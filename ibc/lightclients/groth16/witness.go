package groth16

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/witness"
)

func (h Header) GenerateStateTransitionPublicWitness(trustedStateRoot []byte) (witness.Witness, error) {
	w, err := witness.New(ecc.BN254.ScalarField())
	if err != nil {
		return nil, err
	}

	numInputs := 5

	// Create a channel to send values to the witness
	values := make(chan any, numInputs)
	values <- h.TrustedHeight
	values <- trustedStateRoot
	values <- h.NewHeight
	values <- h.NewStateRoot
	values <- h.DataRoots
	close(values)

	err = w.Fill(numInputs, 0, values)
	if err != nil {
		return nil, fmt.Errorf("failed to fill witness: %w", err)
	}

	return w, nil
}
