package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sitaram/go-rate-limiter/ratelimit"
	"golang.org/x/time/rate"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	// Per-IP limiter: 20 req/sec, burst of 20
	ipLimiters := ratelimit.NewLimiterMap(func() *rate.Limiter {
		return rate.NewLimiter(rate.Limit(20), 20)
	})

	// Build middleware chain
	handler := logging(perIP(ipLimiters)(mux))

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Println("Listening on :8080")
	log.Fatal(srv.ListenAndServe())
}
