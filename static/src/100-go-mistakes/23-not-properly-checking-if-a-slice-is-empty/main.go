package main

// checking if a slice is nil is not equivalent to checking if it is empty; we also need to consider the empty slice

// instead we should check the length using len() which returns 0 for both a nil and empty slice

// when designing interfaces, we should avoid distinguishing nil and empty slices; when returning slices, it should make neither a semantic nor a technical difference if we return a nil or empty slice

func main() {

}