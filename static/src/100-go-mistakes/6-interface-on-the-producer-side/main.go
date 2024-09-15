package main

// an interface can either be defined in the same package as the concrete implementation (producer side) or in an external package where it is used (consumer side)

// the producer should not force an abstraction on all its clients; it should be up to the client to decide whether it needs some form of abstraction

// interfaces are satisfied implicitly which prevents creating circular dependencies

// this allows clients to define the most accurate abstraction for its need

// Interface Segregation Principle: no client should be forced to depend on methods it does not use

// in particular contexts, we may *know* an abstraction will be helpful for consumers, so we can have an interface on the producer side; but we should strive to keep it as minimal as possible

func main() {

}
