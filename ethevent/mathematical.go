package ethevent

import (
	"fmt"
	"math/big"
)

// BigIntAdd 加
func BigIntAdd(numstr string, num int64) string {
	n, _ := new(big.Int).SetString(numstr, 10)
	m := new(big.Int)
	m.SetInt64(num)
	m.Add(n, m)
	return m.String()
}

// BigIntReduce 减
func BigIntReduce(numstr string, num int64) string {
	n, _ := new(big.Int).SetString(numstr, 10)
	m := new(big.Int)
	m.SetInt64(-num)
	m.Add(n, m)
	return m.String()
}

// BigIntMul 乘
func BigIntMul(numstr string, num int64) string {
	n, _ := new(big.Int).SetString(numstr, 10)
	m := new(big.Int)
	m.SetInt64(num)
	m.Mul(n, m)
	return m.String()
}

// BigIntDiv 除
func BigIntDiv(numstr string, num int64) string {
	n, _ := new(big.Int).SetString(numstr, 10)
	m := new(big.Int)
	m.SetInt64(num)
	m.Div(n, m)
	return m.String()
}

// TestMathematical 测试运算
func TestMathematical() {
	numstr := "1515631351536151161464461151511561"
	// 加
	fmt.Println(BigIntAdd(numstr, 99))
	// 减
	fmt.Println(BigIntReduce(numstr, 99))
	// 乘
	fmt.Println(BigIntMul(numstr, 99))
	// 除
	fmt.Println(BigIntDiv(numstr, 99))
}
