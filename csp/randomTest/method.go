package randomtest

import (
	"fmt"
	"github.com/mjibson/go-dsp/fft"
	"math"
)

//TestFunc all test function
var TestFunc = []struct {
	Method func([]byte, float64) bool
	Name   string
}{
	{MonobitFrequencyTest, "单比特频数检测"},
	{FrequencyTestWithABlock, "块内频数检测"},
	{PokerTest, "扑克检测"},
	{SerialTest, "重叠子序列检测"},
	{RunsTest, "游程总数检测"},
	{RunsDistributionTest, "游程分布检测"},
	{TestForTheLongestRunOfOnesInABlock, "块内最大1游程检测"},
	{BinaryDerivativeTest, "二元推导检测"},
	{AutocorrelationTest, "自相关检测"},
	{BinaryMatrixRankTest, "矩阵秩检测"},
	{CumulativeTest, "累加和检测"},
	{ApproximateEntropyTest, "近似熵检测"},
	{LinearComplexityTest, "线性复杂度检测"},
	{MaurersUniversalTest, "Maurer通用统计检测"},
	{DiscreteFourierTransformTest, "离散傅立叶检测"},
}

//MonobitFrequencyTest 单比特频数检测
func MonobitFrequencyTest(data []byte, alpha float64) bool {
	if len(data) == 0 {
		fmt.Println("data cannot be empty")
		return false
	}
	bits := Bytes2Bits(data)
	S := 0
	for i := 0; i < len(bits); i++ {
		if bits[i] == 1 {
			S++
		} else {
			S--
		}
	}
	V := math.Abs(float64(S)) / math.Sqrt(float64(len(bits)))
	P := math.Erfc(V / math.Sqrt(2))
	//fmt.Printf("P = %v\n",P)
	if P >= alpha {
		return true
	}
	return false
}

//FrequencyTestWithABlock 块内频数检测
func FrequencyTestWithABlock(data []byte, alpha float64) bool {
	if len(data) == 0 {
		fmt.Println("data cannot be empty")
		return false
	}
	bits := Bytes2Bits(data)
	n := len(bits)
	m := math.Min(float64(n), 100)
	N := (int)(float64(n) / m)
	var V float64
	var Pi float64

	for i := 0; i < N; i++ {
		Pi = 0
		for j := 0; j < int(m); j++ {
			if bits[i*int(m)+j] == 1 {
				Pi++
			}
		}
		Pi = Pi / m
		V += (Pi - 0.5) * (Pi - 0.5)
	}
	V *= 4.0 * m
	P := igamc(float64(N)/2.0, V/2.0)
	//fmt.Printf("P = %v\n",P)
	if P >= alpha {
		return true
	}
	return false
}

// PokerTest 扑克检测
func PokerTest(data []byte, alpha float64) bool {
	if len(data) == 0 {
		fmt.Println("data cannot be empty")
		return false
	}
	bits := Bytes2Bits(data)
	m := 4 // change m=2 to m=4
	n := len(bits)
	N := n / m
	patterns := make([]int, 1<<uint(m))
	var V float64
	var P float64
	for i := 0; i < N; i++ {
		tmp := 0
		for j := 0; j < m; j++ {
			tmp <<= 1
			if bits[i*m+j] == 1 {
				tmp++
			}
		}
		patterns[tmp]++
	}
	for i := 0; i < (1 << uint(m)); i++ {
		V += float64(patterns[i] * patterns[i])
	}
	V *= float64(uint(1) << uint(m))
	V /= float64(N)
	V -= float64(N)
	P = igamc(float64(((uint(1)<<uint(m))-1)>>1), V/2)
	//fmt.Printf("P = %v\n",P)
	if P >= alpha {
		return true
	}
	return false
}

