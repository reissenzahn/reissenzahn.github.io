package main

import "fmt"

// a slice is empty if its length is equal to 0; a slice is nil if it equals nil

func log(s []string) {
	fmt.Printf("empty=%t\tnil=%t\n", len(s) == 0, s == nil)
}

func main() {
	// a nil slice also has a length of 0
	var s []string
	log(s)

	// this can occasionally be useful to pass a nil slice in a single line
	s = []string(nil)
	log(s)

	s = []string{}
	log(s)

	s = make([]string, 0)
	log(s)

	// empty=true      nil=true
	// empty=true      nil=true
	// empty=true      nil=false
	// empty=true      nil=false

	// append() also works for nil slices
	var s1 []string
	fmt.Println(append(s1, "foo")) // [foo]

	// a nil slice does not require any allocation so we should favor returning a nil slice

	// []string{} should be avoided as it results in an allocation and provides no further benefits
}
