package main

// returning an interface is often considered a bad practice

// this creates a dependency between the implementation package and the cleint interface
// func NewInMemoryStore() client.Store {}

// this restricts flexibility because it forces all clients to use one particular type for abstraction and causes circular dependencies between packages

// if a client needs to abstract an implementation for whatever reason, it can still do that on its side

// "Be conservative in what you do, be liberal in what you accept from others." - Postel's law

// return structs instead of interfaces; accept interfaces if possible

// an exception to this is up-front abstractions like error and io.Reader which are forced by the language designers rather than being defined by clients
// func LimitReader(r Reader, n int64) Reader {}
