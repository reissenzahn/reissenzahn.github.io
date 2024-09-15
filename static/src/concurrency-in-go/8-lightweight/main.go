package main

import (
	"fmt"
	"runtime"
	"sync"
)

// goroutines are lightweight; they are initially only allocated a few kilobytes and the runtime grows and shrinks the memory for storing the stack automatically

// the garbage collector does nothing to cleanup goroutines that have been abandoned somehow

const N = 1e4

func memoryUsage() uint64 {
	runtime.GC()
	var s runtime.MemStats
	runtime.ReadMemStats(&s)
	return s.Sys
}

func main() {
	var c <-chan interface{}

	var wg sync.WaitGroup
	wg.Add(N)

	before := memoryUsage()
	for i := 0; i < N; i++ {
		go func() {
			wg.Done()
			<-c
		}()
	}
	wg.Wait()
	after := memoryUsage()

	fmt.Printf("%.3fkb", float64(after-before)/N/1000)
}
