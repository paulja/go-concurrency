package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)

	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			// fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
}

func doWork(
	done <-chan interface{},
	pulseInterval time.Duration,
) (
	<-chan interface{},
	<-chan time.Time,
) {
	heartbeat := make(chan interface{}, 1)
	results := make(chan time.Time)
	go func() {
		defer close(heartbeat)
		defer close(results)

		pulse := time.Tick(pulseInterval)
		workgen := time.Tick(2 * pulseInterval)

		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
		}
		sendResult := func(r time.Time) {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case results <- r:
					return
				}
			}
		}

		for {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse()
			case r := <-workgen:
				sendResult(r)
			}
		}
	}()
	return heartbeat, results
}
