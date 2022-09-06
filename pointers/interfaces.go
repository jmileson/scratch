package main

import "fmt"

type Fooable interface {
	Foo()
}

func (bs *BigStruct) Foo() {
	fmt.Printf("We're doing foo: %d\n", bs.value)
}

func WorkWithFooable(fooable Fooable) {
	fooable.Foo()
}