//SerialTest 重叠子序列检测, m = 5
func SerialTest(data []byte, alpha float64) bool {
	if len(data) == 0 {
		fmt.Println("data cannot be empty")
		return false
	}
	bits := Bytes2Bits(data)
	n := len(bits)
	var m = 5
	patterns1 := make([]int, 1<<uint(m))
	patterns2 := make([]int, 1<<(uint(m)-1))
	patterns3 := make([]int, 1<<(uint(m)-2))
	Phi1 := 0.0
	Phi2 := 0.0
	Phi3 := 0.0
	var DPhi2, D2Phi2 float64
	var P1, P2 float64
	mask1 := (1 << uint(m)) - 1
	mask2 := (1 << (uint(m) - 1)) - 1
	mask3 := (1 << (uint(m) - 2)) - 1
	tmp := 0

	for i := 0; i < m-1; i++ {
		bits = append(bits, bits[i])
	}
	for i := 0; i < m-1; i++ {
		tmp <<= 1
		if bits[i] == 1 {
			tmp++
		}
	}

	for i := 0; i < n; i++ {
		// tmp should mod 2^m , or it will be overflow,
		// although the overflow don't influence the calculation result.
		tmp <<= 1
		tmp %= 1 << uint(m)
		if bits[m-1+i] == 1 {
			tmp++
		}
		// calculate the number of various subsequences
		patterns1[tmp&mask1]++
		patterns2[tmp&mask2]++
		patterns3[tmp&mask3]++
	}

	// get the square of patterns1[i]
	for i := 0; i <= mask1; i++ {
		Phi1 += math.Pow(float64(patterns1[i]), 2.0)
	}
	Phi1 *= float64(mask1 + 1)
	Phi1 /= float64(n)
	Phi1 -= float64(n)
	for i := 0; i <= mask2; i++ {
		Phi2 += math.Pow(float64(patterns2[i]), 2.0)
	}
	Phi2 *= float64(mask2 + 1)
	Phi2 /= float64(n)
	Phi2 -= float64(n)
	for i := 0; i <= mask3; i++ {
		Phi3 += math.Pow(float64(patterns3[i]), 2.0)
	}
	Phi3 *= float64(mask3 + 1)
	Phi3 /= float64(n)
	Phi3 -= float64(n)

	DPhi2 = Phi1 - Phi2
	D2Phi2 = Phi1 - 2*Phi2 + Phi3

	P1 = igamc(float64(uint(1)<<uint(m))/4.0, DPhi2/2.0)
	P2 = igamc(float64(uint(1)<<uint(m))/8.0, D2Phi2/2.0)
	//fmt.Printf("P1 = %v\n",P1)
	//fmt.Printf("P2 = %v\n",P2)
	if P1 >= alpha && P2 >= alpha {
		return true
	}
	return false
}

//RunsTest 游程总数检测
func RunsTest(data []byte, alpha float64) bool {
	if len(data) == 0 {
		fmt.Println("data cannot be empty")
		return false
	}
	bits := Bytes2Bits(data)
	n := len(bits)
	// Pi represent the proportion of 1 in the sequence
	var Pi float64
	// Vobs represent the sum of runs in the sequence to be checked
	Vobs := 1
	var P float64

	for i := 0; i < n-1; i++ {
		if bits[i] != bits[i+1] {
			Vobs++
		}
		if bits[i] == 1 {
			Pi++
		}
	}
	if bits[n-1] == 1 {
		Pi++
	}
	Pi /= float64(n)
	P = math.Erfc(math.Abs(float64(Vobs)-2.0*float64(n)*Pi*(1.0-Pi)) / (2.0 * math.Sqrt(2.0*float64(n)) * Pi * (1.0 - Pi)))
	//fmt.Printf("P = %v\n",P)
	if P >= alpha {
		return true
	}
	return false
}

//RunsDistributionTest 游程分布检测
func RunsDistributionTest(data []byte, alpha float64) bool {
	if len(data) == 0 {
		fmt.Println("data cannot be empty")
		return false
	}
	bits := Bytes2Bits(data)
	n := len(bits)
	// e[i] represent the expectations of runs with length i, which in a random binary sequence with length n
	e := make([]float64, 50)
	// b[i] record the number of 1-runs in a binary sequence with length i
	b := make([]float64, 50)
	// g[i] record the number of 0-runs in a binary sequence with length i
	g := make([]float64, 50)
	// k is the largest integer satisfying e[i] >= 5
	k := 0
	var V float64
	var cur = bits[0]
	var cnt = 0
	for {
		k++
		e[k] = float64(n-k+3) / float64(uint(1)<<uint(k+2))
		if e[k] <= 5.0 {
			break
		}
	}
	k--
	bits = append(bits, bits[n-1])
	for i := 0; i <= n; i++ {
		if bits[i] == cur {
			cnt++
		} else {
			if cnt <= k {
				if cur == 1 {
					b[cnt]++
				} else {
					g[cnt]++
				}
			}
			cur = bits[i]
			cnt = 1
		}
	}
	bits = bits[:len(bits)-1]
	for i := 1; i <= k; i++ {
		V += (b[i] - e[i]) * (b[i] - e[i]) / e[i]
		V += (g[i] - e[i]) * (g[i] - e[i]) / e[i]
	}
	P := igamc(float64(k)-1, V/2.0)
	//fmt.Printf("P = %v\n",P)
	if P >= alpha {
		return true
	}
	return false
}

