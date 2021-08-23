//+build amd64 arm64

package sm2

import (
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	rand2 "math/rand"
	"testing"
)

func TestRREEDDCC(t *testing.T) {
	p, _ := new(big.Int).SetString("ffffffff00000001000000000000000000000000ffffffffffffffffffffffff", 16)
	R, _ := new(big.Int).SetString("10000000000000000", 16)
	RInv, _ := new(big.Int).SetString("ffffffff0000000100000000000000000000000100000000", 16)
	p1 := new(big.Int)
	p1.Mul(R, RInv).Sub(p1, big.NewInt(1)).Div(p1, p)
	fmt.Println(p1.Text(16)) //1
}

func TestP256Mul(t *testing.T) {
	p, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)
	rInv, _ := new(big.Int).SetString("fffffffb00000005fffffffc00000002fffffffd00000006fffffff900000004", 16)

	for i := 0; i < 0xfffff; i++ {
		a, _ := rand.Int(rand.Reader, p)
		b, _ := rand.Int(rand.Reader, p)
		bba, bbb := a.Bits(), b.Bits()
		ina := &[4]uint64{uint64(bba[0]), uint64(bba[1]), uint64(bba[2]), uint64(bba[3])}
		inb := &[4]uint64{uint64(bbb[0]), uint64(bbb[1]), uint64(bbb[2]), uint64(bbb[3])}

		r := new(big.Int)
		r.Mul(a, b).Mul(r, rInv).Mod(r, p)
		res := new([4]uint64)
		p256Mul(res, ina, inb)
		r2 := new(big.Int).SetBits([]big.Word{big.Word(res[0]), big.Word(res[1]), big.Word(res[2]), big.Word(res[3])})
		if r2.Text(16) != r.Text(16) {
			t.Log("input", a.Text(16))
			t.Log("input", b.Text(16))
			t.Error("want:", r.Text(16))
			bbs1, bbs2 := r.Bits(), r2.Bits()
			t.Errorf("%x %x %x %x \n", bbs1[0], bbs1[1], bbs1[2], bbs1[3])
			t.Error("get :", r2.Text(16))
			t.Errorf("%x %x %x %x \n", bbs2[0], bbs2[1], bbs2[2], bbs2[3])
			return
		}
	}
}

func TestP256Sqr(t *testing.T) {
	sm2P, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)
	e256, _ := new(big.Int).SetString("10000000000000000000000000000000000000000000000000000000000000000", 16)
	e256Inv := new(big.Int).ModInverse(e256, sm2P)
	for i := uint64(0); i < 0xffff; i++ {
		bigA, _ := rand.Int(rand.Reader, sm2P)
		product := new(big.Int).Mul(bigA, bigA)
		h := big.NewInt(0).Rsh(product, 256)
		product.Mod(product, e256)
		product.Mul(product, e256Inv)
		product.Mod(product, sm2P)
		product.Add(product, h)
		product.Mod(product, sm2P)

		fpA := new([4]uint64)
		for i, v := range bigA.Bits() {
			fpA[i] = uint64(v)
		}
		res := new([4]uint64)
		p256Sqr(res, fpA, 1)
		fpProduct := new(big.Int).SetBits(
			[]big.Word{big.Word(res[0]), big.Word(res[1]), big.Word(res[2]), big.Word(res[3])})

		if fpProduct.Cmp(product) != 0 {
			t.Error("input:", bigA.Text(16))
			t.Error("want :", product.Text(16))
			t.Error("get  :", fpProduct.Text(16))
			bbs1, bbs2 := product.Bits(), fpProduct.Bits()
			t.Errorf("%x %x %x %x \n", bbs1[0], bbs1[1], bbs1[2], bbs1[3])
			t.Errorf("%x %x %x %x \n", bbs2[0], bbs2[1], bbs2[2], bbs2[3])
			return
		}
	}

	in := []string{
		"1c1da93958df48a747c1b9ed3c1c3b4181ae3fcab8641ebcf7fdaa532b707c16",
		"be7539f80f4c7d7171de7ad654823d30d947326aa08d78f4caf119433f9a4f2f",
		"dcf839cfe6c92440e479c23d4dfd5236a964d3081818291c58846fa34c365c16",
	}
	for i := range in {
		bigA, _ := big.NewInt(0).SetString(in[i], 16)
		product := new(big.Int).Mul(bigA, bigA)
		h := big.NewInt(0).Rsh(product, 256)
		product.Mod(product, e256)
		product.Mul(product, e256Inv)
		product.Mod(product, sm2P)
		product.Add(product, h)
		product.Mod(product, sm2P)

		fpA := new([4]uint64)
		for i, v := range bigA.Bits() {
			fpA[i] = uint64(v)
		}
		res := new([4]uint64)
		p256Sqr(res, fpA, 1)
		fpProduct := new(big.Int).SetBits(
			[]big.Word{big.Word(res[0]), big.Word(res[1]), big.Word(res[2]), big.Word(res[3])})

		if fpProduct.Cmp(product) != 0 {
			t.Error("input:", bigA.Text(16))
			t.Error("want :", product.Text(16))
			t.Error("get  :", fpProduct.Text(16))
			bbs1, bbs2 := product.Bits(), fpProduct.Bits()
			t.Errorf("%x %x %x %x \n", bbs1[0], bbs1[1], bbs1[2], bbs1[3])
			t.Errorf("%x %x %x %x \n", bbs2[0], bbs2[1], bbs2[2], bbs2[3])
			return
		}
	}
}

