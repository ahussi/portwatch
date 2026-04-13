// Package portstate tracks the lifecycle state of individual monitored ports,
// recording when they were first seen, last seen, and how many times they
// have been observed across scan cycles.
package portstate

import (
	"sync"
	"time"
)

// State represents the observed lifecycle state of a single port binding.
type State struct {
	Key       string
	FirstSeen time.Time
	LastSeen  time.Time
	SeenCount int
}

// Tracker maintains per-binding state across scan cycles.
type Tracker struct {
	mu     sync.RWMutex
	states map[string]*State
	now    func() time.Time
}

// New returns a new Tracker. If now is nil, time.Now is used.
func New(now func() time.Time) *Tracker {
	if now == nil {
		now = time.Now
	}
	return &Tracker{
		states: make(map[string]*State),
		now:    now,
	}
}

// Observe records a sighting of the given binding key, creating a new State
// entry on first observation or updating an existing one.
func (t *Tracker) Observe(key string) {
	now := t.now()
	t.mu.Lock()
	defer t.mu.Unlock()
	if s, ok := t.states[key]; ok {
		s.LastSeen = now
		s.SeenCount++
		return
	}
	t.states[key] = &State{
		Key:       key,
		FirstSeen: now,
		LastSeen:  now,
		SeenCount: 1,
	}
}

// Get returns the State for the given key and whether it was found.
func (t *Tracker) Get(key string) (State, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	s, ok := t.states[key]
	if !ok {
		return State{}, false
	}
	return *s, true
}

// Remove deletes the state entry for the given key.
func (t *Tracker) Remove(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.states, key)
}

// Len returns the number of tracked bindings.
func (t *Tracker) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.states)
}

// Keys returns all currently tracked binding keys.
func (t *Tracker) Keys() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	keys := make([]string, 0, len(t.states))
	for k := range t.states {
		keys = append(keys, k)
	}
	return keys
}
