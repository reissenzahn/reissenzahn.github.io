package main

// we can attach either a value or a pointer receiver to a method

// with a value receiver, a copy of the value will be made and passed to the method; any changes to the object remain local to the method

// with a pointer receiver, the address of the object will be passed to the method

// a receiver must be a pointer if the method needs to mutate the receiver or if the method receiver contains a field that cannot be copied

// a receiver should be a pointer if the receiver is a large object to prevent making an expensive copy

// a receiver must be a value if it is necessary to enforce immutability of the receiver or if the receiver is a map, function or channel

// a receiver should be a value if the receiver is a slice that does not need to be mutated, if the receiver is a small array or struct that is naturally a value type without mutable fields or if the receiver is a primitive type

// while we are allowed to mix receiver types, the consensus tends toward forbidding it

// by default, we should choose to go with a value receiver

type account struct {
	balance float64
}

func (a *account) add(v float64) {
	a.balance += v
}

func main() {

}