func TestOrderMul(t *testing.T) {
	sm2N, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123", 16)
	e256Inv, _ := new(big.Int).SetString("10000000000000000000000000000000000000000000000000000000000000000", 16)
	e256Inv.ModInverse(e256Inv, sm2N)
	for i := uint64(0); i < 0xffff; i++ {
		bigA, _ := rand.Int(rand.Reader, sm2N)
		bigB, _ := rand.Int(rand.Reader, sm2N)
		product := new(big.Int).Mul(bigA, bigB)
		product.Mul(product, e256Inv)
		product.Mod(product, sm2N)
		res := new([4]uint64)
		fpA, fpB := new([4]uint64), new([4]uint64)
		for i, v := range bigA.Bits() {
			fpA[i] = uint64(v)
		}
		for i, v := range bigB.Bits() {
			fpB[i] = uint64(v)
		}
		orderMul(res, fpA, fpB)
		fpProduct := new(big.Int).SetBits(
			[]big.Word{big.Word(res[0]), big.Word(res[1]), big.Word(res[2]), big.Word(res[3])})

		if fpProduct.Cmp(product) != 0 {
			t.Error("input:", bigA.Text(16))
			t.Error("want :", product.Text(16))
			t.Error("get  :", fpProduct.Text(16))
			return
		}
	}
}

func TestOrderSqr(t *testing.T) {
	sm2N, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123", 16)
	e256, _ := new(big.Int).SetString("10000000000000000000000000000000000000000000000000000000000000000", 16)
	e256Inv := new(big.Int).ModInverse(e256, sm2N)
	for i := uint64(0); i < 0xffff; i++ {
		bigA, _ := rand.Int(rand.Reader, sm2N)
		product := new(big.Int).Mul(bigA, bigA)
		h := big.NewInt(0).Rsh(product, 256)
		product.Mod(product, e256)
		product.Mul(product, e256Inv)
		product.Mod(product, sm2N)
		product.Add(product, h)
		product.Mod(product, sm2N)
		res := new([4]uint64)
		fpA := new([4]uint64)
		for i, v := range bigA.Bits() {
			fpA[i] = uint64(v)
		}

		orderSqr(res, fpA, 1)
		fpProduct := new(big.Int).SetBits(
			[]big.Word{big.Word(res[0]), big.Word(res[1]), big.Word(res[2]), big.Word(res[3])})

		if fpProduct.Cmp(product) != 0 {
			t.Error("input:", bigA.Text(16))
			t.Error("want :", product.Text(16))
			t.Error("get  :", fpProduct.Text(16))
			bbs1, bbs2 := product.Bits(), fpProduct.Bits()
			t.Errorf("%x %x %x %x \n", bbs1[0], bbs1[1], bbs1[2], bbs1[3])
			t.Errorf("%x %x %x %x \n", bbs2[0], bbs2[1], bbs2[2], bbs2[3])
			return
		}
	}

	in := []string{
		"1c1da93958df48a747c1b9ed3c1c3b4181ae3fcab8641ebcf7fdaa532b707c16",
		"be7539f80f4c7d7171de7ad654823d30d947326aa08d78f4caf119433f9a4f2f",
		"dcf839cfe6c92440e479c23d4dfd5236a964d3081818291c58846fa34c365c16",
	}
	for i := range in {
		bigA, _ := big.NewInt(0).SetString(in[i], 16)
		product := new(big.Int).Mul(bigA, bigA)
		h := big.NewInt(0).Rsh(product, 256)
		product.Mod(product, e256)
		product.Mul(product, e256Inv)
		product.Mod(product, sm2N)
		product.Add(product, h)
		product.Mod(product, sm2N)
		res := new([4]uint64)
		fpA := new([4]uint64)
		for i, v := range bigA.Bits() {
			fpA[i] = uint64(v)
		}

		orderSqr(res, fpA, 1)
		fpProduct := new(big.Int).SetBits(
			[]big.Word{big.Word(res[0]), big.Word(res[1]), big.Word(res[2]), big.Word(res[3])})

		if fpProduct.Cmp(product) != 0 {
			t.Error("input:", bigA.Text(16))
			t.Error("want :", product.Text(16))
			t.Error("get  :", fpProduct.Text(16))
			return
		}
	}
}

