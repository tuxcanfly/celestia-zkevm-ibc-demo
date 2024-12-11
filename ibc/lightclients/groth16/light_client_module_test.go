package groth16_test

import (
	"crypto/sha256"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v9/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v9/modules/core/24-host"
	ibcerrors "github.com/cosmos/ibc-go/v9/modules/core/errors"
	"github.com/cosmos/ibc-go/v9/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v9/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v9/testing"
)

var (
	tmClientID          = clienttypes.FormatClientIdentifier(exported.Tendermint, 100)
	solomachineClientID = clienttypes.FormatClientIdentifier(exported.Solomachine, 0)
)

func (suite *Groth16TestSuite) TestInitialize() {
	var consensusState exported.ConsensusState
	var clientState exported.ClientState

	testCases := []struct {
		name     string
		malleate func()
		expErr   error
	}{
		{
			"valid consensus & client states",
			func() {},
			nil,
		},
		{
			"invalid client state",
			func() {
				clientState.(*ibctm.ClientState).ChainId = ""
			},
			ibctm.ErrInvalidChainID,
		},
		{
			"invalid client state: solomachine client state",
			func() {
				clientState = ibctesting.NewSolomachine(suite.T(), suite.chainA.Codec, "solomachine", "", 2).ClientState()
			},
			fmt.Errorf("failed to unmarshal client state bytes into client state"),
		},
		{
			"invalid consensus: consensus state is solomachine consensus",
			func() {
				consensusState = ibctesting.NewSolomachine(suite.T(), suite.chainA.Codec, "solomachine", "", 2).ConsensusState()
			},
			fmt.Errorf("failed to unmarshal consensus state bytes into consensus state"),
		},
		{
			"invalid consensus state",
			func() {
				consensusState = ibctm.NewConsensusState(time.Now(), commitmenttypes.NewMerkleRoot([]byte(ibctm.SentinelRoot)), []byte("invalidNextValsHash"))
			},
			fmt.Errorf("next validators hash is invalid"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			path := ibctesting.NewPath(suite.chainA, suite.chainB)

			tmConfig, ok := path.EndpointA.ClientConfig.(*ibctesting.TendermintConfig)
			suite.Require().True(ok)

			clientState = ibctm.NewClientState(
				path.EndpointA.Chain.ChainID,
				tmConfig.TrustLevel, tmConfig.TrustingPeriod, tmConfig.UnbondingPeriod, tmConfig.MaxClockDrift,
				suite.chainA.LatestCommittedHeader.GetHeight().(clienttypes.Height), commitmenttypes.GetSDKSpecs(), ibctesting.UpgradePath,
			)

			consensusState = ibctm.NewConsensusState(time.Now(), commitmenttypes.NewMerkleRoot([]byte(ibctm.SentinelRoot)), suite.chainA.ProposedHeader.ValidatorsHash)

			clientID := suite.chainA.App.GetIBCKeeper().ClientKeeper.GenerateClientIdentifier(suite.chainA.GetContext(), clientState.ClientType())

			lightClientModule, err := suite.chainA.App.GetIBCKeeper().ClientKeeper.Route(suite.chainA.GetContext(), clientID)
			suite.Require().NoError(err)

			tc.malleate()

			clientStateBz := suite.chainA.Codec.MustMarshal(clientState)
			consStateBz := suite.chainA.Codec.MustMarshal(consensusState)

			err = lightClientModule.Initialize(suite.chainA.GetContext(), path.EndpointA.ClientID, clientStateBz, consStateBz)

			store := suite.chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointA.ClientID)

			expPass := tc.expErr == nil
			if expPass {
				suite.Require().NoError(err, "valid case returned an error")
				suite.Require().True(store.Has(host.ClientStateKey()))
				suite.Require().True(store.Has(host.ConsensusStateKey(suite.chainB.LatestCommittedHeader.GetHeight())))
			} else {
				suite.Require().ErrorContains(err, tc.expErr.Error())
				suite.Require().False(store.Has(host.ClientStateKey()))
				suite.Require().False(store.Has(host.ConsensusStateKey(suite.chainB.LatestCommittedHeader.GetHeight())))
			}
		})
	}
}

