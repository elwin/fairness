package main

import (
	"fmt"
	"sync"
	"time"
)

const(
	consumers = 10
	retries = 1000
	timeout = 10 * time.Millisecond
)

// For a large enough timeout (so that each routine is ready for when a new struct{}
// arrives on the channel) it seems that the structs{} get distributed evenly,
// that is, each routine receives equally many structs{} from the channel
func main() {

	results := make([]int, consumers)
	channel := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(consumers)
	for i := 0; i < consumers; i++ {
		go consumer(channel, results, i, &wg)
	}

	producer(channel)
	for _, triggered := range results {
		if triggered != retries / consumers {
			fmt.Printf("Not evenly distributed: %v", results)
			return
		}
	}
	fmt.Printf("Evenly distributed, each routine received %d data-elements", retries / consumers)
}

func producer(channel chan<- struct{}) {
	for i := 0; i < retries; i++ {
		channel<- struct {}{}
	}
}

func consumer(channel <-chan struct{}, results []int, index int, wg *sync.WaitGroup) {
	defer wg.Done()
	for range channel {
		results[index]++
		time.Sleep(timeout)
	}
}