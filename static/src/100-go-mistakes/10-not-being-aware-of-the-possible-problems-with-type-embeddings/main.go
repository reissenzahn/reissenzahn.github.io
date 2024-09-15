package main

import (
	"io"
	"sync"
)

// embedding promotes the fields and methods of the embedded type

// with embedding, the embedded type remains the receiver of a method; with sub-classing, the subclass becomes the receiver of a method

// embedding is also used within interfaces to compose an interface with other

type InMem struct {
	sync.Mutex
	m map[string]int
}

func NewInMem() *InMem {
	return &InMem{m: make(map[string]int)}
}

func (i*InMem) Get(key string) (int, bool) {
	i.Lock()
	v, ok := i.m[key]
	i.Unlock()
	return v, ok
}

func main() {
	i := NewInMem()

	// Lock() and Unlock() are promoted and visible to external clients which is probably not desirable
	i.Lock()
}

type Logger struct {
	io.WriteCloser
}

// embedding avoids implementing forwarding methods like this
// func (l Logger) Write(p []byte) (int, error) {
// 	return l.WriteCloser.Write(p)
// }

// embedding is rarely a necessity
