//+build sjk1862

package csp

/*
#cgo CFLAGS : -I./include  -O3
#cgo LDFLAGS : -L"./" -L"/usr/local/lib64" -lfmapiv100

#include <memory.h>
#include "fm_def.h"
#include "fm_cpc_pub.h"

FM_RET Random(unsigned char *buf, int buf_len) {
	FM_U8 dev_index = 0;
	FM_HANDLE phDev;
    FM_RET dev = 0;
    dev = FM_CPC_OpenDevice(&dev_index, FM_DEV_TYPE_PCIE_1_0X, FM_OPEN_MULTITHREAD | FM_OPEN_MULTIPROCESS, &phDev);
    if ((dev & 0x7ff) != FME_OK) {
        return dev & 0x7ff | 0x1000;
    }
    dev = FM_CPC_GenRandom(phDev, buf_len, buf);
	if ((dev & 0x7ff) != FME_OK) {
        return dev & 0x7ff | 0x2000;
    }
	dev = FM_CPC_CloseDevice(phDev);
    if ((dev & 0x7ff) != FME_OK) {
        return dev & 0x7ff | 0x3000;
    }
    return 0;
}

*/
import "C"
import (
	"fmt"
	"unsafe"
)

func init() {
	fmt.Println("use sjk1862-G")
}

type readerStruct struct {
}

//Reader reader for random
var Reader readerStruct

func (rr readerStruct) Read(out []byte) (int, error) {
	ret := C.Random((*C.uchar)(unsafe.Pointer(&out[0])), C.int(len(out)))
	if ret != 0 {
		return 0, fmt.Errorf("init sjk1862-G error, code:%v\n", ret)
	}
	return len(out), nil
}
