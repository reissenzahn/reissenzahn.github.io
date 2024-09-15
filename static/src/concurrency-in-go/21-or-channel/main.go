package main

import (
	"fmt"
	"time"
)

// the or-channel pattern creates a composite done channel through recursion and goroutines

func or(chs ...<-chan interface{}) <-chan interface{} {
	switch len(chs) {
	case 0:
		return nil
	case 1:
		return chs[0]
	}

	done := make(chan interface{})
	go func() {
		defer close(done)

		switch len(chs) {
		case 2:
			select {
			case <-chs[0]:
			case <-chs[1]:
			}
		default:
			select {
			case <-chs[0]:
			case <-chs[1]:
			case <-chs[2]:
			case <-or(append(chs[3:], done)...):
			}
		}
	}()

	return done
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		ch := make(chan interface{})
		go func() {
			defer close(ch)
			time.Sleep(after)
		}()

		return ch
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(2*time.Second),
	)

	fmt.Printf("%v\n", time.Since(start))
}
