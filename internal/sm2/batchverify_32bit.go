//+build !amd64,!arm64

package sm2

import (
	"container/heap"
	"crypto/rand"
	"errors"
	"math/big"
	"strconv"
)

type point struct {
	x, y, z *sm2FieldElement
}
type params struct {
	s, t big.Int
	q    *uint64
	p    *point
	r    *point
}

type Item struct {
	value    *point // 优先级队列中的数据
	priority uint64 // 优先级队列中节点的优先级
	index    int    // index是该节点在堆中的位置
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
	var r13, r1, r2, r3 point
	r1.x = &sm2FieldElement{}
	r1.y = &sm2FieldElement{}
	r1.z = &sm2FieldElement{}
	r3.x = &sm2FieldElement{}
	r3.y = &sm2FieldElement{}
	r3.z = &sm2FieldElement{}
	r13.x = &sm2FieldElement{}
	r13.y = &sm2FieldElement{}
	r13.z = &sm2FieldElement{}
	err := preStep(&res, publicKey, signature, msg)
	if err != nil {
		return err
	}

	preIsZero := step1BaseScalar(&r1, res)
	/*ta, _ := sm2ToAffine(r1.x, r1.y, r1.z)
	fmt.Println(ta.Bytes())*/

	_, err = step2Scalar(&r2, res)
	if err != nil {
		return err
	}
	/*ta, _ = sm2ToAffine(r2.x, r2.y, r2.z)
	fmt.Println(ta.Bytes())*/

	nowIsZero := step3Scalar(&r3, res)
	/*ta, _ = sm2ToAffine(r3.x, r3.y, r3.z)
	fmt.Println(ta.Bytes())*/

	double := point{}
	double.x, double.y, double.z = &sm2FieldElement{}, &sm2FieldElement{}, &sm2FieldElement{}
	sm2PointAdd(r1.x, r1.y, r1.z, r3.x, r3.y, r3.z, r13.x, r13.y, r13.z)
	sm2PointDouble(double.x, double.y, double.z, r1.x, r1.y, r1.z)
	copyCond(r13, double, resIsEqual(r1, r3))
	copyCond(r13, r3, preIsZero)
	copyCond(r13, r1, nowIsZero)
	/*ta, _ = sm2ToAffine(r13.x, r13.y, r13.z)
	fmt.Println(ta.Bytes())*/

	if resIsEqual(r13, r2) == 1 {
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

		rr, ss, e, t, x, y, xp, yp := new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)
		ss = new(big.Int).SetBytes(s)
		rr = new(big.Int).SetBytes(r)
		if ss.Cmp(Sm2().Params().N) >= 0 || rr.Cmp(Sm2().Params().N) >= 0 {
			return errors.New("invalid signature")
		}
		e = new(big.Int).SetBytes(msg[i][:])
		t = new(big.Int).Add(ss, rr)
		t.Mod(t, Sm2().Params().N)
		if t.Sign() == 0 {
			return errors.New("invalid signature")
		}

		var pp, pr point
		xp, yp, err := getP(publicKey[i])
		//fmt.Println(xp, yp)
		if err != nil {
			return err
		}
		pp.x, pp.y, pp.z = sm2FromAffine(xp, yp)

		x.Sub(rr, e)
		x.Mod(x, sm2.N)
		//	fmt.Println("1", x.Bytes())
		y, err = getY(x, int(isBigger))
		if err != nil {
			return err
		}
		pr.x, pr.y, pr.z = sm2FromAffine(x, y)

		q := big.NewInt(0)
		limit, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFF", 16)
		for q.Cmp(big.NewInt(0)) == 0 {
			q, err = rand.Int(rand.Reader, limit)
			if err != nil {
				return err
			}
		}
		tq, err := strconv.ParseUint(q.String(), 10, 64)
		if err != nil {
			return err
		}

		temp := params{q: &tq, s: *ss, t: *t, p: &pp, r: &pr}
		//		fmt.Println(tq, ss.Bytes(), t.Bytes(), sm2ToBig(pp.x).Bytes(), sm2ToBig(pr.x).Bytes())
		*out = append(*out, temp)
	}
	return nil
}

func getP(k []byte) (x1, y1 *big.Int, e error) {
	if len(k) != 65 {
		return nil, nil, errors.New("invalid publicKey")
	}
	//check is on Curve
	x, y := new(big.Int).SetBytes(k[1:33]), new(big.Int).SetBytes(k[33:])
	if !Sm2().IsOnCurve(x, y) {
		return nil, nil, errors.New("invalid publicKey")
	}
	return x, y, nil
}

