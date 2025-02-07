package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, done := context.WithCancelCause(context.Background())
	defer done(nil)

	ctx, gs := generator(ctx, 1, 12, 3 /*"a",*/, 34, 5)
	ctx, m2 := multiply(ctx, gs, 2)
	ctx, a1 := add(ctx, m2, 1)
	_, pipe := multiply(ctx, a1, 2)

	for {
		select {
		case <-time.After(2 * time.Second):
			done(context.DeadlineExceeded)
		case <-ctx.Done():
			if err := context.Cause(ctx); err != nil && err != context.Canceled {
				fmt.Printf("error %+v\n", err)
				return
			}
		case v, ok := <-pipe:
			if ok {
				fmt.Println(v)
			} else {
				return
			}
		}
	}
}

func generator(ctx context.Context, nums ...interface{}) (context.Context, <-chan int) {
	stream := make(chan int)
	ctx, done := context.WithCancelCause(ctx)

	go func() {
		defer close(stream)

		for _, i := range nums {
			select {
			case <-ctx.Done():
				return
			default:
				v, ok := i.(int)
				if !ok {
					done(fmt.Errorf("generator: value is not a int: %v", i))
					return
				}
				stream <- v
			}
		}
	}()

	return ctx, stream
}

func multiply(ctx context.Context, in <-chan int, value int) (context.Context, <-chan int) {
	stream := make(chan int)
	go func() {
		defer close(stream)
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-in:
				if !ok {
					return
				} else {
					stream <- i * value
				}
			}
		}
	}()
	return ctx, stream
}

func add(ctx context.Context, in <-chan int, value int) (context.Context, <-chan int) {
	ctx, _ = context.WithCancelCause(ctx)

	stream := make(chan int)
	go func() {
		defer close(stream)
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-in:
				if !ok {
					return
				} else {
					// if i > 10 {
					// done(fmt.Errorf("number greater than 10"))
					// return
					// time.Sleep(2200 * time.Millisecond)
					// }
					stream <- i + value
				}
			}
		}
	}()
	return ctx, stream
}
