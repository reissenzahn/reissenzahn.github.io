package main

// interface pollution refers to overwhelming code with unnecessary abstractions

// aim to find the right level of abstraction and interface granularity

// the bigger the interface, the weaker the abstraction; adding methods to an interface can decrease its level of reusability
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

// small, well-defined interfaces can be combined to create higher-level abstractions

type ReadWriter interface {
	Reader
	Writer
}

// three concrete use-cases for interfaces are:

// 1. when multiple types implement a common behavior, we can factor out the behavior into an interface
type Sortable interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

// 2. decoupling from an implementation (Liskov Substitution Principle)
type Customer struct{
	id string
}

type CustomerRepository interface {
	Put(Customer) error
}

type CustomerService struct {
	repository CustomerRepository
}

func (cs CustomerService) CreateCustomer(id string) error {
	c := Customer{id: id}
	return cs.repository.Put(c)
}

// this is particularly helpful for tests

// 3. use interfaces to restrict a type to a specific behavior, such as for semantics enforcement

type Config struct{}

func (c *Config) Get() string {
	return ""
}

func (c *Config) Set(value string) {
}

type ConfigReader interface {
	Get() string
}

func main() {
	c := &Config{}

	var cr ConfigReader = c
	cr.Get()
}

// abstractions should be discovered, not created; we shouldn't start creating abstractions if there is no immediate reason to do so

// overusing interfaces makes the code flow more complex by introducing a useless level of indirection

// "Don’t design with interfaces, discover them." - Rob Pike
