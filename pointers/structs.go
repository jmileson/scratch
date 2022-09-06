package main

import "fmt"

// BigStruct represents a struct with many fields.
type BigStruct struct {
	value int
}

// DoSomethingWithBigStruct avoids copying BigStruct when passing in arg.
func DoSomethingWithBigStruct(bs *BigStruct) {
	fmt.Printf("Value is: %d\n", bs.value)
}

// ModifyBigStruct modifies big struct in place.
func ModifyBigStruct(bs *BigStruct, anotherValue int) {
	bs.value = anotherValue
	DoSomethingWithBigStruct(bs)
}
