package main

import (
	"fmt"
	"sync"
)

func main() {
	// declare a channel of type interface{}
	// var ch chan interface{}

	// instantiate a channel
	// ch = make(chan interface{})

	// declare a unidirectional channel (receive only)
	// var ch <-chan interface{}
	// ch := make(<-chan interface{})

	// declare a unidirectional channel (send only)
	// var ch chan<- interface{}
	// ch := make(chan<- interface{})

	// unidirectional channels are often used as function parameters and return types; go will implicitly convert bidirectional channels to unidirectional channels when needed
	// var ch <-chan interface{}
	// ch = make(chan interface{})

	// channels are typed
	ch := make(chan int)

	// send and receive
	go func() {
		ch <- 42
	}()

	fmt.Println(<-ch)

	// channels are blocking; any goroutine that attempts to write to a channel that is full will wait until the channel has been emptied, and any goroutine that attempts to read from a channel that is empty will wait until at least one item is placed on it

	// closing a channel indicates that no more values will be sent
	// close(ch)

	// the receiving form of the <- operator optionally returns two values; the second return value is a way to indicate whether the read off the channel was a value generated by a write elsewhere in the process or a zero value generated from a closed channel
	// v, ok := <-ch

	ch2 := make(chan int)
	close(ch2)
	v, ok := <-ch2
	fmt.Printf("%v, %v", v, ok)

	// range can be used to iterate over the values received on a channel and break when the channel is closed
	for v := range ch2 {
		fmt.Println(v)
	}

	// closing a channel can be used to unblock multiple goroutines
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-begin
		}()
	}
	close(begin)
	wg.Wait()

	// closing or writing to a closed channel will panic

	// buffered channels are given a capacity when they are instantiated; even if no reads are performed on the channel, a goroutine can still perform n writes where n is the capacity
	// var ch chan interface{}
	// ch = make(chan interface{}, 4)

	// an un-buffered channel is simply a buffered channel with a capacity of zero

	// a buffered channel with no receives and a capacity of four would be full after four writes and block on the fifth write
	ch3 := make(chan rune, 4)
	ch3 <- 'A'
	ch3 <- 'B'
	ch3 <- 'C'
	ch3 <- 'D'
	// ch3 <- 'E' // blocks

	// the zero value for a channel is nil; reading from or writing to a nil channel will block, closing a nil channel will panic

	// we can assign channel ownership; the goroutine that owns a channel should instantiate it, perform writes or pass ownership to another goroutine and close the channel
	owner := func() <-chan int {
		ch := make(chan int, 5)
		go func() {
			defer close(ch)
			for i := 0; i < 6; i++ {
				ch <- i
			}
		}()
		return ch
	}

	// consumer only has to handle blocking reads and channel closes
	for v := range owner() {
		fmt.Println(v)
	}

	// encapsulate the lifecycle f the channel within the owner; keep the scope of channel ownership small
}