package groth16

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/witness"
)

type PublicWitness struct {
	TrustedHeight int64 // Provided by the relayer/user
	TrustedCelestiaHeaderHash []byte // Provided by the ZK IBC Client
	TrustedRollupStateRoot []byte // Provided by the ZK IBC Client
	NewHeight int64 // Provided by the relayer/user
	NewRollupStateRoot []byte // Provided by the relayer/user
	NewCelestiaHeaderHash []byte // Provided by Celestia State Machine
	CodeCommitment []byte // Provided during initialization of the IBC Client
	GenesisStateRoot []byte // Provided during initialization of the IBC Client
}

func (p PublicWitness) GenerateStateTransitionPublicWitness(trustedStateRoot []byte) (witness.Witness, error) {
	w, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return nil, err
	}

	numInputs := 5

	trustedCelestiaHeaderHash := 
	newCelestiaHeaderHash := 

	// Create a channel to send values to the witness
	values := make(chan any, numInputs)
	values <- p.TrustedHeight
	values <- trustedStateRoot
	values <- p.NewHeight
	values <- h.NewStateRoot
	values <- h.DataRoots
	close(values)

	err = w.Fill(numInputs, 0, values)
	if err != nil {
		return nil, fmt.Errorf("failed to fill witness: %w", err)
	}

	return w, nil
}
