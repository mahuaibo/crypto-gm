//+build amd64 arm64

package sm2

import "math/big"

func (curve sm2Curve) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	var res, in1, in2 sm2Point

	fromBig(&in1.xyz[0], maybeReduceModP(x1))
	fromBig(&in1.xyz[1], maybeReduceModP(y1))
	fromBig(&in2.xyz[0], maybeReduceModP(x2))
	fromBig(&in2.xyz[1], maybeReduceModP(y2))
	z1 := zForAffine(x1, y1)
	fromBig(&in1.xyz[2], z1)
	z2 := zForAffine(x2, y2)
	fromBig(&in2.xyz[2], z2)
	in1.toMont()
	in2.toMont()
	sm2PointAdd1Asm(&res.xyz, &in1.xyz, &in2.xyz)
	res.toAffine()
	fromMont(&res.xyz[0], &res.xyz[0])
	fromMont(&res.xyz[1], &res.xyz[1])
	return toBig(&res.xyz[0]), toBig(&res.xyz[1])
}

func (curve sm2Curve) Double(x1, y1 *big.Int) (*big.Int, *big.Int) {
	var res, in sm2Point

	fromBig(&in.xyz[0], maybeReduceModP(x1))
	fromBig(&in.xyz[1], maybeReduceModP(y1))
	z1 := zForAffine(x1, y1)
	fromBig(&in.xyz[2], z1)
	in.toMont()
	sm2PointDouble1Asm(&res.xyz, &in.xyz)
	res.toAffine()
	fromMont(&res.xyz[0], &res.xyz[0])
	fromMont(&res.xyz[1], &res.xyz[1])
	return toBig(&res.xyz[0]), toBig(&res.xyz[1])
}

func (curve sm2Curve) ScalarMult(x1, y1 *big.Int, k []byte) (*big.Int, *big.Int) {
	var (
		res    sm2Point
		scalar [4]uint64
	)
	fromBig(&res.xyz[0], x1)
	fromBig(&res.xyz[1], y1)
	z1 := zForAffine(x1, y1)
	fromBig(&res.xyz[2], z1)
	res.toMont()
	kBig := new(big.Int).SetBytes(k)
	mayReduceByN(kBig)
	fromBig(&scalar, kBig)
	getScalar(&scalar)
	res.sm2ScalarMult(scalar[:])
	res.toAffine()
	fromMont(&res.xyz[0], &res.xyz[0])
	fromMont(&res.xyz[1], &res.xyz[1])
	return toBig(&res.xyz[0]), toBig(&res.xyz[1])
}

func (curve sm2Curve) ScalarBaseMult(k []byte) (*big.Int, *big.Int) {
	var (
		res    sm2Point
		scalar [4]uint64
	)
	kBig := new(big.Int).SetBytes(k)
	mayReduceByN(kBig)
	fromBig(&scalar, kBig)
	getScalar(&scalar)
	res.sm2BaseMult(scalar[:])
	res.toAffine()
	fromMont(&res.xyz[0], &res.xyz[0])
	fromMont(&res.xyz[1], &res.xyz[1])
	return toBig(&res.xyz[0]), toBig(&res.xyz[1])
}

