package groth16_test

// import (
// 	"time"

// 	ics23 "github.com/cosmos/ics23/go"

// 	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
// 	connectiontypes "github.com/cosmos/ibc-go/v9/modules/core/03-connection/types"
// 	channeltypes "github.com/cosmos/ibc-go/v9/modules/core/04-channel/types"
// 	commitmenttypes "github.com/cosmos/ibc-go/v9/modules/core/23-commitment/types"
// 	host "github.com/cosmos/ibc-go/v9/modules/core/24-host"
// 	"github.com/cosmos/ibc-go/v9/modules/core/exported"
// 	types "github.com/cosmos/ibc-go/v9/modules/light-clients/07-tendermint"
// 	ibctesting "github.com/cosmos/ibc-go/v9/testing"
// 	ibcmock "github.com/cosmos/ibc-go/v9/testing/mock"
// )

// const (
// 	testClientID     = "clientidone"
// 	testConnectionID = "connectionid"
// 	testPortID       = "testportid"
// 	testChannelID    = "testchannelid"
// 	testSequence     = 1

// 	// Do not change the length of these variables
// 	fiftyCharChainID    = "12345678901234567890123456789012345678901234567890"
// 	fiftyOneCharChainID = "123456789012345678901234567890123456789012345678901"
// )

// var invalidProof = []byte("invalid proof")

// func (suite *Groth16TestSuite) TestStatus() {
// 	var (
// 		path         *ibctesting.Path
// 		clientState  *types.ClientState
// 		clientModule *types.LightClientModule
// 	)

// 	testCases := []struct {
// 		name      string
// 		malleate  func()
// 		expStatus exported.Status
// 	}{
// 		{"client is active", func() {}, exported.Active},
// 		{"client is frozen", func() {
// 			clientState.FrozenHeight = clienttypes.NewHeight(0, 1)
// 			path.EndpointA.SetClientState(clientState)
// 		}, exported.Frozen},
// 		{"client status without consensus state", func() {
// 			clientState.LatestHeight = clientState.LatestHeight.Increment().(clienttypes.Height)
// 			path.EndpointA.SetClientState(clientState)
// 		}, exported.Expired},
// 		{"client status is expired", func() {
// 			suite.coordinator.IncrementTimeBy(clientState.TrustingPeriod)
// 		}, exported.Expired},
// 	}

// 	for _, tc := range testCases {
// 		path = ibctesting.NewPath(suite.chainA, suite.chainB)
// 		suite.coordinator.SetupClients(path)

// 		clientStore := suite.chainA.App.GetIBCKeeper().ClientKeeper.ClientStore(suite.chainA.GetContext(), path.EndpointA.ClientID)
// 		clientState = path.EndpointA.GetClientState().(*types.ClientState)

// 		tc.malleate()

// 		status := clientState.Status(suite.chainA.GetContext(), clientStore, suite.chainA.App.AppCodec())
// 		suite.Require().Equal(tc.expStatus, status)

// 	}
// }

