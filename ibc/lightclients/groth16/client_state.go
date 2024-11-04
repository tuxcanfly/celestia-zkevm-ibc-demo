package groth16

import (
	"bytes"
	"context"
	fmt "fmt"

	sdkerrors "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/celestiaorg/celestia-zkevm-ibc-demo/ibc/mpt"
	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"

	// connectiontypes "github.com/cosmos/ibc-go/v9/modules/core/03-connection"
	commitmenttypes "github.com/cosmos/ibc-go/v9/modules/core/23-commitment/types"
	commitmenttypesv2 "github.com/cosmos/ibc-go/v9/modules/core/23-commitment/types/v2"
	"github.com/cosmos/ibc-go/v9/modules/core/exported"
)

const (
	Groth16ClientType = "groth16"
)

var _ exported.ClientState = (*ClientState)(nil)

// NewClientState creates a new ClientState instance
func NewClientState(
	latestHeight uint64,
	stateTransitionVerifierKey, stateInclusionVerifierKey []byte,
) *ClientState {
	return &ClientState{
		LatestHeight:               latestHeight,
		StateTransitionVerifierKey: stateTransitionVerifierKey,
	}
}

// ClientType is groth16.
func (cs ClientState) ClientType() string {
	return Groth16ClientType
}

// GetLatestHeight returns latest block height.
func (cs ClientState) GetLatestHeight() exported.Height {
	return clienttypes.Height{
		RevisionNumber: 0,
		RevisionHeight: cs.LatestHeight,
	}
}

// Status returns the status of the groth16 client.
func (cs ClientState) status(
	ctx context.Context,
	clientStore storetypes.KVStore,
	cdc codec.BinaryCodec,
) exported.Status {
	return exported.Active
}

// Validate performs a basic validation of the client state fields.
func (cs ClientState) Validate() error {
	if cs.StateTransitionVerifierKey == nil {
		return sdkerrors.Wrap(clienttypes.ErrInvalidClient, "state transition verifier key is nil")
	}
	return nil
}

// ZeroCustomFields returns a ClientState that is a copy of the current ClientState
// with all client customizable fields zeroed out
func (cs ClientState) ZeroCustomFields() exported.ClientState {
	// Copy over all chain-specified fields
	// and leave custom fields empty
	return &ClientState{
		LatestHeight:               cs.LatestHeight,
		StateTransitionVerifierKey: cs.StateTransitionVerifierKey,
	}
}

// initialize will check that initial consensus state is a Tendermint consensus state
// and will store ProcessedTime for initial consensus state as ctx.BlockTime()
func (cs ClientState) initialize(ctx context.Context, _ codec.BinaryCodec, clientStore storetypes.KVStore, consState exported.ConsensusState) error {
	if _, ok := consState.(*ConsensusState); !ok {
		return sdkerrors.Wrapf(clienttypes.ErrInvalidConsensus, "invalid initial consensus state. expected type: %T, got: %T",
			&ConsensusState{}, consState)
	}
	// set metadata for initial consensus state.
	setConsensusMetadata(ctx, clientStore, cs.GetLatestHeight())
	return nil
}

// verifyDelayPeriodPassed will ensure that at least delayTimePeriod amount of time and delayBlockPeriod number of blocks have passed
// since consensus state was submitted before allowing verification to continue.
func verifyDelayPeriodPassed(ctx sdk.Context, store storetypes.KVStore, proofHeight exported.Height, delayTimePeriod, delayBlockPeriod uint64) error {
	// check that executing chain's timestamp has passed consensusState's processed time + delay time period
	processedTime, ok := GetProcessedTime(store, proofHeight)
	if !ok {
		return sdkerrors.Wrapf(ErrProcessedTimeNotFound, "processed time not found for height: %s", proofHeight)
	}
	currentTimestamp := uint64(ctx.BlockTime().UnixNano())
	validTime := processedTime + delayTimePeriod
	// NOTE: delay time period is inclusive, so if currentTimestamp is validTime, then we return no error
	if currentTimestamp < validTime {
		return sdkerrors.Wrapf(ErrDelayPeriodNotPassed, "cannot verify packet until time: %d, current time: %d",
			validTime, currentTimestamp)
	}
	// check that executing chain's height has passed consensusState's processed height + delay block period
	processedHeight, ok := GetProcessedHeight(store, proofHeight)
	if !ok {
		return sdkerrors.Wrapf(ErrProcessedHeightNotFound, "processed height not found for height: %s", proofHeight)
	}
	currentHeight := clienttypes.GetSelfHeight(ctx)
	validHeight := clienttypes.NewHeight(processedHeight.GetRevisionNumber(), processedHeight.GetRevisionHeight()+delayBlockPeriod)
	// NOTE: delay block period is inclusive, so if currentHeight is validHeight, then we return no error
	if currentHeight.LT(validHeight) {
		return sdkerrors.Wrapf(ErrDelayPeriodNotPassed, "cannot verify packet until height: %s, current height: %s",
			validHeight, currentHeight)
	}
	return nil
}