//TestForTheLongestRunOfOnesInABlock 块内最大"1"游程检测方法
func TestForTheLongestRunOfOnesInABlock(data []byte, alpha float64) bool {
	var pi = []float64{0.0882, 0.2092, 0.2483, 0.1933, 0.1208, 0.0675, 0.0727}
	if len(data) == 0 {
		fmt.Println("Error: data cannot be empty!")
		return false
	}
	bits := Bytes2Bits(data)
	n := len(bits)
	m := 10000
	if n < m {
		fmt.Println("Error: the length of binary sequence must more than m!")
		return false
	}
	var N = n / m
	v := make([]float64, 7)
	var V float64
	var P float64

	// statistics the longest runs of one in a block
	for i := 0; i < N; i++ {
		lr1 := 0
		mlr1 := 0
		for j := 0; j < m; j++ {
			if bits[i*m+j] == 1 {
				lr1++
				mlr1 = int(math.Max(float64(mlr1), float64(lr1)))
			} else {
				lr1 = 0
			}
		}
		if mlr1 <= 10 {
			v[0]++
		}
		if mlr1 >= 16 {
			v[6]++
		}
		if 10 < mlr1 && mlr1 < 16 {
			v[mlr1-10]++
		}
	}

	for i := 0; i < 7; i++ {
		V += (v[i] - float64(N)*pi[i]) * (v[i] - float64(N)*pi[i]) / (float64(N) * pi[i])
	}
	P = igamc(3, V/2.0)
	//fmt.Printf("P = %v\n",P)
	if P >= alpha {
		return true
	}
	return false
}

//BinaryDerivativeTest 二元推导检测
func BinaryDerivativeTest(data []byte, alpha float64) bool {
	if len(data) == 0 {
		fmt.Println("data cannot be empty")
		return false
	}
	bits := Bytes2Bits(data)
	n := len(bits)
	// k can be 3 or 7
	k := 7
	S := 0
	var V float64
	var P float64
	var _bits = make([]byte, len(bits))
	_bits[n-1] = bits[n-1]
	// xor the two adjacent bits in the initial sequence ε
	for i := 0; i < k; i++ {
		for j := 0; j < n-i-1; j++ {
			_bits[j] = bits[j] ^ bits[j+1]
		}
	}
	// transfer the 0 1 to -1 1 of ε'(the new sequence), and accumulate the sum
	for i := 0; i < n-k; i++ {
		if _bits[i] == 1 {
			S++
		} else {
			S--
		}
	}
	V = math.Abs(float64(S)) / math.Sqrt(float64(n-k))
	P = math.Erfc(math.Abs(V) / math.Sqrt(2))
	//fmt.Printf("P = %v\n",P)
	if P >= alpha {
		return true
	}
	return false
}

//AutocorrelationTest 自相关检测
func AutocorrelationTest(data []byte, alpha float64) bool {
	if len(data) == 0 {
		fmt.Println("data cannot be empty")
		return false
	}
	bits := Bytes2Bits(data)
	n := len(bits)
	// d can be 1, 2, 8, 16
	d := 16
	Ad := 0
	var V float64
	var P float64
	// Ad 表示待检序列ε将其左移d位后所得新序列间不同元素的个数，称d为时延
	for i := 0; i < n-d; i++ {
		if bits[i]^bits[i+d] == 1 {
			Ad++
		}
	}

	V = 2.0 * (float64(Ad) - (float64(n-d) / 2.0)) / math.Sqrt(float64(n-d))
	P = math.Sqrt(math.Abs(V) / math.Sqrt(2))
	//fmt.Printf("P = %v\n",P)
	if P >= alpha {
		return true
	}
	return false
}

