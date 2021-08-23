//+build amd64 arm64

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

	a := sm2Point{}
	res[0].q = uint64(1)
	res[1].q = uint64(1)
	step1BaseScalar(&a, res)
}

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

		a, b := sm2Point{}, sm2Point{}
		pre := 1
		res[0].q = uint64(1)
		res[1].q = uint64(1)
		step3Scalar(&a, res)
		for _, v := range res {

			double := sm2Point{}
			temp := [4]uint64{}
			smallOrderMul(&temp, &RRN, &v.q) //toMont
			orderMul(&temp, &temp, &v.t)
			getScalar(&temp)
			v.p.sm2ScalarMult(temp[:])

			now := scalarIsZero(&temp)
			j := b
			isEqual := sm2PointAdd2Asm(&b.xyz, &j.xyz, &v.p.xyz)
			sm2PointDouble2Asm(&double.xyz, &v.p.xyz)
			b.copyConditional(&double, isEqual)
			b.copyConditional(v.p, pre)
			b.copyConditional(&j, now)
			if pre != 1 || now != 1 {
				pre = 0
			}

		}
		a.toAffine()
		b.toAffine()
		assert.Equal(t, a.xyz[0][0], b.xyz[0][0])
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
		a, b := sm2Point{}, sm2Point{}
		pre := 1
		res[0].q = uint64(4)
		res[1].q = uint64(2)
		step2Scalar(&a, res)
		for _, v := range res {
			double := sm2Point{}
			temp := [4]uint64{v.q, 0, 0, 0}
			//getScalar(&temp)

			v.r.sm2ScalarMult(temp[:])
			now := scalarIsZero(&temp)
			j := b
			isEqual := sm2PointAdd2Asm(&b.xyz, &j.xyz, &v.r.xyz)
			sm2PointDouble2Asm(&double.xyz, &v.r.xyz)
			b.copyConditional(&double, isEqual)
			b.copyConditional(v.r, pre)
			b.copyConditional(&j, now)

			if pre != 1 || now != 1 {
				pre = 0
			}
		}

		a.toAffine()
		b.toAffine()
		fromMont(&a.xyz[0], &a.xyz[0])
		fromMont(&b.xyz[0], &b.xyz[0])
		assert.Equal(t, toBig(&b.xyz[0]), toBig(&a.xyz[0]))
	}

}
func BenchmarkStep2Scalar2(b *testing.B) {
	for i := 0; i < 1; i++ {
		X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
		Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
		h := sm3.SignHashSM3(X, Y, []byte(msg))
		X = append(X, Y...)
		X = append([]byte{0}, X...)
		sig2 := []byte{48, 69, 2, 32, 13, 190, 159, 134, 254, 112, 95, 175, 247, 34, 5, 132, 150, 56, 225, 46, 210, 30, 177,
			157, 21, 183, 236, 17, 65, 204, 237, 255, 46, 57, 182, 207, 2, 33, 0, 196, 252, 200, 58, 188, 213, 181, 112, 101,
			211, 201, 31, 210, 140, 96, 168, 47, 81, 168, 169, 229, 100, 44, 65, 148, 114, 181, 46, 68, 141, 26, 225, 0, 0}

		res := []params{}
		preStep(&res, [][]byte{X, X}, [][]byte{sig2, sig2}, [][]byte{h, h})
		a := sm2Point{}
		res[0].q = uint64(4)
		res[1].q = uint64(2)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			step2Scalar(&a, res)
		}
	}

} //61274 ns

