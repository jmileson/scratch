package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkSucceeds(t *testing.T) {
	assert := assert.New(t)

	ok, err := work(1 * time.Second)

	assert.True(ok)
	assert.NoError(err)
}

// TODO: get this test working
func TestWorkFailsNormally(t *testing.T) {
	assert := assert.New(t)

	ok, err := work(1*time.Second + 1)

	assert.False(ok)
	assert.NoError(err)
}

// TODO: get this test working
// NOTE: this test ensures that the timeout is working right
func TestWorkFailsWithTimeout(t *testing.T) {
	assert := assert.New(t)

	ok, err := work(1 * time.Nanosecond)

	assert.False(ok)
	assert.ErrorIs(err, errTimeout)
}
