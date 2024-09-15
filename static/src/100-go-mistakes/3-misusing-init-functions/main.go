package main

import (
	"fmt"
	"log"
)

// an init function takes no arguments and returns no result; it is executed when the package is initialized after all constant and variable declarations are evaluated

var count = func() int {
	fmt.Println("variable declaration evaluated") // 1
	return 0
}()

func init() {
	fmt.Println("init() function called") // 2
}

func main() {
	fmt.Println("main() function called") // 3
}

// if a package depends on another package then the init() function in the other package will be called first

// if there are multiple init() functions defined in a package then they will be executed according to the alphabetical order of the source files

// if multiple init() functions are defined in a single source file then they will be executed in the source order

// init() functions can be used for side effects and a package can be imported for its side effects only
// import (
//   _ "foo"
// )

// init() functions cannot be invoked directly

// initializing a database connection in an init() function has downsides:
// - error management is limited to panic()
// - the init() function will be executed before running any tests
// - this approach requires using a global variable

var db *Database

type Database struct {}

func (d *Database) Ping() error {
	return nil
}

func NewDatabase() (*Database, error) {
	return &Database{}, nil
}

func init() {
	d, err := NewDatabase()
	if err != nil {
		log.Panic(err)
	}

	if err := d.Ping(); err != nil {
		log.Panic(err)
	}

	db = d
}



