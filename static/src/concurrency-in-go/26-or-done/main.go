package main

import "fmt"

// the or-done-channel pattern encapsulates the boilerplate involved in handling the done channel

func orDone(done, in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}

				select {
				case out <- v:
				case <-done:
				}
			}
		}
	}()
	return out
}

func main() {
	done := make(chan interface{})
	ch := make(chan interface{})

	go func() {
		defer close(done)
		ch <- "A"
		ch <- "B"
		
		// TODO: kinda weird that this is not guaranteed to be processed
		ch <- "C"
	}()

	for v := range orDone(done, ch) {
		fmt.Println(v)
	}
}
