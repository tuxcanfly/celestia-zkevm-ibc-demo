package header

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)


type Keeper struct {
	// binaryCodec is used to marshal and unmarshal data from the store.
	binaryCodec codec.BinaryCodec

	// storeKey is key that is used to fetch the signal store from the multi
	// store.
	storeKey storetypes.StoreKey
}

func NewKeeper(
	binaryCodec codec.BinaryCodec, storeKey storetypes.StoreKey,
) *Keeper {
	return &Keeper{
		binaryCodec: binaryCodec,
		storeKey:    storeKey,
	}
}

func (k Keeper) SaveBlockHeader(ctx context.Context, height int64, header tmproto.Header) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
	store.Set(sdk.Uint64ToBigEndian(uint64(height)), header.DataHash)
}

func (k Keeper) GetBlockHeader(ctx context.Context, height int64) ([]byte, bool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
	key := sdk.Uint64ToBigEndian(uint64(height))
	if !store.Has(key) {
		return nil, false
	}
	headerHash := store.Get(key)
	return headerHash, true
}

// PruneHeaders prunes block headers that are older than the retention window.
func (k Keeper) PruneHeaders(ctx sdk.Context) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
    iterator := store.Iterator(nil, nil) // Start from the lowest key
    defer iterator.Close()

	latestHeight, ok := k.GetLatestSavedBlockHeight(ctx)
	if !ok {
		return
	}
    
    // Calculate the minimum height to retain
    minHeightToRetain := latestHeight - retentionPeriod

    for ; iterator.Valid(); iterator.Next() {
        // Convert the key (height) from []byte to int64
		height := sdk.BigEndianToUint64(iterator.Key())

        // If the height is below the minimum height to retain, delete it
        if height < minHeightToRetain {
            store.Delete(iterator.Key())
        } else {
            // Since entries are sorted by height, we can break early
            break
        }
    }
}

func (k Keeper) GetLatestSavedBlockHeight(ctx context.Context) (uint64, bool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
	storeIterator := store.ReverseIterator(nil, nil)
	defer storeIterator.Close()

	if !storeIterator.Valid() {
		return 0, false
	}
	// parse the key to get the height	
	height := sdk.BigEndianToUint64(storeIterator.Key())
	return height, true
}


