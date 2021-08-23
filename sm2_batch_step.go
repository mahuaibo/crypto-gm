package gm

import (
	"github.com/ultramesh/crypto-gm/internal/sm2"
	"sync"
)

//BatchHeapGo in amd64
type BatchHeapGo = sm2.BatchHeapGo

var batchHeapGoPool = &sync.Pool{
	New: func() interface{} {
		return &BatchHeapGo{}
	},
}

//GetHeap get Heap
func GetHeap() *BatchHeapGo {
	heap := batchHeapGoPool.Get().(*BatchHeapGo)
	return heap
}

//CloseHeap close Heap
func CloseHeap(in *BatchHeapGo) {
	batchHeapGoPool.Put(in)
}

//BatchVerifyInit BatchVerify Init
func BatchVerifyInit(ctx *BatchHeapGo, publicKey, signature, msg [][]byte) bool {
	return sm2.BatchVerifyInit(ctx, publicKey, signature, msg)
}

//BatchVerifyEnd BatchVerify End
func BatchVerifyEnd(ctx *BatchHeapGo) bool {
	return sm2.BatchVerifyEnd(ctx)
}