func TestP256invert(t *testing.T) {
	var (
		in, out, in2 [4]uint64
	)
	one := [4]uint64{1, 0, 0, 0}
	for i := 0; i < 4; i++ {
		in[i] = rand2.Uint64()
	}
	copy(in2[:], in[:])
	p, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)
	RR := [4]uint64{0x0000000200000003, 0x00000002ffffffff, 0x0000000100000001, 0x400000002}
	p256Mul(&in, &in, &RR) // to montgomery from
	p256Invert(&out, &in)
	p256Mul(&out, &out, &one) //from montgomery from
	r2 := new(big.Int).SetBits([]big.Word{big.Word(out[0]), big.Word(out[1]), big.Word(out[2]), big.Word(out[3])})

	inBig := new(big.Int).SetBits([]big.Word{big.Word(in2[0]), big.Word(in2[1]), big.Word(in2[2]), big.Word(in2[3])})
	inBig.ModInverse(inBig, p)

	if r2.Cmp(inBig) != 0 {
		t.Errorf("want : %s", inBig.Text(16))
		t.Errorf("get : %s", r2.Text(16))
	}
}

func TestOrderInverse(t *testing.T) {
	var (
		in [4]uint64
	)
	for i := 0; i < 4; i++ {
		in[i] = rand2.Uint64()
	}
	N, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123", 16)
	r := new(big.Int).SetBits([]big.Word{big.Word(in[0]), big.Word(in[1]), big.Word(in[2]), big.Word(in[3])})

	ordInverse(&in)
	orderMul(&in, &in, &one)
	r.ModInverse(r, N)
	r2 := new(big.Int).SetBits([]big.Word{big.Word(in[0]), big.Word(in[1]), big.Word(in[2]), big.Word(in[3])})

	if r2.Text(16) != r.Text(16) {
		t.Error("want:", r.Text(16))
		bbs1, bbs2 := r.Bits(), r2.Bits()
		t.Errorf("%x %x %x %x \n", bbs1[0], bbs1[1], bbs1[2], bbs1[3])
		t.Error("get :", r2.Text(16))
		t.Errorf("%x %x %x %x \n", bbs2[0], bbs2[1], bbs2[2], bbs2[3])
	}
}

func TestP256Add(t *testing.T) {
	var in1, in2 [4]uint64

	for i := 0; i < 4; i++ {
		in1[i] = rand2.Uint64()
		in2[i] = rand2.Uint64()
	}
	p, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)

	a := new(big.Int).SetBits([]big.Word{big.Word(in1[0]), big.Word(in1[1]), big.Word(in1[2]), big.Word(in1[3])})
	b := new(big.Int).SetBits([]big.Word{big.Word(in2[0]), big.Word(in2[1]), big.Word(in2[2]), big.Word(in2[3])})
	a.Add(a, b).Mod(a, p)
	res := new([4]uint64)
	p256Add(res, &in1, &in2)
	r2 := new(big.Int).SetBits([]big.Word{big.Word(res[0]), big.Word(res[1]), big.Word(res[2]), big.Word(res[3])})
	if r2.Text(16) != a.Text(16) {
		t.Error("want:", a.Text(16))
		t.Error("get :", r2.Text(16))
	}

}

func TestP256Sub(t *testing.T) {
	var in1, in2 [4]uint64

	for i := 0; i < 4; i++ {
		in1[i] = rand2.Uint64()
		in2[i] = rand2.Uint64()
	}
	p, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)

	a := new(big.Int).SetBits([]big.Word{big.Word(in1[0]), big.Word(in1[1]), big.Word(in1[2]), big.Word(in1[3])})
	b := new(big.Int).SetBits([]big.Word{big.Word(in2[0]), big.Word(in2[1]), big.Word(in2[2]), big.Word(in2[3])})
	a.Sub(a, b).Mod(a, p)
	res := new([4]uint64)
	p256Sub(res, &in1, &in2)
	r2 := new(big.Int).SetBits([]big.Word{big.Word(res[0]), big.Word(res[1]), big.Word(res[2]), big.Word(res[3])})

	if r2.Text(16) != a.Text(16) {
		t.Error("want:", a.Text(16))
		t.Error("get :", r2.Text(16))
	}
}

