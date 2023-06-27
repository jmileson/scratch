package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	ch := make(chan bool, 1)
	simulateWork(1*time.Second, ch)

	status := http.StatusOK

	w.WriteHeader(status)

	ok := <-ch
	json.NewEncoder(w).Encode(map[string]bool{"result": ok})
}

func TimeoutHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrap := func(ch chan<- bool) {
			next.ServeHTTP(w, r)
			ch <- true
		}

		ch := make(chan bool, 1)
		go wrap(ch)

		select {
		case <-time.After(2 * time.Second):
			// write error response
			// log?
		case <-ch:
			// nothing to do
		}
	})
}

func register() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(http.ResponseWriter, *http.Request) {}))
	mux.Handle("/", TimeoutHandler(http.HandlerFunc(handleRequest)))

	return mux
}
