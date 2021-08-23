//+build !amd64,!arm64

package sm2

import (
	"container/heap"
	"crypto/rand"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"github.com/ultramesh/crypto-gm/internal/sm3"
	"math/big"
	"testing"
)

func TestBatchVerify(t *testing.T) {
	for i := 0; i < 1; i++ {
		Key, _ := hex.DecodeString("6332a6b9f834f5c25df0555ff84b2c0cd278f42457bb95534faa4bae0608f537")
		X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
		Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
		h := sm3.SignHashSM3(X, Y, []byte(msg))

		X = append(X, Y...)
		X = append([]byte{0}, X...)

		var sig [][]byte
		var pk [][]byte
		var msg [][]byte
		for j := 0; j < 64; j++ {
			sig2, flag, _ := Sign(h, rand.Reader, Key)
			/*sig2 := []byte{48, 69, 2 ,32, 13, 190, 159, 134, 254, 112, 95, 175, 247, 34, 5 ,132, 150, 56, 225, 46, 210, 30, 177,
			157, 21, 183, 236, 17, 65, 204, 237, 255, 46, 57, 182, 207, 2 ,33, 0, 196, 252, 200, 58, 188, 213, 181, 112, 101,
			211, 201, 31, 210, 140, 96, 168, 47, 81 ,168, 169, 229, 100, 44, 65, 148, 114, 181, 46, 68, 141, 26, 225}*/
			ret := make([]byte, len(sig2)+1)
			ret[0] = flag
			copy(ret[1:], sig2)
			sig = append(sig, ret)

			pk = append(pk, X)
			msg = append(msg, h)
		}
		err := BatchVerify(pk, sig, msg)
		assert.Nil(t, err)
	}
}

//sum(s*q)*G + sum(t*q*P) = sum(q*R)
func Test2(t *testing.T) {
	X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
	Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
	h := sm3.SignHashSM3(X, Y, []byte(msg))
	X = append(X, Y...)
	X = append([]byte{0}, X...)
	sig2 := []byte{0, 48, 69, 2, 32, 13, 190, 159, 134, 254, 112, 95, 175, 247, 34, 5, 132, 150, 56, 225, 46, 210, 30, 177,
		157, 21, 183, 236, 17, 65, 204, 237, 255, 46, 57, 182, 207, 2, 33, 0, 196, 252, 200, 58, 188, 213, 181, 112, 101,
		211, 201, 31, 210, 140, 96, 168, 47, 81, 168, 169, 229, 100, 44, 65, 148, 114, 181, 46, 68, 141, 26, 225}

	res := []params{}
	err := preStep(&res, [][]byte{X, X}, [][]byte{sig2, sig2}, [][]byte{h, h})
	assert.Nil(t, err)
}

func TestPriorityQueue_Push(t *testing.T) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{
		priority: 4,
		index:    1})
	heap.Push(&pq, &Item{
		priority: 10,
		index:    1})
	heap.Push(&pq, &Item{
		priority: 11,
		index:    1})
	for pq.Len() > 1 {
		item1 := heap.Pop(&pq).(*Item)
		item2 := heap.Pop(&pq).(*Item)
		newq := item1.priority - item2.priority
		heap.Push(&pq, &Item{priority: item2.priority})

		if newq > 0 {
			heap.Push(&pq, &Item{value: item1.value, priority: newq})
		}
	}
	i := heap.Pop(&pq).(*Item)
	assert.Equal(t, uint64(1), i.priority)
}

func TestPreStep(t *testing.T) {
	Key, _ := hex.DecodeString("6332a6b9f834f5c25df0555ff84b2c0cd278f42457bb95534faa4bae0608f537")
	X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
	Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
	h := sm3.SignHashSM3(X, Y, []byte(msg))
	X = append(X, Y...)
	X = append([]byte{0}, X...)
	sig2, flag, err := Sign(h, rand.Reader, Key)
	sig2 = []byte{48, 68, 2, 32, 35, 44, 162, 42, 54, 17, 34, 192, 46, 104, 117, 255, 246, 66, 217, 174, 9, 109, 179, 35, 4, 93, 3, 18, 133,
		132, 217, 114, 2, 178, 180, 212, 2, 32, 15, 80, 157, 49, 108, 74, 97, 132, 254, 112, 20, 254, 115, 171, 129, 130, 217, 183, 7, 54, 154,
		121, 217, 209, 82, 8, 139, 153, 168, 190, 119, 69}
	flag = 1
	assert.Nil(t, err)
	ret := make([]byte, len(sig2)+1)
	ret[0] = flag
	copy(ret[1:], sig2)
	var out []params
	assert.Nil(t, preStep(&out, [][]byte{X}, [][]byte{ret}, [][]byte{h}))
}

