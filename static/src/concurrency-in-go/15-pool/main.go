package main

import (
	"fmt"
	"sync"
)

// Pool is a concurrent-safe implementation of the object pool pattern; this provides a way to create and make available a fixed number of objects

// this is commonly used to constrain the creation of things that are expensive (e.g. database connections) so that only a fixed number of them are ever created but an indeterminate number of operations can still request access to these objects

// Get() will check whether there are any available instances within the pool to return to the caller, and if not, call its New() member variable to create a new one; when finished, callers call Put() to place the instance they were working with back in the pool

// a Pool can also be useful for warming a cache of pre-allocated objects for operations that must run as quickly as possible

// the New member variable should be thread-safe when called

const N = 1024 * 1024

func main() {
	// pool := &sync.Pool{
	// 	New: func() interface{} {
	// 		fmt.Println("Creating new instance")
	// 		return struct{}{}
	// 	},
	// }

	// pool.Get()
	// instance := pool.Get()
	// pool.Put(instance)
	// pool.Get()

	var count int
	pool := &sync.Pool{
		New: func() any {
			count += 1 // TODO: seems like this should be thread safe?
			buf := make([]byte, 1024)
			return &buf
		},
	}

	// seed the pool with 4KB
	pool.Put(pool.New())
	pool.Put(pool.New())
	pool.Put(pool.New())
	pool.Put(pool.New())

	var wg sync.WaitGroup
	wg.Add(N)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()

			buf := pool.Get().(*[]byte)
			defer pool.Put(buf)

			// do something
		}()
	}

	wg.Wait()
	fmt.Println(count)
}
