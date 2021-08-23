//+build  amd64

package sm2

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	rand2 "math/rand"
	"testing"
	_ "unsafe"
)

//go:linkname p256MulForCompare crypto/elliptic.p256Mul
func p256MulForCompare(res, in1, in2 []uint64)

//go:linkname p256SqrForCompare crypto/elliptic.p256Sqr
func p256SqrForCompare(res, in []uint64, n int)

//go:linkname p256InverseForCompare crypto/elliptic.p256Inverse
func p256InverseForCompare(res, in []uint64)

func BenchmarkP256Mul(bb *testing.B) {
	ina := &[4]uint64{0x01, 0x03, 0x05, 0x07}
	res := new([4]uint64)
	bb.ResetTimer()
	for i := 0; i < bb.N; i++ {
		p256Mul(res, ina, ina)
	}
} //24.8ns

func BenchmarkSmallMul(bb *testing.B) {
	ina := &[4]uint64{0x01, 0x03, 0x05, 0x07}
	inb := uint64(1)
	res := new([4]uint64)
	bb.ResetTimer()
	for i := 0; i < bb.N; i++ {
		smallOrderMul(res, ina, &inb)
	}
} //19.7ns

func BenchmarkP256Sqr(bb *testing.B) {
	ina := &[4]uint64{0x01, 0x03, 0x05, 0x07}
	bb.ResetTimer()
	res := new([4]uint64)
	for i := 0; i < bb.N; i++ {
		p256Sqr(res, ina, 1)
	}
} //17.5ns 使用比较传送 16.5

func BenchmarkP256MulForCompare(b *testing.B) {
	ina := []uint64{0x01, 0x03, 0x05, 0x07}
	inb := []uint64{0x01, 0x03, 0x05, 0x07}
	r := make([]uint64, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p256MulForCompare(r, ina, inb)
	}
} // 18.2ns

func BenchmarkP256SqrForCompare(b *testing.B) {
	ina := []uint64{0x01, 0x03, 0x05, 0x07}
	r := make([]uint64, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p256SqrForCompare(r, ina, 1)
	}
} // 13.7ns

func BenchmarkSm2Inverse(b *testing.B) {
	ina := [4]uint64{0x01, 0x03, 0x05, 0x07}
	r := [4]uint64{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p256Invert(&r, &ina)
	}
} // 4581 ns
func BenchmarkP256Inverse(b *testing.B) {
	ina := [4]uint64{0x01, 0x03, 0x05, 0x07}
	r := [4]uint64{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p256InverseForCompare(r[:], ina[:])
	}
} //4108 ns
//--------------------------------------------

//go:linkname p256OrderMulForCompare crypto/elliptic.p256OrdMul
func p256OrderMulForCompare(res, in1, in2 []uint64)

//go:linkname p256OrderSqrForCompare crypto/elliptic.p256OrdSqr
func p256OrderSqrForCompare(res, in []uint64, n int)

//go:linkname invertibleForCompare crypto/elliptic.invertible
type invertibleForCompare interface {
	Inverse(k *big.Int) *big.Int
}

func BenchmarkP256OrderMul(bb *testing.B) {
	ina := &[4]uint64{0x01, 0x03, 0x05, 0x07}
	res := new([4]uint64)
	bb.ResetTimer()
	for i := 0; i < bb.N; i++ {
		orderMul(res, ina, ina)
	}
} //26.8ns

func BenchmarkP256OrderSqr(bb *testing.B) {
	ina := &[4]uint64{0x01, 0x03, 0x05, 0x07}
	bb.ResetTimer()
	res := new([4]uint64)
	for i := 0; i < bb.N; i++ {
		orderSqr(res, ina, 1)
	}
} //28.1ns 使用比较传送 23.8

func BenchmarkP256OrderMulForCompare(b *testing.B) {
	ina := []uint64{0x01, 0x03, 0x05, 0x07}
	inb := []uint64{0x01, 0x03, 0x05, 0x07}
	r := make([]uint64, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p256OrderMulForCompare(r, ina, inb)
	}
} // 26.8ns

