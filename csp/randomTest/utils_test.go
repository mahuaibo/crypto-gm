package randomtest

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBytes2Bits(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{
			"test1",
			args{
				[]byte{11, 22, 33},
			},
			[]byte{0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Bytes2Bits(tt.args.data)
			fmt.Println(len(got))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes2Bits() = %v, want %v", got, tt.want)

			}
		})
	}
}

func Test_rank(t *testing.T) {
	a := [][]int{
		{1, 1, 1},
		{1, 0, 1},
		{1, 1, 0},
	}
	b := [][]int{
		{1, 1, 0},
		{1, 1, 0},
	}
	c := make([][]int, 100)
	for i := range c {
		c[i] = make([]int, 100)
	}
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			if i == j {
				c[i][j] = 1
			} else {
				c[i][j] = 0
			}
		}
	}
	type args struct {
		matrix [][]int
		m      int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			"matrix1",
			args{
				a,
				3,
			},
			3,
		},
		{
			"matrix2",
			args{
				b,
				2,
			},
			1,
		},
		{
			"matrix3",
			args{
				c,
				100,
			},
			100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rank(tt.args.matrix, tt.args.m); got != tt.want {
				t.Errorf("rank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rowEchelon(t *testing.T) {
	a := [][]int{
		{1, 1, 1},
		{1, 0, 1},
		{1, 1, 0},
	}
	b := [][]int{
		{1, 1, 0},
		{1, 1, 0},
	}

	a1 := [][]int{
		{1, 1, 1},
		{0, 1, 0},
		{0, 0, 1},
	}
	b1 := [][]int{
		{1, 1, 0},
		{0, 0, 0},
	}
	type args struct {
		matrix [][]int
		m      int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		// TODO: Add test cases.
		{
			"matrix1",
			args{
				a,
				3,
			},
			a1,
		},
		{
			"matrix2",
			args{
				b,
				2,
			},
			b1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rowEchelon(tt.args.matrix, tt.args.m); !(reflect.DeepEqual(got, a1) || reflect.DeepEqual(got, b1)) {
				t.Errorf("rowEchelon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalCDF(t *testing.T) {
	type args struct {
		x float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			"test1",
			args{2.0},
			0.9772498680518208,
		},
		{
			"test2",
			args{10.0},
			1.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalCDF(tt.args.x); got != tt.want {
				t.Errorf("normalCDF() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试数据来自 matlab
// https://www.mathworks.com/help/matlab/ref/gammainc.html
func Test_igamc(t *testing.T) {
	type args struct {
		a float64
		x float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			"test1",
			args{
				1.5,
				2.35,
			},
			0.19512957232911837,
		},
		{
			"test2",
			args{
				2.0,
				1.55,
			},
			0.5412323332581948,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := igamc(tt.args.a, tt.args.x); got != tt.want {
				t.Errorf("igamc() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 输入的数据必须是二元序列
func Test_linearComplexityTest(t *testing.T) {
	a := []byte{0, 0, 1, 0, 0, 0, 1, 1, 1, 0, 1}
	type args struct {
		a []byte
		M int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			"test1",
			args{
				a,
				11,
			},
			6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range a {
				if a[i] != 0 && a[i] != 1 {
					t.Errorf("the  slice must be a binary sequence !")
					return
				}
			}
			if got := linearComplexityTest(tt.args.a, tt.args.M); got != tt.want {
				t.Errorf("linearComplexityTest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pow2DoubleArr(t *testing.T) {
	type args struct {
		data []float64
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pow2DoubleArr(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pow2DoubleArr() = %v, want %v", got, tt.want)
			}
		})
	}
}