//BinaryMatrixRankTest 矩阵秩检测
func BinaryMatrixRankTest(data []byte, alpha float64) bool {
	if len(data) == 0 {
		fmt.Println("data cannot be empty")
		return false
	}
	bits := Bytes2Bits(data)
	n := len(bits)
	// n >= M*Q, 且 n - N*M*Q 要小
	// M、Q 是序列长度n=1000000比特时的参数推荐值，本规范取M=Q=32
	M := 32
	Q := 32
	N := n / (M * Q)
	// Fm 表示秩为 M 的矩阵的个数
	// Fm1 表示秩为 M-1 的矩阵的个数
	// Fr 表示秩小于 M-1 的矩阵的个数
	var Fm, Fm1, Fr = 0, 0, 0
	// 定义一个 M*Q 的矩阵
	matrix := make([][]int, M)
	for i := range matrix {
		matrix[i] = make([]int, Q)
	}
	var V float64
	var P float64
	var r int

	for i := 0; i < N; i++ {
		// 根据子序列来设置矩阵
		for j := 0; j < M; j++ {
			for k := 0; k < Q; k++ {
				if bits[i*M*Q+j*Q+k] == 1 {
					matrix[j][k] = 1
				} else {
					matrix[j][k] = 0
				}
			}
		}
		r = rank(matrix, M)
		if r == M {
			Fm++
		} else if r == M-1 {
			Fm1++
		} else {
			Fr++
		}
	}
	V += (float64(Fm) - 0.2888*float64(N)) * (float64(Fm) - 0.2888*float64(N)) / (0.2888 * float64(N))
	V += (float64(Fm1) - 0.5776*float64(N)) * (float64(Fm1) - 0.5776*float64(N)) / (0.5776 * float64(N))
	V += (float64(Fr) - 0.1336*float64(N)) * (float64(Fr) - 0.1336*float64(N)) / (0.1336 * float64(N))
	P = igamc(1, V/2.0)
	//fmt.Printf("P = %v\n",P)
	if P >= alpha {
		return true
	}
	return false
}

//CumulativeTest 累加和检测
func CumulativeTest(data []byte, alpha float64) bool {
	if len(data) == 0 {
		fmt.Println("data cannot be empty")
		return false
	}
	bits := Bytes2Bits(data)
	n := len(bits)
	// S代表待检序列ε的累加和，将ε中的 0 1 转换成 -1 1 进行计算
	S := 0
	// Z代表所有S中的绝对值的最大值
	Z := 0
	P := 1.0
	for i := 0; i < n; i++ {
		if bits[i] == 1 {
			S++
		} else {
			S--
		}
		Z = int(math.Max(float64(Z), math.Abs(float64(S))))
	}

	for i := ((-n / Z) + 1) / 4; i <= ((n/Z)-1)/4; i++ {
		P -= normalCDF((4*float64(i)+1)*float64(Z)/math.Sqrt(float64(n))) - normalCDF((4*float64(i)-1)*float64(Z)/math.Sqrt(float64(n)))
	}
	for i := ((-n / Z) - 3) / 4; i <= ((n/Z)-1)/4; i++ {
		P += normalCDF((4*float64(i)+3)*float64(Z)/math.Sqrt(float64(n))) - normalCDF((4*float64(i)+1)*float64(Z)/math.Sqrt(float64(n)))
	}
	//fmt.Printf("P = %v\n",P)
	if P >= alpha {
		return true
	}
	return false
}

