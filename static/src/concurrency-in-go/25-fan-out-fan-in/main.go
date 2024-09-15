package main

import "sync"

// fan-out refers to starting multiple goroutines to handle input from a pipeline; fan-in refers to combining multiple results into one channel

// consider fanning out a stage if it starts running slow and has order-independence

func fanIn(done <-chan interface{}, ins ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	var wg sync.WaitGroup
	wg.Add(len(ins))

	for _, in := range ins {
		go func() {
			defer wg.Done()
			for v := range in {
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	// done := make(chan interface{})
	// defer close(done)
	// start := time.Now()
	// rand := func() interface{} { return rand.Intn(50000000) }
	// randIntStream := toInt(done, repeatFn(done, rand))
	// numFinders := runtime.NumCPU()
	// fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	// finders := make([]<-chan interface{}, numFinders)
	// fmt.Println("Primes:")
	// for i := 0; i < numFinders; i++ {
	// 	finders[i] = primeFinder(done, randIntStream)
	// }
	// for prime := range take(done, fanIn(done, finders...), 10) {
	// 	fmt.Printf("\t%d\n", prime)
	// }
	// fmt.Printf("Search took: %v", time.Since(start))
}
