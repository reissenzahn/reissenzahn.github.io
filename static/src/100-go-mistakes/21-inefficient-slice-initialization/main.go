package main

// func convert(foos []Foo) []Bar {
// 	bars := make([]Bar, 0)
// 	for _, foo := range foos {
// 		bars = append(bars, fooToBar(foo))
// 	}
// 	return bars
// }

// avoid repeatedly allocating a new backing array
// func convert(foos []Foo) []Bar {
// 	n := len(foos)
// 	bars := make([]Bar, 0, n)
// 	for _, foo := range foos {
// 		bars = append(bars, fooToBar(foo))
// 	}
// 	return bars
// }

// avoid repeated calls to append(), which has a small overhead compared to direct assignment
// func convert(foos []Foo) []Bar {
// 	n := len(foos)
// 	bars := make([]Bar, n)
// 	for i, foo := range foos {
// 		bars[i] = fooToBar(foo)
// 	}
// 	return bars
// }

// Converting one slice type into another is a frequent operation for Go developers. As we have seen, if the length of the future slice is already known, there is no good rea- son  to  allocate  an  empty  slice  first.  Our  options  are  to  allocate  a  slice  with  either  a given capacity or a given length.

func main() {

}
