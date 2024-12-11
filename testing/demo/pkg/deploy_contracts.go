package pkg

import (
	"fmt"
	"os/exec"
)

func DeployContracts() error {
	fmt.Println("Deploying IBC smart contracts...")

	// Define the deployment command
	cmd := exec.Command("forge", "script", "E2ETestDeploy.s.sol:E2ETestDeploy", "--rpc-url", "http://localhost:8545", "--private-key", "0x82bfcfadbf1712f6550d8d2c00a39f05b33ec78939d0167be2a737d691f33a6a", "--broadcast")

	// Set the working directory
	cmd.Dir = "./solidity-ibc-eureka/scripts"

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to deploy contracts: %v\nOutput: %s", err, string(output))
	}

	fmt.Println(string(output))
	return nil
}
