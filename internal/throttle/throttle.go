// Package throttle provides a token-bucket style throttle for limiting
// how frequently portwatch emits alerts for a given port or key.
package throttle

import (
	"sync"
	"time"
)

// Throttle limits event emission to at most Burst events per Window per key.
type Throttle struct {
	mu     sync.Mutex
	bucket map[string]*entry
	window time.Duration
	burst  int
}

type entry struct {
	tokens    int
	windowEnd time.Time
}

// New creates a Throttle that allows up to burst events within each window
// duration per key. Subsequent calls within the same window are denied.
func New(window time.Duration, burst int) *Throttle {
	if burst < 1 {
		burst = 1
	}
	return &Throttle{
		bucket: make(map[string]*entry),
		window: window,
		burst:  burst,
	}
}

// Allow returns true if the event identified by key should be allowed through.
// It consumes one token from the bucket for the key. When the window expires
// the bucket is refilled to the burst limit.
func (t *Throttle) Allow(key string) bool {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()

	e, ok := t.bucket[key]
	if !ok || now.After(e.windowEnd) {
		t.bucket[key] = &entry{
			tokens:    t.burst - 1,
			windowEnd: now.Add(t.window),
		}
		return true
	}
	if e.tokens > 0 {
		e.tokens--
		return true
	}
	return false
}

// Remaining returns the number of tokens left in the current window for key.
func (t *Throttle) Remaining(key string) int {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()

	e, ok := t.bucket[key]
	if !ok || now.After(e.windowEnd) {
		return t.burst
	}
	return e.tokens
}

// Prune removes expired entries to reclaim memory.
func (t *Throttle) Prune() {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()

	for k, e := range t.bucket {
		if now.After(e.windowEnd) {
			delete(t.bucket, k)
		}
	}
}
