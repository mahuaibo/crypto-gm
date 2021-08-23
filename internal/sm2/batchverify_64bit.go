//+build amd64 arm64

package sm2

import (
	"container/heap"
	"errors"
	"github.com/ultramesh/crypto-gm/csp/genrand/rdrand"
	"math/big"
	"unsafe"
)

type params struct {
	s, t [4]uint64
	q    uint64
	p    *sm2Point
	r    *sm2Point
}

type Item struct {
	value    *sm2Point // 优先级队列中的数据
	priority uint64    // 优先级队列中节点的优先级
	index    int       // index是该节点在堆中的位置
}

/*
sum(s*q)*G + sum(t*q*P) = sum(q*R)
sig:(s, r)
R = (x, y)
	x = r - e
	y = getY(x)
t = s + r
P = (xp, yp)
q = rdrand.RandUint64()
*/
func BatchVerify(publicKey, signature, msg [][]byte) error {
	res := []params{}
	//r1: step1Result, r2: step2Result, r3: step3Result, r13: r1+r3
	//toVerify: r2 == r13
	double, r13, r1, r2, r3 := sm2Point{}, sm2Point{}, sm2Point{}, sm2Point{}, sm2Point{}
	err := preStep(&res, publicKey, signature, msg)
	if err != nil {
		return err
	}
	/*res[0].q = 2594979037778109695
	res[1].q = 18209446621440883996*/

	r1IsZero := step1BaseScalar(&r1, res)

	_, err = step2Scalar(&r2, res)
	if err != nil {
		return err
	}

	r3IsZero := step3Scalar(&r3, res)

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

	p256Invert(&r2.xyz[2], &r2.xyz[2])
	p256Sqr(&temp2, &r2.xyz[2], 1)
	p256Mul(&r2.xyz[2], &temp2, &r2.xyz[2]) //xyz[2] = 1/z3
	p256Mul(&r2.xyz[0], &r2.xyz[0], &temp2)
	p256Mul(&r2.xyz[1], &r2.xyz[1], &r2.xyz[2])
	p256Mul(&r2.xyz[0], &r2.xyz[0], &one)
	p256Mul(&r2.xyz[1], &r2.xyz[1], &one)

	if resIsEqual(&r13, &r2) {
		//fmt.Println("access：")
		//fmt.Printf("%b\n", res[0].q)
		//fmt.Printf("%b\n", res[1].q)
		return nil
	} else {
		//fmt.Println("error：")
		//fmt.Printf("%b\n", res[0].q)
		//fmt.Printf("%b\n", res[1].q)
		return errors.New("verification failed")
	}
}

func resIsEqual(p1, p2 *sm2Point) bool {
	for i := 0; i < 2; i++ {
		for j := 0; j < 4; j++ {
			if p1.xyz[i][j] != p2.xyz[i][j] {
				return false
			}
		}
	}
	return true
}

func preStep(out *[]params, publicKey, signature, msg [][]byte) error {
	for i := range signature {
		head := 0
		for head < len(signature[i]) && signature[i][head] != 0x30 {
			head++
		}
		if head == 0 {
			return errors.New("invalid signature without batchverify tag")
		}
		isBigger := signature[i][head-1]
		r, s := unMarshal(signature[i][head:])

		var rr, ss, e, t, x, y, xp, yp [4]uint64
		big2little(&rr, r)
		big2little(&ss, s)
		if biggerThan(&rr, n) || biggerThan(&ss, n) {
			return errors.New("invalid signature")
		}
		big2little(&e, msg[i][:])
		orderAdd(&t, &ss, &rr)
		if t[0]|t[1]|t[2]|t[3] == 0 {
			return errors.New("invalid signature")
		}

		err := getP(&xp, &yp, publicKey[i])
		if err != nil {
			return err
		}

		orderSub(&x, &rr, &e)
		//fmt.Println(toBig(&x).Bytes())
		err = getY(&y, &x, int(isBigger))
		if err != nil {
			return err
		}

		var pp, pr sm2Point
		maybeReduceModPASM(&x)
		maybeReduceModPASM(&y)
		p256Mul(&pr.xyz[0], &x, &RR)
		p256Mul(&pr.xyz[1], &y, &RR)
		// This sets Z value to 1, in the Montgomery domain.
		pr.xyz[2] = R

		maybeReduceModPASM(&xp)
		maybeReduceModPASM(&yp)
		p256Mul(&pp.xyz[0], &xp, &RR)
		p256Mul(&pp.xyz[1], &yp, &RR)
		// This sets Z value to 1, in the Montgomery domain.
		pp.xyz[2] = R

		q := rdrand.RandUint64()
		//q := uint64(1)
		temp := params{q: q, s: ss, t: t, p: &pp, r: &pr}
		/*pp.toMont()
		pp.toAffine()
		tp := toBig(&pp.xyz[0])
		fmt.Println(tp.Bytes())

		pr.toMont()
		pr.toAffine()
		tp = toBig(&pr.xyz[0])
		fmt.Println(tp.Bytes())*/
		*out = append(*out, temp)
	}
	return nil
}