func TestGetY(t *testing.T) {
	t.Run("bigger y", func(t *testing.T) {
		//sm2P, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)

		e := Sm2()
		a, b := e.ScalarBaseMult(big.NewInt(4).Bytes()) //可使用内部函数替换
		y := big.Int{}
		y = *b

		//p256Mul(&x, &x, &RR)
		b, _ = getY(a, 1)

		if y.Cmp(b) != 0 {
			t.Errorf("want : %s", y.Text(16))
			t.Errorf("get  : %s", b.Text(16))
		}
	})

	t.Run("smaller y", func(t *testing.T) {
		//sm2P, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)

		e := Sm2()
		a, b := e.ScalarBaseMult(big.NewInt(2).Bytes()) //可使用内部函数替换
		y := big.Int{}
		y = *b

		//p256Mul(&x, &x, &RR)
		b, _ = getY(a, 0)

		if y.Cmp(b) != 0 {
			t.Errorf("want : %s", y.Text(16))
			t.Errorf("get  : %s", b.Text(16))
		}
	})

}

//sum(t*q*P)
//sum(t*q*P)
func TestStep3Scalar(t *testing.T) {
	for i := 0; i < 1; i++ {
		X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
		Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
		h := sm3.SignHashSM3(X, Y, []byte(msg))
		X = append(X, Y...)
		X = append([]byte{0}, X...)
		sig2 := []byte{0, 48, 69, 2, 32, 13, 190, 159, 134, 254, 112, 95, 175, 247, 34, 5, 132, 150, 56, 225, 46, 210, 30, 177,
			157, 21, 183, 236, 17, 65, 204, 237, 255, 46, 57, 182, 207, 2, 33, 0, 196, 252, 200, 58, 188, 213, 181, 112, 101,
			211, 201, 31, 210, 140, 96, 168, 47, 81, 168, 169, 229, 100, 44, 65, 148, 114, 181, 46, 68, 141, 26, 225}

		res := []params{}
		err := preStep(&res, [][]byte{X, X}, [][]byte{sig2, sig2}, [][]byte{h, h})
		assert.Nil(t, err)

		a := point{}
		a.x = &sm2FieldElement{}
		a.y = &sm2FieldElement{}
		a.z = &sm2FieldElement{}
		b := point{}
		b.x = &sm2FieldElement{0}
		b.y = &sm2FieldElement{0}
		b.z = &sm2FieldElement{0}

		*res[0].q = uint64(1)
		*res[1].q = uint64(1)

		step3Scalar(&a, res)
		for i, v := range res {
			d := [8]uint32{}
			temp := big.NewInt(1)
			temp.Mul(new(big.Int).SetUint64(*v.q), &v.t)
			temp.Mod(temp, sm2.N)
			sm2GetScalar2(&d, temp.Bytes())
			//fmt.Println(sm2ToBig(v.p.x).Bytes())
			//			ta, _ := sm2ToAffine(v.p.x, v.p.y, v.p.z)
			//        fmt.Println(sm2ToBig(v.p.x).Bytes())
			//fmt.Println(d)
			sm2PointToAffine(v.p.x, v.p.y, v.p.x, v.p.y, v.p.z)
			sm2ScalarMult2(v.p.x, v.p.y, v.p.z, v.p.x, v.p.y, &d)
			//fmt.Println(sm2ToBig(v.p.x).Bytes())
			//      			j := point{}
			//			j = b
			//ba, _ := sm2ToAffine(b.x, b.y, b.z)
			//fmt.Println(ba.Bytes())
			if i == 0 {
				b.x = v.p.x
				b.y = v.p.y
				b.z = v.p.z
			} else {
				sm2PointDouble(b.x, b.y, b.z, v.p.x, v.p.y, v.p.z)
			}

			//ba, _ = sm2ToAffine(b.x, b.y, b.z)
			//fmt.Println(ba.Bytes())
		}
		aa, _ := sm2ToAffine(a.x, a.y, a.z)
		ba, _ := sm2ToAffine(b.x, b.y, b.z)
		//fmt.Println(aa.Bytes())
		//fmt.Println(ba.Bytes())
		assert.Equal(t, aa, ba)
	}
}