func BenchmarkP256OrderSqrForCompare(b *testing.B) {
	ina := []uint64{0x01, 0x03, 0x05, 0x07}
	r := make([]uint64, 4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p256OrderSqrForCompare(r, ina, 1)
	}
} // 21.5ns

func BenchmarkOrderInverse(b *testing.B) {
	ina := [4]uint64{0x01, 0x03, 0x05, 0x07}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ordInverse(&ina)
	}
} // 7107 ns
func BenchmarkP256OrderInverse(b *testing.B) {
	ina := [4]uint64{0x01, 0x03, 0x05, 0x07}
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	c := key.PublicKey.Curve
	in := c.(invertibleForCompare)
	inBig := toBig(&ina)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in.Inverse(inBig)
	}
} // 7144 ns

//go:linkname p256PointAddForCompare crypto/elliptic.p256PointAddAsm
func p256PointAddForCompare(res, in1, in2 []uint64) int

//go:linkname p256PointDoubleForCompare crypto/elliptic.p256PointDoubleAsm
func p256PointDoubleForCompare(res, in1 []uint64)

//go:linkname p256PointAddAffineForCompare crypto/elliptic.p256PointAddAffineAsm
func p256PointAddAffineForCompare(res, in1, in2 []uint64, sign, sel, zero int)

func BenchmarkSm2PointAdd1(b *testing.B) {
	var in1, in2, res sm2Point
	// R mod p
	R := [4]uint64{0x0000000000000001, 0x00000000ffffffff, 0x0000000000000000, 0x100000000}

	for i := 0; i < 4; i++ {
		in1.xyz[0][i] = rand2.Uint64()
		in1.xyz[1][i] = rand2.Uint64()
		in2.xyz[0][i] = rand2.Uint64()
		in2.xyz[1][i] = rand2.Uint64()
	}
	copy(in1.xyz[2][:], R[:])
	copy(in2.xyz[2][:], R[:])
	all := new([12]uint64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm2PointAdd1(&res.xyz, &in1.xyz, &in2.xyz, all)
	}
} // 244
func BenchmarkSm2PointAdd1Asm(b *testing.B) {
	var in1, in2, res sm2Point
	// R mod p
	R := [4]uint64{0x0000000000000001, 0x00000000ffffffff, 0x0000000000000000, 0x100000000}

	for i := 0; i < 4; i++ {
		in1.xyz[0][i] = rand2.Uint64()
		in1.xyz[1][i] = rand2.Uint64()
		in2.xyz[0][i] = rand2.Uint64()
		in2.xyz[1][i] = rand2.Uint64()
	}
	copy(in1.xyz[2][:], R[:])
	copy(in2.xyz[2][:], R[:])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm2PointAdd1Asm(&res.xyz, &in1.xyz, &in2.xyz)
	}
} // 218
func BenchmarkSm2PointAdd2(b *testing.B) {
	var in1, in2, res sm2Point
	// R mod p
	R := [4]uint64{0x0000000000000001, 0x00000000ffffffff, 0x0000000000000000, 0x100000000}

	for i := 0; i < 4; i++ {
		in1.xyz[0][i] = rand2.Uint64()
		in1.xyz[1][i] = rand2.Uint64()
		in2.xyz[0][i] = rand2.Uint64()
		in2.xyz[1][i] = rand2.Uint64()
	}
	copy(in1.xyz[2][:], R[:])
	copy(in2.xyz[2][:], R[:])
	all := new([44]uint64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm2PointAdd2(&res.xyz, &in1.xyz, &in2.xyz, all)
	}
} // 350
func BenchmarkSm2PointAdd2Asm(b *testing.B) {
	var in1, in2, res sm2Point
	// R mod p
	R := [4]uint64{0x0000000000000001, 0x00000000ffffffff, 0x0000000000000000, 0x100000000}

	for i := 0; i < 4; i++ {
		in1.xyz[0][i] = rand2.Uint64()
		in1.xyz[1][i] = rand2.Uint64()
		in2.xyz[0][i] = rand2.Uint64()
		in2.xyz[1][i] = rand2.Uint64()
	}
	copy(in1.xyz[2][:], R[:])
	copy(in2.xyz[2][:], R[:])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm2PointAdd2Asm(&res.xyz, &in1.xyz, &in2.xyz)
	}
} //307
func BenchmarkP256PointAdd(b *testing.B) {
	var in1, in2, res [12]uint64
	// R mod p
	R := [4]uint64{0x0000000000000001, 0xffffffff00000000, 0xffffffffffffffff, 0x00000000fffffffe}

	for i := 0; i < 8; i++ {
		in1[i] = rand2.Uint64()
		in2[i] = rand2.Uint64()
	}
	copy(in1[8:12], R[:])
	copy(in2[8:12], R[:])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p256PointAddForCompare(res[:], in1[:], in2[:])
	}
} // 273 ns