// func (suite *Groth16TestSuite) TestValidate() {
// 	testCases := []struct {
// 		name        string
// 		clientState *types.ClientState
// 		expPass     bool
// 	}{
// 		{
// 			name:        "valid client",
// 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// 			expPass:     true,
// 		},
// 		{
// 			name:        "valid client with nil upgrade path",
// 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), nil),
// 			expPass:     true,
// 		},
// 		{
// 			name:        "invalid chainID",
// 			clientState: types.NewClientState("  ", types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// 			expPass:     false,
// 		},
// 		{
// 			// NOTE: if this test fails, the code must account for the change in chainID length across tendermint versions!
// 			// Do not only fix the test, fix the code!
// 			// https://github.com/cosmos/ibc-go/issues/177
// 			name:        "valid chainID - chainID validation failed for chainID of length 50! ",
// 			clientState: types.NewClientState(fiftyCharChainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// 			expPass:     true,
// 		},
// 		{
// 			// NOTE: if this test fails, the code must account for the change in chainID length across tendermint versions!
// 			// Do not only fix the test, fix the code!
// 			// https://github.com/cosmos/ibc-go/issues/177
// 			name:        "invalid chainID - chainID validation did not fail for chainID of length 51! ",
// 			clientState: types.NewClientState(fiftyOneCharChainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// 			expPass:     false,
// 		},
// 		{
// 			name:        "invalid trust level",
// 			clientState: types.NewClientState(chainID, types.Fraction{Numerator: 0, Denominator: 1}, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// 			expPass:     false,
// 		},
// 		{
// 			name:        "invalid zero trusting period",
// 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, 0, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// 			expPass:     false,
// 		},
// 		{
// 			name:        "invalid negative trusting period",
// 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, -1, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// 			expPass:     false,
// 		},
// 		{
// 			name:        "invalid zero unbonding period",
// 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, 0, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// 			expPass:     false,
// 		},
// 		{
// 			name:        "invalid negative unbonding period",
// 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, -1, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// 			expPass:     false,
// 		},
// 		{
// 			name:        "invalid zero max clock drift",
// 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, 0, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// 			expPass:     false,
// 		},
// 		{
// 			name:        "invalid negative max clock drift",
// 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, -1, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// 			expPass:     false,
// 		},
// 		{
// 			name:        "invalid revision number",
// 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, clienttypes.NewHeight(1, 1), commitmenttypes.GetSDKSpecs(), upgradePath),
// 			expPass:     false,
// 		},
// 		{
// 			name:        "invalid revision height",
// 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, clienttypes.ZeroHeight(), commitmenttypes.GetSDKSpecs(), upgradePath),
// 			expPass:     false,
// 		},
// 		{
// 			name:        "trusting period not less than unbonding period",
// 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, ubdPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// 			expPass:     false,
// 		},
// 		{
// 			name:        "proof specs is nil",
// 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, ubdPeriod, ubdPeriod, maxClockDrift, height, nil, upgradePath),
// 			expPass:     false,
// 		},
// 		{
// 			name:        "proof specs contains nil",
// 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, ubdPeriod, ubdPeriod, maxClockDrift, height, []*ics23.ProofSpec{ics23.TendermintSpec, nil}, upgradePath),
// 			expPass:     false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		err := tc.clientState.Validate()
// 		if tc.expPass {
// 			suite.Require().NoError(err, tc.name)
// 		} else {
// 			suite.Require().Error(err, tc.name)
// 		}
// 	}
// }

// func (suite *Groth16TestSuite) TestInitialize() {
// 	testCases := []struct {
// 		name           string
// 		consensusState exported.ConsensusState
// 		clientState    exported.ClientState
// 		expPass        bool
// 	}{
// 		{
// 			name:           "valid consensus",
// 			consensusState: &types.ConsensusState{},
// 			clientState:    &types.ClientState{},
// 			expPass:        true,
// 		},
// 		{
// 			name:           "invalid consensus: consensus state is solomachine consensus",
// 			consensusState: ibctesting.NewSolomachine(suite.T(), suite.chainA.Codec, "solomachine", "", 2).ConsensusState(),
// 			clientState:    ibctesting.NewSolomachine(suite.T(), suite.chainA.Codec, "solomachine", "", 2).ClientState(),
// 			expPass:        false,
// 		},
// 	}

// 	path := ibctesting.NewPath(suite.chainA, suite.chainB)
// 	err := path.EndpointA.CreateClient()
// 	suite.Require().NoError(err)

// 	// TODO: figure out client id
// 	lightClientModule, err := suite.chainA.App.GetIBCKeeper().ClientKeeper.Route(suite.chainA.GetContext(), "clientID")

// 	for _, tc := range testCases {
// 		clientStateBz := suite.chainA.Codec.MustMarshal(tc.clientState)
// 		consensusStateBz := suite.chainA.Codec.MustMarshal(tc.consensusState)
// 		err := lightClientModule.Initialize(suite.chainA.GetContext(), "CLIENDID", clientStateBz, consensusStateBz)
// 		if tc.expPass {
// 			suite.Require().NoError(err, "valid case returned an error")
// 		} else {
// 			suite.Require().Error(err, "invalid case didn't return an error")
// 		}
// 	}
// }

