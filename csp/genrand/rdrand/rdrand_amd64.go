package rdrand

import (
	"crypto/rand"
	"golang.org/x/sys/cpu"
	"unsafe"
)

//Rand generate random 8 bytes every time
func Rand(out []byte) (n int)

//randUint64ASM generate random uint64
func randUint64ASM() uint64

//RandUint64 generate random uint64
func RandUint64() uint64 {
	if cpu.X86.HasRDRAND {
		return randUint64ASM()
	}
	var a uint64
	_, _ = rand.Read((*[8]byte)(unsafe.Pointer(&a))[:])
	return a
}
