---
title: "Go"
date: "2024-08-03"
tags: ["Programming"]
---


## Fundamentals

### Value Types
```go
package main

import "fmt"

func main() {
  // boolean
  var t bool = true
  fmt.Printf("%v\n", t)

  // signed integer (int8, int16, int32, int64)
  var i int32 = -1
  fmt.Printf("%v\n", i)

  // unsigned integer (uint8, uint16, uint32, uint64)
  var u uint32 = 42
  fmt.Printf("%v\n", u)

  // floating-point (float32, float64)
  var d float64 = 3.14
  fmt.Printf("%v\n", d)

  // machine-sized integer (int, uint)
  var j int = 7
  fmt.Printf("%v\n", j)

  // byte (alias for uint8)
  var b byte = 128
  fmt.Printf("%v\n", b)

  // rune (alias for uint32)
  var r rune = 'A'
  fmt.Printf("%q", r)

  // string
  var s string = "hello"
}
```

### Variables
```go
package main

import "fmt"

func main() {
  // uninitialized variables are assigned an appropriate zero value (false, 0, 0.0, "")

  // declare a variable
  var a bool
  fmt.Printf("%v\n", a)

  // declare multiple variables of the same type
  var b, c int
  fmt.Printf("%v\n", b)
  fmt.Printf("%v\n", c)

  // declare and initialize a variable
  var d string = "hello"
  fmt.Printf("%v\n", d)

  // declare and initialize multiple variables
  var e, f rune = 'A', 'B'
  fmt.Printf("%v\n", e)
  fmt.Printf("%v\n", f)

  // declare and initialize a variable with type inference
  var g = 3.14
  fmt.Printf("%v\n", g)

  // declare and initialize multiple variables with type inference
  var h, i = true, -1
  fmt.Printf("%v\n", h)
  fmt.Printf("%v\n", i)

  // declare and initialize a variable using short form (within functions only)
  j := 7
  fmt.Printf("%v\n", j)

  // declare and initialize multiple variables using short form
  k, l := "test", 2.71
  fmt.Printf("%v\n", k)
  fmt.Printf("%v\n", l)
}
```

### Constants
```go
package main

// constants must be initialized at declaration and are immutable

// declare and initialize a constant
const max int = 10

// declare and initialize a constant with type inference
const min = 10

func main() {
}
```

### Blank Identifier
```go
package main

import "fmt"

func main() {
  // the blank identifier does not introduce a binding
  _ = true

  // this can be used to selectively ignore unneeded values in an assignment
  _, err := fmt.Println()
  if err != nil {
    panic(err)
  }
}
```

### Strings
```go
package main

import "fmt"

// a string represents an immutable sequence of unicode code points

func main() {
  // interpreted string literal
  s := "hello\nworld"
  fmt.Println(s)

  // raw string literal
  t := `hello
world`
  fmt.Println(t)

  // string can be freely converted to and from byte slices
  fmt.Println(string([]byte("hello")))
}
```

### Arrays
```go
package main

import "fmt"

// an array is a fixed-length sequence of elements of a particular type

func main() {
  // declare an array
  var arr [5]int

  // the zero value of an array is an array of the specified length of zero-valued elements
  fmt.Println(arr)

  // initialize an array
  arr = [5]int{1, 2, 3, 4, 5}

  // initialize an array omitting an explicit length
  arr = [...]int{1, 2, 3, 4, 5}

  // initialize an array by specifying only some of the elements
  arr = [...]int{1, 3: 4, 5}

  // initialize multi-dimensional array
  var mat = [2][3]int{
    {1, 2, 3},
    {1, 2, 3},
  }
  fmt.Println(mat[1][2])

  // assign a value to an array element
  arr[1] = 2

  // retrieve array element at specified index
  fmt.Println(arr[1])

  // determine the length of an array
  fmt.Println(len(arr))
}
```

### Slices
```go
package main

import "fmt"

// a slice is a resizable sequence of elements of a particular type

func main() {
  var slc []int

  // initialize a slice of zero values
  slc = make([]int, 3)

  // initialize slice
  slc = []int{1, 2, 3}

  // append values to a slice
  slc = append(slc, 4)
  slc = append(slc, 5, 6)

  // determine length of slice
  fmt.Println(len(slc))
}
```

// array capacity
cap(s)

// slice operator creates a new slice backed by the given array or slice
s[1:3]
s[:3]
s[1:]
s[:]

// changes to a slice will be reflected in the underlying array
s[3] = 42

