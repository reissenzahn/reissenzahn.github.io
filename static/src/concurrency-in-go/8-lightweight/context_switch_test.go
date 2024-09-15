package main

import (
	"sync"
	"testing"
)

// context switching is also much cheaper in a software-defined scheduler
func BenchmarkContextSwitch(b *testing.B) {
	begin := make(chan struct{})
	c := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(2)

	var token struct{}
	go func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			c <- token
		}
	}()

	go func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			<-c
		}
	}()

	b.StartTimer()
	close(begin)
	wg.Wait()
}

// go test -bench=. -cpu=1
