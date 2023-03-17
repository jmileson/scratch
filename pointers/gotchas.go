package main

import "fmt"

type Obj struct {
	value int
}

func printSlice(vals []*int) {
	for i, v := range vals {
		fmt.Printf("%d is %d\n", i, *v)
	}
}

// LoopWithPointerCopy illustrates a common issue when working
// with loops: taking the pointer to the loop variable.
func LoopWithPointerCopy(vals []int) {
	copy := make([]*int, len(vals))
	// Here we're looping over the range and on each iteration
	// assigning a value from vals into the loop variable v
	for i, v := range vals {
		// This causes problems, but why?
		copy[i] = &v
	}

	printSlice(copy)
}

// LoopWithPointerCopyExpanded is a rewrite of LoopWithPointerCopy
// that "desugars" the loop in that function to explain what's
// happening. Though not exactly what the Go compiled is doing
// conceptually this function behaves the same.
func LoopWithPointerCopyExpanded(vals []int) {
	copy := make([]*int, len(vals))

	// v here is the same as the loop variable v from the original function.
	// It is declared once, given it's own address and zero-value
	var v int
	for i := 0; i < len(vals); i++ {
		// inside the loop, we assign a value from the original loop
		// into our v, changing it from it's last assigned value, into
		// the new value. However, the address for v remains the same
		// on each iteration, regardless of it's value.
		v = vals[i]
		// Thus, each iteration we assign the same address of v into
		// different elements of the copy slice
		copy[i] = &v
	}

	// And when we print it out, we get n repetitions of the
	// last value of v.
	printSlice(copy)
}

// LoopWithPointerCopyFixed shows you how to avoid the issue described above.
func LoopWithPointerCopyFixed(vals []int) {
	copy := make([]*int, len(vals))
	// We can still use the sugared syntax from the original function
	for i, v := range vals {
		// v2 is declared on each iteration, not just once, which means
		// it has a distinct address on each iteration of the loop
		v2 := v
		// So this works
		copy[i] = &v2
	}

	printSlice(copy)
}

func printSliceObj(vals []Obj) {
	for _, v := range vals {
		fmt.Printf("Obj is: %v\n", v)
	}
}

// LoopWithStruct shows a related issue with operations that modify in place.
func LoopWithStruct(vals []Obj) {
	// A similar problem arises when you try to modify values in place
	for _, v := range vals {
		v.value += 1
	}

	printSliceObj(vals)
}

// LoopWithStructFixed shows that it isn't really possible to do modification
// in place with the sugared syntax because the sugared syntax includes an
// implicit copy.
func LoopWithStructFixed(vals []Obj) {
	// A similar problem arises when you try to modify values in place
	for i := range vals {
		vals[i].value += 1
	}

	printSliceObj(vals)
}
