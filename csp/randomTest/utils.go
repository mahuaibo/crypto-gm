package randomtest

import (
	"math"
)

//Bytes2Bits 将bytes切片转换为bit切片，这里是大端序
func Bytes2Bits(data []byte) []byte {
	buf := make([]byte, len(data)*8)
	for i := range data {
		buf[7+i*8] = data[i] & 0x01 >> 0
		buf[6+i*8] = data[i] & 0x02 >> 1
		buf[5+i*8] = data[i] & 0x04 >> 2
		buf[4+i*8] = data[i] & 0x08 >> 3
		buf[3+i*8] = data[i] & 0x10 >> 4
		buf[2+i*8] = data[i] & 0x20 >> 5
		buf[1+i*8] = data[i] & 0x40 >> 6
		buf[0+i*8] = data[i] & 0x80 >> 7
	}

	return buf
}

// 求矩阵的秩,这里的矩阵很特别，是只含有 0 1 的序列
// 所以只需要判断关键位置是否为 1, 即可判断矩阵的秩
func rank(matrix [][]int, m int) int {
	temp := make([][]int, m)
	for i := range temp {
		temp[i] = make([]int, m)
		copy(temp[i], matrix[i])
	}

	rowEchelon(temp, m)
	rank := 0
	// 判断经过初等变换过的矩阵每一行是否至少有一个 1
	for i := 0; i < m; i++ {
		notZero := false
		for j := 0; j < m; j++ {
			if temp[i][j] != 0 {
				notZero = true
				break
			}
		}
		if notZero {
			rank++
		}
	}
	return rank
}

// 对矩阵进行初等变换得到阶梯型矩阵
func rowEchelon(matrix [][]int, m int) (a [][]int) {
	pivotstartrow := 0
	pivotstartcol := 0
	pivotrow := 0
	for i := 0; i < m; i++ {
		found := false
		// 寻找第一个 matrix[k][pivotstartcol] 元素为 1 的行数 k ，找到则将 found 置为 true
		// 并将行数 k 赋值给 pivotrow ，同时跳出循环， 此时为 1 的元素是 matrix[pivotrow][pivotstartcol]
		for k := pivotstartrow; k < m; k++ {
			if matrix[k][pivotstartcol] == 1 {
				found = true
				pivotrow = k
				break
			}
		}
		if found {
			// 如果找到的行数 pivotrow 不是 pivotstartrow，那么上下交换这两行的所有元素
			// 此时为 1 的元素是 matrix[pivotstartrow][pivotstartcol]
			if pivotrow != pivotstartrow {
				for k := 0; k < m; k++ {
					matrix[pivotrow][k] ^= matrix[pivotstartrow][k]
					matrix[pivotstartrow][k] ^= matrix[pivotrow][k]
					matrix[pivotrow][k] ^= matrix[pivotstartrow][k]
				}
			}
			// 判断矩阵同一列中第 pivotstartrow 行后是否还有1，如果有 1 则将该行所有元素与第 pivotstartrow 行做异或以消除重复的1
			for j := pivotstartrow + 1; j < m; j++ {
				if matrix[j][pivotstartcol] == 1 {
					for k := 0; k < m; k++ {
						matrix[j][k] = matrix[pivotstartrow][k] ^ matrix[j][k]
					}
				}
			}

			pivotstartcol++
			pivotstartrow++
		} else {
			pivotstartcol++
		}

	}
	return matrix
}

// 标准正态分布函数，利用其和误差函数的关系得到
// https://www.cnblogs.com/htj10/p/8621771.html
func normalCDF(x float64) float64 {
	return (1.0 + math.Erf(x/math.Sqrt(2))) / 2
}

var mAXLOG = 7.09782712893383996732224e2
var biginv = 2.22044604925031308085e-16
var big = 4.503599627370496e15
var mACHEP = 1.11022302462515654042e-16

func igam(a float64, x float64) float64 {
	var ans, ax, c, r float64

	if (x <= 0.0) || (a <= 0.0) {
		return 0.0
	}

	if (x > 1.0) && (x > a) {
		return 1.e0 - igamc(a, x)
	}
	lg, _ := math.Lgamma(a)
	ax = a*math.Log(x) - x - lg
	if ax < -mAXLOG {
		return 0.0
	}
	ax = math.Exp(ax)

	r = a
	c = 1.0
	ans = 1.0

	for {
		r += 1.0
		c *= x / r
		ans += c
		if c/ans <= mACHEP {
			break
		}
	}

	return ans * ax / a
}

