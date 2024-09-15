package main

import (
	"fmt"
	"time"
)

// heartbeats provide a way for concurrent processes to signal life to outside parties

// heartbeats can either occur on a time interval or at the beginning of a unit of work

func doWork(done <-chan interface{}, interval time.Duration) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)

	go func() {
		defer close(heartbeat)
		defer close(results)

		pulse := time.Tick(interval)
		workGen := time.Tick(2*interval)

		sendPulse := func() {
				select {
				case heartbeat <-struct{}{}:
				// there may be no one listening to the heartbeat
				default:
				}
		}

		for {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				for {
					select {
					case <-done:
						return
					case <-pulse:  // note
						sendPulse()
					case results <- r:
						return
					}
				}
			}
		}
	}()

	return heartbeat, results
}

func main() {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2*time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Println(r.Second())
		case <-time.After(timeout):
			return
		}
	}
}
