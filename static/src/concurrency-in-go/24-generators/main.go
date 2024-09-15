package main

import (
	"fmt"
	"math/rand"
)

// a generator for a pipeline is any function that converts a set of discrete values into a stream of values on a channel

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case ch <- v:
				}
			}
		}
	}()
	return ch
}

func take(done <-chan interface{}, in <-chan interface{}, n int) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case out <- <-in:
			}
		}
	}()
	return out
}

func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case ch <- fn():
			}
		}
	}()
	return ch
}

// separate stage for type assertions
func toString(done <-chan interface{}, in <-chan interface{}) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-done:
				return
			case out <- v.(string):
			}
		}
	}()
	return out
}

func main() {
	done := make(chan interface{})
	defer close(done)

	for v := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", v)
	}

	rand := func() interface{} { return rand.Int() }
	for v := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(v)
	}

	var msg string
	for v := range toString(done, take(done, repeat(done, "hello", "world"), 5)) {
		msg += v
	}
	fmt.Println(msg)
}
