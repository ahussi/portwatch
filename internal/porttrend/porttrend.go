// Package porttrend tracks binding frequency over time to detect anomalous port activity.
package porttrend

import (
	"sync"
	"time"
)

// Entry holds trend data for a single port key.
type Entry struct {
	Key       string
	Count     int
	FirstSeen time.Time
	LastSeen  time.Time
}

// Tracker records how often each port binding is observed.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	clock   func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]*Entry),
		clock:   time.Now,
	}
}

// WithClock replaces the clock used for timestamps (useful in tests).
func WithClock(fn func() time.Time) func(*Tracker) {
	return func(t *Tracker) { t.clock = fn }
}

// NewWithOptions returns a Tracker with functional options applied.
func NewWithOptions(opts ...func(*Tracker)) *Tracker {
	t := New()
	for _, o := range opts {
		o(t)
	}
	return t
}

// Record increments the observation count for key.
func (t *Tracker) Record(key string) {
	now := t.clock()
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[key]
	if !ok {
		t.entries[key] = &Entry{Key: key, Count: 1, FirstSeen: now, LastSeen: now}
		return
	}
	e.Count++
	e.LastSeen = now
}

// Get returns the Entry for key and whether it exists.
func (t *Tracker) Get(key string) (Entry, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[key]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot of all entries.
func (t *Tracker) All() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, *e)
	}
	return out
}

// Reset removes all tracked data.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[string]*Entry)
}
