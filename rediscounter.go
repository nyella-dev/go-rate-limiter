package main

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisCounter struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCounter(addr string) *RedisCounter {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisCounter{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisCounter) Increment(key string) int {
	count, _ := r.client.Incr(r.ctx, key).Result()
	return int(count)
}

func (r *RedisCounter) Get(key string) int {
	val, _ := r.client.Get(r.ctx, key).Result()
	count, _ := strconv.Atoi(val) // convert string to int
	return count
}

func (r *RedisCounter) Reset(key string) {
	r.client.Del(r.ctx, key)
}
