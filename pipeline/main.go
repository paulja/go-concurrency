package main

import "fmt"

func main() {
	done := make(chan interface{})
	defer close(done)

	stream := generator(done, 1, 12, 3, 34, 5)
	pipeline := multiply(done, add(done, multiply(done, stream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}

func generator(done chan interface{}, nums ...int) <-chan int {
	stream := make(chan int)
	go func() {
		defer close(stream)
		for _, i := range nums {
			select {
			case <-done:
				return
			case stream <- i:
			}
		}
	}()
	return stream
}

func multiply(done chan interface{}, in <-chan int, value int) <-chan int {
	stream := make(chan int)
	go func() {
		defer close(stream)
		for i := range in {
			select {
			case <-done:
				return
			case stream <- i * value:
			}
		}
	}()
	return stream
}

func add(done chan interface{}, in <-chan int, value int) <-chan int {
	stream := make(chan int)
	go func() {
		defer close(stream)
		for i := range in {
			select {
			case <-done:
				return
			case stream <- i + value:
			}
		}
	}()
	return stream
}