// z = 1
func TestSm2PointAdd1(t *testing.T) {
	var in1, in2, res1, res2 sm2Point
	// R mod p
	R := [4]uint64{0x0000000000000001, 0x00000000ffffffff, 0x0000000000000000, 0x100000000}
	for i := 0; i < 4; i++ {
		in1.xyz[0][i] = rand2.Uint64()
		in1.xyz[1][i] = rand2.Uint64()
		in2.xyz[0][i] = rand2.Uint64()
		in2.xyz[1][i] = rand2.Uint64()
	}
	p, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)

	x1 := new(big.Int).SetBits([]big.Word{big.Word(in1.xyz[0][0]), big.Word(in1.xyz[0][1]), big.Word(in1.xyz[0][2]), big.Word(in1.xyz[0][3])})
	x2 := new(big.Int).SetBits([]big.Word{big.Word(in2.xyz[0][0]), big.Word(in2.xyz[0][1]), big.Word(in2.xyz[0][2]), big.Word(in2.xyz[0][3])})

	y1 := new(big.Int).SetBits([]big.Word{big.Word(in1.xyz[1][0]), big.Word(in1.xyz[1][1]), big.Word(in1.xyz[1][2]), big.Word(in1.xyz[1][3])})
	y2 := new(big.Int).SetBits([]big.Word{big.Word(in2.xyz[1][0]), big.Word(in2.xyz[1][1]), big.Word(in2.xyz[1][2]), big.Word(in2.xyz[1][3])})

	x1.Mod(x1, p)
	x2.Mod(x2, p)
	y1.Mod(y1, p)
	y2.Mod(y2, p)

	fromBig(&in1.xyz[0], x1)
	fromBig(&in1.xyz[1], y1)
	fromBig(&in2.xyz[0], x2)
	fromBig(&in2.xyz[1], y2)
	// R*R mod p
	RR := [4]uint64{0x0000000200000003, 0x00000002ffffffff, 0x0000000100000001, 0x400000002}

	p256Mul(&in1.xyz[0], &in1.xyz[0], &RR)
	p256Mul(&in1.xyz[1], &in1.xyz[1], &RR)
	p256Mul(&in2.xyz[0], &in2.xyz[0], &RR)
	p256Mul(&in2.xyz[1], &in2.xyz[1], &RR)

	copy(in1.xyz[2][:], R[:])
	copy(in2.xyz[2][:], R[:])
	all1 := new([12]uint64)
	sm2PointAdd1(&res1.xyz, &in1.xyz, &in2.xyz, all1)

	sm2PointAdd1Asm(&res2.xyz, &in1.xyz, &in2.xyz)

	res1.toAffine()
	fromMont(&res1.xyz[0], &res1.xyz[0])
	fromMont(&res1.xyz[1], &res1.xyz[1])

	res2.toAffine()

	fromMont(&res2.xyz[0], &res2.xyz[0])
	fromMont(&res2.xyz[1], &res2.xyz[1])

	a := new(big.Int).SetBits([]big.Word{big.Word(res1.xyz[0][0]), big.Word(res1.xyz[0][1]), big.Word(res1.xyz[0][2]), big.Word(res1.xyz[0][3])})
	b := new(big.Int).SetBits([]big.Word{big.Word(res1.xyz[1][0]), big.Word(res1.xyz[1][1]), big.Word(res1.xyz[1][2]), big.Word(res1.xyz[1][3])})

	j := new(big.Int).SetBits([]big.Word{big.Word(res2.xyz[0][0]), big.Word(res2.xyz[0][1]), big.Word(res2.xyz[0][2]), big.Word(res2.xyz[0][3])})
	k := new(big.Int).SetBits([]big.Word{big.Word(res2.xyz[1][0]), big.Word(res2.xyz[1][1]), big.Word(res2.xyz[1][2]), big.Word(res2.xyz[1][3])})

	if j.Cmp(a) != 0 || k.Cmp(b) != 0 {
		t.Error("get : ", a.Text(16), b.Text(16))
		t.Error("get : ", j.Text(16), k.Text(16))

	}

}

