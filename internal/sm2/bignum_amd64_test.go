//+build amd64

package sm2

import (
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	rand2 "math/rand"
	"testing"
	"time"
)

//test for mulxq
func TestMul(t *testing.T) {
	var (
		a, b uint64
		res  [2]uint64
	)
	for i := 0; i < 0xff; i++ {
		a = rand2.Uint64()
		b = rand2.Uint64()
		mul(&res, a, b)
		aBig := new(big.Int).SetUint64(a)
		bBig := new(big.Int).SetUint64(b)
		aBig.Mul(aBig, bBig)
		r2 := new(big.Int).SetBits([]big.Word{big.Word(res[0]), big.Word(res[1])})
		assert.Equal(t, aBig.Text(10), r2.Text(10))
	}
}

func TestMRInv(t *testing.T) {
	t.Skip("mRInv is demo, never used, so skip")
	util := big.NewInt(0)
	product := new(big.Int)
	RInv := new(big.Int).SetBits([]big.Word{0xffffffff00000001, 0xffffffffffffffff, 0xfffffffeffffffff})
	bigA := big.NewInt(0)
	fp := new([4]uint64)
	tt := time.Now()
	for i := uint64(0); i < 0xfffffffff; i++ {
		bigA.SetUint64(i)
		mRInv(fp, i)
		util.SetBits([]big.Word{big.Word(fp[0]), big.Word(fp[1]), big.Word(fp[2]), big.Word(fp[3])})
		product.Mul(bigA, RInv)
		if util.Cmp(product) != 0 {
			fmt.Println(i) //249108103168
			return
		}
		if i%0x10000000 == 0 { //e28
			fmt.Printf("finish 0x%x, %v\n", i, time.Now().Sub(tt).String())
		}
	}
}

func TestREDC64(t *testing.T) {
	fp := new([4]uint64)
	p, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)
	R, _ := new(big.Int).SetString("10000000000000000", 16)
	RInv := new(big.Int)
	RInv.ModInverse(R, p)
	get := big.NewInt(0)
	expect := new(big.Int)
	for i := uint64(1); i < 0xfffffff; i++ {
		bigA := big.NewInt(int64(i))
		bbs := bigA.Bits()
		fp[0], fp[1], fp[2], fp[3] = uint64(bbs[0]), 0, 0, 0
		REDC64(fp)
		get.SetBits([]big.Word{big.Word(fp[0]), big.Word(fp[1]), big.Word(fp[2]), big.Word(fp[3])})

		expect.Mul(bigA, RInv).Mod(expect, p)
		if expect.Cmp(get) != 0 {
			t.Error("input:", bigA.Text(16))
			t.Error("want :", expect.Text(16))
			t.Error("get  :", get.Text(16))
			return
		}
	}

	for i := uint64(0); i < 0xffff; i++ {
		bigA, _ := rand.Int(rand.Reader, p)
		bbs := bigA.Bits()
		fp[0], fp[1], fp[2], fp[3] = uint64(bbs[0]), uint64(bbs[1]), uint64(bbs[2]), uint64(bbs[3])
		REDC64(fp)
		get.SetBits([]big.Word{big.Word(fp[0]), big.Word(fp[1]), big.Word(fp[2]), big.Word(fp[3])})

		expect.Mul(bigA, RInv).Mod(expect, p)
		if expect.Cmp(get) != 0 {
			t.Error("input:", bigA.Text(16))
			t.Error("want :", expect.Text(16))
			t.Error("get  :", get.Text(16))
			return
		}
	}
	in := []string{
		"17d2712f32ff9ab5b53232bfd220acc74d841f33c1761b78662dac5781ef263d",
		"fffffffe00000001ffffffff00000000fffffffe00000001fffffffe00000000",
		"0",
		"1",
	}
	for i := range in {
		bigA, _ := big.NewInt(0).SetString(in[i], 16)
		bbs := bigA.Bits()
		fp[0], fp[1], fp[2], fp[3] = 0, 0, 0, 0
		for i := range bbs {
			fp[i] = uint64(bbs[i])
		}
		REDC64(fp)
		get.SetBits([]big.Word{big.Word(fp[0]), big.Word(fp[1]), big.Word(fp[2]), big.Word(fp[3])})

		expect.Mul(bigA, RInv).Mod(expect, p)
		if expect.Cmp(get) != 0 {
			t.Error("input:", bigA.Text(16))
			t.Error("want :", expect.Text(16))
			t.Error("get  :", get.Text(16))
			return
		}
	}
}

