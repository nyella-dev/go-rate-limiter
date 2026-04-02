package main

type Limiter interface {
	Allow(key string) bool
}

type RateLimiter struct {
	counter Counter
	limit   int
}

func NewRateLimiter(limit int, counter Counter) *RateLimiter {
	return &RateLimiter{
		counter: counter,
		limit:   limit,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	count := rl.counter.Increment(key)
	if count > rl.limit {
		return false
	}

	return true
}
