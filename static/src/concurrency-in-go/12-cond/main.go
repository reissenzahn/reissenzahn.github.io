package main

import (
	"fmt"
	"sync"
	"time"
)

// a Cond provides a rendezvous point for goroutines waiting for or announcing the occurrence of an event

// this allows a goroutine to efficiently sleep until it was signaled to wake and check its conditions

// the call to Wait() doesn't just block but also suspends the current goroutine

// upon entering Wait(), Unlock() is called on the Locker and upon existing Wait(), Lock() is called on the Locker

// Signal() notifies the goroutine that has been waiting the longest whereas Broadcast sends a signal to all waiting goroutines

func main() {
	// instantiate Cond providing a sync.Locker
	// c := sync.NewCond(&sync.Mutex{})

	// lock the Locker for the condition
	// c.L.Lock()

	// for conditionTrue() == false {
	// wait to be notified that the condition had occurred; Unlock() is called and goroutine will be suspended until condition is occurred at which point Lock() is called again
	// 	c.Wait()
	// }

	// Unlock() the sync.Locker
	// c.L.Unlock()

	// queue has a capacity of 2
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0)

	pop := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Popped")
		c.L.Unlock()
		c.Signal()
	}

	// only enqueue item when there is room
	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
		}
		fmt.Println("Pushed")
		queue = append(queue, struct{}{})
		go pop(1 * time.Second)
		c.L.Unlock()
	}
}
