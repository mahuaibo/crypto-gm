//+build !amd64,!arm64

package sm2

const heapBatchSize = 64

//BatchHeapGo
type BatchHeapGo struct {
	Size      uint64
	leftPoint point
	midout    []params
}

//BatchVerifyInit BatchVerify Init
func BatchVerifyInit(ctx *BatchHeapGo, publicKey, signature, msg [][]byte) bool {
	err := preStep(&ctx.midout, publicKey, signature, msg)
	if err != nil {
		return false
	}
	double, r13, r1, r3 := point{}, point{}, point{}, point{}
	r1.x = &sm2FieldElement{}
	r1.y = &sm2FieldElement{}
	r1.z = &sm2FieldElement{}
	r3.x = &sm2FieldElement{}
	r3.y = &sm2FieldElement{}
	r3.z = &sm2FieldElement{}
	r13.x = &sm2FieldElement{}
	r13.y = &sm2FieldElement{}
	r13.z = &sm2FieldElement{}
	double.x, double.y, double.z = &sm2FieldElement{}, &sm2FieldElement{}, &sm2FieldElement{}

	r1IsZero := step1BaseScalar(&r1, ctx.midout)
	r3IsZero := step3Scalar(&r3, ctx.midout)

	sm2PointAdd(r1.x, r1.y, r1.z, r3.x, r3.y, r3.z, r13.x, r13.y, r13.z)
	sm2PointDouble(double.x, double.y, double.z, r1.x, r1.y, r1.z)
	copyCond(r13, double, resIsEqual(r1, r3))
	copyCond(r13, r3, r1IsZero)
	copyCond(r13, r1, r3IsZero)

	ctx.leftPoint = r13
	ctx.Size = uint64(len(ctx.midout))
	return true
}

//BatchVerifyEnd BatchVerifyEnd
func BatchVerifyEnd(ctx *BatchHeapGo) bool {
	r2 := point{}
	r2.x, r2.y, r2.z = &sm2FieldElement{}, &sm2FieldElement{}, &sm2FieldElement{}

	_, err := step2Scalar(&r2, ctx.midout[:ctx.Size])
	if err != nil {
		return false
	}

	return resIsEqual(ctx.leftPoint, r2) == 1
}
