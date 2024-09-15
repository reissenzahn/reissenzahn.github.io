package main

import (
	"fmt"
	"sync"
)

// Signal() notifies the goroutine that has been waiting the longest whereas Broadcast sends a signal to all waiting goroutines

type Button struct {
	Clicked *sync.Cond
}

func subscribe(c *sync.Cond, f func()) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		wg.Done()
		c.L.Lock()
		defer c.L.Unlock()
		c.Wait()
		f()
	}()

	// prevent subscribe from existing until the goroutine is confirmed to be running
	wg.Wait()
}

func main() {
	button := Button{
		Clicked: sync.NewCond(&sync.Mutex{}),
	}

	var wg sync.WaitGroup
	wg.Add(3)

	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window")
		wg.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("Displaying dialog box")
		wg.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("Submit form")
		wg.Done()
	})

	button.Clicked.Broadcast()

	wg.Wait()
}
