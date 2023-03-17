package main

import (
	"fmt"
	"io"
)

type BigStruct struct {
	// pretend lots of fields
	Field string
}

func write(w io.Writer) error {
	_, err := w.Write([]byte("test"))

	return err
}

func (bs *BigStruct) Write(val []byte) (int, error) {
	bs.Field = string(val)

	return len(val), nil
}

func main() {
	bs := BigStruct{}

	write(&bs)
	// write(bs)

	fmt.Println(bs.Field)
}
