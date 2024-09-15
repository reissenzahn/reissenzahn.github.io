package main

// assigning a value to an any type results in the loss of all type information and requires a type assertion to get anything useful out of the variable

// any can be helpful if there is a genuine need for accepting or returning any possible type (e.g. marshalling, formatting)
// func Marshal(v any) ([]byte, error) {}

// in general, we should avoid overgeneralized code; a little bit of duplication is better than a bad abstraction
