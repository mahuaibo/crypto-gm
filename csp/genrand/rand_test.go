package csp

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandom(t *testing.T) {
	out := make([]byte, 128)
	n, err := Reader.Read(out)
	assert.Nil(t, err)
	assert.Equal(t, n, 128)
	fmt.Println(hex.EncodeToString(out))
}

func TestRandUint64(t *testing.T) {
	_, err := Reader.ReadUint64()
	assert.Nil(t, err)
}

// generate 8 Bytes: BenchmarkRandUint64-4    598928	      1857 ns/op
func BenchmarkRandUint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Reader.ReadUint64()
		assert.Nil(b, err)
	}
}

// Mac
// generate   8 Bytes: BenchmarkRandGo-4   	  417153	      2727 ns/op
// generate 128 Bytes: BenchmarkRandGo-4   	  328885	      3322 ns/op
// generate 256 Bytes: BenchmarkRandGo-4   	  307812	      3494 ns/op
// 需要注意的是，使用Go的crypto包生成随机数的方法在不同平台上速度差异明显，而用rdrand汇编则速度差不多
// ino4
// generate 128 Bytes: BenchmarkRandGo-8                 187671              6393 ns/op
func BenchmarkRandGo(b *testing.B) {
	out := make([]byte, 256)
	for i := 0; i < b.N; i++ {
		n, err := Reader.ReadGo(out)
		assert.Nil(b, err)
		assert.Equal(b, n, 256)
	}
}

// generate 128 Bytes: BenchmarkRand-4   	    9027	    156415 ns/op
// generate 256 Bytes: BenchmarkRand-4   	    4345	    308171 ns/op
// 性能低下的主要原因是因为RDRAND指令执行本身就慢，尽管循环会影响性能，但不是主要原因
// 将RDRAND指令替换为MOV指令后进行测试：
// generate 256 Bytes: BenchmarkRandToMov-4    430555        2830 ns/op
func BenchmarkRand(b *testing.B) {
	out := make([]byte, 256)
	for i := 0; i < b.N; i++ {
		n, err := Reader.Read(out)
		assert.Nil(b, err)
		assert.Equal(b, n, 256)
	}
}