//------------------------------------

// The following are modified methods from the v9 IBC Client interface. The idea is to make
// it easy to update this client once Celestia moves to v9 of IBC
func (cs ClientState) verifyMembership(
	ctx context.Context,
	clientStore storetypes.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	delayTimePeriod uint64,
	delayBlockPeriod uint64,
	proof []byte,
	path exported.Path,
	value []byte,
) error {
	// Path validation
	merklePath, ok := path.(commitmenttypesv2.MerklePath)
	if !ok {
		return sdkerrors.Wrapf(commitmenttypes.ErrInvalidProof, "expected %T, got %T", commitmenttypesv2.MerklePath{}, path)
	}

	consensusState, err := GetConsensusState(clientStore, cdc, height)
	if err != nil {
		return fmt.Errorf("failed to get consensus state: %w", err)
	}

	// Inclusion verification only supports MPT tries currently
	// KeyPath structure is arbitrary and we could structure it as we see fit
	// or we could add a new type that's more compatible with ethereum
	verifiedValue, err := mpt.VerifyMerklePatriciaTrieProof(consensusState.StateRoot, merklePath.KeyPath[0], proof)
	if err != nil {
		return fmt.Errorf("inclusion verification failed: %w", err)
	}

	if !bytes.Equal(value, verifiedValue) {
		return fmt.Errorf("retrieved value does not match the value passed to the client")
	}

	return nil
}

// VerifyNonMembership verifies a proof of the absence of a key in the Merkle tree.
// It's the same as VerifyMembership, but the value is nil
func (cs ClientState) verifyNonMembership(
	ctx context.Context,
	clientStore storetypes.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	delayTimePeriod uint64,
	delayBlockPeriod uint64,
	proof []byte,
	path exported.Path,
) error {
	consensusState, err := GetConsensusState(clientStore, cdc, height)
	if err != nil {
		return fmt.Errorf("failed to get consensus state: %w", err)
	}

	merklePath, ok := path.(commitmenttypesv2.MerklePath)
	if !ok {
		return sdkerrors.Wrapf(commitmenttypes.ErrInvalidProof, "expected %T, got %T", commitmenttypesv2.MerklePath{}, path)
	}

	// Inclusion verification only supports MPT tries currently
	// KeyPath structure is arbitrary and we could structure it as we see fit
	// or we could add a new type that's more compatible with ethereum
	verifiedValue, err := mpt.VerifyMerklePatriciaTrieProof(consensusState.StateRoot, merklePath.KeyPath[0], proof)
	if err != nil {
		fmt.Errorf("inclusion verification failed", err)
	}

	// if verifiedValue is not nil error
	if verifiedValue != nil {
		fmt.Errorf("the value for the specified key exists", err)
	}

	return nil
}

func (cs ClientState) getTimestampAtHeight(
	clientStore storetypes.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
) (uint64, error) {
	consensusState, err := GetConsensusState(clientStore, cdc, height)
	if err != nil {
		return 0, fmt.Errorf("failed to get consensus state: %w", err)
	}

	return consensusState.GetTimestamp(), nil
}

//------------------------------------

// Checking for misbehaviour is a noop for groth16
// TODO: removed in favor of
// The `CheckMisbehaviourAndUpdateState` function has been removed from `ClientState` interface.
// This functionality is now encapsulated by the usage of `VerifyClientMessage`, `CheckForMisbehaviour`, `UpdateStateOnMisbehaviour`.

// func (cs ClientState) CheckMisbehaviourAndUpdateState(
// 	ctx sdk.Context,
// 	cdc codec.BinaryCodec,
// 	clientStore storetypes.KVStore,
// 	misbehaviour exported.Misbehaviour,
// ) (exported.ClientState, error) {
// 	return &cs, nil
// }

// CheckForMisbehaviour is a noop for groth16
func (ClientState) CheckForMisbehaviour(ctx context.Context, cdc codec.BinaryCodec, clientStore storetypes.KVStore, msg exported.ClientMessage) bool {
	return false
}

func (cs ClientState) CheckSubstituteAndUpdateState(
	ctx context.Context, cdc codec.BinaryCodec, subjectClientStore,
	substituteClientStore storetypes.KVStore, substituteClient exported.ClientState,
) error {
	return sdkerrors.Wrap(clienttypes.ErrUpdateClientFailed, "cannot update groth16 client with a proposal")
}
