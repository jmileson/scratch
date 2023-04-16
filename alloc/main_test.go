package main

import "testing"

func BenchmarkF1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t, err := f1(i)
		_ = t
		_ = err
	}
}

func BenchmarkF2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t, err := f2(i)
		_ = t
		_ = err
	}
}
