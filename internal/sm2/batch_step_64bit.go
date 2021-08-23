//+build amd64 arm64

package sm2

const heapBatchSize = 64

//BatchHeapGo
type BatchHeapGo struct {
	Size      uint64
	leftPoint sm2Point
	midout    []params
}

//BatchVerifyInit BatchVerify Init
func BatchVerifyInit(ctx *BatchHeapGo, publicKey, signature, msg [][]byte) bool {
	err := preStep(&ctx.midout, publicKey, signature, msg)
	if err != nil {
		return false
	}
	double, r13, r1, r3 := sm2Point{}, sm2Point{}, sm2Point{}, sm2Point{}
	r1IsZero := step1BaseScalar(&r1, ctx.midout)
	r3IsZero := step3Scalar(&r3, ctx.midout)

	isEqual := sm2PointAdd2Asm(&r13.xyz, &r1.xyz, &r3.xyz)
	sm2PointDouble2Asm(&double.xyz, &r1.xyz)
	r13.copyConditional(&double, isEqual)
	r13.copyConditional(&r1, r3IsZero)
	r13.copyConditional(&r3, r1IsZero)

	temp2 := [4]uint64{} //1/z2
	p256Invert(&r13.xyz[2], &r13.xyz[2])
	p256Sqr(&temp2, &r13.xyz[2], 1)
	p256Mul(&r13.xyz[2], &temp2, &r13.xyz[2]) //xyz[2] = 1/z3
	p256Mul(&r13.xyz[0], &r13.xyz[0], &temp2)
	p256Mul(&r13.xyz[1], &r13.xyz[1], &r13.xyz[2])
	p256Mul(&r13.xyz[0], &r13.xyz[0], &one)
	p256Mul(&r13.xyz[1], &r13.xyz[1], &one)

	ctx.leftPoint = r13
	ctx.Size = uint64(len(ctx.midout))
	return true
}

//BatchVerifyEnd BatchVerifyEnd
func BatchVerifyEnd(ctx *BatchHeapGo) bool {
	r2 := sm2Point{}
	_, err := step2Scalar(&r2, ctx.midout[:ctx.Size])
	if err != nil {
		return false
	}

	temp2 := [4]uint64{} //1/z2
	p256Invert(&r2.xyz[2], &r2.xyz[2])
	p256Sqr(&temp2, &r2.xyz[2], 1)
	p256Mul(&r2.xyz[2], &temp2, &r2.xyz[2]) //xyz[2] = 1/z3
	p256Mul(&r2.xyz[0], &r2.xyz[0], &temp2)
	p256Mul(&r2.xyz[1], &r2.xyz[1], &r2.xyz[2])
	p256Mul(&r2.xyz[0], &r2.xyz[0], &one)
	p256Mul(&r2.xyz[1], &r2.xyz[1], &one)

	return resIsEqual(&ctx.leftPoint, &r2)
}