func getP(xp, yp *[4]uint64, k []byte) error {
	if len(k) != 65 {
		return errors.New("invalid publicKey")
	}
	//check is on Curve
	x, y := new(big.Int).SetBytes(k[1:33]), new(big.Int).SetBytes(k[33:])
	if !Sm2().IsOnCurve(x, y) {
		return errors.New("invalid publicKey")
	}
	fromBig(xp, x)
	fromBig(yp, y)
	return nil
}

//in, out: not Mont
func getY(out, in *[4]uint64, flag int) error {
	x3, threex, b, y, y2, a3 := [4]uint64{}, [4]uint64{}, [4]uint64{}, [4]uint64{}, [4]uint64{}, [4]uint64{3}
	yy, inn := [4]uint64{}, [4]uint64{}
	p256Mul(&inn, in, &RR)
	p256Mul(&a3, &a3, &RR)
	p256Mul(&threex, &inn, &a3)
	p256Mul(&x3, &inn, &inn)
	p256Mul(&x3, &x3, &inn)
	p256Sub(&x3, &x3, &threex)
	fromBig(&b, Sm2().Params().B)
	p256Mul(&b, &b, &RR)
	p256Add(&x3, &x3, &b)

	p256Mul(&b, &x3, &one)

	//x3: g
	invertForY(&y, &x3)

	p256Mul(&y2, &y, &y)

	if y2 != x3 {
		return errors.New("invalid x value")
	}
	yy = y
	p256NegCond(&y, 1)

	p256Mul(&y, &y, &one)
	p256Mul(&yy, &yy, &one)
	if (flag == 1 && biggerThan(&yy, &y)) || (flag == 0 && !biggerThan(&yy, &y)) {
		y = yy
	}
	//p256Mul(&y, &y, &RR)
	*out = y
	return nil
}

/*

u+1:
1111111111111111111111111111111		x31:31
011111111111111111111111111111111	x32:33
11111111111111111111111111111111	x32:32
11111111111111111111111111111111	x32:32
11111111111111111111111111111111	x32:32
00000000000000000000000000000001	x1:32
00000000000000000000000000000000000000000000000000000000000000		0:62

x31 = 2^31-1
*/
// mod P , out and in are in Montgomery form
func invertForY(out, in *[4]uint64) {
	var all [40]uint64
	x1 := (*[4]uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(&all))))
	x2 := (*[4]uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(x1)) + 32))
	x4 := (*[4]uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(x2)) + 32))
	x6 := (*[4]uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(x4)) + 32))
	x7 := (*[4]uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(x6)) + 32))
	x8 := (*[4]uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(x7)) + 32))
	x15 := (*[4]uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(x8)) + 32))
	x30 := (*[4]uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(x15)) + 32))
	x31 := (*[4]uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(x30)) + 32))
	x32 := (*[4]uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(x31)) + 32))
	copy(x1[:], in[:])
	p256Sqr(x2, in, 1)
	p256Mul(x2, x2, in)

	p256Sqr(x4, x2, 2)
	p256Mul(x4, x4, x2)

	p256Sqr(x6, x4, 2)
	p256Mul(x6, x6, x2)

	p256Sqr(x7, x6, 1)
	p256Mul(x7, x7, in)

	p256Sqr(x8, x7, 1)
	p256Mul(x8, x8, in)

	p256Sqr(x15, x8, 7)
	p256Mul(x15, x15, x7)

	p256Sqr(x30, x15, 15)
	p256Mul(x30, x30, x15)

	p256Sqr(x31, x30, 1)
	p256Mul(x31, x31, in) //x31

	p256Sqr(x32, x31, 1)
	p256Mul(x32, x32, in)

	p256Sqr(out, x31, 33)
	p256Mul(out, out, x32)

	p256Sqr(out, out, 32)
	p256Mul(out, out, x32)

	p256Sqr(out, out, 32)
	p256Mul(out, out, x32)

	p256Sqr(out, out, 32)
	p256Mul(out, out, x32)

	p256Sqr(out, out, 32)
	p256Mul(out, out, x1)

	p256Sqr(out, out, 62)
}

