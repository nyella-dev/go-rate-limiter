package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
)

func main() {
	log.SetFlags(0)
	runtime.GOMAXPROCS(1)
	rateLimiter := NewRateLimiter(10, NewRedisCounter("localhost:6379"))
	http.HandleFunc("/", handleRequest(rateLimiter))

	// log every panic
	defer func() {
		if r := recover(); r != nil {
			log.Println("Server crashed:", r)
		}
	}()

	fmt.Println("Server running on: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(rl *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		// If Token is Not Valid We Want to Rate Limit Based on IP
		var key string
		if token == "nil" {
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip, _, _ = net.SplitHostPort(r.RemoteAddr)
			}
			key = ip
		} else {
			key = token
		}

		if !rl.Allow(key) {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintln(w, "Rate limit exceeded")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Request allowed")
	}

}