func BenchmarkPointDouble1(b *testing.B) {
	var in1, res sm2Point
	// R mod p
	R := [4]uint64{0x0000000000000001, 0x00000000ffffffff, 0x0000000000000000, 0x100000000}

	for i := 0; i < 4; i++ {
		in1.xyz[0][i] = rand2.Uint64()
		in1.xyz[1][i] = rand2.Uint64()
	}
	copy(in1.xyz[2][:], R[:])
	all := new([16]uint64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm2PointDouble1(&res.xyz, &in1.xyz, all)
	}
} // 164 ns p256:139

func BenchmarkPointDouble2(b *testing.B) {
	var in1, res sm2Point
	// R mod p
	R := [4]uint64{0x0000000000000001, 0x00000000ffffffff, 0x0000000000000000, 0x100000000}

	for i := 0; i < 4; i++ {
		in1.xyz[0][i] = rand2.Uint64()
		in1.xyz[1][i] = rand2.Uint64()
	}
	copy(in1.xyz[2][:], R[:])
	all := new([24]uint64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm2PointDouble2(&res.xyz, &in1.xyz, all)
	}
} // 222

func BenchmarkCoordinateExchange(b *testing.B) {
	e := sm2Point{}
	for i := 0; i < 4; i++ {
		e.xyz[0][i] = rand2.Uint64()
		e.xyz[1][i] = rand2.Uint64()
		e.xyz[2][i] = rand2.Uint64()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.toAffine()
		fromMont(&e.xyz[0], &e.xyz[0])
		fromMont(&e.xyz[1], &e.xyz[1])
		z1 := zForAffine(toBig(&e.xyz[0]), toBig(&e.xyz[1]))
		fromBig(&e.xyz[2], z1)
		e.toMont()
	}
} //4868 ns

func BenchmarkPointDouble1Asm(b *testing.B) {
	var in1, res sm2Point
	// R mod p
	R := [4]uint64{0x0000000000000001, 0x00000000ffffffff, 0x0000000000000000, 0x100000000}

	for i := 0; i < 4; i++ {
		in1.xyz[0][i] = rand2.Uint64()
		in1.xyz[1][i] = rand2.Uint64()
	}
	copy(in1.xyz[2][:], R[:])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm2PointDouble1Asm(&res.xyz, &in1.xyz)
	}
} // 114ns
func BenchmarkPointDouble2Asm(b *testing.B) {
	var in1, res sm2Point
	// R mod p
	R := [4]uint64{0x0000000000000001, 0x00000000ffffffff, 0x0000000000000000, 0x100000000}

	for i := 0; i < 4; i++ {
		in1.xyz[0][i] = rand2.Uint64()
		in1.xyz[1][i] = rand2.Uint64()
	}
	copy(in1.xyz[2][:], R[:])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm2PointDouble2Asm(&res.xyz, &in1.xyz)
	}
} // 161ns

func BenchmarkP256PointDouble(b *testing.B) {
	var in1, res [12]uint64
	// R mod p
	R := [4]uint64{0x0000000000000001, 0xffffffff00000000, 0xffffffffffffffff, 0x00000000fffffffe}

	for i := 0; i < 8; i++ {
		in1[i] = rand2.Uint64()
	}
	copy(in1[8:12], R[:])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p256PointDoubleForCompare(res[:], in1[:])
	}
} // 146ns