func TestSm2PointAdd2(t *testing.T) {
	var in1, in2, res2, res3 sm2Point

	for i := 0; i < 4; i++ {
		in1.xyz[0][i] = rand2.Uint64()
		in1.xyz[1][i] = rand2.Uint64()
		in1.xyz[2][i] = rand2.Uint64()

		in2.xyz[0][i] = rand2.Uint64()
		in2.xyz[1][i] = rand2.Uint64()
		in2.xyz[2][i] = rand2.Uint64()

	}
	p, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)

	x1 := new(big.Int).SetBits([]big.Word{big.Word(in1.xyz[0][0]), big.Word(in1.xyz[0][1]), big.Word(in1.xyz[0][2]), big.Word(in1.xyz[0][3])})
	x2 := new(big.Int).SetBits([]big.Word{big.Word(in2.xyz[0][0]), big.Word(in2.xyz[0][1]), big.Word(in2.xyz[0][2]), big.Word(in2.xyz[0][3])})

	y1 := new(big.Int).SetBits([]big.Word{big.Word(in1.xyz[1][0]), big.Word(in1.xyz[1][1]), big.Word(in1.xyz[1][2]), big.Word(in1.xyz[1][3])})
	y2 := new(big.Int).SetBits([]big.Word{big.Word(in2.xyz[1][0]), big.Word(in2.xyz[1][1]), big.Word(in2.xyz[1][2]), big.Word(in2.xyz[1][3])})

	z1 := new(big.Int).SetBits([]big.Word{big.Word(in1.xyz[2][0]), big.Word(in1.xyz[2][1]), big.Word(in1.xyz[2][2]), big.Word(in1.xyz[2][3])})
	z2 := new(big.Int).SetBits([]big.Word{big.Word(in2.xyz[2][0]), big.Word(in2.xyz[2][1]), big.Word(in2.xyz[2][2]), big.Word(in2.xyz[2][3])})

	x1.Mod(x1, p)
	x2.Mod(x2, p)
	y1.Mod(y1, p)
	y2.Mod(y2, p)

	z1.Mod(z1, p)
	z2.Mod(z2, p)

	fromBig(&in1.xyz[0], x1)
	fromBig(&in1.xyz[1], y1)
	fromBig(&in1.xyz[2], z1)

	fromBig(&in2.xyz[0], x2)
	fromBig(&in2.xyz[1], y2)
	fromBig(&in2.xyz[2], z2)

	// R*R mod p
	RR := [4]uint64{0x0000000200000003, 0x00000002ffffffff, 0x0000000100000001, 0x400000002}

	p256Mul(&in1.xyz[0], &in1.xyz[0], &RR)
	p256Mul(&in1.xyz[1], &in1.xyz[1], &RR)
	p256Mul(&in1.xyz[2], &in1.xyz[2], &RR)

	p256Mul(&in2.xyz[0], &in2.xyz[0], &RR)
	p256Mul(&in2.xyz[1], &in2.xyz[1], &RR)
	p256Mul(&in2.xyz[2], &in2.xyz[2], &RR)

	var all2 = new([44]uint64)
	sm2PointAdd2(&res2.xyz, &in1.xyz, &in2.xyz, all2)

	sm2PointAdd2Asm(&res3.xyz, &in1.xyz, &in2.xyz)

	res2.toAffine()

	fromMont(&res2.xyz[0], &res2.xyz[0])
	fromMont(&res2.xyz[1], &res2.xyz[1])

	res3.toAffine()

	fromMont(&res3.xyz[0], &res3.xyz[0])
	fromMont(&res3.xyz[1], &res3.xyz[1])

	c := new(big.Int).SetBits([]big.Word{big.Word(res2.xyz[0][0]), big.Word(res2.xyz[0][1]), big.Word(res2.xyz[0][2]), big.Word(res2.xyz[0][3])})
	d := new(big.Int).SetBits([]big.Word{big.Word(res2.xyz[1][0]), big.Word(res2.xyz[1][1]), big.Word(res2.xyz[1][2]), big.Word(res2.xyz[1][3])})

	h := new(big.Int).SetBits([]big.Word{big.Word(res3.xyz[0][0]), big.Word(res3.xyz[0][1]), big.Word(res3.xyz[0][2]), big.Word(res3.xyz[0][3])})
	i := new(big.Int).SetBits([]big.Word{big.Word(res3.xyz[1][0]), big.Word(res3.xyz[1][1]), big.Word(res3.xyz[1][2]), big.Word(res3.xyz[1][3])})

	if c.Cmp(h) != 0 || d.Cmp(i) != 0 {
		t.Error("get : ", c.Text(16), d.Text(16))
		t.Error("get : ", h.Text(16), i.Text(16))
	}

}
func TestSm2PointDouble1(t *testing.T) {
	var in, res1, res4 sm2Point
	// R mod p
	R := [4]uint64{0x0000000000000001, 0x00000000ffffffff, 0x0000000000000000, 0x100000000}
	for i := 0; i < 4; i++ {
		in.xyz[0][i] = rand2.Uint64()
		in.xyz[1][i] = rand2.Uint64()
	}
	p, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)

	x1 := new(big.Int).SetBits([]big.Word{big.Word(in.xyz[0][0]), big.Word(in.xyz[0][1]), big.Word(in.xyz[0][2]), big.Word(in.xyz[0][3])})
	y1 := new(big.Int).SetBits([]big.Word{big.Word(in.xyz[1][0]), big.Word(in.xyz[1][1]), big.Word(in.xyz[1][2]), big.Word(in.xyz[1][3])})

	x1.Mod(x1, p)
	y1.Mod(y1, p)

	fromBig(&in.xyz[0], x1)
	fromBig(&in.xyz[1], y1)
	// R*R mod p
	RR := [4]uint64{0x0000000200000003, 0x00000002ffffffff, 0x0000000100000001, 0x400000002}

	p256Mul(&in.xyz[0], &in.xyz[0], &RR)
	p256Mul(&in.xyz[1], &in.xyz[1], &RR)

	copy(in.xyz[2][:], R[:])
	all1 := new([16]uint64)
	sm2PointDouble1(&res1.xyz, &in.xyz, all1)

	sm2PointDouble1Asm(&res4.xyz, &in.xyz)

	res1.toAffine()
	fromMont(&res1.xyz[0], &res1.xyz[0])
	fromMont(&res1.xyz[1], &res1.xyz[1])

	res4.toAffine()

	fromMont(&res4.xyz[0], &res4.xyz[0])
	fromMont(&res4.xyz[1], &res4.xyz[1])

	a := new(big.Int).SetBits([]big.Word{big.Word(res1.xyz[0][0]), big.Word(res1.xyz[0][1]), big.Word(res1.xyz[0][2]), big.Word(res1.xyz[0][3])})
	b := new(big.Int).SetBits([]big.Word{big.Word(res1.xyz[1][0]), big.Word(res1.xyz[1][1]), big.Word(res1.xyz[1][2]), big.Word(res1.xyz[1][3])})

	j := new(big.Int).SetBits([]big.Word{big.Word(res4.xyz[0][0]), big.Word(res4.xyz[0][1]), big.Word(res4.xyz[0][2]), big.Word(res4.xyz[0][3])})
	k := new(big.Int).SetBits([]big.Word{big.Word(res4.xyz[1][0]), big.Word(res4.xyz[1][1]), big.Word(res4.xyz[1][2]), big.Word(res4.xyz[1][3])})

	if a.Cmp(j) != 0 || b.Cmp(k) != 0 {
		t.Error("get : ", a.Text(16), b.Text(16))
		t.Error("get : ", j.Text(16), k.Text(16))
	}
}
func TestSm2PointDouble2(t *testing.T) {
	var in, res2, res3 sm2Point
	for i := 0; i < 4; i++ {
		in.xyz[0][i] = rand2.Uint64()
		in.xyz[1][i] = rand2.Uint64()
		in.xyz[2][i] = rand2.Uint64()
	}
	p, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)

	x1 := new(big.Int).SetBits([]big.Word{big.Word(in.xyz[0][0]), big.Word(in.xyz[0][1]), big.Word(in.xyz[0][2]), big.Word(in.xyz[0][3])})
	y1 := new(big.Int).SetBits([]big.Word{big.Word(in.xyz[1][0]), big.Word(in.xyz[1][1]), big.Word(in.xyz[1][2]), big.Word(in.xyz[1][3])})
	z1 := new(big.Int).SetBits([]big.Word{big.Word(in.xyz[2][0]), big.Word(in.xyz[2][1]), big.Word(in.xyz[2][2]), big.Word(in.xyz[2][3])})

	x1.Mod(x1, p)
	y1.Mod(y1, p)
	z1.Mod(z1, p)

	fromBig(&in.xyz[0], x1)
	fromBig(&in.xyz[1], y1)
	fromBig(&in.xyz[2], z1)

	// R*R mod p
	RR := [4]uint64{0x0000000200000003, 0x00000002ffffffff, 0x0000000100000001, 0x400000002}

	p256Mul(&in.xyz[0], &in.xyz[0], &RR)
	p256Mul(&in.xyz[1], &in.xyz[1], &RR)
	p256Mul(&in.xyz[2], &in.xyz[2], &RR)

	all2 := new([24]uint64)
	sm2PointDouble2(&res2.xyz, &in.xyz, all2)
	sm2PointDouble2Asm(&res3.xyz, &in.xyz)

	res2.toAffine()

	fromMont(&res2.xyz[0], &res2.xyz[0])
	fromMont(&res2.xyz[1], &res2.xyz[1])

	res3.toAffine()
	fromMont(&res3.xyz[0], &res3.xyz[0])
	fromMont(&res3.xyz[1], &res3.xyz[1])

	c := new(big.Int).SetBits([]big.Word{big.Word(res2.xyz[0][0]), big.Word(res2.xyz[0][1]), big.Word(res2.xyz[0][2]), big.Word(res2.xyz[0][3])})
	d := new(big.Int).SetBits([]big.Word{big.Word(res2.xyz[1][0]), big.Word(res2.xyz[1][1]), big.Word(res2.xyz[1][2]), big.Word(res2.xyz[1][3])})

	h := new(big.Int).SetBits([]big.Word{big.Word(res3.xyz[0][0]), big.Word(res3.xyz[0][1]), big.Word(res3.xyz[0][2]), big.Word(res3.xyz[0][3])})
	i := new(big.Int).SetBits([]big.Word{big.Word(res3.xyz[1][0]), big.Word(res3.xyz[1][1]), big.Word(res3.xyz[1][2]), big.Word(res3.xyz[1][3])})

	if c.Cmp(h) != 0 || d.Cmp(i) != 0 {
		t.Error("get : ", c.Text(16), d.Text(16))
		t.Error("get : ", h.Text(16), i.Text(16))
	}
}

