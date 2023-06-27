package implrefactor

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeGetter struct {
	get func(int) (string, error)
}

func (r *fakeGetter) Get(id int) (string, error) {
	return r.get(id)
}

type fakeFetcher struct {
	fetch func(int) (string, error)
}

func (f *fakeFetcher) Fetch(id int) (string, error) {
	return f.fetch(id)
}

func TestImplRefactor1Get(t *testing.T) {
	// now we're unit testing, but with lots of setup
	repo := fakeGetter{
		func(id int) (string, error) {
			return "repo: " + strconv.Itoa(id), nil
		},
	}
	fetcher := fakeFetcher{
		func(id int) (string, error) {
			return "fetch: " + strconv.Itoa(id), nil
		},
	}

	sut := ImplRefactor1{&repo, &fetcher}

	got, err := sut.Get(1)

	assert.Equal(t, "repo: 1 fetch: 1", got)
	assert.NoError(t, err)
}