func (suite *Groth16TestSuite) TestVerifyClientMessage() {
	var path *ibctesting.Path

	testCases := []struct {
		name     string
		malleate func()
		expErr   error
	}{
		{
			"success",
			func() {},
			nil,
		},
		{
			"failure: client state not found",
			func() {
				store := suite.chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointA.ClientID)
				store.Delete(host.ClientStateKey())
			},
			clienttypes.ErrClientNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			path.SetupClients()

			lightClientModule, err := suite.chainA.App.GetIBCKeeper().ClientKeeper.Route(suite.chainA.GetContext(), path.EndpointA.ClientID)
			suite.Require().NoError(err)

			// ensure counterparty state is committed
			suite.coordinator.CommitBlock(suite.chainB)
			trustedHeight, ok := path.EndpointA.GetClientLatestHeight().(clienttypes.Height)
			suite.Require().True(ok)
			header, err := path.EndpointA.Counterparty.Chain.IBCClientHeader(path.EndpointA.Counterparty.Chain.LatestCommittedHeader, trustedHeight)
			suite.Require().NoError(err)

			tc.malleate()

			err = lightClientModule.VerifyClientMessage(suite.chainA.GetContext(), path.EndpointA.ClientID, header)

			expPass := tc.expErr == nil
			if expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().ErrorIs(err, tc.expErr)
			}
		})
	}
}

func (suite *Groth16TestSuite) TestCheckForMisbehaviourPanicsOnClientStateNotFound() {
	suite.SetupTest()

	path := ibctesting.NewPath(suite.chainA, suite.chainB)
	path.SetupClients()

	lightClientModule, err := suite.chainA.App.GetIBCKeeper().ClientKeeper.Route(suite.chainA.GetContext(), path.EndpointA.ClientID)
	suite.Require().NoError(err)

	// ensure counterparty state is committed
	suite.coordinator.CommitBlock(suite.chainB)
	trustedHeight, ok := path.EndpointA.GetClientLatestHeight().(clienttypes.Height)
	suite.Require().True(ok)
	header, err := path.EndpointA.Counterparty.Chain.IBCClientHeader(path.EndpointA.Counterparty.Chain.LatestCommittedHeader, trustedHeight)
	suite.Require().NoError(err)

	// delete client state
	store := suite.chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointA.ClientID)
	store.Delete(host.ClientStateKey())

	suite.Require().PanicsWithError(errorsmod.Wrap(clienttypes.ErrClientNotFound, path.EndpointA.ClientID).Error(),
		func() {
			lightClientModule.CheckForMisbehaviour(suite.chainA.GetContext(), path.EndpointA.ClientID, header)
		},
	)
}

func (suite *Groth16TestSuite) TestUpdateStatePanicsOnClientStateNotFound() {
	suite.SetupTest()

	path := ibctesting.NewPath(suite.chainA, suite.chainB)
	path.SetupClients()

	lightClientModule, err := suite.chainA.App.GetIBCKeeper().ClientKeeper.Route(suite.chainA.GetContext(), path.EndpointA.ClientID)
	suite.Require().NoError(err)

	// ensure counterparty state is committed
	suite.coordinator.CommitBlock(suite.chainB)
	trustedHeight, ok := path.EndpointA.GetClientLatestHeight().(clienttypes.Height)
	suite.Require().True(ok)
	header, err := path.EndpointA.Counterparty.Chain.IBCClientHeader(path.EndpointA.Counterparty.Chain.LatestCommittedHeader, trustedHeight)
	suite.Require().NoError(err)

	// delete client state
	store := suite.chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointA.ClientID)
	store.Delete(host.ClientStateKey())

	suite.Require().PanicsWithError(
		errorsmod.Wrap(clienttypes.ErrClientNotFound, path.EndpointA.ClientID).Error(),
		func() {
			lightClientModule.UpdateState(suite.chainA.GetContext(), path.EndpointA.ClientID, header)
		},
	)
}

