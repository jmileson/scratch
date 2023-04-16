package main

import (
	"errors"
	"math/rand"
	"strings"
)

var errOhNo = errors.New("oh no")

type T struct {
	ID          int
	Description string
	Number      float64

	Vals []string
}

func f1(i int) (T, error) {
	if i%2 == 0 {
		return T{}, errOhNo
	}

	vals := make([]string, 0)
	for j := 0; j < i%10; j++ {
		vals = append(vals, "something")
	}

	return T{
		ID:          i,
		Description: strings.Repeat("a", i%10),
		Number:      rand.Float64(),
		Vals:        vals,
	}, nil
}

func f2(i int) (T, error) {
	var t T
	if i%2 == 0 {
		return t, errOhNo
	}

	vals := make([]string, 0)
	for j := 0; j < i%10; j++ {
		vals = append(vals, "something")
	}

	t.ID = i
	t.Description = strings.Repeat("a", i%10)
	t.Number = rand.Float64()
	t.Vals = vals

	return t, nil
}
