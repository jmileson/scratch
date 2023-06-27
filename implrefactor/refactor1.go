package implrefactor

// The SUT

type Getter interface {
	Get(id int) (string, error)
}

type Fetcher interface {
	Fetch(id int) (string, error)
}

type ImplRefactor1 struct {
	repo    Getter
	fetcher Fetcher
}

func (i *ImplRefactor1) Get(id int) (string, error) {
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
