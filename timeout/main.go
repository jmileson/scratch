package main

import (
	"errors"
	"time"
)

var (
	errTimeout = errors.New("timeout")
)

func work() (bool, error) {
	// TODO: support a failure case without a timeout, i.e. return false, nil
	// TODO: support timeout on this function

	// simulates work being done
	time.Sleep(100 * time.Millisecond)

	return true, nil
}