func TestBaseMult(t *testing.T) {
	scalar := new([4]uint64)
	for i := 0; i < 4; i++ {
		scalar[i] = rand2.Uint64()
	}
	scalarBig := new(big.Int).SetBits([]big.Word{big.Word(scalar[0]), big.Word(scalar[1]), big.Word(scalar[2]), big.Word(scalar[3])})
	N := Sm2().Params().N
	scalarBig.Mod(scalarBig, N)
	fromBig(scalar, scalarBig)

	res1 := new(sm2Point)
	res1.sm2BaseMult(scalar[:])

	res1.toAffine()
	one := [4]uint64{1, 0, 0, 0}
	p256Mul(&res1.xyz[0], &res1.xyz[0], &one)
	p256Mul(&res1.xyz[1], &res1.xyz[1], &one)

	a := new(big.Int).SetBits([]big.Word{big.Word(res1.xyz[0][0]), big.Word(res1.xyz[0][1]), big.Word(res1.xyz[0][2]), big.Word(res1.xyz[0][3])})
	b := new(big.Int).SetBits([]big.Word{big.Word(res1.xyz[1][0]), big.Word(res1.xyz[1][1]), big.Word(res1.xyz[1][2]), big.Word(res1.xyz[1][3])})

	fmt.Println("get : ", a.Text(16), b.Text(16))

}