### Maps
```go
package main

import "fmt"

// a map is an associative data structure that stores key-value pairs; any type that is comparable using == may be used as a key

func main() {
  // declare a map of string keys to integers
  var tbl map[string]int

  // initialize an empty map
  tbl = make(map[string]int)

  // initialize a map
  tbl = map[string]int{
    "foo": 42,
    "bar": -1,
  }

  // assign key-value pair
  tbl["foo"] = 42

  // retrieve value associated with a given key
  fmt.Println(tbl["foo"])

  // the zero value of the value type is returned if the specified key does not exist
  fmt.Println(tbl["baz"])

  // test the membership of a key in a map
  _, ok := tbl["baz"]
  if !ok {
    fmt.Println("missing")
  }

  // determine the number of key-value pairs in a map
  fmt.Println(len(tbl))

  // delete a key-value pair
  delete(tbl, "foo")

  // delete all key-value pairs
  clear(tbl)
}
```

### Pointers
```go
package main

import "fmt"

func main() {
  // a pointer stores the address of a variable
  var p *int

  // the zero value of a pointer is nil
  fmt.Println(p)
  
  // retrieve the address of a variable
  v := 1
  p = &v

  // dereference a pointer to retrieve the corresponding value
  fmt.Println(*p)

  // assign the value of a variable via a pointer
  *p = 2
  fmt.Println(v)

  // pass a variable to a function by pointer
  func (p *int) {
    *p++
  }(&v)
  fmt.Println(v)

  // a pointer can be safely returned from a function as the runtime will transfer values to the heap as required
  q := func () *int {
    v := 2
    return &v
  }()

  // pointers are equal only if they contain the same address or are both nil
  fmt.Println(p == q)

  // a pointer to a struct will be automatically dereferenced when accessing a field
  r := &struct {
    v int
  }{1}
  fmt.Println((*r).v)
  fmt.Println(r.v)
}
```

### If/Else
```go
package main

import "fmt"

func main() {
  // if statement
  if 8%4 == 0 {
    fmt.Println("8 is divisible by 4")
  }

  // if statement with else branch
  if 7 % 2 == 0 {
    fmt.Println("7 is even")
  } else {
    fmt.Println("7 is odd")
  }

  // a statement can precede the conditional and any variables declared in this statement are available in the subsequent branches
  if _, err := fmt.Println("hello"); err != nil {
    panic(err)
  }
}
```

### For
```go
package main

import "fmt"

func main() {
  // iterate using initial/condition/after
  for i := 0; i < 3; i++ {
    fmt.Println(i)
  }

  // iterate while condition holds true
  i := 0
  for i < 3 {
    fmt.Println(i)
    i++
  }

  // iterate forever
  for {
    // break out of loop
    break
  }

  // iterate N times
  for i := range 6 {
    // continue to next iteration of loop
    if i%2 == 0 {
      continue
    }
    fmt.Println(i)
  }

  // iterate over slice
  for i, v := range []int{1, 2, 3} {
    fmt.Println(i, "->", v)
  }

  // iterate over map
  for k, v := range map[string]int{"foo": 42, "bar": -1} {
    fmt.Println(k, "->", v)
  }
}
```

### Switch
```go
// switch statement
i := 0
switch i % 3 {
  case 0:
    fmt.Println("zero")
    // explicit keyword required to fallthrough to next case
    fallthrough
  case 1:
    fmt.Println("one")
  case 2:
    fmt.Println("two")
  default:
    fmt.Println("???")
}

// switch statement with statement and expressions in case statements
switch hour := time.Now().Hour(); {
case hour >= 5 && hour < 9:
  fmt.Println("I'm writing")
case hour >= 9 && hour < 18:
  fmt.Println("I'm working")
default:
  fmt.Println("I'm sleeping")
}
```

### Functions


### Structs
```go
package main

import "fmt"

// a struct is a typed collection of fields

// declare a struct
type Rectangle struct {
  Length int
  Width  int
}

// declare a constructor function for the struct
func NewRectangle(length, width int) *Rectangle {
  return &Rectangle{
    Length: length,
    Width: width,
  }
}

func main() {
  var r Rectangle

  // initialize a struct using ordered fields
  r = Rectangle{3, 4}

  // initialize a struct using named fields
  r = Rectangle{Width: 4, Length: 3}

  // omitted fields will be zero-valued
  fmt.Println(Rectangle{Width: 4})

  // assign and read struct fields
  r.Length = 3
  fmt.Println(r.Length)

  // declare and initialize anonymous struct
  s := struct {
    Length int
    Width  int
  }{3, 4}
  fmt.Println(s)
}
```

