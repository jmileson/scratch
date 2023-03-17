package thrdprty

import (
	"os"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
)

const (
	envAPIKey  = "TEACHABLE_PAYMENTS_STRIPE_CONNECT_API_KEY"
	envAccount = "TEACHABLE_PAYMENTS_STRIPE_CONNECT_ACCOUNT_ID"
)

type Service struct {
	sc      *client.API
	apiKey  string
	Account string
}

type Balance struct {
	Amount   int64
	Currency string
}

func NewService() Service {
	// setup client with API key
	sc := client.API{}
	sc.Init(os.Getenv(envAPIKey), nil)

	return Service{
		sc:      &sc,
		Account: os.Getenv(envAccount),
	}
}

func (s *Service) GetBalance(account string) ([]Balance, error) {
	params := stripe.BalanceParams{}
	params.SetStripeAccount(account)

	balResp, err := s.sc.Balance.Get(&params)
	if err != nil {
		return nil, err
	}

	b := make([]Balance, len(balResp.Available))
	for i, amt := range balResp.Available {
		b[i] = Balance{
			Amount:   amt.Amount,
			Currency: string(amt.Currency),
		}
	}

	return b, nil
}