// // func (suite *Groth16TestSuite) TestVerifyClientConsensusState() {
// // 	testCases := []struct {
// // 		name           string
// // 		clientState    *types.ClientState
// // 		consensusState *types.ConsensusState
// // 		prefix         commitmenttypes.MerklePrefix
// // 		proof          []byte
// // 		expPass        bool
// // 	}{
// // 		// FIXME: uncomment
// // 		// {
// // 		// 	name:        "successful verification",
// // 		// 	clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height,  commitmenttypes.GetSDKSpecs()),
// // 		// 	consensusState: types.ConsensusState{
// // 		// 		Root: commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()),
// // 		// 	},
// // 		// 	prefix:  commitmenttypes.NewMerklePrefix([]byte("ibc")),
// // 		// 	expPass: true,
// // 		// },
// // 		{
// // 			name:        "ApplyPrefix failed",
// // 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// // 			consensusState: &types.ConsensusState{
// // 				Root: commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()),
// // 			},
// // 			prefix:  commitmenttypes.MerklePrefix{},
// // 			expPass: false,
// // 		},
// // 		{
// // 			name:        "latest client height < height",
// // 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// // 			consensusState: &types.ConsensusState{
// // 				Root: commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()),
// // 			},
// // 			prefix:  commitmenttypes.NewMerklePrefix([]byte("ibc")),
// // 			expPass: false,
// // 		},
// // 		{
// // 			name:        "proof verification failed",
// // 			clientState: types.NewClientState(chainID, types.DefaultTrustLevel, trustingPeriod, ubdPeriod, maxClockDrift, height, commitmenttypes.GetSDKSpecs(), upgradePath),
// // 			consensusState: &types.ConsensusState{
// // 				Root:               commitmenttypes.NewMerkleRoot(suite.header.Header.GetAppHash()),
// // 				NextValidatorsHash: suite.valsHash,
// // 			},
// // 			prefix:  commitmenttypes.NewMerklePrefix([]byte("ibc")),
// // 			proof:   []byte{},
// // 			expPass: false,
// // 		},
// // 	}

// // 	for i, tc := range testCases {
// // 		tc := tc

// // 		err := tc.clientState.VerifyClientConsensusState(
// // 			nil, suite.cdc, height, "chainA", tc.clientState.LatestHeight, tc.prefix, tc.proof, tc.consensusState,
// // 		)

// // 		if tc.expPass {
// // 			suite.Require().NoError(err, "valid test case %d failed: %s", i, tc.name)
// // 		} else {
// // 			suite.Require().Error(err, "invalid test case %d passed: %s", i, tc.name)
// // 		}
// // 	}
// // }

// // test verification of the connection on chainB being represented in the
// // light client on chainA
// func (suite *Groth16TestSuite) TestVerifyConnectionState() {
// 	var (
// 		clientState *types.ClientState
// 		proof       []byte
// 		proofHeight exported.Height
// 	)

// 	testCases := []struct {
// 		name     string
// 		malleate func()
// 		expPass  bool
// 	}{
// 		{
// 			"successful verification", func() {}, true,
// 		},
// 		{
// 			"latest client height < height", func() {
// 				proofHeight = clientState.LatestHeight.Increment()
// 			}, false,
// 		},
// 		{
// 			"proof verification failed", func() {
// 				proof = invalidProof
// 			}, false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc

// 		suite.Run(tc.name, func() {
// 			suite.SetupTest() // reset

// 			// setup testing conditions
// 			path := ibctesting.NewPath(suite.chainA, suite.chainB)
// 			suite.coordinator.Setup(path)
// 			connection := path.EndpointB.GetConnection()

// 			var ok bool
// 			clientStateI := suite.chainA.GetClientState(path.EndpointA.ClientID)
// 			clientState, ok = clientStateI.(*types.ClientState)
// 			suite.Require().True(ok)

// 			// make connection proof
// 			connectionKey := host.ConnectionKey(path.EndpointB.ConnectionID)
// 			proof, proofHeight = suite.chainB.QueryProof(connectionKey)

// 			tc.malleate() // make changes as necessary

// 			connectionKeeper := suite.chainA.App.GetIBCKeeper().ConnectionKeeper
// 			// TODO: make sure we're setting ChainA and ChainB correctly
// 			err := connectionKeeper.VerifyConnectionState(suite.chainA.GetContext(), connection, proofHeight, proof, path.EndpointB.ConnectionID, path.EndpointB.GetConnection())

