package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	consumers = 10
	retries   = 1000
	timeout   = 10 * time.Millisecond
)

// In this scenario we want to find out, if the channel operates using a FIFO channel,
// that is if the first routine to listen on the channel will be the first one to receive.
// For that we let all threads queue up after each iteration randomly (as opposed to stay in queue).
// If the output is once again evenly distributed this suggests it is not a FIFO queue, instead having
// state that keeps track of each routine. If the output is roughly evenly distributed but seems to contain
// randomness (e.g. ƒollowing a gaussian normal distribution), this suggests that a FIFO could have been used.
func main() {

	results := make([]int, consumers)
	channel := make(chan struct{})
	release := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(consumers)
	for i := 0; i < consumers; i++ {
		go consumer(channel, release, results, i, &wg)
	}

	producer(channel, release)
	fmt.Println(results)
}

func producer(channel, release chan<- struct{}) {
	for i := 0; i < retries; i++ {
		channel <- struct{}{}
		for j := 0; j < consumers - 1; j++ {
			release <- struct{}{}
		}
		time.Sleep(timeout)
	}
}

func consumer(channel, release <-chan struct{}, results []int, index int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-channel:
			results[index]++
		case <-release:
		}
		time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
	}
}
