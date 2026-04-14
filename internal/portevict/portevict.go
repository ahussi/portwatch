// Package portevict tracks port eviction events — bindings that were
// present in a previous scan but are no longer observed.
package portevict

import (
	"sync"
	"time"
)

// Event represents a single port eviction (a binding that disappeared).
type Event struct {
	Key       string
	Port      int
	Proto     string
	Addr      string
	EvictedAt time.Time
}

// Tracker records and queries port eviction events.
type Tracker struct {
	mu     sync.RWMutex
	events []Event
	cap    int
}

// New creates a new Tracker with the given maximum capacity.
// If cap is <= 0 it defaults to 256.
func New(capacity int) *Tracker {
	if capacity <= 0 {
		capacity = 256
	}
	return &Tracker{cap: capacity}
}

// Record appends an eviction event, evicting the oldest entry when full.
func (t *Tracker) Record(e Event) {
	if e.EvictedAt.IsZero() {
		e.EvictedAt = time.Now()
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.events) >= t.cap {
		t.events = t.events[1:]
	}
	t.events = append(t.events, e)
}

// Since returns all eviction events that occurred at or after the given time.
func (t *Tracker) Since(ts time.Time) []Event {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var out []Event
	for _, e := range t.events {
		if !e.EvictedAt.Before(ts) {
			out = append(out, e)
		}
	}
	return out
}

// All returns a copy of all recorded eviction events in insertion order.
func (t *Tracker) All() []Event {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Event, len(t.events))
	copy(out, t.events)
	return out
}

// Len returns the current number of recorded events.
func (t *Tracker) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.events)
}

// Clear removes all recorded events.
func (t *Tracker) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = t.events[:0]
}
