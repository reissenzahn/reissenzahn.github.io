package main

import (
	"errors"
	"fmt"
)

// the %w directive allows for wrapping errors conveniently

// error wrapping refers to packing an error inside a wrapper container that also makes the source error available

// this can be used to add additional context to an error or mark an error as being of a particular type

// a caller can handle the error by unwrapping it and checking the source error

// wrapping an error by creating a custom error type is cumbersome
type BarError struct {
	Err error
}

func (b BarError) Error() string {
	return "bar failed: " + b.Err.Error()
}

// the %w allows for wrapping a source error to add additional context without having to create another error type

// the source error remains available

// a client can unwrap the parent error and check whether the source error was of a specific type or value

// however, making the source error available to callers also introduces potential coupling; to avoid this we can use %v to transform the error instead of wrapping it

func main() {
	err := errors.New("oops")

	fmt.Println(BarError{Err: err})

	fmt.Println(fmt.Errorf("bar failed: %w", err))
}

// To summarize, when handling an error, we can decide to wrap it. Wrapping is about adding additional context to an error and/or marking an error as a specific type. If we need to mark an error, we should create a custom error type. However, if we just want to  add  extra  context,  we  should  use  fmt.Errorf  with  the  %w  directive  as  it  doesn’t require creating a new error type. Yet, error wrapping creates potential coupling as it makes the source error available for the caller. If we want to prevent it, we shouldn’t use error wrapping but error transformation, for example, using fmt.Errorf with the %v directive.