func (suite *Groth16TestSuite) TestStatus() {
	var (
		path        *ibctesting.Path
		clientState *ibctm.ClientState
	)

	testCases := []struct {
		name      string
		malleate  func()
		expStatus exported.Status
	}{
		{
			"client is active",
			func() {},
			exported.Active,
		},
		{
			"client is frozen",
			func() {
				clientState.FrozenHeight = clienttypes.NewHeight(0, 1)
				path.EndpointA.SetClientState(clientState)
			},
			exported.Frozen,
		},
		{
			"client status without consensus state",
			func() {
				newLatestHeight, ok := clientState.LatestHeight.Increment().(clienttypes.Height)
				suite.Require().True(ok)
				clientState.LatestHeight = newLatestHeight
				path.EndpointA.SetClientState(clientState)
			},
			exported.Expired,
		},
		{
			"client status is expired",
			func() {
				suite.coordinator.IncrementTimeBy(clientState.TrustingPeriod)
			},
			exported.Expired,
		},
		{
			"client state not found",
			func() {
				store := suite.chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointA.ClientID)
				store.Delete(host.ClientStateKey())
			},
			exported.Unknown,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			path.SetupClients()

			lightClientModule, err := suite.chainA.App.GetIBCKeeper().ClientKeeper.Route(suite.chainA.GetContext(), path.EndpointA.ClientID)
			suite.Require().NoError(err)

			var ok bool
			clientState, ok = path.EndpointA.GetClientState().(*ibctm.ClientState)
			suite.Require().True(ok)

			tc.malleate()

			status := lightClientModule.Status(suite.chainA.GetContext(), path.EndpointA.ClientID)
			suite.Require().Equal(tc.expStatus, status)
		})

	}
}

func (suite *Groth16TestSuite) TestLatestHeight() {
	var (
		path   *ibctesting.Path
		height exported.Height
	)

	testCases := []struct {
		name      string
		malleate  func()
		expHeight exported.Height
	}{
		{
			"success",
			func() {},
			clienttypes.Height{RevisionNumber: 0x1, RevisionHeight: 0x4},
		},
		{
			"client state not found",
			func() {
				store := suite.chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointA.ClientID)
				store.Delete(host.ClientStateKey())
			},
			clienttypes.ZeroHeight(),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			path.SetupClients()

			lightClientModule, err := suite.chainA.App.GetIBCKeeper().ClientKeeper.Route(suite.chainA.GetContext(), path.EndpointA.ClientID)
			suite.Require().NoError(err)

			tc.malleate()

			height = lightClientModule.LatestHeight(suite.chainA.GetContext(), path.EndpointA.ClientID)
			suite.Require().Equal(tc.expHeight, height)
		})
	}
}

