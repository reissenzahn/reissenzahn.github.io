package main

import "fmt"

// the copy() built-in allows copying elements from a source slice into a destination slice

// the destination argument is the first argument while the source is the second argument

// the number of elements copied to the destination slice corresponds to the minimum between the source slice's length and the destination slice's length

func main() {
	src := []int{0, 1, 2}
	var dst []int
	copy(dst, src)
	fmt.Println(dst) // []

	// correct
	dst2 := make([]int, len(src))
	copy(dst2, src)
	fmt.Println(dst2)

	// append() can also be used to copy elements, though copy() is usually more idiomatic
	dst3 := append([]int(nil), src...)
	fmt.Println(dst3)
}
