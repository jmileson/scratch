package main

import (
	"encoding/json"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	ok, err := work()

	status := http.StatusOK
	if err != nil {
		status = http.StatusGatewayTimeout
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]bool{"result": ok})
}

func register() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRequest)

	return mux
}