//sum(q*R)
func TestStep2Scalar(t *testing.T) {
	for i := 0; i < 1; i++ {
		X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
		Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
		h := sm3.SignHashSM3(X, Y, []byte(msg))
		X = append(X, Y...)
		X = append([]byte{0}, X...)
		sig2 := []byte{0, 48, 69, 2, 32, 13, 190, 159, 134, 254, 112, 95, 175, 247, 34, 5, 132, 150, 56, 225, 46, 210, 30, 177,
			157, 21, 183, 236, 17, 65, 204, 237, 255, 46, 57, 182, 207, 2, 33, 0, 196, 252, 200, 58, 188, 213, 181, 112, 101,
			211, 201, 31, 210, 140, 96, 168, 47, 81, 168, 169, 229, 100, 44, 65, 148, 114, 181, 46, 68, 141, 26, 225}

		res := []params{}
		err := preStep(&res, [][]byte{X, X}, [][]byte{sig2, sig2}, [][]byte{h, h})
		assert.Nil(t, err)
		a, b := point{}, point{}

		b.x, b.y, b.z = &sm2FieldElement{}, &sm2FieldElement{}, &sm2FieldElement{}
		pre := 1
		*res[0].q = 5
		*res[1].q = 6
		step2Scalar(&a, res)
		for _, v := range res {
			double := point{}

			double.x, double.y, double.z = &sm2FieldElement{}, &sm2FieldElement{}, &sm2FieldElement{}
			d := [8]uint32{}
			temp := new(big.Int).SetUint64(*v.q)
			temp.Mod(temp, sm2.N)
			sm2GetScalar2(&d, temp.Bytes())

			sm2PointToAffine(v.r.x, v.r.y, v.r.x, v.r.y, v.r.z)

			sm2ScalarMult2(v.r.x, v.r.y, v.r.z, v.r.x, v.r.y, &d)

			now := scalarIsZero(&d)
			j := point{}
			j.x, j.y, j.z = &sm2FieldElement{}, &sm2FieldElement{}, &sm2FieldElement{}
			*j.x = *b.x
			*j.z = *b.z
			*j.y = *b.y
			sm2PointAdd(v.r.x, v.r.y, v.r.z, j.x, j.y, j.z, b.x, b.y, b.z)
			sm2PointDouble(double.x, double.y, double.z, j.x, j.y, j.z)
			copyCond(b, double, resIsEqual(j, *v.r))
			copyCond(b, *v.r, pre)
			copyCond(b, j, now)
			if pre != 1 || now != 1 {
				pre = 0
			}
		}
		aa, _ := sm2ToAffine(a.x, a.y, a.z)
		ba, _ := sm2ToAffine(b.x, b.y, b.z)
		//fmt.Println(aa.Bytes())
		//fmt.Println(ba.Bytes())
		assert.Equal(t, aa, ba)
	}
}

func TestStep1(t *testing.T) {
	X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
	Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
	h := sm3.SignHashSM3(X, Y, []byte(msg))
	X = append(X, Y...)
	X = append([]byte{0}, X...)
	sig2 := []byte{0, 48, 69, 2, 32, 13, 190, 159, 134, 254, 112, 95, 175, 247, 34, 5, 132, 150, 56, 225, 46, 210, 30, 177,
		157, 21, 183, 236, 17, 65, 204, 237, 255, 46, 57, 182, 207, 2, 33, 0, 196, 252, 200, 58, 188, 213, 181, 112, 101,
		211, 201, 31, 210, 140, 96, 168, 47, 81, 168, 169, 229, 100, 44, 65, 148, 114, 181, 46, 68, 141, 26, 225}

	res := []params{}
	err := preStep(&res, [][]byte{X, X}, [][]byte{sig2, sig2}, [][]byte{h, h})
	assert.Nil(t, err)

	a := point{}
	*res[0].q = uint64(1)
	*res[1].q = uint64(1)
	step1BaseScalar(&a, res)
	//aa, _ := sm2ToAffine(a.x, a.y, a.z)
	//fmt.Println(aa.Bytes())
}
