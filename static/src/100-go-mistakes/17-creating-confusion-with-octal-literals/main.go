package main

import "os"

// numbers can include underscores for clarity
const billion = 1_000_000_000

// an integer literal starting with 0 is considered an octal integer
// 100 + 010 == 108

func main() {
	// using 0o as a prefix means the same thing but is more readable
	_, _ = os.OpenFile("foo", os.O_RDONLY, 0o644)
}
