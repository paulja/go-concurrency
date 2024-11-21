package main

import "fmt"

func main() {
	var fib func(int) <-chan int
	fib = func(n int) <-chan int {
		res := make(chan int)
		go func() {
			defer close(res)
			if n <= 2 {
				res <- 1
				return
			}
			res <- <-fib(n-1) + <-fib(n-2)
		}()
		return res
	}
	x := 30
	fmt.Printf("fib(%d) = %d", x, <-fib(x))
}