func (curve sm2Curve) combinedMult(X, Y *[4]uint64, baseScalar, scalar *[4]uint64) {
	var r1, r2 sm2Point
	getScalar(baseScalar)
	r1IsInfinity := scalarIsZero(baseScalar)
	r1.sm2BaseMult(baseScalar[:])

	getScalar(scalar)
	r2IsInfinity := scalarIsZero(scalar)
	maybeReduceModPASM(X)
	maybeReduceModPASM(Y)
	p256Mul(&r2.xyz[0], X, &RR)
	p256Mul(&r2.xyz[1], Y, &RR)

	// This sets r2's Z value to 1, in the Montgomery domain.
	r2.xyz[2] = R
	r2.sm2ScalarMult(scalar[:])

	var sum, double sm2Point
	pointsEqual := sm2PointAdd2Asm(&sum.xyz, &r1.xyz, &r2.xyz)
	sm2PointDouble1Asm(&double.xyz, &r1.xyz)
	sum.copyConditional(&double, pointsEqual)
	sum.copyConditional(&r1, r2IsInfinity)
	sum.copyConditional(&r2, r1IsInfinity)

	p256Invert(&sum.xyz[2], &sum.xyz[2])
	p256Sqr(&sum.xyz[2], &sum.xyz[2], 1)
	p256Mul(&sum.xyz[0], &sum.xyz[0], &sum.xyz[2])
	p256Mul(&sum.xyz[0], &sum.xyz[0], &one)
	*X = sum.xyz[0]
}
func (p *sm2Point) sm2BaseMult(scalar []uint64) {
	//precomputeOnce.Do(initTable)

	wvalue := (scalar[0] << 1) & 0x7f
	sel, sign := boothW6(uint(wvalue))
	p256SelectBase(&p.xyz, 0, sel)
	p256NegCond(&p.xyz[1], sign)

	// (This is one, in the Montgomery domain.)
	copy(p.xyz[2][:], R[:])

	t := new(sm2Point)
	// (This is one, in the Montgomery domain.
	copy(t.xyz[2][:], R[:])

	index := uint(5)
	zero := sel

	for i := 1; i < 43; i++ {
		if index < 192 {
			wvalue = ((scalar[index/64] >> (index % 64)) + (scalar[index/64+1] << (64 - (index % 64)))) & 0x7f
		} else {
			wvalue = (scalar[index/64] >> (index % 64)) & 0x7f
		}
		index += 6
		sel, sign = boothW6(uint(wvalue))
		p256SelectBase(&t.xyz, i, sel)
		p256PointAddAffineAsm(&p.xyz, &p.xyz, &t.xyz, sign, sel, zero)
		zero |= sel
	}

}

type sm2Point struct {
	xyz [3][4]uint64
}

func (p *sm2Point) p256StorePoint(r *[16 * 4 * 3]uint64, index int) {
	copy(r[index*12:index*12+4], p.xyz[0][:])
	copy(r[index*12+4:index*12+8], p.xyz[1][:])
	copy(r[index*12+8:index*12+12], p.xyz[2][:])

}

func (p *sm2Point) toMont() {
	p256Mul(&p.xyz[0], &p.xyz[0], &RR)
	p256Mul(&p.xyz[1], &p.xyz[1], &RR)
	p256Mul(&p.xyz[2], &p.xyz[2], &RR)

}

// CopyConditional copies overwrites p with src if v == 1, and leaves p
// unchanged if v == 0.
func (p *sm2Point) copyConditional(src *sm2Point, v int) {
	pMask := uint64(v) - 1
	srcMask := ^pMask

	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			p.xyz[i][j] = (p.xyz[i][j] & pMask) | (src.xyz[i][j] & srcMask)
		}
	}
}
func (p *sm2Point) toAffine() {
	zz := &[4]uint64{}
	zInv := &[4]uint64{}
	p256Invert(zInv, &p.xyz[2])
	copy(zz[:], zInv[:])
	p256Sqr(zz, zz, 1)
	p256Mul(&p.xyz[0], &p.xyz[0], zz)
	p256Mul(zz, zz, zInv)
	p256Mul(&p.xyz[1], &p.xyz[1], zz)
}

