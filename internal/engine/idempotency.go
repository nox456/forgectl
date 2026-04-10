package engine

import "sync"

type IdempotencyGuard struct {
	mu      sync.Mutex
	results map[string]map[string]any
}

func NewIdempotencyGuard() *IdempotencyGuard {
	return &IdempotencyGuard{
		results: make(map[string]map[string]any),
	}
}

func (g *IdempotencyGuard) CheckOrClaim(key string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.results[key]; exists {
		return true
	}

	g.results[key] = nil

	return false
}
