package main

// the tee-channel takes a channel to read from and returns two channels that will get the same value

func orDone(done, in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}

				select {
				case out <- v:
				case <-done:
				}
			}
		}
	}()
	return out
}

func tee(done <-chan interface{}, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1, out2 := make(chan interface{}), make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)

		for v := range orDone(done, in) {
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1 <- v:
					out1 = nil
				case out2 <- v:
					out2 = nil
				}
			}
		}
	}()

	return out1, out2
}