func (suite *Groth16TestSuite) TestGetTimestampAtHeight() {
	var (
		path   *ibctesting.Path
		height exported.Height
	)
	expectedTimestamp := time.Unix(1, 0)

	testCases := []struct {
		name     string
		malleate func()
		expErr   error
	}{
		{
			"success",
			func() {},
			nil,
		},
		{
			"failure: client state not found",
			func() {
				store := suite.chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointA.ClientID)
				store.Delete(host.ClientStateKey())
			},
			clienttypes.ErrClientNotFound,
		},
		{
			"failure: consensus state not found for height",
			func() {
				clientState, ok := path.EndpointA.GetClientState().(*ibctm.ClientState)
				suite.Require().True(ok)
				height = clientState.LatestHeight.Increment()
			},
			clienttypes.ErrConsensusStateNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			path.SetupClients()

			clientState, ok := path.EndpointA.GetClientState().(*ibctm.ClientState)
			suite.Require().True(ok)
			height = clientState.LatestHeight

			// grab consensusState from store and update with a predefined timestamp
			consensusState := path.EndpointA.GetConsensusState(height)
			tmConsensusState, ok := consensusState.(*ibctm.ConsensusState)
			suite.Require().True(ok)

			tmConsensusState.Timestamp = expectedTimestamp
			path.EndpointA.SetConsensusState(tmConsensusState, height)

			lightClientModule, err := suite.chainA.App.GetIBCKeeper().ClientKeeper.Route(suite.chainA.GetContext(), path.EndpointA.ClientID)
			suite.Require().NoError(err)

			tc.malleate()

			timestamp, err := lightClientModule.TimestampAtHeight(suite.chainA.GetContext(), path.EndpointA.ClientID, height)

			expPass := tc.expErr == nil
			if expPass {
				suite.Require().NoError(err)

				expectedTimestamp := uint64(expectedTimestamp.UnixNano())
				suite.Require().Equal(expectedTimestamp, timestamp)
			} else {
				suite.Require().ErrorIs(err, tc.expErr)
			}
		})
	}
}

func (suite *Groth16TestSuite) TestRecoverClient() {
	var (
		subjectClientID, substituteClientID string
		subjectClientState                  exported.ClientState
	)

	testCases := []struct {
		name     string
		malleate func()
		expErr   error
	}{
		{
			"success",
			func() {
			},
			nil,
		},
		{
			"cannot parse malformed substitute client ID",
			func() {
				substituteClientID = ibctesting.InvalidID
			},
			host.ErrInvalidID,
		},
		{
			"substitute client ID does not contain 07-tendermint prefix",
			func() {
				substituteClientID = solomachineClientID
			},
			clienttypes.ErrInvalidClientType,
		},
		{
			"cannot find subject client state",
			func() {
				subjectClientID = tmClientID
			},
			clienttypes.ErrClientNotFound,
		},
		{
			"cannot find substitute client state",
			func() {
				substituteClientID = tmClientID
			},
			clienttypes.ErrClientNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			ctx := suite.chainA.GetContext()

			subjectPath := ibctesting.NewPath(suite.chainA, suite.chainB)
			subjectPath.SetupClients()
			subjectClientID = subjectPath.EndpointA.ClientID
			subjectClientState = suite.chainA.GetClientState(subjectClientID)

			substitutePath := ibctesting.NewPath(suite.chainA, suite.chainB)
			substitutePath.SetupClients()
			substituteClientID = substitutePath.EndpointA.ClientID

			tmClientState, ok := subjectClientState.(*ibctm.ClientState)
			suite.Require().True(ok)
			tmClientState.FrozenHeight = tmClientState.LatestHeight
			suite.chainA.App.GetIBCKeeper().ClientKeeper.SetClientState(ctx, subjectPath.EndpointA.ClientID, tmClientState)

			lightClientModule, err := suite.chainA.App.GetIBCKeeper().ClientKeeper.Route(suite.chainA.GetContext(), subjectClientID)
			suite.Require().NoError(err)

			tc.malleate()

			err = lightClientModule.RecoverClient(ctx, subjectClientID, substituteClientID)

			expPass := tc.expErr == nil
			if expPass {
				suite.Require().NoError(err)

				// assert that status of subject client is now Active
				lightClientModule, err := suite.chainA.App.GetIBCKeeper().ClientKeeper.Route(suite.chainA.GetContext(), subjectClientID)
				suite.Require().NoError(err)
				suite.Require().Equal(lightClientModule.Status(suite.chainA.GetContext(), subjectClientID), exported.Active)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, tc.expErr)
			}
		})
	}
}

