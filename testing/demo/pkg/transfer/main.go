package main

import (
	"fmt"
	"os"

	"github.com/celestiaorg/celestia-zkevm-ibc-demo/testing/demo/pkg/utils"
	"github.com/cosmos/cosmos-sdk/client"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v9/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
)

const (
	// relayer is an address registered in simapp
	relayer = "cosmos1ltvzpwf3eg8e9s7wzleqdmw02lesrdex9jgt0q"

	// denom is the denomination of the token on simapp
	denom = "stake"
)

func main() {
	err := SubmitMsgTransfer()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func SubmitMsgTransfer() error {
	fmt.Printf("Submitting MsgTransfer...\n")

	clientCtx, err := utils.SetupClientContext()
	if err != nil {
		return fmt.Errorf("failed to setup client context: %v", err)
	}

	submitMsgTransfer(clientCtx)

	return nil
}

func submitMsgTransfer(clientCtx client.Context) error {
	portID := "transfer"
	channelID := "channel-0"
	version := "ics20-2"
	tokens := sdktypes.NewCoins(sdktypes.NewInt64Coin(denom, 100))
	sender := relayer
	recipient := relayer
	timeoutHeight := clienttypes.NewHeight(0, 100)
	timeoutTimestamp := uint64(0)
	memo := "test transfer"
	forwarding := &transfertypes.Forwarding{}
	msgTransfer, err := createMsgTransfer(portID, channelID, version, tokens, sender, recipient, timeoutHeight, timeoutTimestamp, memo, forwarding)

	if err != nil {
		return fmt.Errorf("failed to create MsgTransfer: %w", err)
	}

	response, err := utils.BroadcastMessages(clientCtx, relayer, 200_000, msgTransfer)
	if err != nil {
		return fmt.Errorf("failed to broadcast MsgTransfer", err)
	}

	if response.Code != 0 {
		return fmt.Errorf("failed to execute MsgTransfer", response.RawLog)
	}

	return nil
}

// createMsgTransfer returns a MsgTransfer that is constructed based on the version of the IBC transfer module in-use.
// Note: the IBC transfer module version should be specified by the channel metadata: v1 is ics20-1. v2 is ics20-2.
// TODO: I thought IBC Eureka eliminated the notion of portID so maybe this uses an outdated type.
func createMsgTransfer(portID string, channelID string, version string, tokens sdktypes.Coins, sender string, receiver string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, memo string, forwarding *transfertypes.Forwarding) (*transfertypes.MsgTransfer, error) {
	switch version {
	case transfertypes.V1:
		return &transfertypes.MsgTransfer{
			SourcePort:       portID,
			SourceChannel:    channelID,
			Token:            tokens[0],
			Sender:           sender,
			Receiver:         receiver,
			TimeoutHeight:    timeoutHeight,
			TimeoutTimestamp: timeoutTimestamp,
			Memo:             memo,
			Tokens:           sdktypes.NewCoins(),
		}, nil
	case transfertypes.V2:
		return transfertypes.NewMsgTransfer(portID, channelID, tokens, sender, receiver, timeoutHeight, timeoutTimestamp, memo, forwarding), nil
	default:
		return nil, fmt.Errorf("unsupported transfer version: %s", version)
	}
}
