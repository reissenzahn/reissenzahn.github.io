package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	worker := func(ch <-chan string) <-chan interface{} {
		done := make(chan interface{})
		go func() {
			defer close(done)

			for v := range ch {
				fmt.Println(v)
			}
		}()

		return done
	}

	// goroutine will never exit and will leak
	go worker(nil)

	// establish a signal between the parent goroutine and its children that allows the parent to signal cancellation to its children; by convention, this signal is usually a read-only channel named done
	worker2 := func(done <-chan interface{}, ch <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer close(completed)
			for {
				select {
				case v := <-ch:
					fmt.Println(v)
				case <-done:
					return
				}
			}
		}()
		return completed
	}

	done := make(chan interface{})
	completed := worker2(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		close(done)
	}()
	<-completed

	// a similar pattern can be used if a goroutine becomes blocked on a write
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randStream
	}

	done2 := make(chan interface{})
	randStream := newRandStream(done2)
	for i := 0; i < 3; i++ {
		fmt.Println(<-randStream)
	}
	close(done)
}

// by convention, if a goroutine is responsible for creating a goroutine, it is also responsible for ensuring it can stop the goroutine