func TestREDC1111(t *testing.T) {
	sm2N, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123", 16)
	e64, _ := new(big.Int).SetString("10000000000000000", 16)
	e64Inv := new(big.Int).ModInverse(e64, sm2N)

	in := []string{
		"1c1da93958df48a747c1b9ed3c1c3b4181ae3fcab8641ebcf7fdaa532b707c16",
		"be7539f80f4c7d7171de7ad654823d30d947326aa08d78f4caf119433f9a4f2f",
		"dcf839cfe6c92440e479c23d4dfd5236a964d3081818291c58846fa34c365c16",
	}
	for i := range in {
		bigA, _ := new(big.Int).SetString(in[i], 16)
		bbs0 := bigA.Bits()
		fpA := new([4]uint64)
		for i, v := range bigA.Bits() {
			fpA[i] = uint64(v)
		}
		redc1111(fpA)
		fpProduct := new(big.Int).SetBits(
			[]big.Word{big.Word(fpA[0]), big.Word(fpA[1]), big.Word(fpA[2]), big.Word(fpA[3])})

		bigA.Mul(bigA, e64Inv).Mod(bigA, sm2N)
		bbs1, bbs2 := bigA.Bits(), fpProduct.Bits()
		if bigA.Cmp(fpProduct) != 0 {
			t.Log(fmt.Sprintf("input : %x %x %x %x \n", bbs0[0], bbs0[1], bbs0[2], bbs0[3]))
			t.Log(fmt.Sprintf("expect: %x %x %x %x \n", bbs1[0], bbs1[1], bbs1[2], bbs1[3]))
			t.Log(fmt.Sprintf("get   : %x %x %x %x \n", bbs2[0], bbs2[1], bbs2[2], bbs2[3]))
		}
	}
}

func TestREDC2222(t *testing.T) {
	sm2P, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)
	e64, _ := new(big.Int).SetString("10000000000000000", 16)
	e64Inv := new(big.Int).ModInverse(e64, sm2P)

	in := []string{
		"1c1da93958df48a747c1b9ed3c1c3b4181ae3fcab8641ebcf7fdaa532b707c16",
		"be7539f80f4c7d7171de7ad654823d30d947326aa08d78f4caf119433f9a4f2f",
		"dcf839cfe6c92440e479c23d4dfd5236a964d3081818291c58846fa34c365c16",
	}
	for i := range in {
		bigA, _ := new(big.Int).SetString(in[i], 16)
		bbs0 := bigA.Bits()
		fpA := new([4]uint64)
		for i, v := range bigA.Bits() {
			fpA[i] = uint64(v)
		}
		redc2222(fpA)
		fpProduct := new(big.Int).SetBits(
			[]big.Word{big.Word(fpA[0]), big.Word(fpA[1]), big.Word(fpA[2]), big.Word(fpA[3])})

		bigA.Mul(bigA, e64Inv).Mod(bigA, sm2P)
		bbs1, bbs2 := bigA.Bits(), fpProduct.Bits()
		if bigA.Cmp(fpProduct) != 0 {
			t.Log(fmt.Sprintf("input : %x %x %x %x \n", bbs0[0], bbs0[1], bbs0[2], bbs0[3]))
			t.Log(fmt.Sprintf("expect: %x %x %x %x \n", bbs1[0], bbs1[1], bbs1[2], bbs1[3]))
			t.Log(fmt.Sprintf("get   : %x %x %x %x \n", bbs2[0], bbs2[1], bbs2[2], bbs2[3]))
			return
		}
	}
}
