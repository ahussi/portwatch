// Package suppress provides a mechanism to temporarily suppress alerts
// for specific port bindings, preventing repeated notifications during
// known maintenance windows or expected transient activity.
package suppress

import (
	"sync"
	"time"
)

// Entry represents a single suppression rule.
type Entry struct {
	Key       string
	ExpiresAt time.Time
}

// IsExpired reports whether the suppression entry has passed its expiry time.
func (e Entry) IsExpired(now time.Time) bool {
	return now.After(e.ExpiresAt)
}

// Manager holds active suppressions keyed by binding key (e.g. "tcp:8080").
type Manager struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New creates a Manager. If nowFn is nil, time.Now is used.
func New(nowFn func() time.Time) *Manager {
	if nowFn == nil {
		nowFn = time.Now
	}
	return &Manager{
		entries: make(map[string]Entry),
		now:     nowFn,
	}
}

// Suppress adds or refreshes a suppression for the given key lasting duration d.
func (m *Manager) Suppress(key string, d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[key] = Entry{
		Key:       key,
		ExpiresAt: m.now().Add(d),
	}
}

// IsSuppressed reports whether the given key is currently suppressed.
// Expired entries are pruned lazily on read.
func (m *Manager) IsSuppressed(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entries[key]
	if !ok {
		return false
	}
	if e.IsExpired(m.now()) {
		delete(m.entries, key)
		return false
	}
	return true
}

// Remove explicitly lifts a suppression before it expires.
func (m *Manager) Remove(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, key)
}

// Prune removes all expired entries and returns the count removed.
func (m *Manager) Prune() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.now()
	removed := 0
	for k, e := range m.entries {
		if e.IsExpired(now) {
			delete(m.entries, k)
			removed++
		}
	}
	return removed
}

// Len returns the number of active (non-expired) suppressions.
func (m *Manager) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	now := m.now()
	count := 0
	for _, e := range m.entries {
		if !e.IsExpired(now) {
			count++
		}
	}
	return count
}
