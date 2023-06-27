package implrefactor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImplGet(t *testing.T) {
	// using the concrete implementations makes these
	// tests integration tests and requires a lot more env setup
	repo := ConcreteRepository{}
	fetcher := ConcreteFetcher{}

	sut := Impl{&repo, &fetcher}

	got, err := sut.Get(1)

	assert.Equal(t, "repo: 1 fetch: 1", got)
	assert.NoError(t, err)
}