// 不完全伽马函数
func igamc(a float64, x float64) float64 {
	var ans, ax, c, yc, r, t, y, z float64
	var pk, pkm1, pkm2, qk, qkm1, qkm2 float64

	if (x <= 0) || (a <= 0) {
		return 1.0
	}
	if (x < 1.0) || (x < a) {
		return 1.e0 - igam(a, x)
	}
	lg, _ := math.Lgamma(a)
	ax = a*math.Log(x) - x - lg

	if ax < -mAXLOG {
		return 0.0
	}
	ax = math.Exp(ax)

	y = 1.0 - a
	z = x + y + 1.0
	c = 0.0
	pkm2 = 1.0
	qkm2 = x
	pkm1 = x + 1.0
	qkm1 = z * x
	ans = pkm1 / qkm1

	for {
		c += 1.0
		y += 1.0
		z += 2.0
		yc = y * c
		pk = pkm1*z - pkm2*yc
		qk = qkm1*z - qkm2*yc
		if qk != 0 {
			r = pk / qk
			t = math.Abs((ans - r) / r)
			ans = r
		} else {
			t = 1.0
		}
		pkm2 = pkm1
		pkm1 = pk
		qkm2 = qkm1
		qkm1 = qk
		if math.Abs(pk) > big {
			pkm2 *= biginv
			pkm1 *= biginv
			qkm2 *= biginv
			qkm1 *= biginv
		}
		if t <= mACHEP {
			break
		}
	}
	return ans * ax

}

// 利用 Berlekamp-Massey 算法计算子序列的线性复杂度，实际上的输入是一个 0 1 的序列
// 二元序列 S 的线性复杂度定义为生成它的最短线性反馈移位寄存器的长度
// https://en.wikipedia.org/wiki/Berlekamp%E2%80%93Massey_algorithm?utm_source=hacpai.com
func linearComplexityTest(s []byte, M int) int {
	// Nn 代表算法中的 n，实际是指元素的下标
	Nn := 0
	// 分配初始值 L=0, m=-1
	L := 0
	m := -1
	d := 0
	B := make([]int, M)
	C := make([]int, M)
	P := make([]int, M)
	T := make([]int, M)
	// 分配初始值 b[0]=1, c[0]=1, 除此之外均为 0
	C[0] = 1
	B[0] = 1
	for Nn < M {
		// 求 d 的值，异或操作相当于两数相加再对 2 取余
		d = int(s[Nn])
		for i := 1; i <= L; i++ {
			d += C[i] * int(s[Nn-i])
		}
		d = d % 2
		// 如果计算出的 d=0 ，那么 c 已经是一个多项式了，可以消除 n-L 至 n
		// 计算出的 d!=0, 则进行以下操作
		if d == 1 {
			// 让 t 变成 c 的副本
			for i := 0; i < M; i++ {
				T[i] = C[i]
				P[i] = 0
			}
			// 取临时的 P 是为了方便计算，当然也可以不使用 P，而直接根据 C[i]、B[i] 计算出新的 C[i]
			for j := 0; j < M; j++ {
				if B[j] == 1 {
					P[j+Nn-m] = 1
				}
			}
			// 计算新的 C[i] 的值
			for i := 0; i < M; i++ {
				C[i] = (C[i] + P[i]) % 2
			}
			// L <= n/2， 则进行以下操作
			if L <= Nn/2 {
				L = Nn + 1 - L
				m = Nn
				for i := 0; i < M; i++ {
					B[i] = T[i]
				}
			}
		}
		Nn++
	}
	return L
}

// 影响因子 c(L,K)
func mutFactorC(L int, K int) float64 {
	var v float64
	v = 0.7
	v -= 0.8 / float64(L)
	v += (4.0 + 32.0/float64(L)) * (math.Pow(float64(K), -3.0/float64(L)) / 15.0)
	return v
}

// 将新序列长度扩展到2的指数倍，长度不够的补0
func pow2DoubleArr(data []float64) []float64 {
	var newData []float64
	dataLength := len(data)

	sumNum := 2
	for {
		sumNum *= 2
		if sumNum >= dataLength {
			break
		}
	}
	addLength := sumNum - dataLength

	if addLength != 0 {
		newData = make([]float64, sumNum)
		for i := 0; i < dataLength; i++ {
			newData[i] = data[i]
		}
		for i := dataLength; i < sumNum; i++ {
			newData[i] = 0
		}
	} else {
		newData = data
	}
	return newData
}