### Methods
```go
package main

import "fmt"

// methods can be defined on struct types

type Rectangle struct {
  Length int
  Width  int
}

// define a method on a value receiver type
func (r Rectangle) Area() int {
  return r.Length * r.Width
}

// define a method on a pointer receiver type

func (r *Rectangle) Scale(factor int) {
  r.Length *= factor
  r.Width *= factor
}

func main() {
  // conversion between values and pointers are handled automatically for method calls
  r := Rectangle{Length: 2, Width: 3}
  r.Scale(2)
  fmt.Println(r.Area())

  p := &r
  p.Scale(2)
  fmt.Println(r.Area())
}
```

### Interfaces
```go

```

### Errors
```go
package main

import (
  "errors"
  "fmt"
)

// errors are represented by the built-in error interface type which defines a single Error() method which returns a string

// create an error from a string
var divideByZeroErr = errors.New("divide by zero")

// an error is returned as the last return value by convention
func divide(x, y int) (int, error) {
  if y == 0 {
    return 0, divideByZeroErr
  }

  return x / y, nil
}

func main() {
  // check for the presence of an error and handle it
  _, err := divide(1, 0)
  if err != nil {
    fmt.Println(err)
  }
}
```

### Goroutines
```go
package main

import (
  "fmt"
  "time"
)

func greet() {
  fmt.Println("hello")
}

// a goroutine is a function that executes concurrently with other goroutines

// the main goroutine is created automatically when the process runs
func main() {
  // start a goroutine from a function call
  go greet()

  // start a goroutine from an anonymous function call
  go func() {
    fmt.Println("hello")
  }()

  // goroutines can reference variables even if they have fallen out of scope as the runtime will transfer those variables to the heap as required
  s := "hello"
  go func() {
    fmt.Println(s)
  }()

  time.Sleep(time.Second)
}
```

### Channels
```go
package main

// a channel serves as a conduit for a stream of values of a particular type

func main() {
  // declare a channel
  var ch chan int

  // initialize an empty channel
  ch = make(chan int)

  // declare a unidirectional channel (receive only)
  var rch 

  var wch
}
```



## Standard Library

### sync

#### WaitGroup
```go
package main

import (
  "fmt"
  "sync"
)

// a wait group can be used to wait for a set of concurrent operations to complete provided the results of those operations do not need to be collected

const N = 5

func main() {
  var wg sync.WaitGroup

  // increment counter by N
  wg.Add(N)

  for i := 0; i < N; i++ {
    go func() {
      // decrement counter
      defer wg.Done()
      fmt.Println(i)
    }()
  }

  // block until counter is zero
  wg.Wait()
}
```

#### Once
```go
package main

import (
  "fmt"
  "sync"
)

// sync.Once ensures that only one call to Do() ever calls the provided callback

const N = 100

func main() {
  var i int

  var wg sync.WaitGroup
  wg.Add(N)

  var once sync.Once

  for i := 0; i < N; i++ {
    go func() {
      defer wg.Done()

      once.Do(func() {
        i++
      })
    }()
  }

  // sync.Once only counts the number of times Do() is called, not how many times unique functions passed into Do() are called
  once.Do(func() {
    i--
  })

  fmt.Println(i)
}
```

#### Mutex
```go
package main

import (
  "fmt"
  "sync"
)

// a mutex allows for synchronizing access to a shared resource

const N = 1000

func main() {
  var i int
  var mu sync.Mutex
  var wg sync.WaitGroup

  for i := 0; i < N; i++ {
    wg.Add(1)
    go func() {
      defer wg.Done()

      mu.Lock()

      // Unlock() is commonly called within a defer statement to ensure it will be called even if the function panics
      defer mu.Unlock()

      i++
    }()

    wg.Add(1)
    go func() {
      defer wg.Done()
      mu.Lock()
      defer mu.Unlock()
      i--
    }()
  }

  wg.Wait()
  fmt.Println(i)
}
```


#### RWMutex
```go
package main

import (
  "time"
  "sync"
)

// a RWMutex allows for an arbitrary number of reader to hold a reader lock or a single writer to hold a writer lock

const N = 5
const M = 1000

func writer(wg *sync.WaitGroup, l sync.Locker) {
  defer wg.Done()
  for i := 0; i < N; i++ {
    l.Lock()
    l.Unlock()
    time.Sleep(time.Nanosecond)
  }
}

func reader(wg *sync.WaitGroup, l sync.Locker) {
  defer wg.Done()
  l.Lock()
  l.Unlock()
}

func main() {
  var mu sync.RWMutex

  var wg sync.WaitGroup

  wg.Add(1)
  go writer(&wg, &mu)

  wg.Add(M)
  for i := 0; i < M; i++ {
    go reader(&wg, mu.RLocker())
  }

  wg.Wait()
}
```


