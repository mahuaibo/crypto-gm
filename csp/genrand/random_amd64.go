//+build amd64

package csp

import (
	"crypto/rand"
	"github.com/ultramesh/crypto-gm/csp/genrand/rdrand"
)

type readerStruct struct {
}

//Reader reader for random
var Reader readerStruct

//Read 使用rdrand指令生成随机数到输入的slice中
func (r readerStruct) Read(out []byte) (int, error) {
	return rdrand.Rand(out), nil
}

//ReadUint64 使用rdrand指令生成随机数到uint64
func (r readerStruct) ReadUint64() (uint64, error) {
	return rdrand.RandUint64(), nil
}

//ReadGo 调用go的crypto来生成随机数到输入的slice中
func (r readerStruct) ReadGo(out []byte) (int, error) {
	return rand.Read(out)
}
