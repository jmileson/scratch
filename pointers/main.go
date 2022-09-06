package main

func main() {
	// STRUCTS
	bs := BigStruct{1}
	DoSomethingWithBigStruct(&bs)
	ModifyBigStruct(&bs, 100)

	// SLICES
	vals := []int{1, 2, 3, 4}
	WorkWithSlice(vals)
	WorkWithSlicePointer(&vals) // in general, don't do this

	// INTERFACES
	WorkWithFooable(&bs)

	// GOTCHAS
	// Prints 4, 4, 4, 4; bad news
	LoopWithPointerCopy(vals)
	LoopWithPointerCopyExpanded(vals)
	// Prints 1, 2, 3, 4; nice!
	LoopWithPointerCopyFixed(vals)

	vals2 := []Obj{
		{1},
		{2},
		{3},
	}
	// Prints 1, 2, 3; oh no...
	LoopWithStruct(vals2)
	// Prints 2, 3, 4; hell yes
	LoopWithStructFixed(vals2)
}
