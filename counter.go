package main

type Counter interface {
	Increment(key string) int
	Get(key string) int
	Reset(key string)
}

type MemoryCounter struct {
	counts map[string]int
}

func NewMemoryCounter() *MemoryCounter {
	return &MemoryCounter{
		counts: make(map[string]int), // initializes the map
	}
}

func (m *MemoryCounter) Increment(key string) int {
	m.counts[key]++
	return m.counts[key]
}

func (m *MemoryCounter) Get(key string) int {
	return m.counts[key]
}

func (m *MemoryCounter) Reset(key string) {
	m.counts[key] = 0
}
