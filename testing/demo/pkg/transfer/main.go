package main

import (
	"fmt"
	"os"
	"time"

	"cosmossdk.io/math"
	"github.com/celestiaorg/celestia-zkevm-ibc-demo/testing/demo/pkg/utils"
	"github.com/cosmos/cosmos-sdk/client"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v9/modules/apps/transfer/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v9/modules/core/04-channel/v2/types"
	ibctesting "github.com/cosmos/ibc-go/v9/testing"
	"github.com/cosmos/solidity-ibc-eureka/abigen/ics20lib"
)

const (
	// relayer is an address registered in simapp
	relayer = "cosmos1ltvzpwf3eg8e9s7wzleqdmw02lesrdex9jgt0q"

	// ethereumUserAddress is an address registered in the ethereum chain
	ethereumUserAddress = "0x7f39c581f595b53c5cb19b5a6e5b8f3a0b1f2f6e"

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
	msgTransfer, err := createMsgTransfer()
	if err != nil {
		return fmt.Errorf("failed to create MsgTransfer: %w", err)
	}

	response, err := utils.BroadcastMessages(clientCtx, relayer, 200_000, &msgTransfer)
	if err != nil {
		return fmt.Errorf("failed to broadcast MsgTransfer %w", err)
	}

	if response.Code != 0 {
		return fmt.Errorf("failed to execute MsgTransfer %v", response.RawLog)
	}

	return nil
}

func createMsgTransfer() (channeltypesv2.MsgSendPacket, error) {
	timeout := uint64(time.Now().Add(30 * time.Minute).Unix())
	coin := sdktypes.NewCoin(denom, math.NewInt(100))

	transferPayload := ics20lib.ICS20LibFungibleTokenPacketData{
		Denom:    coin.Denom,
		Amount:   coin.Amount.BigInt(),
		Sender:   relayer,
		Receiver: ethereumUserAddress,
		Memo:     "test transfer",
	}
	transferBz, err := ics20lib.EncodeFungibleTokenPacketData(transferPayload)
	if err != nil {
		return channeltypesv2.MsgSendPacket{}, err
	}
	payload := channeltypesv2.Payload{
		SourcePort:      transfertypes.PortID,
		DestinationPort: transfertypes.PortID,
		Version:         transfertypes.V1,
		Encoding:        transfertypes.EncodingABI,
		Value:           transferBz,
	}
	msg := channeltypesv2.MsgSendPacket{
		SourceChannel:    ibctesting.FirstChannelID,
		TimeoutTimestamp: timeout,
		Payloads:         []channeltypesv2.Payload{payload},
		Signer:           relayer,
	}
	return msg, nil
}
