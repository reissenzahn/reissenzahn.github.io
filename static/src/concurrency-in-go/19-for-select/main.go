package main

// for { // Either loop infinitely or range over something
// 	select {
// 	// Do some work with channels
// 	}
// }


// Sending iteration variables out on a channel
// for _, s := range []string{"a", "b", "c"} {
//     select {
//     case <-done:
//         return
//     case stringStream <- s:
//     }
// }

// Looping infinitely waiting to be stopped
// for {
//     select {
//     case <-done:
//         return
//     default:
//     }
//     // Do non-preemptable work
// }

// for {
//     select {
//     case <-done:
//         return
//     default:
//         // Do non-preemptable work
//     }
// }
