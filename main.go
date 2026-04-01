package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	rateLimiter := NewRateLimiter(10)

	http.HandleFunc("/", handleRequest(rateLimiter))
	log.SetFlags(0)
	fmt.Println("Server running on: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
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
