package main

import (
	"fmt"
	genrand "github.com/ultramesh/crypto-gm/csp/genrand"
	random "github.com/ultramesh/crypto-gm/csp/randomTest"
	"io"
	"os"
	"sync"
)

// 1000组128*1024Byte随机数：
//crypto/rand
//983 987 963 985 989 987 991 989 999 989 984 990 988 989 983
//sjk1862
//990 992 973 979 985 974 986 985 996 987 989 986 995 990 989
//990 988 971 981 992 992 990 992 1000 984 989 986 989 991 994
//rdrand
//987 994 984 983 992 987 991 993 999 987 985 991 987 988 990
func main() {
	result := make([]int, len(random.TestFunc))
	for i := 0; i < 1000; i++ {
		data, err := generateRandomnessToFile(genrand.Reader, fmt.Sprintf("random%03v.bin", i), 128*1024)
		//data := make([]byte, 128*1024)
		//_, err := genrand.Reader.Read(data)
		if err != nil {
			panic(err)
		}
		var wg sync.WaitGroup
		wg.Add(len(random.TestFunc))
		for j := range random.TestFunc {
			go func(k int) {
				if random.TestFunc[k].Method(data, 0.01) {
					result[k]++
				}
				wg.Done()
			}(j)
		}
		wg.Wait()
		fmt.Println(i)
	}
	fmt.Println(result)
}

// 生成随机数到文件
func generateRandomnessToFile(reader io.Reader, filename string, length int) ([]byte, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	//defer file.Close()

	data := make([]byte, length)
	n, err := reader.Read(data)
	if n != length || err != nil {
		return nil, fmt.Errorf("generate random error")
	}
	n, err = file.Write(data)
	if n != length || err != nil {
		return nil, fmt.Errorf("generate random error")
	}
	return data, nil
}
