package main

import (
	"errors"
	"time"
)

var (
	errTimeout = errors.New("timeout")
)

func simulateWork(timeout time.Duration, ch chan<- bool) {
	time.Sleep(100 * time.Millisecond)

	ch <- (timeout % 2) == 0
}

func work(timeout time.Duration) (bool, error) {
	// TODO: support timeout on this function

	// simulates work being done
	ch := make(chan bool, 1)
	go simulateWork(timeout, ch)

	select {
	case <-time.After(timeout):
		return false, errTimeout
	case ok := <-ch:
		return ok, nil
	}
}
