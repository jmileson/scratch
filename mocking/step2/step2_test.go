// Things are testable now too, since we can provide a
// destination, we can create a buffer, call Logger.Info
// and check that the result is what we expected.
package step2_test

import (
	"bytes"
	"testing"

	"github.com/jmileson/scratch/mocking/step2"
	"github.com/stretchr/testify/assert"
)

func TestLogToBuffer(t *testing.T) {
	assert := assert.New(t)

	// create an in memory buffer
	// and assign that to the logger
	buf := bytes.Buffer{}
	logger := step2.Logger{&buf}

	logger.Info("sup")

	// out buffer should have something in it
	assert.Equal(4, buf.Len())
	assert.Equal("sup\n", buf.String())
}
