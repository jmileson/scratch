package thrdprty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBalanceReturnsError(t *testing.T) {
	assert := assert.New(t)

	assert.FailNow("not implemented")
}

func TestGetBalance(t *testing.T) {
	assert := assert.New(t)

	type getBalanceTC struct {
		input string
		exp   []Balance
	}

	tests := []getBalanceTC{}
	for _, tc := range tests {
		// TODO: write some tests
		_ = tc
	}

	assert.FailNow("not implemented")
}
