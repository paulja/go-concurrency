package main

import "fmt"

func main() {
	// tee
	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2), 4))
	for v := range out1 {
		fmt.Printf("Out1: %v, Out2: %v\n", v, <-out2)
	}

	// bridge
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}
	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}
}

func orDone(done, c <-chan interface{}) <-chan interface{} {
	stream := make(chan interface{})
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case stream <- v:
				case <-done:
				}
			}
		}
	}()
	return stream
}

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	stream := make(chan interface{})
	go func() {
		defer close(stream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case stream <- v:
				}
			}
		}
	}()
	return stream
}

func take(done, instream <-chan interface{}, num int) <-chan interface{} {
	outstream := make(chan interface{})
	go func() {
		defer close(outstream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case outstream <- <-instream:
			}
		}
	}()
	return outstream
}

func tee(done, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})
	go func() {
		defer close(out1)
		defer close(out2)
		for val := range orDone(done, in) {
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()
	return out1, out2
}

func bridge(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if !ok {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range orDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}
