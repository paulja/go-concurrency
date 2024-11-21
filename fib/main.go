package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	var num int
	flag.IntVar(&num, "n", 0, "number param")
	flag.Parse()

	done := make(chan interface{})
	defer close(done)

	select {
	case n := <-fib(done, num):
		fmt.Printf("fib(%d) = %d\n", num, n)
	case <-time.Tick(5 * time.Second):
		done <- true
		fmt.Printf("request timed out (%d)\n", num)
	}
}

func fib(done <-chan interface{}, n int) <-chan int {
	res := make(chan int)
	go func() {
		defer close(res)

		select {
		case <-done:
			return
		default:
		}

		if n <= 2 {
			res <- 1
			return
		}
		res <- <-fib(done, n-1) + <-fib(done, n-2)
	}()
	return res
}