// 			if tc.expPass {
// 				suite.Require().NoError(err)
// 			} else {
// 				suite.Require().Error(err)
// 			}
// 		})
// 	}
// }

// // test verification of the channel on chainB being represented in the light
// // client on chainA
// func (suite *Groth16TestSuite) TestVerifyChannelState() {
// 	var (
// 		clientState *types.ClientState
// 		proof       []byte
// 		proofHeight exported.Height
// 	)

// 	testCases := []struct {
// 		name     string
// 		malleate func()
// 		expPass  bool
// 	}{
// 		{
// 			"successful verification", func() {}, true,
// 		},
// 		{
// 			"latest client height < height", func() {
// 				proofHeight = clientState.LatestHeight.Increment()
// 			}, false,
// 		},
// 		{
// 			"proof verification failed", func() {
// 				proof = invalidProof
// 			}, false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc

// 		suite.Run(tc.name, func() {
// 			suite.SetupTest() // reset

// 			// setup testing conditions
// 			path := ibctesting.NewPath(suite.chainA, suite.chainB)
// 			suite.coordinator.Setup(path)
// 			channel := path.EndpointB.GetChannel()

// 			var ok bool
// 			clientStateI := suite.chainA.GetClientState(path.EndpointA.ClientID)
// 			clientState, ok = clientStateI.(*types.ClientState)
// 			suite.Require().True(ok)

// 			// make channel proof
// 			channelKey := host.ChannelKey(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID)
// 			proof, proofHeight = suite.chainB.QueryProof(channelKey)

// 			tc.malleate() // make changes as necessary

// 			// TODO make sure this is supposed to be ChainB
// 			connection := path.EndpointB.GetConnection()
// 			connectionKeeper := suite.chainA.App.GetIBCKeeper().ConnectionKeeper

// 			err := connectionKeeper.VerifyChannelState(suite.chainA.GetContext(), connection, proofHeight, proof, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, channel)

// 			if tc.expPass {
// 				suite.Require().NoError(err)
// 			} else {
// 				suite.Require().Error(err)
// 			}
// 		})
// 	}
// }

// // test verification of the packet commitment on chainB being represented
// // in the light client on chainA. A send from chainB to chainA is simulated.
// func (suite *Groth16TestSuite) TestVerifyPacketCommitment() {
// 	var (
// 		clientState      *types.ClientState
// 		proof            []byte
// 		delayTimePeriod  uint64
// 		delayBlockPeriod uint64
// 		proofHeight      exported.Height
// 	)

// 	testCases := []struct {
// 		name     string
// 		malleate func()
// 		expPass  bool
// 	}{
// 		{
// 			"successful verification", func() {}, true,
// 		},
// 		{
// 			name: "delay time period has passed",
// 			malleate: func() {
// 				delayTimePeriod = uint64(time.Second.Nanoseconds())
// 			},
// 			expPass: true,
// 		},
// 		{
// 			name: "delay time period has not passed",
// 			malleate: func() {
// 				delayTimePeriod = uint64(time.Hour.Nanoseconds())
// 			},
// 			expPass: false,
// 		},
// 		{
// 			name: "delay block period has passed",
// 			malleate: func() {
// 				delayBlockPeriod = 1
// 			},
// 			expPass: true,
// 		},
// 		{
// 			name: "delay block period has not passed",
// 			malleate: func() {
// 				delayBlockPeriod = 1000
// 			},
// 			expPass: false,
// 		},
// 		{
// 			"latest client height < height", func() {
// 				proofHeight = clientState.LatestHeight.Increment()
// 			}, false,
// 		},
// 		{
// 			"proof verification failed", func() {
// 				proof = invalidProof
// 			}, false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc

// 		suite.Run(tc.name, func() {
// 			suite.SetupTest() // reset

// 			// setup testing conditions
// 			path := ibctesting.NewPath(suite.chainA, suite.chainB)
// 			suite.coordinator.Setup(path)
// 			timeoutHeight := clienttypes.NewHeight(0, 100)

// 			// send packet
// 			sequence, err := path.EndpointB.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
// 			suite.Require().NoError(err)

// 			var ok bool
// 			packet := channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, timeoutHeight, 0)
// 			clientStateI := suite.chainA.GetClientState(path.EndpointA.ClientID)
// 			clientState, ok = clientStateI.(*types.ClientState)
// 			suite.Require().True(ok)