//ApproximateEntropyTest 近似熵检测
func ApproximateEntropyTest(data []byte, alpha float64) bool {
	if len(data) == 0 {
		fmt.Println("data cannot be empty")
		return false
	}
	bits := Bytes2Bits(data)
	// m can be 2 or 5
	m := 5
	n := len(bits)
	var Cjm float64
	phim := 0.0
	phim1 := 0.0
	tmp := 0
	var V float64
	var P float64

	//Round1
	for i := 0; i < m-1; i++ {
		bits = append(bits, bits[i])
	}

	mask := (1 << uint(m)) - 1
	pattern := make([]int, 1<<uint(m))

	for i := 0; i < m-1; i++ {
		tmp <<= 1
		if bits[i] == 1 {
			tmp++
		}
	}
	for i := 0; i < n; i++ {
		tmp <<= 1
		if bits[m-1+i] == 1 {
			tmp++
		}
		pattern[tmp&mask]++
	}
	for i := 0; i < (1 << uint(m)); i++ {
		Cjm = float64(pattern[i]) / float64(n)
		if Cjm != 0 {
			phim += Cjm * math.Log(Cjm)
		}
	}
	bits = bits[:len(bits)-m+1]
	//Round2
	m++
	for i := 0; i < m-1; i++ {
		bits = append(bits, bits[i])
	}

	mask = (1 << uint(m)) - 1
	pattern1 := make([]int, 1<<uint(m))

	for i := 0; i < m-1; i++ {
		tmp <<= 1
		if bits[i] == 1 {
			tmp++
		}
	}
	for i := 0; i < n; i++ {
		tmp <<= 1
		//tmp %= 1 << uint(m)
		if bits[m-1+i] == 1 {
			tmp++
		}
		pattern1[tmp&mask]++
	}
	for i := 0; i < (1 << uint(m)); i++ {
		Cjm = float64(pattern1[i]) / float64(n)
		//fmt.Println(Cjm)
		if Cjm != 0 {
			phim1 += Cjm * math.Log(Cjm)
		}
	}
	bits = bits[:len(bits)-(m-1)]

	//Final
	m--
	ApEn := phim - phim1
	V = 2.0 * float64(n) * (math.Log(2) - ApEn)
	P = igamc(float64(uint(1)<<uint(m))/2.0, V/2.0)
	//fmt.Printf("P = %v\n",P)

	if P >= alpha {
		return true
	}
	return false
}

//LinearComplexityTest 线性复杂度检测
func LinearComplexityTest(data []byte, alpha float64) bool {
	if len(data) == 0 {
		fmt.Println("data cannot be empty")
		return false
	}
	bits := Bytes2Bits(data)
	// 本规范取 m = 500
	m := 500
	n := len(bits)
	N := n / m

	var v = []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	// pi 对应标准中规定的 𝛑 值
	var pi = []float64{0.010417, 0.03125, 0.12500, 0.5000, 0.25000, 0.06250, 0.020833}
	var V = 0.0
	var P float64

	arr := make([]byte, m)
	var complexity int // 线性复杂度
	var T, miu float64

	num := math.Pow(-1.0, float64(m+1))
	// 计算 μ 值
	miu = float64(m)/2.0 + (9.0+num)/36.0 - (float64(m)/3.0+2.0/9.0)/math.Pow(2.0, float64(m))
	for i := 0; i < N; i++ {
		for j := 0; j < m; j++ {
			arr[j] = bits[i*m+j]
		}
		complexity = linearComplexityTest(arr, m)
		num = math.Pow(-1.0, float64(m))
		T = num*(float64(complexity)-miu) + 2.0/9.0
		if T <= -2.5 {
			v[0]++
		} else if T <= -1.5 {
			v[1]++
		} else if T <= -0.5 {
			v[2]++
		} else if T <= 0.5 {
			v[3]++
		} else if T <= 1.5 {
			v[4]++
		} else if T <= 2.5 {
			v[5]++
		} else {
			v[6]++
		}
	}

	for i := 0; i < 7; i++ {
		V += math.Pow(v[i]-float64(N)*pi[i], 2.0) / (float64(N) * pi[i])
	}

	P = igamc(3.0, V/2.0)
	//fmt.Printf("P = %v\n",P)

	if P >= alpha {
		return true
	}
	return false
}

