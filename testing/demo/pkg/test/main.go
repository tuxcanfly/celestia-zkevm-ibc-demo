package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func DeployContracts() error {
	fmt.Println("Deploying IBC smart contracts...")

	// Define the deployment command
	cmd := exec.Command("forge", "script", "E2ETestDeploy.s.sol:E2ETestDeploy", "--rpc-url", "http://localhost:8545", "--private-key", "0x82bfcfadbf1712f6550d8d2c00a39f05b33ec78939d0167be2a737d691f33a6a", "--broadcast")

	// Export PRIVATE_KEY
	cmd.Env = append(cmd.Env, "PRIVATE_KEY=0x82bfcfadbf1712f6550d8d2c00a39f05b33ec78939d0167be2a737d691f33a6a")

	// Set the working directory
	cmd.Dir = "./solidity-ibc-eureka/scripts"

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to deploy contracts: %v\nOutput: %s", err, string(output))
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
		fmt.Errorf("Error unmarshalling JSON: %v", err)
	}

	// Print the contract addresses
	// Access specific fields
	fmt.Println("Contract Addresses:")

	// Access the transactions field
	transactions, ok := runLatest["transactions"].([]interface{})
	if !ok {
		fmt.Errorf("Error asserting transactions field")
	}

	// Extract contract addresses
	for _, tx := range transactions {
		transaction, ok := tx.(map[string]interface{})
		if !ok {
			fmt.Errorf("Error asserting transaction to map[string]interface{}")
		}

		// Access the contractAddress
		if contractAddress, ok := transaction["contractAddress"].(string); ok {
			fmt.Println("Contract Address:", contractAddress)
		} else {
			fmt.Println("No valid contractAddress found")
		}
	}

	return nil
}