func BenchmarkStep2Scalar(b *testing.B) {
	X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
	Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
	h := sm3.SignHashSM3(X, Y, []byte(msg))
	X = append(X, Y...)
	X = append([]byte{0}, X...)
	sig2 := []byte{48, 69, 2, 32, 13, 190, 159, 134, 254, 112, 95, 175, 247, 34, 5, 132, 150, 56, 225, 46, 210, 30, 177,
		157, 21, 183, 236, 17, 65, 204, 237, 255, 46, 57, 182, 207, 2, 33, 0, 196, 252, 200, 58, 188, 213, 181, 112, 101,
		211, 201, 31, 210, 140, 96, 168, 47, 81, 168, 169, 229, 100, 44, 65, 148, 114, 181, 46, 68, 141, 26, 225, 0, 0}

	res := []params{}
	sig := [][]byte{}
	pk := [][]byte{}
	msg := [][]byte{}
	for j := 0; j < 64; j++ {
		sig = append(sig, sig2)
		pk = append(pk, X)
		msg = append(msg, h)
	}
	preStep(&res, pk, sig, msg)
	a := sm2Point{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		step2Scalar(&a, res)
	}
} // 782512 ns

func BenchmarkStep1(b *testing.B) {
	X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
	Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
	h := sm3.SignHashSM3(X, Y, []byte(msg))
	X = append(X, Y...)
	X = append([]byte{0}, X...)
	sig2 := []byte{48, 69, 2, 32, 13, 190, 159, 134, 254, 112, 95, 175, 247, 34, 5, 132, 150, 56, 225, 46, 210, 30, 177,
		157, 21, 183, 236, 17, 65, 204, 237, 255, 46, 57, 182, 207, 2, 33, 0, 196, 252, 200, 58, 188, 213, 181, 112, 101,
		211, 201, 31, 210, 140, 96, 168, 47, 81, 168, 169, 229, 100, 44, 65, 148, 114, 181, 46, 68, 141, 26, 225, 0, 0}

	res := []params{}
	sig := [][]byte{}
	pk := [][]byte{}
	msg := [][]byte{}
	for j := 0; j < 64; j++ {
		sig = append(sig, sig2)
		pk = append(pk, X)
		msg = append(msg, h)
	}
	preStep(&res, pk, sig, msg)
	a := sm2Point{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		step1BaseScalar(&a, res)
	}
} // 16576 ns

func BenchmarkStep3(b *testing.B) {
	X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
	Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
	h := sm3.SignHashSM3(X, Y, []byte(msg))
	X = append(X, Y...)
	X = append([]byte{0}, X...)
	sig2 := []byte{48, 69, 2, 32, 13, 190, 159, 134, 254, 112, 95, 175, 247, 34, 5, 132, 150, 56, 225, 46, 210, 30, 177,
		157, 21, 183, 236, 17, 65, 204, 237, 255, 46, 57, 182, 207, 2, 33, 0, 196, 252, 200, 58, 188, 213, 181, 112, 101,
		211, 201, 31, 210, 140, 96, 168, 47, 81, 168, 169, 229, 100, 44, 65, 148, 114, 181, 46, 68, 141, 26, 225, 0, 0}

	res := []params{}
	sig := [][]byte{}
	pk := [][]byte{}
	msg := [][]byte{}
	for j := 0; j < 64; j++ {
		sig = append(sig, sig2)
		pk = append(pk, X)
		msg = append(msg, h)
	}
	preStep(&res, pk, sig, msg)
	a := sm2Point{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		step3Scalar(&a, res)
	}
} // 72053 ns

func BenchmarkPreStep(b *testing.B) {
	X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
	Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
	h := sm3.SignHashSM3(X, Y, []byte(msg))
	X = append(X, Y...)
	X = append([]byte{0}, X...)
	sig2 := []byte{48, 69, 2, 32, 13, 190, 159, 134, 254, 112, 95, 175, 247, 34, 5, 132, 150, 56, 225, 46, 210, 30, 177,
		157, 21, 183, 236, 17, 65, 204, 237, 255, 46, 57, 182, 207, 2, 33, 0, 196, 252, 200, 58, 188, 213, 181, 112, 101,
		211, 201, 31, 210, 140, 96, 168, 47, 81, 168, 169, 229, 100, 44, 65, 148, 114, 181, 46, 68, 141, 26, 225, 0, 0}

	res := []params{}
	sig := [][]byte{}
	pk := [][]byte{}
	msg := [][]byte{}
	for j := 0; j < 64; j++ {
		sig = append(sig, sig2)
		pk = append(pk, X)
		msg = append(msg, h)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		preStep(&res, pk, sig, msg)
	}
} // 601154 ns

