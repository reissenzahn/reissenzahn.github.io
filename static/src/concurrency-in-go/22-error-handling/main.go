package main

import (
	"fmt"
	"net/http"
)

// in general, concurrent processes should send their errors to another part of the program that can make a more informed decision about what to do

// errors should be considered first-class citizens when constructing values to return from goroutines; if your goroutine can produce errors, those errors should be tightly coupled with your result type and passed through the same lines of communication

type Result struct {
	Err      error
	Response *http.Response
}

func fetchAll(done <-chan interface{}, urls ...string) <-chan Result {
	results := make(chan Result)
	go func() {
		defer close(results)

		for _, url := range urls {
			response, err := http.Get(url)
			result := Result{Err: err, Response: response}
			select {
			case <-done:
				return
			case results <- result:
			}
		}
	}()

	return results
}

func main() {
	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	errCount := 0

	for r := range fetchAll(done, urls...) {
		if r.Err != nil {
			fmt.Printf("error: %v\n", r.Err)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors!")
				break
			}
			continue
		}
		fmt.Printf("%v\n", r.Response.Status)
	}
}