// 			// make packet commitment proof
// 			packetKey := host.PacketCommitmentKey(packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
// 			proof, proofHeight = path.EndpointB.QueryProof(packetKey)

// 			// reset time and block delays to 0, malleate may change to a specific non-zero value.
// 			delayTimePeriod = 0
// 			delayBlockPeriod = 0
// 			tc.malleate() // make changes as necessary

// 			ctx := suite.chainA.GetContext()

// 			// TODO: check this is correct
// 			connection := path.EndpointA.GetConnection()
// 			connection.DelayPeriod = delayTimePeriod

// 			// set time per block param
// 			if delayBlockPeriod != 0 {
// 				suite.chainA.App.GetIBCKeeper().ConnectionKeeper.SetParams(suite.chainA.GetContext(), connectiontypes.NewParams(delayBlockPeriod))
// 			}

// 			commitment := channeltypes.CommitPacket(suite.chainA.App.GetIBCKeeper().Codec(), packet)
// 			// TODO: check this too
// 			err = suite.chainB.App.GetIBCKeeper().ConnectionKeeper.VerifyPacketCommitment(ctx, connection, proofHeight, proof, packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence(), commitment)

// 			if tc.expPass {
// 				suite.Require().NoError(err)
// 			} else {
// 				suite.Require().Error(err)
// 			}
// 		})
// 	}
// }

// // test verification of the acknowledgement on chainB being represented
// // in the light client on chainA. A send and ack from chainA to chainB
// // is simulated.
// func (suite *Groth16TestSuite) TestVerifyPacketAcknowledgement() {
// 	var (
// 		clientState      *types.ClientState
// 		proof            []byte
// 		delayTimePeriod  uint64
// 		delayBlockPeriod uint64
// 		proofHeight      exported.Height
// 	)

// 	testCases := []struct {
// 		name     string
// 		malleate func()
// 		expPass  bool
// 	}{
// 		{
// 			"successful verification", func() {}, true,
// 		},
// 		{
// 			name: "delay time period has passed",
// 			malleate: func() {
// 				delayTimePeriod = uint64(time.Second.Nanoseconds())
// 			},
// 			expPass: true,
// 		},
// 		{
// 			name: "delay time period has not passed",
// 			malleate: func() {
// 				delayTimePeriod = uint64(time.Hour.Nanoseconds())
// 			},
// 			expPass: false,
// 		},
// 		{
// 			name: "delay block period has passed",
// 			malleate: func() {
// 				delayBlockPeriod = 1
// 			},
// 			expPass: true,
// 		},
// 		{
// 			name: "delay block period has not passed",
// 			malleate: func() {
// 				delayBlockPeriod = 10
// 			},
// 			expPass: false,
// 		},
// 		{
// 			"latest client height < height", func() {
// 				proofHeight = clientState.LatestHeight.Increment()
// 			}, false,
// 		},
// 		{
// 			"proof verification failed", func() {
// 				proof = invalidProof
// 			}, false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc

// 		suite.Run(tc.name, func() {
// 			suite.SetupTest() // reset

// 			// setup testing conditions
// 			path := ibctesting.NewPath(suite.chainA, suite.chainB)
// 			suite.coordinator.Setup(path)
// 			timeoutHeight := clienttypes.NewHeight(0, 100)

// 			// send packet
// 			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
// 			suite.Require().NoError(err)

// 			// write receipt and ack
// 			packet := channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
// 			err = path.EndpointB.RecvPacket(packet)
// 			suite.Require().NoError(err)

// 			var ok bool
// 			clientStateI := suite.chainA.GetClientState(path.EndpointA.ClientID)
// 			clientState, ok = clientStateI.(*types.ClientState)
// 			suite.Require().True(ok)

// 			// make packet acknowledgement proof
// 			acknowledgementKey := host.PacketAcknowledgementKey(packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())
// 			proof, proofHeight = suite.chainB.QueryProof(acknowledgementKey)

// 			// reset time and block delays to 0, malleate may change to a specific non-zero value.
// 			delayTimePeriod = 0
// 			delayBlockPeriod = 0
// 			tc.malleate() // make changes as necessary

