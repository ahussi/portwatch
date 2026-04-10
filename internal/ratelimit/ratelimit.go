// Package ratelimit provides a simple token-bucket rate limiter
// used to suppress duplicate alerts for the same port binding.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter suppresses repeated events for the same key within a cooldown window.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	seen     map[string]time.Time
}

// New creates a Limiter with the given cooldown duration.
// Events with the same key will be suppressed until the cooldown elapses.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		seen:     make(map[string]time.Time),
	}
}

// Allow returns true if the event for key should be processed,
// or false if it is within the cooldown window of a previous event.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if last, ok := l.seen[key]; ok && now.Sub(last) < l.cooldown {
		return false
	}
	l.seen[key] = now
	return true
}

// Reset removes the cooldown record for key, allowing the next event through immediately.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.seen, key)
}

// Prune removes all expired entries from the internal map.
// Call periodically to prevent unbounded memory growth.
func (l *Limiter) Prune() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for k, t := range l.seen {
		if now.Sub(t) >= l.cooldown {
			delete(l.seen, k)
		}
	}
}

// Len returns the number of tracked keys currently in the limiter.
func (l *Limiter) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.seen)
}
