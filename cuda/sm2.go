//+build cuda

package cuda

/*
#cgo LDFLAGS: -L/usr/local/lib -L/usr/lib -L./ -lsm2cuda
#cgo CFLAGS: -std=c99
#include "./sm2cuda.h"
*/
import "C"
import (
	"encoding/asn1"
	"errors"
	"math/big"
	"unsafe"
)

func init() {
	C.init_sm2cuda()
}

//[]uint32 0==success other number=fail
func VerifySignatureGPUM(sig, dgest, pubX [][]byte) ([]uint32, []error) {
	type ECDSASignature struct {
		R, S *big.Int
	}
	if !(len(sig) == len(dgest) && len(dgest) == len(pubX)) {
		return []uint32{1}, []error{errors.New("num of sig not equal dgest")}
	}
	data := make([]byte, 0, len(sig)*128)
	ret := make([]uint32, len(sig))
	rerr := make([]error, len(sig))
	for i := 0; i < len(sig); i++ {
		ecdsaSignature := new(ECDSASignature)
		_, err := asn1.Unmarshal(sig[i], ecdsaSignature)
		if err != nil || sig[i] == nil { //sig ==nil means transaction is Invaild
			BadData := make([]byte, 32)
			data = append(data, BadData...)
			rerr[i] = errors.New("cannot unmarshal this signature")
			continue
		}
		r := bigToBytes(ecdsaSignature.R, 256)
		s := bigToBytes(ecdsaSignature.S, 256)
		ph := pubX[i]
		e := dgest[i]

		//(s,r,ph,e)
		data = append(data, s...)
		data = append(data, r...)
		data = append(data, ph...)
		data = append(data, e...)

	}
	C.sm2ver_cuda((*C.xint)(unsafe.Pointer(&data[0])), C.int(len(sig)), (*C.xint)(unsafe.Pointer(&ret[0])))
	return ret, rerr
}

func bigToBytes(num *big.Int, base int) []byte {
	ret := make([]byte, base/8)

	if len(num.Bytes()) > base/8 {
		return num.Bytes()
	}

	return append(ret[:len(ret)-len(num.Bytes())], num.Bytes()...)
}
