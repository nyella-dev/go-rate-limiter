package main

import "time"

type Limiter interface {
	Allow(key string) bool
}

type RateLimiter struct {
	counter Counter
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, counter Counter, window time.Duration) *RateLimiter {
	return &RateLimiter{
		counter: counter,
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	count := rl.counter.Increment(key, rl.window)
	if count > rl.limit {
		return false
	}

	return true
}
