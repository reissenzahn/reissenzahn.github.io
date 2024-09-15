package main

import "fmt"

// a slice is backed by an array and handles the logic of adding an element if the backing array is full or shrinking the backing array if it is almost empty

// internally, a slice holds a pointer to the backing array as well as a length and capacity

// the length is the number of elements the slice contains while the capacity is the number of elements in the backing array

// accessing an element outside the length range is forbidden

// slicing creates another slice that references the same backing array

func main() {
	s := make([]int, 3, 6)
	// [0 0 0 _ _ _]
	//  ^s (3)

	s[1] = 1
	// [0 1 0 _ _ _]
	//  ^s (3)

	s = append(s, 2)
	// [0 1 0 2 _ _]
	//  ^s (3)

	s1 := make([]int, 3, 6)
	s2 := s1[1:3]
	// [0 0 0 _ _ _]
	//  ^s1 (3, 6)
	//    ^s2 (2, 5)

	s1[1] = 1
	// [0 1 0 _ _ _]
	//  ^s1 (3, 6)
	//    ^s2 (2, 5)

	s2 = append(s2, 2)
	// [0 1 0 2 _ _]
	//  ^s1 (3, 6)
	//    ^s2 (3, 5)

	// if append() has sufficient capacity, the destination is resliced to accommodate the new elements. If it does not, a new underlying array will be allocated.
	s2 = append(s2, 3)
	s2 = append(s2, 4)
	s2 = append(s2, 5)
	// [0 1 0 2 3 4]
	//  ^s1 (3, 6)
	// [1 0 2 3 4 5 _ _ _ _]
	//  ^s2 (6, 10)

	s1[1] = 0

	fmt.Println(s1) // [0 0 0]
	fmt.Println(s2) // [1 0 2 3 4 5]
}