func getY(in *big.Int, flag int) (out *big.Int, e error) {
	a := big.NewInt(3)
	var yy, x, x3, threex, b, y, y2, a3 sm2FieldElement
	sm2FromBig(&x, in) // x = in * R mod P
	sm2FromBig(&a3, a)
	sm2FromBig(&b, Sm2().Params().B)
	//	fmt.Println(x)
	sm2Mul(&threex, &x, &a3)
	sm2Mul(&x3, &x, &x)
	sm2Mul(&x3, &x3, &x)
	sm2Sub(&x3, &x3, &threex)
	sm2Add(&x3, &x3, &b)
	//	fmt.Println(x3)
	//x3: g
	invertForY(&y, &x3)

	sm2Mul(&y2, &y, &y)

	r1 := sm2ToBig(&y2)
	r2 := sm2ToBig(&x3)
	if r1.Cmp(r2) != 0 {
		return nil, errors.New("invalid x value")
	}
	yy = y

	negCond(&y, 1)
	out1 := sm2ToBig(&y)
	o := sm2ToBig(&yy)
	if (flag == 1 && o.Cmp(out1) > 0) || (flag == 0 && o.Cmp(out1) < 0) {
		out1 = o
	}
	//	fmt.Println(out.Bytes())
	return out1, nil
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
func invertForY(out, in *sm2FieldElement) {
	var x1, x2, x4, x6, x7, x8, x15, x30, x31, x32 sm2FieldElement
	copy(x1[:], in[:])
	sm2Square(&x2, in)
	sm2Mul(&x2, &x2, in)

	sm2SquareTimes(&x4, &x2, 2)
	sm2Mul(&x4, &x4, &x2)

	sm2SquareTimes(&x6, &x4, 2)
	sm2Mul(&x6, &x6, &x2)

	sm2Square(&x7, &x6)
	sm2Mul(&x7, &x7, in)

	sm2SquareTimes(&x8, &x7, 1)
	sm2Mul(&x8, &x8, in)

	sm2SquareTimes(&x15, &x8, 7)
	sm2Mul(&x15, &x15, &x7)

	sm2SquareTimes(&x30, &x15, 15)
	sm2Mul(&x30, &x30, &x15)

	sm2SquareTimes(&x31, &x30, 1)
	sm2Mul(&x31, &x31, in) //x31

	sm2SquareTimes(&x32, &x31, 1)
	sm2Mul(&x32, &x32, in)

	sm2SquareTimes(out, &x31, 33)
	sm2Mul(out, out, &x32)

	sm2SquareTimes(out, out, 32)
	sm2Mul(out, out, &x32)

	sm2SquareTimes(out, out, 32)
	sm2Mul(out, out, &x32)

	sm2SquareTimes(out, out, 32)
	sm2Mul(out, out, &x32)

	sm2SquareTimes(out, out, 32)
	sm2Mul(out, out, &x1)

	sm2SquareTimes(out, out, 62)
}

// out in Jacobian mode
func step1BaseScalar(out *point, in []params) int {
	sum := big.Int{}
	for i := range in {
		tmp := big.Int{}
		tmp.Mul(new(big.Int).SetUint64(*in[i].q), &in[i].s)
		tmp.Mod(&tmp, sm2.N)

		if i == 0 {
			sum = tmp
		} else {
			sum.Add(&sum, &tmp)
			sum.Mod(&sum, sm2.N)
		}
	}
	summ := [8]uint32{}
	sm2GetScalar2(&summ, sum.Bytes())
	out.x, out.y, out.z = &sm2FieldElement{}, &sm2FieldElement{}, &sm2FieldElement{}
	sm2BaseMult2(out.x, out.y, out.z, &summ)
	return scalarIsZero(&summ)
}

// in and out in Mont and Jacobian mode
func step2Scalar(out *point, in []params) (int, error) {
	pq := make(PriorityQueue, len(in))
	for i, v := range in {
		pq[i] = &Item{
			value:    v.r,
			priority: *v.q,
			index:    i}
		i++
	}
	heap.Init(&pq)
	for pq.Len() > 1 {
		item1 := heap.Pop(&pq).(*Item)
		item2 := heap.Pop(&pq).(*Item)
		newq := item1.priority - item2.priority

		var sum point
		sum.x, sum.y, sum.z = &sm2FieldElement{}, &sm2FieldElement{}, &sm2FieldElement{}
		double := point{}
		double.x, double.y, double.z = &sm2FieldElement{}, &sm2FieldElement{}, &sm2FieldElement{}
		sm2PointAdd(item1.value.x, item1.value.y, item1.value.z, item2.value.x, item2.value.y, item2.value.z, sum.x, sum.y, sum.z)
		sm2PointDouble(double.x, double.y, double.z, item1.value.x, item1.value.y, item1.value.z)
		copyCond(sum, double, resIsEqual(*item1.value, *item2.value))
		flag := 0
		if item2.priority == 0 {
			flag = 1
		}
		copyCond(sum, *item1.value, flag)
		if item1.priority > 0 {
			flag = 0
		}
		copyCond(sum, *item2.value, flag)

		if item2.priority > 0 {
			heap.Push(&pq, &Item{value: &sum, priority: item2.priority})
			/*ta, _ := sm2ToAffine(sum.x, sum.y, sum.z)
			fmt.Println(item2.priority, ta.Bytes())*/
		}

		if newq > 0 {
			heap.Push(&pq, &Item{value: item1.value, priority: newq})
			/*ta, _ := sm2ToAffine(item1.value.x, item1.value.y, item1.value.z)
			fmt.Println(newq, ta.Bytes())*/
		}
	}
	if pq.Len() < 1 {
		return 0, errors.New("step2 error!")
	}
	item := heap.Pop(&pq).(*Item)

	//fmt.Println(item.priority)
	d := [8]uint32{}
	temp := new(big.Int).SetUint64(item.priority)
	temp.Mod(temp, sm2.N)
	sm2GetScalar2(&d, temp.Bytes())
	out.x, out.y, out.z = &sm2FieldElement{}, &sm2FieldElement{}, &sm2FieldElement{}
	sm2PointToAffine(item.value.x, item.value.y, item.value.x, item.value.y, item.value.z)

	sm2ScalarMult2(out.x, out.y, out.z, item.value.x, item.value.y, &d)
	return scalarIsZero(&d), nil
}

func step3Scalar(out *point, in []params) int {
	cl := map[point]big.Int{}
	sum, double := point{}, point{}
	sum.x, sum.y, sum.z = &sm2FieldElement{}, &sm2FieldElement{}, &sm2FieldElement{}

	double.x, double.y, double.z = &sm2FieldElement{}, &sm2FieldElement{}, &sm2FieldElement{}
	for i, v := range in {
		res, temp := big.Int{}, big.Int{}
		res.Mul(new(big.Int).SetUint64(*in[i].q), &in[i].t)
		res.Mod(&res, sm2.N)
		_, ok := cl[*v.p]
		if !ok {
			cl[*v.p] = res
		} else {
			temp = cl[*v.p]
			res.Add(&res, &temp)
			res.Mod(&res, sm2.N)
			cl[*v.p] = temp
		}
	}
	preIsZero := 1
	for ii, v := range cl {
		i := point{}
		i.x, i.y, i.z = &sm2FieldElement{}, &sm2FieldElement{}, &sm2FieldElement{}
		vv := [8]uint32{}
		sm2GetScalar2(&vv, v.Bytes())
		nowIsZero := scalarIsZero(&vv)
		sm2PointToAffine(i.x, i.y, ii.x, ii.y, ii.z)
		sm2ScalarMult2(i.x, i.y, i.z, i.x, i.y, &vv)

		j := point{}
		j.x, j.y, j.z = &sm2FieldElement{}, &sm2FieldElement{}, &sm2FieldElement{}
		*j.x = *sum.x
		*j.z = *sum.z
		*j.y = *sum.y
		sm2PointAdd(sum.x, sum.y, sum.z, i.x, i.y, i.z, sum.x, sum.y, sum.z)
		sm2PointDouble(double.x, double.y, double.z, i.x, i.y, i.z)
		copyCond(sum, double, resIsEqual(j, i))
		copyCond(sum, i, preIsZero)
		copyCond(sum, j, nowIsZero)

		if preIsZero != 1 || nowIsZero != 1 {
			preIsZero = 0
		}
	}
	out.x = sum.x
	out.y = sum.y
	out.z = sum.z
	return preIsZero
}

func copyCond(out, in point, f int) {
	if f == 1 {
		*out.x = *in.x
		*out.y = *in.y
		*out.z = *in.z
	}
}

func resIsEqual(p1, p2 point) int {
	p1x, p1y := sm2ToAffine(p1.x, p1.y, p1.z)
	p2x, p2y := sm2ToAffine(p2.x, p2.y, p2.z)
	if p1x.Cmp(p2x) == 0 && p1y.Cmp(p2y) == 0 {
		return 1
	}
	return 0
	/*for j := 0; j < 8; j++ {
		if p1.x[j] != p2.x[j] || p1.y[j] != p2.y[j] || p1.z[j] != p2.z[j]{
			return 0
		}
	}
	return 1*/
}