// out in Mont and Jacobian mode
func step1BaseScalar(out *sm2Point, in []params) int {
	sum := [4]uint64{}
	for i := range in {
		res := [4]uint64{}
		smallOrderMul(&res, &RRN, &in[i].q) //toMont
		orderMul(&res, &res, &in[i].s)
		if i == 0 {
			copy(sum[:], res[:])
		} else {
			orderAdd(&sum, &sum, &res)
		}
	}
	getScalar(&sum)         //mode N
	out.sm2BaseMult(sum[:]) //out in Jacobian
	//out.toAffine()		//out in Affine
	return scalarIsZero(&sum)
}

// in and out in Mont and Jacobian mode
func step2Scalar(out *sm2Point, in []params) (int, error) {
	pq := make(PriorityQueue, len(in))
	for i, v := range in {
		pq[i] = &Item{
			value:    v.r,
			priority: v.q,
			index:    i}
		i++
	}
	heap.Init(&pq)
	for pq.Len() > 1 {
		item1 := heap.Pop(&pq).(*Item)
		item2 := heap.Pop(&pq).(*Item)

		i1IsInfinity := uint64IsZero(item1.priority)
		i2IsInfinity := uint64IsZero(item2.priority)
		newq := item1.priority - item2.priority

		var sum, double sm2Point
		pointsEqual := sm2PointAdd2Asm(&sum.xyz, &item1.value.xyz, &item2.value.xyz)
		sm2PointDouble2Asm(&double.xyz, &item2.value.xyz)
		sum.copyConditional(&double, pointsEqual)
		sum.copyConditional(item1.value, i2IsInfinity)
		sum.copyConditional(item2.value, i1IsInfinity)
		if item2.priority > 0 {
			heap.Push(&pq, &Item{value: &sum, priority: item2.priority})
		}

		if newq > 0 {
			heap.Push(&pq, &Item{value: item1.value, priority: newq})
		}
	}
	if pq.Len() < 1 {
		return 0, errors.New("step2 error!")
	}
	item := heap.Pop(&pq).(*Item)
	e := *item.value

	//fmt.Println(item.priority)
	temp := [4]uint64{uint64(item.priority), 0, 0, 0}
	getScalar(&temp)
	e.sm2ScalarMult(temp[:])

	out.xyz = e.xyz
	return scalarIsZero(&[4]uint64{item.priority}), nil
}

func step3Scalar(out *sm2Point, in []params) int {
	cl := map[sm2Point][4]uint64{}
	sum, double := sm2Point{}, sm2Point{}
	for i, v := range in {
		res, temp := [4]uint64{}, [4]uint64{}
		smallOrderMul(&res, &RRN, &in[i].q) //toMont
		orderMul(&res, &res, &in[i].t)
		_, ok := cl[*v.p]
		//fmt.Println(toBig(&v.p.xyz[0]).Bytes())
		if !ok {
			cl[*v.p] = res
		} else {
			temp = cl[*v.p]
			orderAdd(&temp, &res, &temp)
			cl[*v.p] = temp
		}
	}
	preIsZero := 1
	for i, v := range cl {
		j := sm2Point{}
		nowIsZero := scalarIsZero(&v)
		i.sm2ScalarMult(v[:])
		j = sum
		isEqual := sm2PointAdd2Asm(&sum.xyz, &j.xyz, &i.xyz)
		sm2PointDouble2Asm(&double.xyz, &i.xyz)
		sum.copyConditional(&double, isEqual)
		sum.copyConditional(&i, preIsZero)
		sum.copyConditional(&j, nowIsZero)
		if preIsZero != 1 || nowIsZero != 1 {
			preIsZero = 0
		}
	}
	out.xyz = sum.xyz
	return preIsZero
}
