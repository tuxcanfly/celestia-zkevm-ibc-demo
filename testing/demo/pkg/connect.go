package pkg

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

func EstablishIBCConnection() error {
	// Connect to local geth node
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		return fmt.Errorf("failed to connect to ethereum client: %v", err)
	}

	// Create wallet from mnemonic
	mnemonic := "swallow practice city pear alert scale endless service rather clever salmon toss tenant law antenna garage order helmet host disorder innocent proof crunch length"
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %v", err)
	}

	// Derive account
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, true)
	if err != nil {
		return fmt.Errorf("failed to derive account: %v", err)
	}

	// Get private key
	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return fmt.Errorf("failed to get private key: %v", err)
	}

	publicKey, err := wallet.PublicKey(account)
	if err != nil {
		return fmt.Errorf("failed to get public key: %v", err)
	}

	address := crypto.PubkeyToAddress(*publicKey)
	fmt.Printf("Address: %s\n", address.Hex())

	// Create transaction
	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %v", err)
	}

	toAddress := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
	value := big.NewInt(1000000000000000000) // 1 ETH

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      21000,
		GasPrice: gasPrice,
	})

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain id: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
	return nil
}
