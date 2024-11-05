package mpt

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	gethtrie "github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/require"

	"encoding/gob"
	"fmt"
)

func TestVerifyMerklePatriciaTrieProof(t *testing.T) {
	trie, vals := RandomTrie(500)
	root := trie.Hash()
	for i, prover := range makeProvers(trie) {

		for _, kv := range vals {
			proof := prover(kv.k)
			proofBytes, err := proofListToBytes(*proof)
			require.NoError(t, err)

			if proof == nil {
				t.Fatalf("prover %d: missing key %x while constructing proof", i, kv.k)
			}

			val, err := VerifyMerklePatriciaTrieProof(root.Bytes(), kv.k, proofBytes)
			if err != nil {
				t.Fatalf("prover %d: failed to verify proof for key %x: %v\nraw proof: %x", i, kv.k, err, proof)
			}
			if !bytes.Equal(val, kv.v) {
				t.Fatalf("prover %d: verified value mismatch for key %x: have %x, want %x", i, kv.k, val, kv.v)
			}
		}
	}
}

// makeProvers creates Merkle trie provers based on different implementations to
// test all variations.
func makeProvers(trie *gethtrie.Trie) []func(key []byte) *ProofList {
	var provers []func(key []byte) *ProofList

	// Create a direct trie based Merkle prover
	provers = append(provers, func(key []byte) *ProofList {
		var proof ProofList
		trie.Prove(key, &proof)
		return &proof
	})
	// Create a leaf iterator based Merkle prover
	provers = append(provers, func(key []byte) *ProofList {
		var proof ProofList
		if it := gethtrie.NewIterator(trie.MustNodeIterator(key)); it.Next() && bytes.Equal(key, it.Key) {
			for _, p := range it.Prove() {
				proof.Put(crypto.Keccak256(p), p)
			}
		}
		return &proof
	})
	return provers
}

// Converts proofList to []byte by encoding it with gob
func proofListToBytes(proof ProofList) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(proof); err != nil {
		return nil, fmt.Errorf("failed to encode proofList to bytes: %w", err)
	}
	return buf.Bytes(), nil
}
