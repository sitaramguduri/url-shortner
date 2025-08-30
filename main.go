package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"

	"tinyurl/ratelimit"
)

func main(){
	mux := http.NewServeMux()

	mux.HandleFunc("health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintln(w, "pong")
	})

	ipLimiters := ratelimit.NewLimiterMap(func() *rate.Limiter{
		return rate.NewLimiter(rate.Limit(20), 20)
	})

	handler :=logging(perIP(ipLimiters)(mux))


	srv := &http.Server{
		Addr:	":8080",
		Handler: handler,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Println("listening on :8080")
	log.Fatal(srv.ListenAndServe())
}

func logging(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func perIP(lm *ratelimit.LimiterMap) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			key := clientIP(r)
			if !lm.Allow(key){
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string{
	if xf := r.Header.Get("X-forwarded-for"); xf != ""{
		return strings.TrimSpace(strings.Split(xf, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}