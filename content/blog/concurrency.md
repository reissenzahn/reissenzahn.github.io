

## Overview

Concurrency is a property of the code; parallelism is a property of the running program.

### Atomicity
- An operation is considered to be atomic if it is indivisible or uninterruptible within a given context.
- Combining atomic operations does not necessarily produce a larger atomic operation.

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	// the increment operation is not atomic within a concurrent context: retrieve the value, increment the value, store the value
	x := 0

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for i := 0; i < 1_000_000; i++ {
			x++
		}
		wg.Done()
	}()

	for i := 0; i < 1_000_000; i++ {
		x++
	}

	wg.Wait()
	fmt.Println(x)
}

```


### Memory Access Synchronization
- A critical section is a section of a program that needs exclusive access to a shared resource.
- A mutex can be used to synchronize access to a shared resource.
- Synchronizing access to a shared resource does not necessarily solve for correctness as the order of execution of critical sections is not deterministic.

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	x := 0

	var mu sync.Mutex

	go func() {
		// critical section 1
		mu.Lock()
		// exclusive access to the shared variable
		x++
		mu.Unlock()
	}()

	// critical section 2
	mu.Lock()
	if x == 0 {
		fmt.Println("0")
	} else {
		fmt.Println(x)
	}
	mu.Unlock()
}
```

### Race Conditions
- A race condition occurs when two or more operations must execute in a certain order but the program has not been written such that this order is guaranteed to be maintained.


#### Data Races
- A data race is a type of race condition that occurs when one concurrent operation attempts to read a variable while at some undetermined time another concurrent operation is attempting to write to that variable.

```go
package main

import "fmt"

// three possible outcomes are possible: nothing is printed, "0" is printed or "1" is printed

func main() {
  x := 0

  go func() {
    x++
  }()

  if x == 0 {
    fmt.Printf("%v\n", x)
  }
}
```


### Deadlock
- A deadlocked program is one in which all concurrent processes are waiting on one another.
- The Coffman conditions must be present for deadlocks to arise:
  1. Mutual exclusion: A concurrent process holds exclusive rights to a resource at any one time.
  2. Wait-for: A concurrent process must simultaneously hold a resource and be waiting for an additional resource.
  3. No preemption: A resource held by a concurrent process can only be released by that process.
  4. Circular wait: A concurrent process must be waiting on a chain of other processes which are in turn waiting on it.

```go
package main

import (
  "sync"
  "time"
)

func main() {
  // mutual exclusion
  var m1, m2 sync.Mutex

  var wg sync.WaitGroup
  wg.Add(2)

  go func() {
    defer wg.Done()

    m1.Lock()
    defer m1.Unlock()

    time.Sleep(time.Second)

    // wait-for
    m2.Lock()
    defer m2.Unlock()
  }()

  go func() {
    defer wg.Done()

    m2.Lock()
    defer m2.Unlock()

    time.Sleep(time.Second)

    // wait-for
    m1.Lock()
    defer m1.Unlock()
  }()

  wg.Wait()
}
```

### Starvation
- Starvation refers to a scenario where a concurrent process cannot obtain all the resources it needs to perform work.
- This usually implies that there are one or more greedy concurrent processes that are unfairly preventing one or more other processes from accomplishing work as efficiently as possible, or even at all.
- Memory access synchronization requires finding a balance between preferring coarse-grained synchronization for performance and fine-grained synchronization for fairness.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex

  var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		count := 0
		for begin := time.Now(); time.Since(begin) <= 1*time.Second; {
			mu.Lock()
			time.Sleep(3*time.Nanosecond)
			mu.Unlock()
			count++
		}

		fmt.Printf("greedy = %v\n", count)
	}()

	go func() {
		defer wg.Done()

		count := 0
		for begin := time.Now(); time.Since(begin) <= 1*time.Second; {
			for i := 0; i < 3; i++ {
				mu.Lock()
				time.Sleep(1 * time.Nanosecond)
				mu.Unlock()
			}
			count++
		}

		fmt.Printf("polite = %v\n", count)
	}()

	wg.Wait()
}
```

#### Livelocks
- A livelock is a type of starvation that occurs when a program is actively performing concurrent operations but these operations do nothing to advance the state of the program.
- All te concurrent processes are starved equally and no work is accomplished.
- This commonly occurs when concurrent processes are attempting to prevent a deadlock without coordination.

```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// simulate a constant cadence
	cadence := sync.NewCond(&sync.Mutex{})
	go func() {
		for range time.Tick(1 * time.Millisecond) {
			cadence.Broadcast()
		}
	}()

	takeStep := func() {
		cadence.L.Lock()
		cadence.Wait()
		cadence.L.Unlock()
	}

	// attempt to move in a direction and return whether the attempt was successful
	tryDirection := func(direction *int32) bool {
		atomic.AddInt32(direction, 1)

		takeStep()
		if atomic.LoadInt32(direction) == 1 {
			return true
		}

		takeStep()
		atomic.AddInt32(direction, -1)
		return false
	}

	// maintain intention to move in a direction
	var left, right int32

	walk := func(wg *sync.WaitGroup, name string) {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			fmt.Printf("%s scoots left\n", name)
			if tryDirection(&left) {
				fmt.Printf("%s was successful\n", name)
				return
			}

			fmt.Printf("%s scoots right\n", name)
			if tryDirection(&right) {
				fmt.Printf("%s was successful\n", name)
				return
			}
		}

		fmt.Printf("%v gives up\n", name)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go walk(&wg, "Alice")
	go walk(&wg, "Barbara")

	wg.Wait()
}
```

