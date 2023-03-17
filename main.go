package main

import (
	"log"

	"github.com/jmileson/scratch/thrdprty"
)

func main() {
	srv := thrdprty.NewService()

	bal, err := srv.GetBalance(srv.Account)
	if err != nil {
		log.Fatalf("unable to get account balance: %v\n", err)
	}

	log.Printf("Account %s balances:", srv.Account)
	for _, b := range bal {
		log.Printf("%s %f", b.Currency, float64(b.Amount)/100.0)
	}

}
