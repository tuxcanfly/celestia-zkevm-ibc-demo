package main

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"cosmossdk.io/x/tx/signing"
	"github.com/celestiaorg/celestia-zkevm-ibc-demo/ibc/lightclients/groth16"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	legacysigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/gogoproto/proto"
	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v9/modules/core/04-channel/v2/types"
	commitmenttypesv2 "github.com/cosmos/ibc-go/v9/modules/core/23-commitment/types/v2"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitializeLightClient() error {
	fmt.Println("Initializing the groth16 light client")
	// ETH SET UP
	ethClient, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		return fmt.Errorf("failed to connect to ethereum client: %v", err)
	}

	// retrieve the genesis block
	genesisBlock, err := ethClient.BlockByNumber(context.Background(), big.NewInt(0))
	if err != nil {
		return fmt.Errorf("failed to get genesis block: %v", err)
	}

	latestBlock, err := ethClient.BlockByNumber(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %v", err)
	}

	clientState := groth16.NewClientState(latestBlock.Number().Uint64(), []byte{}, []byte{}, []byte{}, genesisBlock.Root().Bytes())
	clientStateAny, err := cdctypes.NewAnyWithValue(clientState)
	if err != nil {
		return fmt.Errorf("failed to create client state any: %v", err)
	}
	latestBlockTime := time.Unix(int64(latestBlock.Time()), 0)
	consensusState := groth16.NewConsensusState(latestBlockTime, latestBlock.Root().Bytes())
	consensusStateAny, err := cdctypes.NewAnyWithValue(consensusState)
	if err != nil {
		return fmt.Errorf("failed to create consensus state any: %v", err)
	}

	// relayer address registered in simapp
	relayer := "cosmos1ltvzpwf3eg8e9s7wzleqdmw02lesrdex9jgt0q"
	msg := &clienttypes.MsgCreateClient{
		ClientState:    clientStateAny,
		ConsensusState: consensusStateAny,
		Signer:         relayer,
	}

	// simapp address
	clientCtx, err := SetupClientContext()
	if err != nil {
		return fmt.Errorf("failed to setup client context: %v", err)
	}

	// broadcast the initial client creation message
	_, err = BroadcastMessages(clientCtx, relayer, 200, msg)
	if err != nil {
		return err
	}

	// establish transfer channel for groth16 client on simapp
	cosmosMerklePathPrefix := commitmenttypesv2.NewMerklePath([]byte("simd"))
	_, err = BroadcastMessages(clientCtx, relayer, 200_000, &channeltypesv2.MsgCreateChannel{
		ClientId:         groth16.ModuleName,
		MerklePathPrefix: cosmosMerklePathPrefix,
		Signer:           relayer,
	})
	if err != nil {
		return err
	}

	return nil
}

// SetupClientContext initializes a Cosmos SDK ClientContext
func SetupClientContext() (client.Context, error) {
	// Get the user's home directory
	home, err := os.Getwd()
	if err != nil {
		return client.Context{}, fmt.Errorf("failed to initialize keyring: %v", err)
	}

	// Chain-specific configurations
	chainID := "zkibc-demo"
	cometNodeURI := "http://localhost:5123"                                // Comet RPC endpoint
	appName := "celestia-zkevm-ibc-demo"                                   // Name of the application from the genesis file
	grpcAddr := "localhost:9190"                                           // gRPC endpoint
	homeDir := filepath.Join(home, "testing", "files", "simapp-validator") // Path to the keyring directory

	// Initialise codec with the necessary registerers
	interfaceRegistry, _ := cdctypes.NewInterfaceRegistryWithOptions(cdctypes.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32AccountAddrPrefix(),
			},
			ValidatorAddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
			},
		},
	})
	std.RegisterInterfaces(interfaceRegistry)
	// Register Account interface
	authtypes.RegisterInterfaces(interfaceRegistry)
	appCodec := codec.NewProtoCodec(interfaceRegistry)

	// Keyring setup
	kr, err := keyring.New(appName, keyring.BackendTest, homeDir, nil, appCodec)
	if err != nil {
		return client.Context{}, fmt.Errorf("failed to initialize keyring: %v", err)
	}

	rec, err := kr.List()
	if err != nil {
		fmt.Println(err, "keyring list error")
	}
	addr, err := rec[0].GetAddress()
	if err != nil {
		fmt.Println(err, "keyring address error")
	}

	for _, r := range rec {
		add, _ := r.GetAddress()
		fmt.Println(r.Name, "record name", add.String(), "record address")
		fmt.Println(add, "ADDRESS")
	}

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return client.Context{}, fmt.Errorf("failed to create gRPC connection: %v", err)
	}

	if err != nil {
		return client.Context{}, fmt.Errorf("failed to create tendermint client: %v", err)
	}

	txConfigOpts := authtx.ConfigOptions{
		EnabledSignModes:           append(authtx.DefaultSignModes, legacysigning.SignMode_SIGN_MODE_TEXTUAL),
		TextualCoinMetadataQueryFn: txmodule.NewGRPCCoinMetadataQueryFn(conn),
	}

	txConfig, err := authtx.NewTxConfigWithOptions(appCodec, txConfigOpts)
	if err != nil {
		return client.Context{}, fmt.Errorf("failed to create tx config: %v", err)
	}

	cometNode, err := client.NewClientFromNode(cometNodeURI)
	if err != nil {
		return client.Context{}, err
	}

	// Initialize ClientContext
	clientCtx := client.Context{}.
		WithChainID(chainID).
		WithKeyring(kr).
		WithHomeDir(homeDir).
		WithGRPCClient(conn).
		WithFromAddress(addr).
		WithFromName(rec[0].Name).
		WithSkipConfirmation(true).
		WithInterfaceRegistry(interfaceRegistry).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithKeyring(kr).
		WithTxConfig(txConfig).
		WithBroadcastMode("sync").
		WithClient(cometNode).
		WithCodec(appCodec)

	return clientCtx, nil
}