func BenchmarkBatchVerify(b *testing.B) {
	sig2 := []byte{48, 69, 2, 32, 13, 190, 159, 134, 254, 112, 95, 175, 247, 34, 5, 132, 150, 56, 225, 46, 210, 30, 177,
		157, 21, 183, 236, 17, 65, 204, 237, 255, 46, 57, 182, 207, 2, 33, 0, 196, 252, 200, 58, 188, 213, 181, 112, 101,
		211, 201, 31, 210, 140, 96, 168, 47, 81, 168, 169, 229, 100, 44, 65, 148, 114, 181, 46, 68, 141, 26, 225, 0, 0}
	X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
	Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
	h := sm3.SignHashSM3(X, Y, []byte(msg))

	X = append(X, Y...)
	X = append([]byte{0}, X...)
	sig := [][]byte{}
	pk := [][]byte{}
	msg := [][]byte{}
	for j := 0; j < 64; j++ {
		sig = append(sig, sig2)
		pk = append(pk, X)
		msg = append(msg, h)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BatchVerify(pk, sig, msg)
	}
} //1462635 ns

func BenchmarkBatchVerify2(b *testing.B) {
	sig2 := []byte{48, 69, 2, 32, 13, 190, 159, 134, 254, 112, 95, 175, 247, 34, 5, 132, 150, 56, 225, 46, 210, 30, 177,
		157, 21, 183, 236, 17, 65, 204, 237, 255, 46, 57, 182, 207, 2, 33, 0, 196, 252, 200, 58, 188, 213, 181, 112, 101,
		211, 201, 31, 210, 140, 96, 168, 47, 81, 168, 169, 229, 100, 44, 65, 148, 114, 181, 46, 68, 141, 26, 225, 0, 0}
	X, _ := hex.DecodeString("86d3205ed0c3db8ef35a74b6bf924cbef75988e835f65f422884e3b1c8cdbde1")
	Y, _ := hex.DecodeString("ea7eee5e7ff177622c3081aea9375d3cfec41867298261aae8f8e1434c9e81f0")
	h := sm3.SignHashSM3(X, Y, []byte(msg))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 64; j++ {
			VerifySignature(sig2, h, X, Y) //79499 ns
		}
	}
} //4798976 ns

func TestBatchVerify(t *testing.T) {
	for i := 0; i < 1000; i++ {
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
			sig1, flag, _ := Sign(h, rand.Reader, Key)
			ret := make([]byte, len(sig1)+1)
			ret[0] = flag
			copy(ret[1:], sig1)
			sig = append(sig, ret)
			/*sig2 := []byte{48, 69, 2 ,32, 13, 190, 159, 134, 254, 112, 95, 175, 247, 34, 5 ,132, 150, 56, 225, 46, 210, 30, 177,
			157, 21, 183, 236, 17, 65, 204, 237, 255, 46, 57, 182, 207, 2 ,33, 0, 196, 252, 200, 58, 188, 213, 181, 112, 101,
			211, 201, 31, 210, 140, 96, 168, 47, 81 ,168, 169, 229, 100, 44, 65, 148, 114, 181, 46, 68, 141, 26, 225, 0, 0}
			*/
			pk = append(pk, X)
			msg = append(msg, h)
		}
		err := BatchVerify(pk, sig, msg)
		assert.Nil(t, err)
	}
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

