//+build cuda

package gm

import (
	"github.com/ultramesh/crypto-gm/cuda"
)

// Verify verify the signature by SM2PublicKey self, so the first parameter will be ignored.
func (key *SM2PublicKey) Verify(_ []byte, signature, digest []byte) (bool, error) {
	r, err := cuda.VerifySignatureGPUM([][]byte{signature}, [][]byte{digest}, [][]byte{key.X[:]})
	return r[0] == 0, err[0]
}
