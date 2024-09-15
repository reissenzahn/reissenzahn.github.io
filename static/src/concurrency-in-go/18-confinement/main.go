package main

import (
	"bytes"
	"fmt"
	"sync"
)

// confinement refers to ensuring information is only ever available from one concurrent process; when this is achieved, a concurrent program is implicitly safe and no synchronization is needed

// there are two types of confinement: ad hoc (achieved through convention) and lexical (using lexical scope to expose only the correct data and concurrency primitives for multiple concurrent processes to use)

// avoiding the need for synchronization avoids creating critical sections and sidesteps complexity

func main() {
	// ad hoc confinement
	data := make([]int, 4)

	// by convention, we only access data from the loopData function
	loopData := func(ch chan<- int) {
		defer close(ch)
		for _, v := range data {
			ch <- v
		}
	}

	ch := make(chan int)
	go loopData(ch)

	for v := range ch {
		fmt.Println(v)
	}

	// lexical confinement
	owner := func() <-chan int {
		// instantiate the channel within the lexical scope of the owner function; this limits the scope of the write aspect of the results channel to the closure defined below it; it confines the write aspect of this channel to prevent other goroutines from writing to it
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i < 6; i++ {
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Println(result)
		}
	}

	// receive the read aspect f the channel and pass it to the consumer; this confines the main goroutine to a read-only view of the channel
	results := owner()
	consumer(results)

	// use confinement for a data structure that is not concurrent-safe, bytes.Buffer
	print := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buf bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buf, "%c", b)
		}
		fmt.Println(buf.String())
	}

	// print doesn't close around the data slice and so cannot access it; it needs to take in a slice of byte to operate on
	var wg sync.WaitGroup
	wg.Add(2)
	d := []byte("golang")

	// we pass in different subsets of the slice, thus constraining the goroutines to the part of the slice being passed in--there is no need to synchronize memory access
	go print(&wg, d[:3])
	go print(&wg, d[3:])
	wg.Wait()
}