//MaurersUniversalTest Maurer通用统计检测
func MaurersUniversalTest(data []byte, alpha float64) bool {
	if len(data) == 0 {
		fmt.Println("data cannot be empty")
		return false
	}
	bits := Bytes2Bits(data)
	n := len(bits)
	// 将待检序列ε分为两部分：初始序列和测试序列
	// 初始序列包括Q个L位的非重叠子序列，测试序列包括K个L位的非重叠子序列
	// 将多余的位（不够组成一个完整的L位子序列）舍弃
	// 本规范取 L = 7, Q = 1280
	L := 7
	Q := 1280
	K := n/L - Q

	// 针对初始序列，创建一个表，它以L位值作为表中的索引值
	// T[j] 表示表中第j个元素的值，实际代表每个L位子序列出现的频数
	T := make([]int, 1<<uint(L))
	mask := (1 << uint(L)) - 1
	sum := 0.0
	V := 0.0
	var P, sigma float64
	// E 表示期望值, variance 表示方差,这里不推荐 L 的值在 6 以下，故设为 0 值，实际可计算出
	var E = []float64{0, 0, 0, 0, 0, 0, 5.2177052, 6.1962507, 7.1836656,
		8.1764248, 9.1723243, 10.170032, 11.168765,
		12.168070, 13.167693, 14.167488, 15.167379}
	var variance = []float64{0, 0, 0, 0, 0, 0, 2.954, 3.125, 3.238, 3.311, 3.356, 3.384,
		3.401, 3.410, 3.416, 3.419, 3.421}
	tmp := 0

	for i := 1; i <= Q; i++ {
		for j := 0; j < L; j++ {
			tmp <<= 1
			if bits[(i-1)*L+j] == 1 {
				tmp++
			}
		}
		T[tmp&mask] = i
	}

	index := Q * L
	for i := Q + 1; i <= Q+K; i++ {
		for j := 0; j < L; j++ {
			tmp <<= 1
			if bits[(i-1-Q)*L+j+index] == 1 {
				tmp++
			}
		}
		sum += math.Log(float64(i-T[tmp&mask])) / math.Log(2.0)
		T[tmp&mask] = i
	}

	// 计算方差 σ
	sigma = math.Sqrt(variance[L]/float64(K)) * mutFactorC(L, K)
	V = (sum/float64(K) - E[L]) / sigma
	P = math.Erfc(math.Abs(V) / math.Sqrt(2.0))
	//fmt.Printf("P = %v\n",P)
	if P >= alpha {
		return true
	}
	return false
}

//DiscreteFourierTransformTest 离散傅里叶检测
func DiscreteFourierTransformTest(data []byte, alpha float64) bool {
	if len(data) == 0 {
		fmt.Println("data cannot be empty")
		return false
	}
	bits := Bytes2Bits(data)
	n := len(bits)
	r := make([]float64, n)
	// T 表示门限值
	T := math.Sqrt(2.995732274 * float64(n))
	N0 := 0.95 * float64(n) / 2
	// N1 表示系数f[i]中小于门限值T的复数个数
	N1 := 0

	// 将待检序列ε中的 0 1 分别转换成 -1 1 , 得到新序列
	for i := 0; i < n; i++ {
		if bits[i] == 1 {
			r[i] = 1.0
		} else {
			r[i] = -1.0
		}
	}
	// 将新序列长度扩展到2的指数倍，长度不够的补0，得到的序列前面是 -1 1, 后面是0   ？ why
	r = pow2DoubleArr(r)
	// 对新序列进行快速傅立叶变换，得到一系列的复数f[i]
	// FFT是离散傅立叶变换(DFT)的快速算法
	f := fft.FFTReal(r)
	//fmt.Println(f)

	// 对每一个f[i]，计算其系数mod[i], i 取值为 n/2 - 1 以降低计算量
	var mod = make([]float64, len(f))
	//for i := 0; i < len(f); i++{
	for i := 0; i < n/2-1; i++ {
		// a 表示复数的实部, b 表示复数的虚部
		a := real(f[i])
		b := imag(f[i])
		// mod[i]表示求复数的模
		mod[i] = math.Abs(math.Sqrt(a*a + b*b))
	}
	for i := 0; i < n/2-1; i++ {
		if mod[i] < T {
			N1++
		}
	}
	//if math.Abs(r[0]) < T{
	//	N1++
	//}
	// GMT 0005-2012 中是 0.95*0.05*n/4
	V := (float64(N1) - N0) / math.Sqrt(0.95*0.05*float64(n)/4)
	P := math.Erfc(math.Abs(V) / math.Sqrt(2))
	//fmt.Printf("P = %v\n",P)
	if P >= alpha {
		return true
	}
	return false
}
