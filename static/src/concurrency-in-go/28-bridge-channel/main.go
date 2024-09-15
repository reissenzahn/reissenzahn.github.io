package main

// the bridge-channel destructures a channel of channels into a single channel

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

func bridge(done <-chan interface{}, in <-chan <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			var ch <-chan interface{}
			select {
			case maybeCh, ok := <-in:
				if !ok {
					return
				}
				ch = maybeCh
			case <-done:
				return
			}
			
			for v := range orDone(done, ch) {
				select {
				case out <- v:
				case <-done:
				}
			}
		}
	}()
	return out
}

func main() {

}