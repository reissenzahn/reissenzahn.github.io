
A B C D E F G H I J K L M N O P Q R S T U V W X Y Z

## Abstraction
- Abstractions should be discovered, not created.

## Control flow
- Avoid shadowing variables.
- Avoid deep nesting and keep the happy path aligned on the left; return as early as possible.

## Comments

## Coupling & Cohesion

## Data Transfer Objects

## Dependencies

## Dependency Injection

## Documentation

## Dynamic Configuration

- Feature flags
- Percent migrations

## Encapsulation
- Avoid excessive use of getters and setters.

## Exceptions

## Formatters & Linters
- To improve code quality and consistency, use linters and formatters.

## Generics

## Inheritance

## Interfaces
- Create an interface when it is needed, not when you forsee needing it, or at least prove the abstraction to be a valid one.
- Keep interfaces on the client side.
- Functions should accept interfaces and return concrete implementations.

## Interviews


## Noisy Neighbors
- The basic problem this  term  refers  to  is  that  other  applications  running  on  the  same  physical  system  as yours can have a noticeable impact on your performance and resource availability.

## Perfect Accounting & Janitors


## Project structure
- Standardize a layout.
- Give packages meaningful and specific names.
- Avoid creating packages like common, util or shared.

## Primitive Types
- Integer overflows (https://en.wikipedia.org/wiki/Ariane_flight_V88).
- Comparing floats using a delta.


## Polymorphism

## Queuing
- Queue depth allows for trading off latency for fault tolerance during large load spikes.
- Increasing queue depth due to a continuous rate of rejections is almost always bad.

## Repositories

## Requirements
- Functional requirements
- Non-functional requirements
  - Scalability: The ability of a system to continue to behave as expected in the face of significant upward or downward changes in demand.
    - Vertical scaling
    - Horizontal scaling
  - Loose coupling: The components of a system have minimal knowledge of each other; changes to one component generally do not require changes to other components.
    - Service contracts, protocols
    - Distributed monoliths
  - Resilience: A measure of how well a system withstands and recovers from errors and faults; the degree to which it can continue to operate correctly in the face of errors and faults.
    - Redundancy
    - Partial failures
    - Circuit breakers and retries
  - Reliability: The ability of a system to behave as expected for a given time interval.
  - Manageability: The ease (or lack thereof) with which changes can be made to the behavior of a running system, up to and including deploying system components.
    - Configuration changes
    - Feature flags
    - Credential rotation
    - Certificate renewal
    - Deployments
    - Patching
  - Maintainability: The ease with which changes can be made to the functionality of a system, most often its code.
  - Observability: A measure of how well the internal states of a system can be inferred from its external outputs.
    - Metrics, logging, tracing
    - "Data is not information"

## SOLID

## Stack & Heap

## State Machines

## Sweepers


<!-- BREAK -->


## Dynamic Configuration

It should always be possible to configure feature flags statically on a per-dimension basis. In addition, dynamic configuration can be used.

- Audit trail.

## Deployment Mechanism

Ideally, all changes should be deployed as part of the CI/CD release process with safety provided by automated tests and monitoring. In reality, this is not always feasible:

- Monitoring gaps result in a bad change being promoted through the pipeline.
- Rollbacks cause significant disruption to the team (along with confusion when trying to recall the bad commit).
- One box environments do not receive sufficient traffic to fully exercise the new code. The smallest deployment unit receives too much traffic to limit blast radius.
- The change is being deployed ahead of its dependencies in which case it cannot be enabled by default.
- The change is being released along with unrelated changes that cannot be easily rolled back.
- The change is too risky for automated deployments and requires manual enablement and validation steps.
- The change is not compatible with all deployment units as a result of a snowflake architecture.

Feature flag based deployments can sometimes strike a better balance between automation and manual action and validation. It will often only be necessary to manually test the feature in a subset of deployment environments after which a static configuration change can be deployed to enable the feature flag by default in all environments.

## Percent Migration

To further reduce the blast radius of a potentially bad change, a percent migration applied to the new code path in a single deployment unit.

This can be implemented in a few different ways for a chosen migration rate, \(r \in [0, 1]\):

- Pick a random number \(x \in [0, 1]\) and check \(x \le p\).
- Maintain a count, \(c\). For each request, check \(c \mod \frac{1}{r} < 1\) and increment \(c\).
- Concatenate the user ID and flag name to obtain \(s\), check \(h(s) \mod 100 < r * 100\).

```go
type RateCounter struct {
	rateProvider func() float64
	counter      atomic.Uint64
}

func (rc *RateCounter) IsOpen() (bool, error) {
	rate := rc.rateProvider()
	if !(0.0 <= rate && rate <= 1.0) {
		return false, fmt.Errorf("Invalid rate %f", rate)
	}

	count := rc.counter.Add(1)
	period := 1.0 / rate

	return math.Mod(float64(count) - 1, period) < 1, nil
}

func main() {
	rc := RateCounter{rateProvider: func() float64 {
		return 0.7
	}}

	total := 0
	for i := 0; i < 100; i++ {
		isOpen, err := rc.IsOpen()
		if err != nil {
			log.Fatal(err)
		}

		if isOpen {
			total++
		}

		log.Println(isOpen)
	}

	log.Println(total)
}
```

```go
// isInPercentage check if the user is in the cohort for the toggle.
func (f *FlagData) isInPercentage(flagName string, user ffcontext.Context) bool {
	percentage := int32(f.getActualPercentage())
	maxPercentage := uint32(100 * percentageMultiplier)

	// <= 0%
	if percentage <= 0 {
		return false
	}
	// >= 100%
	if uint32(percentage) >= maxPercentage {
		return true
	}

	hashID := utils.Hash(flagName+user.GetKey()) % maxPercentage
	return hashID < uint32(percentage)
}
```

<!-- https://github.com/thomaspoignant/go-feature-flag/blob/c84e9326f895c67913f04949c4a76645c18da48f/testutils/flagv1/flag_data.go#L125 -->

## Private Beta

Sometimes we can choose the test cohort explicitly. If this is the case, provide access to the feature as a "private beta" experience for select customers with the understanding that it is not in its final state.

## Scheduled Enablement/Expiration

It should be possible to specify "active after" and "inactive after" dates for feature flags.


## Feature Flag Hygiene

Feature flags should be regularly reviewed and removed from code where possible.

Feature flags require visibility to be useful; otherwise, they become a liability.


## Toothless Mode

There exists a trade-off between the extent to which the new code path runs and how toothless that new code path actually is.


## De-fanged Side-effects


## Common Pitfalls

- Don't read the value of a feature flag twice on the same code path (unless this is explicitly intended). If the flag is set dynamically, then it's value may have changed between reads with unpredictable results.



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

## Overview

- More visibility is required as system become more complicated.
- Observability is a system property that reflects how well the internal states of a system can be inferred from knowledge of its external outputs. A system can be considered observable when it is possible to quickly and consistently ask novel questions about it with minimal prior knowledge and without having to re-instrument or build new code. An observable system lets you ask it questions that you have not thought of yet.
  - Traditionally, monitoring focuses on asking questions in the hope of identifying or predicting some expected or previously observed failure modes.
  - At a certain level of complexity, the number of "unknown unknowns" starts to overwhelm the number of "known unknowns". Failures are more often unpredicted and monitoring for every possible failure mode becomes effectively impossible.
  - Understanding all possible failure (or non-failure) states in a complex system is pretty much impossible. The first step to achieving observability is to stop looking for specific, expected failure modes.
  - The three pillars of observability:
    - Tracing: Follows a request as it propagates through a system, allowing the entire end-to-end request flow to be reconstructed as a DAG called a trace.
    - Metrics: The collection of numerical data points representing the state of various aspects of a system at specific points in time.  Collections  of  data  points, representing  observations  of  the  same  subject  at  various  times,  are  particularly useful  for  visualization  and  mathematical  analysis,  and  can  be  used  to  highlight trends, identify anomalies, and predict future behavior
    - Logging: Logging is the process of appending records of noteworthy events to an immuta‐ ble  record—the  log—for  later  review  or  analysis.  A  log  can  take  a  variety  of forms, from a continuously appended file on disk to a full-text search engine like Elasticsearch. Logs provides valuable, context-rich insight into application- specific events emitted by processes. However, it’s important that log entries are properly structured; not doing so can sharply limit their utility
  - A truly observable system will interweave tracing, metrics and logging so that each can reference the others. For instance, metrics might be used to track down a subset of misbehaving traces and those traces might highlight logs that could help to find the underlying cause of the behavior.
  -

## Logs

Logging  is  the  act  of  recording  events  that  occur  during  the  running  of  a  program.
It is often an undervalued activity in programming because it is additional work that
has little immediate payback for the programmer.
During  the  normal  operations  of  a  program,  logging  is  an  overhead,  taking  up
processing  cycles  to  write  to  a  file,  database,  or  even  to  the  screen.  In  addition,
unmanaged  logs  can  cause  problems.  The  classic  case  of  logfiles  getting  so  big  that
they take up all the available disk space and crash the server is too real and happens
too often.
However, when something happens, and you want to find out the sequence of events
that led to it, logs become an invaluable diagnostic resource. Logs can also be moni‐
tored in real time, and alerts can be sent out when needed.

### Log Rotation

### Log Retention

### Sensitive User Data

- Logging complex types

> My analogy for user-data is that its like a pool of toxic sludge. You might be willing to poke it with a stick to see how deep it is, but you definitely wouldn't dive into it.


### Compliance



## Metrics

### Faults
Zero-ed metrics.

### Errors

### Latency

### Saturation

#### Starvation

Note our technique here for identifying the starvation: a metric. Starvation makes for
a good argument for recording and sampling metrics. One of the ways you can detect
and solve starvation is by logging when work is accomplished, and then determining
if your rate of work is as high as you expect it.

### Cache Metrics

https://guava.dev/releases/snapshot/api/docs/com/google/common/cache/CacheStats.html

### Queue Metrics


### Database Query Metrics


## Alarms

### Severity levels

### Aggregate hierarchies

### Alarm fatigue

### Anomaly detection
Very skeptical.


## Tracing

By tracking requests as they propagate through the system (even across process, network and security boundaries) tracing can help you to pinpoint component failures, identify performance bottlenecks and analyze service dependencies.

There are two fundamental concepts:

- Spans: A span describes a unit of work performed by a request, such as a fork in the execution flow or hop across the network, as it propagates through a system. Each span has an associated name, a start time and a duration. They can be and typically are nested and ordered to model casual relationships.
- Traces: 

