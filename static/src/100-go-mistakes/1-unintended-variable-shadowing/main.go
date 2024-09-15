package main

import "fmt"

// variable shadowing allows a name declared in a block to be re-declared in an inner block

type Counter struct {
	count int
}

func (c *Counter) Increment() {
	// this will panic
	c.count += 1
}

func NewCounter() (*Counter, error) {
	return &Counter{}, nil
}

func NewCounterWithTracing() (*Counter, error) {
	return &Counter{}, nil
}

func problem(tracing bool) error {
	var counter *Counter
	if tracing {
		// assign result of function call to inner client variable
		counter, err := NewCounterWithTracing()
		if err != nil {
			return err
		}
		fmt.Println(counter)
	} else {
		counter, err := NewCounter()
		if err != nil {
			return err
		}
		fmt.Println(counter)
	}

	// counter is always nil
	counter.Increment()

	return nil
}

func solution(tracing bool) error {
	var counter *Counter
	if tracing {
		// create an error variable and use the assignment operator
		var err error
		counter, err = NewCounterWithTracing()
		if err != nil {
			return err
		}
	} else {
		// or introduce a temporary variable and assign that variable to the outer variable
		c, err := NewCounter()
		if err != nil {
			return err
		}
		counter = c
	}

	counter.Increment()

	return nil
}

func main() {
	problem(true)
}
