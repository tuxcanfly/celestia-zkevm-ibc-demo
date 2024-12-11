package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strings"
	"time"

	ibcexported "github.com/cosmos/ibc-go/v9/modules/core/exported"
	"github.com/cosmos/solidity-ibc-eureka/abigen/icscore"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const channelId = "channel-0"

type ContractAddresses struct {
	ERC20           string `json:"erc20"`
	Escrow          string `json:"escrow"`
	IBCStore        string `json:"ibcstore"`
	ICS07Tendermint string `json:"ics07Tendermint"`
	ICS20Transfer   string `json:"ics20Transfer"`
	ICS26Router     string `json:"ics26Router"`
	ICSCore         string `json:"icsCore"`
}

func DeployContracts() error {
	fmt.Println("Deploying IBC smart contracts...")

	// Define the deployment command
	cmd := exec.Command("forge", "script", "E2ETestDeploy.s.sol:E2ETestDeploy", "--rpc-url", "http://localhost:8545", "--private-key", "0x82bfcfadbf1712f6550d8d2c00a39f05b33ec78939d0167be2a737d691f33a6a", "--broadcast")

	// Export PRIVATE_KEY
	cmd.Env = append(cmd.Env, "PRIVATE_KEY=0x82bfcfadbf1712f6550d8d2c00a39f05b33ec78939d0167be2a737d691f33a6a")

	// Set the working directory
	cmd.Dir = "./solidity-ibc-eureka/scripts"

	// Direct the command's output to the standard output and standard error of the Go program
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to deploy contracts: %v", err)
	}

	// print contract addresses from broadcast run-latest.json
	filePath := "./solidity-ibc-eureka/broadcast/E2ETestDeploy.s.sol/80087/run-latest.json"

	// Read the JSON file
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Parse JSON into a map
	var runLatest map[string]interface{}
	if err := json.Unmarshal(file, &runLatest); err != nil {
		return fmt.Errorf("Error unmarshalling JSON: %w", err)
	}

	// Extract and print the contract addresses
	returns, ok := runLatest["returns"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no valid returns found")
	}

	returnValue, ok := returns["0"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no valid return value found")
	}

	value, ok := returnValue["value"].(string)
	if !ok {
		return fmt.Errorf("no valid value found")
	}

	// Unescape the JSON string
	unescapedValue := strings.ReplaceAll(value, "\\\"", "\"")

	var addresses ContractAddresses
	if err := json.Unmarshal([]byte(unescapedValue), &addresses); err != nil {
		return fmt.Errorf("error unmarshalling contract addresses: %w", err)
	}

	fmt.Println("Contract Addresses:", addresses)

	ethChainId := big.NewInt(80087)
	ethPrivateKey := "0x82bfcfadbf1712f6550d8d2c00a39f05b33ec78939d0167be2a737d691f33a6a"

	ethClient, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		return fmt.Errorf("failed to connect to ethereum client: %v", err)
	}

	// get ics core contract
	icsCoreContract, err := icscore.NewContract(ethcommon.HexToAddress(addresses.ICSCore), ethClient)
	if err != nil {
		return fmt.Errorf("failed to instantiate ICS Core contract: %v", err)
	}

	channel := icscore.IICS04ChannelMsgsChannel{
		CounterpartyId: channelId,
		MerklePrefix:   [][]byte{[]byte("ibc"), []byte("")},
	}

	tmLightClientAddress := ethcommon.HexToAddress(addresses.ICS07Tendermint)

	key, err := crypto.ToECDSA(ethcommon.FromHex(ethPrivateKey))
	if err != nil {
		return fmt.Errorf("failed to convert private key: %v", err)
	}

	tx, err := icsCoreContract.AddChannel(GetTransactOpts(key, ethChainId, ethClient), ibcexported.Tendermint, channel, tmLightClientAddress)
	if err != nil {
		return fmt.Errorf("failed to add channel: %v", err)
	}

	receipt := GetTxReciept(context.Background(), ethClient, tx.Hash())

	event, err := GetEvmEvent(receipt, icsCoreContract.ParseICS04ChannelAdded)
	if err != nil {
		return fmt.Errorf("failed to get event: %v", err)
	}

	if event.Channel.CounterpartyId == channelId {
		fmt.Println("counterparty channel added successfully")
	}

	if event.ChannelId == ibcexported.Tendermint {
		fmt.Println("channel added successfully")
	}

	return nil
}

func GetTransactOpts(key *ecdsa.PrivateKey, chainID *big.Int, ethClient *ethclient.Client) *bind.TransactOpts {
	fromAddress := crypto.PubkeyToAddress(key.PublicKey)
	nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}

	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}

	txOpts, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		panic(err)
	}

	txOpts.Nonce = big.NewInt(int64(nonce))
	txOpts.GasPrice = gasPrice

	return txOpts
}

func GetTxReciept(ctx context.Context, ethClient *ethclient.Client, hash ethcommon.Hash) *ethtypes.Receipt {
	var receipt *ethtypes.Receipt
	var err error
	err = WaitForCondition(time.Second*30, time.Second, func() (bool, error) {
		receipt, err = ethClient.TransactionReceipt(ctx, hash)
		if err != nil {
			return false, nil
		}
		fmt.Println(receipt, "receipt")
		fmt.Println(receipt != nil, "IS NOT NIL")
		return receipt != nil, nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("made it here")
	return receipt
}

// GetEvmEvent parses the logs in the given receipt and returns the first event that can be parsed
func GetEvmEvent[T any](receipt *ethtypes.Receipt, parseFn func(log ethtypes.Log) (*T, error)) (event *T, err error) {
	for _, l := range receipt.Logs {
		event, err = parseFn(*l)
		if err == nil && event != nil {
			break
		}
	}

	if event == nil {
		err = fmt.Errorf("event not found")
	}

	return
}

// WaitForCondition periodically executes the given function fn based on the provided pollingInterval.
// The function fn should return true of the desired condition is met. If the function never returns true within the timeoutAfter
// period, or fn returns an error, the condition will not have been met.
func WaitForCondition(timeoutAfter, pollingInterval time.Duration, fn func() (bool, error)) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutAfter)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("failed waiting for condition after %f seconds", timeoutAfter.Seconds())
		case <-time.After(pollingInterval):
			reachedCondition, err := fn()
			if err != nil {
				return fmt.Errorf("error occurred while waiting for condition: %s", err)
			}

			if reachedCondition {
				return nil
			}
		}
	}
}

func main() {
	err := DeployContracts()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := InitializeLightClient(); err != nil {
		fmt.Println(err)
	}
}