// 			ctx := suite.chainA.GetContext()
// 			// TODO: make sure this is correct
// 			connection := path.EndpointA.GetConnection()
// 			connection.DelayPeriod = delayTimePeriod

// 			// set time per block param
// 			if delayBlockPeriod != 0 {
// 				suite.chainA.App.GetIBCKeeper().ConnectionKeeper.SetParams(suite.chainA.GetContext(), connectiontypes.NewParams(delayBlockPeriod))
// 			}

// 			// TODO: verify ChainA and ChainB are set correctly
// 			connectionKeeper := suite.chainA.App.GetIBCKeeper().ConnectionKeeper
// 			err = connectionKeeper.VerifyPacketAcknowledgement(ctx, connection, proofHeight, proof, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, packet.GetSequence(), ibcmock.MockAcknowledgement.Acknowledgement())
// 			if tc.expPass {
// 				suite.Require().NoError(err)
// 			} else {
// 				suite.Require().Error(err)
// 			}
// 		})
// 	}
// }

// // test verification of the absent acknowledgement on chainB being represented
// // in the light client on chainA. A send from chainB to chainA is simulated, but
// // no receive.
// func (suite *Groth16TestSuite) TestVerifyPacketReceiptAbsence() {
// 	var (
// 		clientState      *types.ClientState
// 		proof            []byte
// 		delayTimePeriod  uint64
// 		delayBlockPeriod uint64
// 		proofHeight      exported.Height
// 	)

// 	testCases := []struct {
// 		name     string
// 		malleate func()
// 		expPass  bool
// 	}{
// 		{
// 			"successful verification", func() {}, true,
// 		},
// 		{
// 			name: "delay time period has passed",
// 			malleate: func() {
// 				delayTimePeriod = uint64(time.Second.Nanoseconds())
// 			},
// 			expPass: true,
// 		},
// 		{
// 			name: "delay time period has not passed",
// 			malleate: func() {
// 				delayTimePeriod = uint64(time.Hour.Nanoseconds())
// 			},
// 			expPass: false,
// 		},
// 		{
// 			name: "delay block period has passed",
// 			malleate: func() {
// 				delayBlockPeriod = 1
// 			},
// 			expPass: true,
// 		},
// 		{
// 			name: "delay block period has not passed",
// 			malleate: func() {
// 				delayBlockPeriod = 10
// 			},
// 			expPass: false,
// 		},
// 		{
// 			"latest client height < height", func() {
// 				proofHeight = clientState.LatestHeight.Increment()
// 			}, false,
// 		},
// 		{
// 			"proof verification failed", func() {
// 				proof = invalidProof
// 			}, false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc

// 		suite.Run(tc.name, func() {
// 			suite.SetupTest() // reset

// 			// setup testing conditions
// 			path := ibctesting.NewPath(suite.chainA, suite.chainB)
// 			suite.coordinator.Setup(path)
// 			timeoutHeight := clienttypes.NewHeight(0, 100)

// 			// send packet, but no recv
// 			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
// 			suite.Require().NoError(err)

// 			var ok bool
// 			packet := channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
// 			clientStateI := suite.chainA.GetClientState(path.EndpointA.ClientID)
// 			clientState, ok = clientStateI.(*types.ClientState)
// 			suite.Require().True(ok)

// 			// make packet receipt absence proof
// 			receiptKey := host.PacketReceiptKey(packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence())
// 			proof, proofHeight = path.EndpointB.QueryProof(receiptKey)

// 			// reset time and block delays to 0, malleate may change to a specific non-zero value.
// 			delayTimePeriod = 0
// 			delayBlockPeriod = 0
// 			tc.malleate() // make changes as necessary

// 			ctx := suite.chainA.GetContext()

// 			connection := path.EndpointA.GetConnection()
// 			// TODO: come back to this setup
// 			connection.DelayPeriod = delayTimePeriod

// 			// set time per block param
// 			// TODO check is this is correct
// 			if delayBlockPeriod != 0 {
// 				suite.chainA.App.GetIBCKeeper().ConnectionKeeper.SetParams(suite.chainA.GetContext(), connectiontypes.NewParams(delayBlockPeriod))
// 			}

