// I want to write a simple logger that writes somewhere
// maybe a terminal, a file, or an in-memory buffer.
// Ultimately, I want to be able to test that my
// logger writes where and what I expect, and that
// the destination I'm writing to doesn't have any
// potential issues holding on to memory.
package step1

import "fmt"

// To start with, we'll just create a simple logger
type Logger struct{}

// And one method for logging
func (l *Logger) Info(msg string) {
	// Here we can print to stdout, but nothing else,
	// which doesn't match our requirements, but is a
	// good place to start.
	// But: this is wholly untestable, without somehow
	// capturing stdout. We can do better
	fmt.Println(msg)
}
