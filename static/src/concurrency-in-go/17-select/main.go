package main

import (
	"fmt"
	"time"
)

// a select block encompasses a series of case statements that guard a series of statements; cases consists of channel reads and writes and are considered simultaneously

// if any of the cases are ready (populated or closed channels, channels not at capacity) the corresponding statements will be executed; if none of the cases are ready then the select statement blocks

func main() {
	ch := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(ch)
	}()

	select {
	case <-ch:
		fmt.Println("Hello!")
	}

	// the runtime will perform pseudorandom uniform over the set of case statements that are ready
	ch1 := make(chan interface{})
	close(ch1)
	ch2 := make(chan interface{})
	close(ch2)

	var count1, count2 int
	for i := 0; i < 1000; i++ {
		select {
		case <-ch1:
			count1++
		case <-ch2:
			count2++
		}
	}

	fmt.Println(count1, count2)

	// if none of the cases are ready then the select statement blocks
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out")
	}

	// the default clause will be executed if all the channels you're selecting against are blocking
	var ch3, ch4 <-chan int
	select {
	case <-ch3:
	case <-ch4:
	default:
		fmt.Println("Nothing to do")
	}

	// a default clause is usually used in conjunction with for-select loop to allow the goroutine to make progress on work while waiting for another goroutine to report a result
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	counter := 0
	loop:
		for {
			select {
			case <-done:
				break loop
			default:
			}

			counter++
			time.Sleep(1 * time.Second)
		}

	// a select statement with no case clauses will block forever
	// select {}
}
