// Package cooldown provides a per-key cooldown tracker that enforces a minimum
// duration between successive events for the same key. Unlike ratelimit, cooldown
// does not count events — it simply gates repeated triggers.
package cooldown

import (
	"sync"
	"time"
)

// Cooldown tracks the last-seen timestamp for arbitrary string keys and
// reports whether enough time has elapsed since the previous trigger.
type Cooldown struct {
	mu       sync.Mutex
	duration time.Duration
	last     map[string]time.Time
}

// New returns a Cooldown with the given minimum interval between triggers.
// If d is zero or negative it defaults to one second.
func New(d time.Duration) *Cooldown {
	if d <= 0 {
		d = time.Second
	}
	return &Cooldown{
		duration: d,
		last:     make(map[string]time.Time),
	}
}

// Ready reports whether key is eligible to trigger again.
// It returns true the first time a key is seen, and again only after the
// configured duration has elapsed since the last successful call to Record.
func (c *Cooldown) Ready(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	t, ok := c.last[key]
	if !ok {
		return true
	}
	return time.Since(t) >= c.duration
}

// Record marks key as having triggered right now.
func (c *Cooldown) Record(key string) {
	c.mu.Lock()
	c.last[key] = time.Now()
	c.mu.Unlock()
}

// ReadyAndRecord atomically checks readiness and, if ready, records the
// trigger. It returns true when the caller should proceed.
func (c *Cooldown) ReadyAndRecord(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if t, ok := c.last[key]; ok && now.Sub(t) < c.duration {
		return false
	}
	c.last[key] = now
	return true
}

// Reset removes the cooldown state for key so the next call to Ready returns true.
func (c *Cooldown) Reset(key string) {
	c.mu.Lock()
	delete(c.last, key)
	c.mu.Unlock()
}

// Prune removes all entries whose last-trigger time is older than the cooldown
// duration, freeing memory for keys that are no longer active.
func (c *Cooldown) Prune() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for k, t := range c.last {
		if now.Sub(t) >= c.duration {
			delete(c.last, k)
		}
	}
}
