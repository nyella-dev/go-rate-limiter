package main

import "sync"

type Counter interface {
	Increment(key string) int
	Get(key string) int
	Reset(key string)
}

type MemoryCounter struct {
	counts map[string]int
	mu     sync.Mutex
}

func NewMemoryCounter() *MemoryCounter {
	return &MemoryCounter{
		counts: make(map[string]int), // initializes the map
	}
}

func (m *MemoryCounter) Increment(key string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counts[key]++
	return m.counts[key]
}

func (m *MemoryCounter) Get(key string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.counts[key]
}

func (m *MemoryCounter) Reset(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counts[key] = 0
}