func TestP256NegCond(t *testing.T) {
	in := new([4]uint64)
	for i := 0; i < 4; i++ {
		in[i] = rand2.Uint64()
	}
	out := new([4]uint64)

	copy(out[:], in[:])
	p256NegCond(out, 1)

	p256Add(out, out, in)
	fmt.Println(out)
}

func TestSaclarMult(t *testing.T) {
	in := new([8]uint64)
	scalar := new([4]uint64)
	p := new(sm2Point)

	for i := 0; i < 8; i++ {
		in[i] = rand2.Uint64()
	}
	for i := 0; i < 4; i++ {
		scalar[i] = rand2.Uint64()
	}
	x1 := new(big.Int).SetBits([]big.Word{big.Word(in[0]), big.Word(in[1]), big.Word(in[2]), big.Word(in[3])})
	y1 := new(big.Int).SetBits([]big.Word{big.Word(in[4]), big.Word(in[5]), big.Word(in[6]), big.Word(in[7])})
	scalarBig := new(big.Int).SetBits([]big.Word{big.Word(scalar[0]), big.Word(scalar[1]), big.Word(scalar[2]), big.Word(scalar[3])})

	P := Sm2().Params().P
	N := Sm2().Params().N
	x1.Mod(x1, P)
	y1.Mod(y1, P)
	scalarBig.Mod(scalarBig, N)

	fromBig(&p.xyz[0], x1)
	fromBig(&p.xyz[1], y1)
	fromBig(scalar, scalarBig)
	R := [4]uint64{0x0000000000000001, 0x00000000ffffffff, 0x0000000000000000, 0x100000000}
	copy(p.xyz[2][:], R[:])

	RR := [4]uint64{0x0000000200000003, 0x00000002ffffffff, 0x0000000100000001, 0x400000002}
	p256Mul(&p.xyz[0], &p.xyz[0], &RR)
	p256Mul(&p.xyz[1], &p.xyz[1], &RR)
	p.sm2ScalarMult(scalar[:])
	one := [4]uint64{1, 0, 0, 0}

	p.toAffine()
	p256Mul(&p.xyz[0], &p.xyz[0], &one)
	p256Mul(&p.xyz[1], &p.xyz[1], &one)

	a := new(big.Int).SetBits([]big.Word{big.Word(p.xyz[0][0]), big.Word(p.xyz[0][1]), big.Word(p.xyz[0][2]), big.Word(p.xyz[0][3])})
	b := new(big.Int).SetBits([]big.Word{big.Word(p.xyz[1][0]), big.Word(p.xyz[1][1]), big.Word(p.xyz[1][2]), big.Word(p.xyz[1][3])})
	fmt.Println("get : ", a.Text(16), b.Text(16))

}

func TestSm2Curve_CombinedMult1(t *testing.T) {
	e := Sm2()
	para := e.Params()
	a, b := e.ScalarBaseMult(big.NewInt(6).Bytes())
	a, b = e.Double(a, b)

	gxBits := para.Gx.Bits()
	gx := &[4]uint64{uint64(gxBits[0]), uint64(gxBits[1]), uint64(gxBits[2]), uint64(gxBits[3])}
	gyBits := para.Gy.Bits()
	gy := &[4]uint64{uint64(gyBits[0]), uint64(gyBits[1]), uint64(gyBits[2]), uint64(gyBits[3])}

	s := big.NewInt(6).Bytes()
	six := [4]uint64{}
	big2little(&six, s)
	tt := big.NewInt(12).Bytes()
	sm2.combinedMult(gx, gy, &six, &six)

	x, _ := sm2.ScalarBaseMult(tt)

	assert.Equal(t, x.Text(16), toBig(gx).Text(16))
	assert.True(t, a.Cmp(x) == 0)
}

