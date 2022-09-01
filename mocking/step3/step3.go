// There's an issue with our implementation though:
// if the destination is some kind of in-memory buffer
// and the logger instance is long lived, we might
// eventually consume all the memory on the machine.
// To be safer, we should ask that the destination
// to flush it's contents to ensure that in-memory
// destinations can periodically clean out.
// Since our package doesn't control the destination
// or the implementation of Flush, we expect callers
// to think about their needs and implement accordingly.
package step3

import (
	"fmt"
	"io"
)

// Instead of using the io.Writer exclusively, we
// ask for some additional behavior from our destination.
type WriteFlusher interface {
	io.Writer
	Flush() int
}

type Logger struct {
	// And we update our struct accordingly
	Out WriteFlusher
}

func (l *Logger) Info(msg string) {
	fmt.Fprintln(l.Out, msg)
	l.Out.Flush()
}

// We know that destinations like stdout and files handle
// flushing on their own when Writing, so we can provide
// a simple adaptor to support conversions for io.Writers
// that don't need to worry about memory consumption
type NoopWriteFlusher struct {
	io.Writer
}

func (n *NoopWriteFlusher) Flush() int {
	// it's OK to have a noop here, the writer
	// is expected to handle it's memory consumption
	return 0
}

// Helper function to make the interface nicer
func AddFlush(w io.Writer) WriteFlusher {
	return &NoopWriteFlusher{w}
}
