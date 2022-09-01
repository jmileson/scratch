// Now that we know we can't really test fmt.Println,
// we can use a little bit of indirection and ask callers
// to provide a destination for where logs should be written.
// This makes things testable, but also helps support our requirement
// that we can write to terminals, files, or in-memory buffers!
package step2

import (
	"fmt"
	"io"
)

type Logger struct {
	// This is an interface provided by the standard library.
	// Any type that can Write([]byte) (int, error) satisfies it
	// including files, stdout and bytes.Buffer
	Out io.Writer
}

func (l *Logger) Info(msg string) {
	// Now we use Fprintln to specify our out
	fmt.Fprintln(l.Out, msg)
}
