package main

import "fmt"

// a pipeline is a series of things that take data in, perform an operation on it, and pass the data back out; we call each of these operations a stage of the pipeline

// this allows for separating the concerns of each stage; you can modify stages independently of one another, you can mix and match how stages are combined, you can process each stage concurrently to upstream or downstream stages, you can fan-out or rate-limit portions of your pipeline

// a stage is something that takes data in, performs a transformation on it and sends the data back out

// a pipeline stage consumes and returns the same type

// batch processing refers to operating on chunks of data all at once whereas stream processing refers to receiving and emitting one element at a time

func generator(done <-chan interface{}, integers ...int) <-chan int {
	ch := make(chan int, len(integers))

	go func() {
		defer close(ch)
		for _, v := range integers {
			select {
			case <-done:
				return
			case ch <- v:
			}
		}
	}()
	return ch
}

func multiply(done <-chan interface{}, in <-chan int, multiplier int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-done:
				return
			case out <- v * multiplier:
			}
		}
	}()
	return out
}

func add(done <-chan interface{}, in <-chan int, additive int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-done:
				return
			case out <- v * additive:
			}
		}
	}()
	return out
}

func main() {
	done := make(chan interface{})
	defer close(done)

	pipeline := multiply(done, add(done, multiply(done, generator(done, 1, 2, 3, 4), 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}
}