func BenchmarkSm2PointAffine(b *testing.B) {
	var in1, in2, res sm2Point
	// R mod p

	R := [4]uint64{0x0000000000000001, 0x00000000ffffffff, 0x0000000000000000, 0x100000000}

	for i := 0; i < 4; i++ {
		in1.xyz[0][i] = rand2.Uint64()
		in1.xyz[1][i] = rand2.Uint64()
		in2.xyz[0][i] = rand2.Uint64()
		in2.xyz[1][i] = rand2.Uint64()
	}
	copy(in1.xyz[2][:], R[:])
	copy(in2.xyz[2][:], R[:])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p256PointAddAffineAsm(&res.xyz, &in1.xyz, &in2.xyz, 1, 1, 1)
	}
} // 226 ns
func BenchmarkP256PonitAddAffine(b *testing.B) {
	var in1, in2, res [12]uint64
	// R mod p
	R := [4]uint64{0x0000000000000001, 0xffffffff00000000, 0xffffffffffffffff, 0x00000000fffffffe}

	for i := 0; i < 8; i++ {
		in1[i] = rand2.Uint64()
		in2[i] = rand2.Uint64()
	}
	copy(in1[8:12], R[:])
	copy(in2[8:12], R[:])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p256PointAddAffineForCompare(res[:], in1[:], in2[:], 1, 1, 1)
	}
} // 223 ns
func BenchmarkBaseMult(b *testing.B) {
	scalar := new([4]uint64)
	for i := 0; i < 4; i++ {
		scalar[i] = rand2.Uint64()
	}
	scalarBig := new(big.Int).SetBits([]big.Word{big.Word(scalar[0]), big.Word(scalar[1]), big.Word(scalar[2]), big.Word(scalar[3])})
	N := Sm2().Params().N
	scalarBig.Mod(scalarBig, N)
	fromBig(scalar, scalarBig)

	res1 := new(sm2Point)
	for i := 0; i < b.N; i++ {
		res1.sm2BaseMult(scalar[:])
	}
} // 11210 ns p256:10347 ns

func BenchmarkScalarMult(b *testing.B) {
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.sm2ScalarMult(scalar[:])
	}
} //63862 ns p256:56180 ns

//go:linkname combinedMultForCompare crypto/elliptic.combinedMult
type combinedMultForCompare interface {
	CombinedMult(bigX, bigY *big.Int, baseScalar, scalar []byte) (x, y *big.Int)
}

func BenchmarkCombinedMult(b *testing.B) {
	//X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
	//Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
	//Xbig := new(big.Int).SetBytes(X)
	//Ybig := new(big.Int).SetBytes(Y)
	//ss, _ := new(big.Int).SetString("7a144441d80d2fc7f348a9ab1026b59fae03697b916554af2953472ec5f21469", 16)
	//tt, _ := new(big.Int).SetString("b75fd90371e9eb21e3858e427e0038fadc20b9f364e4bcf973a6febe3d81c768", 16)
	//b.StartTimer()
	//for i := 0; i < b.N; i++ {
	//	sm2.combinedMult(Xbig, Ybig, ss.Bytes(), tt.Bytes())
	//}
	//b.StopTimer()
} // 72169

func BenchmarkP256CombinedMult(b *testing.B) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	c := key.PublicKey.Curve
	in := c.(combinedMultForCompare)
	ss, _ := new(big.Int).SetString("7a144441d80d2fc7f348a9ab1026b59fae03697b916554af2953472ec5f21469", 16)
	tt, _ := new(big.Int).SetString("b75fd90371e9eb21e3858e427e0038fadc20b9f364e4bcf973a6febe3d81c768", 16)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		in.CombinedMult(key.X, key.Y, ss.Bytes(), tt.Bytes())
	}
	b.StopTimer()
} // 71996

//go:linkname fromMontForCompare crypto/elliptic.p256FromMont
func fromMontForCompare(res, in []uint64)

func BenchmarkFromMont(b *testing.B) {
	ina := &[4]uint64{0x01, 0x03, 0x05, 0x07}
	res := new([4]uint64)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fromMontForCompare(res[:], ina[:])
	}
}
