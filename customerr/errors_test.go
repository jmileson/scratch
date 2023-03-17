package main

import (
	"fmt"
	"testing"
	"time"
)

func TestThing(t *testing.T) {
	go main()
	fmt.Println("hello")
	time.Sleep(5 * time.Second)
	t.Fail()
}
