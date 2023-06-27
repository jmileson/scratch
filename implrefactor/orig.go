package implrefactor

import "strconv"

// Concrete dependencies that require IO, and make it hard to unit test:
type ConcreteRepository struct{}

func (r *ConcreteRepository) Get(id int) (string, error) {
	return "repo: " + strconv.Itoa(id), nil
}

type ConcreteFetcher struct{}

func (f *ConcreteFetcher) Fetch(id int) (string, error) {
	return "fetch: " + strconv.Itoa(id), nil
}

// The SUT
type Impl struct {
	repo    *ConcreteRepository
	fetcher *ConcreteFetcher
}

func (i *Impl) Get(id int) (string, error) {
	acct, err := i.repo.Get(id)
	if err != nil {
		return "", err
	}

	fetched, err := i.fetcher.Fetch(id)
	if err != nil {
		return "", err
	}

	return acct + " " + fetched, nil
}