func (suite *Groth16TestSuite) TestVerifyUpgradeAndUpdateState() {
	var (
		clientID                                              string
		path                                                  *ibctesting.Path
		upgradedClientState                                   exported.ClientState
		upgradedClientStateAny, upgradedConsensusStateAny     *codectypes.Any
		upgradedClientStateProof, upgradedConsensusStateProof []byte
	)

	testCases := []struct {
		name     string
		malleate func()
		expErr   error
	}{
		{
			"success",
			func() {
				// upgrade height is at next block
				upgradeHeight := clienttypes.NewHeight(0, uint64(suite.chainB.GetContext().BlockHeight()+1))

				// zero custom fields and store in upgrade store
				zeroedUpgradedClient := upgradedClientState.(*ibctm.ClientState).ZeroCustomFields()
				zeroedUpgradedClientAny, err := codectypes.NewAnyWithValue(zeroedUpgradedClient)
				suite.Require().NoError(err)

				err = suite.chainB.GetSimApp().UpgradeKeeper.SetUpgradedClient(suite.chainB.GetContext(), int64(upgradeHeight.GetRevisionHeight()), suite.chainB.Codec.MustMarshal(zeroedUpgradedClientAny))
				suite.Require().NoError(err)

				err = suite.chainB.GetSimApp().UpgradeKeeper.SetUpgradedConsensusState(suite.chainB.GetContext(), int64(upgradeHeight.GetRevisionHeight()), suite.chainB.Codec.MustMarshal(upgradedConsensusStateAny))
				suite.Require().NoError(err)

				// commit upgrade store changes and update clients
				suite.coordinator.CommitBlock(suite.chainB)
				err = path.EndpointA.UpdateClient()
				suite.Require().NoError(err)

				upgradedClientStateProof, _ = suite.chainB.QueryUpgradeProof(upgradetypes.UpgradedClientKey(int64(upgradeHeight.GetRevisionHeight())), path.EndpointA.GetClientLatestHeight().GetRevisionHeight())
				upgradedConsensusStateProof, _ = suite.chainB.QueryUpgradeProof(upgradetypes.UpgradedConsStateKey(int64(upgradeHeight.GetRevisionHeight())), path.EndpointA.GetClientLatestHeight().GetRevisionHeight())
			},
			nil,
		},
		{
			"cannot find client state",
			func() {
				clientID = tmClientID
			},
			clienttypes.ErrClientNotFound,
		},
		{
			"upgraded client state is not for tendermint client state",
			func() {
				upgradedClientStateAny = &codectypes.Any{
					Value: []byte("invalid client state bytes"),
				}
			},
			clienttypes.ErrInvalidClient,
		},
		{
			"upgraded consensus state is not tendermint consensus state",
			func() {
				upgradedConsensusStateAny = &codectypes.Any{
					Value: []byte("invalid consensus state bytes"),
				}
			},
			clienttypes.ErrInvalidConsensus,
		},
		{
			"upgraded client state height is not greater than current height",
			func() {
				// upgrade height is at next block
				upgradeHeight := clienttypes.NewHeight(1, uint64(suite.chainB.GetContext().BlockHeight()+1))

				// zero custom fields and store in upgrade store
				zeroedUpgradedClient := upgradedClientState.(*ibctm.ClientState).ZeroCustomFields()
				zeroedUpgradedClientAny, err := codectypes.NewAnyWithValue(zeroedUpgradedClient)
				suite.Require().NoError(err)

				err = suite.chainB.GetSimApp().UpgradeKeeper.SetUpgradedClient(suite.chainB.GetContext(), int64(upgradeHeight.GetRevisionHeight()), suite.chainB.Codec.MustMarshal(zeroedUpgradedClientAny))
				suite.Require().NoError(err)

				err = suite.chainB.GetSimApp().UpgradeKeeper.SetUpgradedConsensusState(suite.chainB.GetContext(), int64(upgradeHeight.GetRevisionHeight()), suite.chainB.Codec.MustMarshal(upgradedConsensusStateAny))
				suite.Require().NoError(err)

				// change upgraded client state height to be lower than current client state height
				tmClient, ok := upgradedClientState.(*ibctm.ClientState)
				suite.Require().True(ok)

				newLatestheight, ok := path.EndpointA.GetClientLatestHeight().Decrement()
				suite.Require().True(ok)

				tmClient.LatestHeight, ok = newLatestheight.(clienttypes.Height)
				suite.Require().True(ok)
				upgradedClientStateAny, err = codectypes.NewAnyWithValue(tmClient)
				suite.Require().NoError(err)

				suite.coordinator.CommitBlock(suite.chainB)
				err = path.EndpointA.UpdateClient()
				suite.Require().NoError(err)

				upgradedClientStateProof, _ = suite.chainB.QueryUpgradeProof(upgradetypes.UpgradedClientKey(int64(upgradeHeight.GetRevisionHeight())), path.EndpointA.GetClientLatestHeight().GetRevisionHeight())
				upgradedConsensusStateProof, _ = suite.chainB.QueryUpgradeProof(upgradetypes.UpgradedConsStateKey(int64(upgradeHeight.GetRevisionHeight())), path.EndpointA.GetClientLatestHeight().GetRevisionHeight())
			},
			ibcerrors.ErrInvalidHeight,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			path = ibctesting.NewPath(suite.chainA, suite.chainB)
			path.SetupClients()

			clientID = path.EndpointA.ClientID
			clientState, ok := path.EndpointA.GetClientState().(*ibctm.ClientState)
			suite.Require().True(ok)
			revisionNumber := clienttypes.ParseChainID(clientState.ChainId)

			newUnbondindPeriod := ubdPeriod + trustingPeriod
			newChainID, err := clienttypes.SetRevisionNumber(clientState.ChainId, revisionNumber+1)
			suite.Require().NoError(err)

			upgradedClientState = ibctm.NewClientState(newChainID, ibctm.DefaultTrustLevel, trustingPeriod, newUnbondindPeriod, maxClockDrift, clienttypes.NewHeight(revisionNumber+1, clientState.LatestHeight.GetRevisionHeight()+1), commitmenttypes.GetSDKSpecs(), upgradePath)
			upgradedClientStateAny, err = codectypes.NewAnyWithValue(upgradedClientState)
			suite.Require().NoError(err)

			nextValsHash := sha256.Sum256([]byte("new-nextValsHash"))
			upgradedConsensusState := ibctm.NewConsensusState(time.Now(), commitmenttypes.NewMerkleRoot([]byte("new-hash")), nextValsHash[:])

			upgradedConsensusStateAny, err = codectypes.NewAnyWithValue(upgradedConsensusState)
			suite.Require().NoError(err)

			tc.malleate()

			lightClientModule, err := suite.chainA.App.GetIBCKeeper().ClientKeeper.Route(suite.chainA.GetContext(), clientID)
			suite.Require().NoError(err)

			err = lightClientModule.VerifyUpgradeAndUpdateState(
				suite.chainA.GetContext(),
				clientID,
				upgradedClientStateAny.Value,
				upgradedConsensusStateAny.Value,
				upgradedClientStateProof,
				upgradedConsensusStateProof,
			)

			expPass := tc.expErr == nil
			if expPass {
				suite.Require().NoError(err)

				expClientState := path.EndpointA.GetClientState()
				expClientStateBz := suite.chainA.Codec.MustMarshal(expClientState)
				suite.Require().Equal(upgradedClientStateAny.Value, expClientStateBz)

				expConsensusState := ibctm.NewConsensusState(upgradedConsensusState.Timestamp, commitmenttypes.NewMerkleRoot([]byte(ibctm.SentinelRoot)), upgradedConsensusState.NextValidatorsHash)
				expConsensusStateBz := suite.chainA.Codec.MustMarshal(expConsensusState)

				consensusStateBz := suite.chainA.Codec.MustMarshal(path.EndpointA.GetConsensusState(path.EndpointA.GetClientLatestHeight()))
				suite.Require().Equal(expConsensusStateBz, consensusStateBz)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, tc.expErr)
			}
		})
	}
}
