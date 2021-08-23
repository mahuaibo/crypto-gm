package sm9

import (
	"io"
	"unsafe"
)

//#cgo amd64 CFLAGS: -D_ARCH_PAINTER=64
//#cgo arm arm64 386 CFLAGS: -D_ARCH_PAINTER=32
//#cgo LDFLAGS: -lm
//#include "sm9_sv.h"
//#include "miracl.h"
//#include <stdlib.h>
import "C"

func init() {
	C.SM9_Init()
}

//Check self check
func Check() {
	if C.SM9_SelfCheck() != 9987 {
		panic("check error")
	}
}

//GenMasterKeyPair generate master key pair
func GenMasterKeyPair(rand io.Reader) (ks, pub []byte) {
	ks = make([]byte, 32)
	_, _ = rand.Read(ks)
	r := C.SM9_Setup((*C.uchar)(C.CBytes(ks)))
	defer C.free(unsafe.Pointer(r))
	pub = C.GoBytes(unsafe.Pointer(r), 128)
	return
}

//GenMasterPub generate master pub key
func GenMasterPub(ks []byte) (pub []byte) {
	r := C.SM9_Setup((*C.uchar)(C.CBytes(ks)))
	defer C.free(unsafe.Pointer(r))
	pub = C.GoBytes(unsafe.Pointer(r), 128)
	return
}

//GenSignKey generate sm9 sign key
func GenSignKey(ID, ks []byte) (dsa []byte) {
	Cid, Cks := C.CBytes(ID), C.CBytes(ks)
	defer C.free(Cid)
	defer C.free(Cks)
	var Cdsa *C.uchar
	Cdsa = C.SM9_GenerateSignKey((*C.char)(Cid), C.int(len(ID)), (*C.uchar)(Cks))
	defer C.free(unsafe.Pointer(Cdsa))
	return C.GoBytes(unsafe.Pointer(Cdsa), 64)
}

//Sign generate signature
func Sign(msg, dsa, Ppub []byte) (sign []byte, err error) {
	Cmsg, Cdsa, CPpub := C.CBytes(msg), C.CBytes(dsa), C.CBytes(Ppub)
	defer C.free(Cmsg)
	defer C.free(Cdsa)
	defer C.free(CPpub)
	Csign, _ := C.SM9_Sign((*C.uchar)(Cmsg), C.int(len(msg)), (*C.uchar)(Cdsa), (*C.uchar)(CPpub))
	defer C.free(unsafe.Pointer(Csign))
	return C.GoBytes(unsafe.Pointer(Csign), 96), nil
}

func sm9SignTest(msg, dsa, Ppub, rand []byte) (sign []byte, err error) {
	Cmsg, Cdsa, CPpub, Cr := C.CBytes(msg), C.CBytes(dsa), C.CBytes(Ppub), C.CBytes(rand)
	defer C.free(Cmsg)
	defer C.free(Cdsa)
	defer C.free(CPpub)
	defer C.free(Cr)
	r := C.mirvar(0)
	C.bytes_to_big(C.BNLEN, (*C.char)(Cr), r)
	Csign, _ := C.SM9_Sign_With_Rand((*C.uchar)(Cmsg), C.int(len(msg)), r, (*C.uchar)(Cdsa), (*C.uchar)(CPpub))
	defer C.free(unsafe.Pointer(Csign))
	defer C.mirkill(r)
	return C.GoBytes(unsafe.Pointer(Csign), 96), nil
}

//Verify verify sm9 signature
func Verify(IDA string, sign, msg, Ppub []byte) bool {
	Csign, Cid, Cmsg, CPpub := C.CBytes(sign), C.CString(IDA), C.CBytes(msg), C.CBytes(Ppub)
	r := C.SM9_Verify((*C.uchar)(Csign), Cid, C.int(len(IDA)), (*C.uchar)(Cmsg), C.int(len(msg)), (*C.uchar)(CPpub))
	defer C.free(Csign)
	defer C.free(unsafe.Pointer(Cid))
	defer C.free(Cmsg)
	defer C.free(CPpub)
	return int(r) == 0
}
