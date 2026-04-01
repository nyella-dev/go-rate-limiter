package main

type Limiter interface {
	Allow(key string) bool
}

type RateLimiter struct {
	counter *MemoryCounter
	limit   int
}

func NewRateLimiter(limit int) *RateLimiter {
	return &RateLimiter{
		counter: NewMemoryCounter(),
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