// 			// TODO double check with chainA and chainB
// 			connectionKeeper := suite.chainA.App.GetIBCKeeper().ConnectionKeeper
// 			err = connectionKeeper.VerifyPacketReceiptAbsence(ctx, connection, proofHeight, proof, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, packet.GetSequence())

// 			if tc.expPass {
// 				suite.Require().NoError(err)
// 			} else {
// 				suite.Require().Error(err)
// 			}
// 		})
// 	}
// }

// // test verification of the next receive sequence on chainB being represented
// // in the light client on chainA. A send and receive from chainB to chainA is
// // simulated.
// func (suite *Groth16TestSuite) TestVerifyNextSeqRecv() {
// 	var (
// 		clientState      *types.ClientState
// 		proof            []byte
// 		delayTimePeriod  uint64
// 		delayBlockPeriod uint64
// 		proofHeight      exported.Height
// 	)

// 	testCases := []struct {
// 		name     string
// 		malleate func()
// 		expPass  bool
// 	}{
// 		{
// 			"successful verification", func() {}, true,
// 		},
// 		{
// 			name: "delay time period has passed",
// 			malleate: func() {
// 				delayTimePeriod = uint64(time.Second.Nanoseconds())
// 			},
// 			expPass: true,
// 		},
// 		{
// 			name: "delay time period has not passed",
// 			malleate: func() {
// 				delayTimePeriod = uint64(time.Hour.Nanoseconds())
// 			},
// 			expPass: false,
// 		},
// 		{
// 			name: "delay block period has passed",
// 			malleate: func() {
// 				delayBlockPeriod = 1
// 			},
// 			expPass: true,
// 		},
// 		{
// 			name: "delay block period has not passed",
// 			malleate: func() {
// 				delayBlockPeriod = 10
// 			},
// 			expPass: false,
// 		},

// 		{
// 			"latest client height < height", func() {
// 				proofHeight = clientState.LatestHeight.Increment()
// 			}, false,
// 		},
// 		{
// 			"proof verification failed", func() {
// 				proof = invalidProof
// 			}, false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc

// 		suite.Run(tc.name, func() {
// 			suite.SetupTest() // reset

// 			// setup testing conditions
// 			path := ibctesting.NewPath(suite.chainA, suite.chainB)
// 			path.SetChannelOrdered()
// 			suite.coordinator.Setup(path)
// 			timeoutHeight := clienttypes.NewHeight(0, 100)

// 			// send packet
// 			sequence, err := path.EndpointA.SendPacket(timeoutHeight, 0, ibctesting.MockPacketData)
// 			suite.Require().NoError(err)

// 			// next seq recv incremented
// 			packet := channeltypes.NewPacket(ibctesting.MockPacketData, sequence, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, timeoutHeight, 0)
// 			err = path.EndpointB.RecvPacket(packet)
// 			suite.Require().NoError(err)

// 			var ok bool
// 			clientStateI := suite.chainA.GetClientState(path.EndpointA.ClientID)
// 			clientState, ok = clientStateI.(*types.ClientState)
// 			suite.Require().True(ok)

// 			// make next seq recv proof
// 			nextSeqRecvKey := host.NextSequenceRecvKey(packet.GetDestPort(), packet.GetDestChannel())
// 			proof, proofHeight = suite.chainB.QueryProof(nextSeqRecvKey)

// 			// reset time and block delays to 0, malleate may change to a specific non-zero value.
// 			delayTimePeriod = 0
// 			delayBlockPeriod = 0
// 			tc.malleate() // make changes as necessary

// 			ctx := suite.chainA.GetContext()
// 			// TODO: check if correct
// 			connection := path.EndpointA.GetConnection()
// 			connection.DelayPeriod = delayTimePeriod

// 			// set time per block param
// 			if delayBlockPeriod != 0 {
// 				suite.chainA.App.GetIBCKeeper().ConnectionKeeper.SetParams(suite.chainA.GetContext(), connectiontypes.NewParams(delayBlockPeriod))
// 			}
// 			err = suite.chainA.App.GetIBCKeeper().ConnectionKeeper.VerifyNextSequenceRecv(
// 				ctx, connection, proofHeight, proof,
// 				packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence()+1,
// 			)

// 			if tc.expPass {
// 				suite.Require().NoError(err)
// 			} else {
// 				suite.Require().Error(err)
// 			}
// 		})
// 	}
// }
