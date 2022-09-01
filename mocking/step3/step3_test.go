// Now that we altered Logger, we can update our test to
// see if our NoopWriteFlusher works, but we should also ensure
// that our Info method is actually trying to flush the destination
package step3_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/jmileson/scratch/mocking/step3"
	"github.com/stretchr/testify/assert"
)

func TestLogToBufferNoopFlush(t *testing.T) {
	assert := assert.New(t)

	// create an in memory buffer
	// and assign that to the logger
	// using our new adaptor
	buf := bytes.Buffer{}
	logger := step3.Logger{
		Out: step3.AddFlush(&buf),
	}

	logger.Info("sup")

	// out buffer should have something in it
	assert.Equal(4, buf.Len())
	assert.Equal("sup\n", buf.String())
}

// Like we said, our package doesn't care about the implementation
// details of the destination, so we can create a test object that
// has very simple behavior, but implments the required interface.
type fakeWriteFlusher struct {
	// any io.Writer will do
	io.Writer
	// this field allows us to provide behavior to the object
	// in our test setup, Flush delegates to the function stored here.
	flush func() int
}

// Implement the rest of the WriteFlusher interface, with delegation to
// the profided behavior
func (f *fakeWriteFlusher) Flush() int {
	return f.flush()
}

func TestLogToFlushableBuffer(t *testing.T) {
	assert := assert.New(t)

	// so we'll still create an in memory buffer
	// but instead of using our adaptor, we'll use a special
	// type that we can control the behavior of. We really only
	// care _that_ Flush was called, not what it does.
	buf := bytes.Buffer{}
	called := false
	writeFlusher := fakeWriteFlusher{
		Writer: &buf,
		flush: func() int {
			// closure captures the variable and updates the state
			// ONLY if the function is called.
			called = true
			return 0
		},
	}
	logger := step3.Logger{
		Out: &writeFlusher,
	}

	logger.Info("sup")

	// out buffer should have something in it
	assert.Equal(4, buf.Len())
	assert.Equal("sup\n", buf.String())
	// and Flush should have been called
	assert.True(called)
}
