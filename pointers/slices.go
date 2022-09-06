package main

import "fmt"

// WorkWithSlice is preferred, slices are already pointers to the underlying array
// so passing a copy of a slice is always very cheap, and doesn't require any
// additional work.
func WorkWithSlice(vals []int) {
	for i, v := range vals {
		fmt.Printf("i: %d value: %d\n", i, v)
	}
}

// WorkWithSlice is NOT preferred, the indirection doesn't give us any additional
// performance and requires extra work.
func WorkWithSlicePointer(vals *[]int) {
	if vals == nil {
		return
	}

	for i, v := range *vals {
		fmt.Printf("i: %d value: %d\n", i, v)
	}
}