func (p *sm2Point) sm2ScalarMult(scalar []uint64) {
	// precomp is a table of precomputed points that stores powers of p
	// from p^1 to p^16.
	var precomp [16 * 4 * 3]uint64
	var t0, t1, t2, t3 sm2Point
	// Prepare the table
	p.p256StorePoint(&precomp, 0) // 1

	sm2PointDouble2Asm(&t0.xyz, &p.xyz)
	sm2PointDouble2Asm(&t1.xyz, &t0.xyz)
	sm2PointDouble2Asm(&t2.xyz, &t1.xyz)
	sm2PointDouble2Asm(&t3.xyz, &t2.xyz)
	t0.p256StorePoint(&precomp, 1)  // 2
	t1.p256StorePoint(&precomp, 3)  // 4
	t2.p256StorePoint(&precomp, 7)  // 8
	t3.p256StorePoint(&precomp, 15) // 16

	sm2PointAdd2Asm(&t0.xyz, &t0.xyz, &p.xyz)

	sm2PointAdd2Asm(&t1.xyz, &t1.xyz, &p.xyz)
	sm2PointAdd2Asm(&t2.xyz, &t2.xyz, &p.xyz)
	t0.p256StorePoint(&precomp, 2) // 3
	t1.p256StorePoint(&precomp, 4) // 5
	t2.p256StorePoint(&precomp, 8) // 9

	sm2PointDouble2Asm(&t0.xyz, &t0.xyz)
	sm2PointDouble2Asm(&t1.xyz, &t1.xyz)
	t0.p256StorePoint(&precomp, 5) // 6
	t1.p256StorePoint(&precomp, 9) // 10

	sm2PointAdd2Asm(&t2.xyz, &t0.xyz, &p.xyz)
	sm2PointAdd2Asm(&t1.xyz, &t1.xyz, &p.xyz)
	t2.p256StorePoint(&precomp, 6)  // 7
	t1.p256StorePoint(&precomp, 10) // 11

	sm2PointDouble2Asm(&t0.xyz, &t0.xyz)
	sm2PointDouble2Asm(&t2.xyz, &t2.xyz)
	t0.p256StorePoint(&precomp, 11) // 12
	t2.p256StorePoint(&precomp, 13) // 14

	sm2PointAdd2Asm(&t0.xyz, &t0.xyz, &p.xyz)
	sm2PointAdd2Asm(&t2.xyz, &t2.xyz, &p.xyz)
	t0.p256StorePoint(&precomp, 12) // 13
	t2.p256StorePoint(&precomp, 14) // 15

	// Start scanning the window from top bit
	index := uint(254)
	var sel, sign int

	wvalue := (scalar[index/64] >> (index % 64)) & 0x3f
	sel, _ = boothW5(uint(wvalue))

	p256Select(&p.xyz, &precomp, sel)
	zero := sel

	for index > 4 {
		index -= 5
		sm2PointDouble2Asm(&p.xyz, &p.xyz)
		sm2PointDouble2Asm(&p.xyz, &p.xyz)
		sm2PointDouble2Asm(&p.xyz, &p.xyz)
		sm2PointDouble2Asm(&p.xyz, &p.xyz)
		sm2PointDouble2Asm(&p.xyz, &p.xyz)

		if index < 192 {
			wvalue = ((scalar[index/64] >> (index % 64)) + (scalar[index/64+1] << (64 - (index % 64)))) & 0x3f
		} else {
			wvalue = (scalar[index/64] >> (index % 64)) & 0x3f
		}

		sel, sign = boothW5(uint(wvalue))

		p256Select(&t0.xyz, &precomp, sel)
		p256NegCond(&t0.xyz[1], sign)
		sm2PointAdd2Asm(&t1.xyz, &t0.xyz, &p.xyz)
		p256MovCond(&t1.xyz, &t1.xyz, &p.xyz, sel)
		p256MovCond(&p.xyz, &t1.xyz, &t0.xyz, zero)
		zero |= sel
	}

	sm2PointDouble2Asm(&p.xyz, &p.xyz)
	sm2PointDouble2Asm(&p.xyz, &p.xyz)
	sm2PointDouble2Asm(&p.xyz, &p.xyz)
	sm2PointDouble2Asm(&p.xyz, &p.xyz)
	sm2PointDouble2Asm(&p.xyz, &p.xyz)

	wvalue = (scalar[0] << 1) & 0x3f
	sel, sign = boothW5(uint(wvalue))

	p256Select(&t0.xyz, &precomp, sel)
	p256NegCond(&t0.xyz[1], sign)
	sm2PointAdd2Asm(&t1.xyz, &p.xyz, &t0.xyz)
	p256MovCond(&t1.xyz, &t1.xyz, &p.xyz, sel)
	p256MovCond(&p.xyz, &t1.xyz, &t0.xyz, zero)
}