// GetFactory returns an instance of tx.Factory that is configured with this Broadcaster's CosmosChain
// and the provided user. ConfigureFactoryOptions can be used to specify arbitrary options to configure the returned
// factory.
func GetFactory(clientContext client.Context, user string, factoryOptions tx.Factory) (tx.Factory, error) {
	sdkAdd, err := sdk.AccAddressFromBech32(user)
	if err != nil {
		return tx.Factory{}, err
	}

	account, err := clientContext.AccountRetriever.GetAccount(clientContext, sdkAdd)
	if err != nil {
		return tx.Factory{}, err
	}

	return defaultTxFactory(clientContext, account), nil
}

// defaultTxFactory creates a new Factory with default configuration.
func defaultTxFactory(clientCtx client.Context, account client.Account) tx.Factory {
	fmt.Println(clientCtx.TxConfig, "TX CONFIG CLIENTCTX")
	return tx.Factory{}.
		WithAccountNumber(account.GetAccountNumber()).
		WithSequence(account.GetSequence()).
		WithSignMode(legacysigning.SignMode_SIGN_MODE_DIRECT).
		WithGas(flags.DefaultGasLimit).
		WithGasPrices("0.0001stake").
		WithMemo("interchaintest").
		WithTxConfig(clientCtx.TxConfig).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithKeybase(clientCtx.Keyring).
		WithChainID(clientCtx.ChainID).
		WithSimulateAndExecute(false)
}

// BroadcastMessages broadcasts the provided messages to the given chain and signs them on behalf of the provided user.
// Once the broadcast response is returned, we wait for two blocks to be created on chain.
func BroadcastMessages(clientContext client.Context, user string, gas uint64, msgs ...interface {
	ProtoMessage()
	Reset()
	String() string
}) (*sdk.TxResponse, error) {
	// Create a tx.Factory instance and apply the factoryOptions function
	factory := tx.Factory{}
	factoryOptions := factory.WithGas(gas)

	f, err := GetFactory(clientContext, user, factoryOptions)
	if err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}
	clientContext.Output = buffer
	clientContext.WithOutput(buffer)

	fmt.Println("BEFORE BROADCASTING TX")
	if err := tx.BroadcastTx(clientContext, f, msgs...); err != nil {
		return &sdk.TxResponse{}, err
	}

	if buffer.Len() == 0 {
		return nil, fmt.Errorf("empty buffer, transaction has not been executed yet")
	}
	// unmarshall buffer into txresponse
	var txResp sdk.TxResponse
	if err := clientContext.Codec.UnmarshalJSON(buffer.Bytes(), &txResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tx response: %v", err)
	}

	fmt.Println(txResp.Code, "TX RESPONSE")
	fmt.Println(txResp.RawLog, "TX RAW LOG")
	return &txResp, nil
}

type User interface {
	KeyName() string
	FormattedAddress() string
}