func TestSm2Curve_CombinedMult2(t *testing.T) {
	e := Sm2()
	para := e.Params()
	a, b := e.ScalarBaseMult(big.NewInt(5).Bytes())
	c, d := e.ScalarMult(para.Gx, para.Gy, big.NewInt(6).Bytes())
	a, b = e.Add(a, b, c, d)

	gxBits := para.Gx.Bits()
	gx := &[4]uint64{uint64(gxBits[0]), uint64(gxBits[1]), uint64(gxBits[2]), uint64(gxBits[3])}
	gyBits := para.Gy.Bits()
	gy := &[4]uint64{uint64(gyBits[0]), uint64(gyBits[1]), uint64(gyBits[2]), uint64(gyBits[3])}

	f := big.NewInt(5).Bytes()
	s := big.NewInt(6).Bytes()
	five, six := [4]uint64{}, [4]uint64{}
	big2little(&five, f)
	big2little(&six, s)
	eleven := big.NewInt(11).Bytes()
	sm2.combinedMult(gx, gy, &five, &six)

	x, _ := sm2.ScalarBaseMult(eleven)

	assert.Equal(t, x, toBig(gx))
	assert.True(t, a.Cmp(x) == 0)
}

func TestOrderAdd(t *testing.T) {
	p, _ := new(big.Int).SetString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	for i := 0; i < 0xffff; i++ {
		a, _ := rand.Int(rand.Reader, p)
		b, _ := rand.Int(rand.Reader, p)
		bba, bbb := a.Bits(), b.Bits()
		ina := &[4]uint64{uint64(bba[0]), uint64(bba[1]), uint64(bba[2]), uint64(bba[3])}
		inb := &[4]uint64{uint64(bbb[0]), uint64(bbb[1]), uint64(bbb[2]), uint64(bbb[3])}
		res := &[4]uint64{}
		orderAdd(res, ina, inb)
		n, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123", 16)
		r := new(big.Int).SetBits([]big.Word{big.Word(res[0]), big.Word(res[1]), big.Word(res[2]), big.Word(res[3])})
		assert.Equal(t, a.Add(a, b).Mod(a, n), r)
	}
}

func TestOrderSub(t *testing.T) {
	p, _ := new(big.Int).SetString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	for i := 0; i < 0xffff; i++ {
		a, _ := rand.Int(rand.Reader, p)
		b, _ := rand.Int(rand.Reader, p)
		bba, bbb := a.Bits(), b.Bits()
		ina := &[4]uint64{uint64(bba[0]), uint64(bba[1]), uint64(bba[2]), uint64(bba[3])}
		inb := &[4]uint64{uint64(bbb[0]), uint64(bbb[1]), uint64(bbb[2]), uint64(bbb[3])}
		res := &[4]uint64{}
		orderSub(res, ina, inb)
		n, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123", 16)
		r := new(big.Int).SetBits([]big.Word{big.Word(res[0]), big.Word(res[1]), big.Word(res[2]), big.Word(res[3])})
		assert.Equal(t, a.Sub(a, b).Mod(a, n), r)
	}
}

func TestBiggerThan(t *testing.T) {
	p, _ := new(big.Int).SetString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	for i := 0; i < 0xfffff; i++ {
		a, _ := rand.Int(rand.Reader, p)
		b, _ := rand.Int(rand.Reader, p)
		bba, bbb := a.Bits(), b.Bits()
		ina := &[4]uint64{uint64(bba[0]), uint64(bba[1]), uint64(bba[2]), uint64(bba[3])}
		inb := &[4]uint64{uint64(bbb[0]), uint64(bbb[1]), uint64(bbb[2]), uint64(bbb[3])}
		assert.Equal(t, a.Cmp(b) > 0, biggerThan(ina, inb))
	}

	a := big.NewInt(0x3333333333333333)
	b := big.NewInt(0x3333333333333333)
	bba, bbb := a.Bits(), b.Bits()
	ina := &[4]uint64{uint64(bba[0])}
	inb := &[4]uint64{uint64(bbb[0])}
	assert.Equal(t, a.Cmp(b) > 0, biggerThan(ina, inb))

	a = big.NewInt(33)
	b = big.NewInt(0xffffff)
	bba, bbb = a.Bits(), b.Bits()
	ina = &[4]uint64{uint64(bba[0])}
	inb = &[4]uint64{uint64(bbb[0])}
	assert.Equal(t, a.Cmp(b) > 0, biggerThan(ina, inb))

	assert.Equal(t, b.Cmp(a) > 0, biggerThan(inb, ina))
}
