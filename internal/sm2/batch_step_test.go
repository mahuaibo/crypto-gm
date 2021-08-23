package sm2

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"github.com/ultramesh/crypto-gm/internal/sm3"
	"testing"
)

func TestBatchVerifyStep(t *testing.T) {
	sig2 := []byte{0, 48, 69, 2, 32, 13, 190, 159, 134, 254, 112, 95, 175, 247, 34, 5, 132, 150, 56, 225, 46, 210, 30, 177,
		157, 21, 183, 236, 17, 65, 204, 237, 255, 46, 57, 182, 207, 2, 33, 0, 196, 252, 200, 58, 188, 213, 181, 112, 101,
		211, 201, 31, 210, 140, 96, 168, 47, 81, 168, 169, 229, 100, 44, 65, 148, 114, 181, 46, 68, 141, 26, 225}
	X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
	Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
	h := sm3.SignHashSM3(X, Y, []byte(msg))

	X = append(X, Y...)
	X = append([]byte{0}, X...)
	var sig [][]byte
	var pk [][]byte
	var msg [][]byte
	for j := 0; j < 64; j++ {
		sig = append(sig, sig2)
		pk = append(pk, X)
		msg = append(msg, h)
	}
	ctx := new(BatchHeapGo)
	err := BatchVerifyInit(ctx, pk, sig, msg)
	assert.True(t, err)
	err = BatchVerifyEnd(ctx)
	assert.True(t, err)
}
