package main

import "fmt"

// floating-point arithmetic is an approximation of real arithmetic

// using the == operator to compare two floating-point numbers can lead to inaccuracies; instead we should compare their difference to see if it less than some small error value

// the result of floating-point calculations depends on the actual processor

func main() {
	var n float32 = 1.001
	fmt.Println(n * n)

	// there are three special kinds of floating-point numbers: positive infinite, negative infinite and NaN (the result of an undefined or un-representable operation)
	var f float64
	pInf := 1 / f
	nInf := -1 / f
	nan := f / f
	fmt.Println(pInf, nInf, nan)

	// the error can accumulate in a sequence of floating-point operations
	for n := range []int{10, 1_000, 1_000_000} {
		r1 := 10_000.
		for i := 0; i < n; i++ {
			r1 += 1.001
		}

		r2 := 0.
		for i := 0; i < n; i++ {
			r2 += 1.001
		}

		fmt.Printf("%d\t%.12f\t%.12f\n", n, r1, r2+10_000)
	}

	// ^this doesn't seem to work on my machine
}
