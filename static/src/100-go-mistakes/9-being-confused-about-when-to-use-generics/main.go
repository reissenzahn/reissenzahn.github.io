package main

import "io"

func getKeys[K comparable, V any](m map[K]V) []K {
	var keys []K
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

type Node[T any] struct {
	Val  T
	next *Node[T]
}

func (n *Node[T]) Add(next *Node[T]) {
	n.next = next
}

// type parameters cannot be used with method arguments, only with function arguments or method receivers

// generics are recommended for:
// - data structures
// - functions manipulating slices, maps and channels of any type
// - factoring out behaviors instead of types

// generics are not recommended when calling a method of the type argument or when they make code more complex

// we should just make the w argument an io.Writer directly
func foo[T io.Writer](w T) {
	w.Write([]byte("hello"))
}

// generics introduce a form of abstraction and unnecessary abstractions introduce complexity

func main() {
}

// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#no-parameterized-methods
// https://stackoverflow.com/questions/64189810/how-to-use-a-type-parameter-in-an-interface-method
