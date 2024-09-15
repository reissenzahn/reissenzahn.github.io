package main

import (
	"fmt"
	"math"
)

// an integer overflow that can be detected at compile time generates a compilation error
// var j int32 = math.MaxInt32 + 1

func main() {
	var i int32 = math.MaxInt32

	// at run time, an integer overflow is silent and does not cause a panic
	i++
	fmt.Println(i)
}

// detect an integer overflow
func Inc32(counter int32) int32 {
	if counter == math.MaxInt32 {
		panic("int32 overflow")
	}
	return counter + 1
}

func AddInt(a, b int) int {
	if a > math.MaxInt-b {
		panic("int overflow")
	}
	return a + b
}

func MultiplyInt(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	result := a * b
	if a == 1 || b == 1 {
		return result
	}
	if a == math.MinInt || b == math.MinInt {
		panic("integer overflow")
	}
	if result/b != a {
		panic("integer overflow")
	}
	return result
}

// the math/big package is provided for handling large numbers