func TestSmallOrderMul(t *testing.T) {
	sm2N, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFF7203DF6B21C6052B53BBF40939D54123", 16)
	limit64, _ := new(big.Int).SetString("FFFFFFFFFFFFFFFF", 16)
	e256Inv, _ := new(big.Int).SetString("10000000000000000000000000000000000000000000000000000000000000000", 16)
	e256Inv.ModInverse(e256Inv, sm2N)
	for i := uint64(0); i < 0xffff; i++ {
		bigA, _ := rand.Int(rand.Reader, sm2N)
		bigB, _ := rand.Int(rand.Reader, limit64)
		product := new(big.Int).Mul(bigA, bigB)
		product.Mul(product, e256Inv)
		product.Mod(product, sm2N)
		res := new([4]uint64)
		fpA, fpB := new([4]uint64), new([1]uint64)
		for i, v := range bigA.Bits() {
			fpA[i] = uint64(v)
		}
		for i, v := range bigB.Bits() {
			fpB[i] = uint64(v)
		}
		smallOrderMul(res, fpA, &fpB[0])
		fpProduct := new(big.Int).SetBits(
			[]big.Word{big.Word(res[0]), big.Word(res[1]), big.Word(res[2]), big.Word(res[3])})

		if fpProduct.Cmp(product) != 0 {
			t.Error("input:", bigA.Text(16), " ", bigB.Text(16))
			t.Error("want :", product.Text(16))
			t.Error("get  :", fpProduct.Text(16))
			return
		}
	}
}

func TestInvertForY(t *testing.T) {
	sm2P, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)
	//uPlusOne, _ := new(big.Int).SetString("3fffffffbfffffffffffffffffffffffffffffffc00000004000000000000000", 16)
	e256Inv, _ := new(big.Int).SetString("10000000000000000000000000000000000000000000000000000000000000000", 16)
	e256Inv.ModInverse(e256Inv, sm2P)

	bigA := big.NewInt(100)

	fpA := [4]uint64{}
	for i, v := range bigA.Bits() {
		fpA[i] = uint64(v)
	}
	p256Mul(&fpA, &fpA, &RR)
	invertForY(&fpA, &fpA)
	p256Mul(&fpA, &fpA, &one)
	fp := new(big.Int).SetBits(
		[]big.Word{big.Word(fpA[0]), big.Word(fpA[1]), big.Word(fpA[2]), big.Word(fpA[3])})

	if fp.Cmp(big.NewInt(10)) != 0 {
		t.Errorf("want : %s", big.NewInt(10).Text(16))
		t.Errorf("get  : %s", fp.Text(16))
	}
}

func TestGetY(t *testing.T) {
	t.Run("bigger y", func(t *testing.T) {
		//sm2P, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)

		e := Sm2()
		a, b := e.ScalarBaseMult(big.NewInt(4).Bytes()) //可使用内部函数替换
		x, y := [4]uint64{}, [4]uint64{}
		fromBig(&x, maybeReduceModP(a))

		//p256Mul(&x, &x, &RR)
		getY(&y, &x, 1)

		fp := new(big.Int).SetBits(
			[]big.Word{big.Word(y[0]), big.Word(y[1]), big.Word(y[2]), big.Word(y[3])})

		if fp.Cmp(b) != 0 {
			t.Errorf("want : %s", b.Text(16))
			t.Errorf("get  : %s", fp.Text(16))
		}
	})

	t.Run("smaller y", func(t *testing.T) {
		//sm2P, _ := new(big.Int).SetString("FFFFFFFEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00000000FFFFFFFFFFFFFFFF", 16)

		e := Sm2()
		a, b := e.ScalarBaseMult(big.NewInt(2).Bytes()) //可使用内部函数替换
		x, y := [4]uint64{}, [4]uint64{}
		fromBig(&x, maybeReduceModP(a))

		//p256Mul(&x, &x, &RR)
		getY(&y, &x, 0)

		fp := new(big.Int).SetBits(
			[]big.Word{big.Word(y[0]), big.Word(y[1]), big.Word(y[2]), big.Word(y[3])})

		if fp.Cmp(b) != 0 {
			t.Errorf("want : %s", b.Text(16))
			t.Errorf("get  : %s", fp.Text(16))
		}
	})

}
