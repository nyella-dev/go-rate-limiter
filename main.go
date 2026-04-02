package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {

	// Environment variable setup
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	limit := 100
	if limitStr := os.Getenv("RATE_LIMIT"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	window := time.Minute
	if windowStr := os.Getenv("WINDOW"); windowStr != "" {
		if parsed, err := time.ParseDuration(windowStr); err == nil {
			window = parsed
		}
	}

	hostname, _ := os.Hostname()

	// End Environmnet Variable Prepation

	log.SetFlags(0)
	runtime.GOMAXPROCS(1)

	// Create Rate Limiter
	rateLimiter := NewRateLimiter(limit, NewRedisCounter(redisAddr), window)

	// Setup HTTP Endpoint
	http.HandleFunc("/", handleRequest(rateLimiter))

	// Setup HTTP Server
	fmt.Println("Server running on: " + hostname + " " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleRequest(rl *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		// If Token is Not Valid We Want to Rate Limit Based on IP
		var key string
		if token == "" {
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip, _, _ = net.SplitHostPort(r.RemoteAddr)
			}
			key = ip
		} else {
			key = token
		}

		if !rl.Allow(key) {
			log.Printf("RATE LIMITED | key=%s count=%d", key, rl.counter.Get(key))
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		log.Printf("ALLOWED | key=%s count=%d", key, rl.counter.Get(key))

		// Instead of Write header this would be a proxy to our API server once we have it up
		// target, _ := url.Parse("http://api:9090")
		// proxy := httputil.NewSingleHostReverseProxy(target)
		w.WriteHeader(http.StatusOK)
	}

}
