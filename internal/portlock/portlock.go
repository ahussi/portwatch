// Package portlock tracks ports that have been explicitly locked (reserved)
// by the operator so that any unexpected process binding to them triggers
// an immediate high-priority alert.
package portlock

import (
	"fmt"
	"sync"
)

// Entry describes a locked port reservation.
type Entry struct {
	Port     int
	Protocol string // "tcp" or "udp"
	Owner    string // expected process name, empty means "any is suspicious"
	Reason   string // human-readable note
}

// Key returns a canonical string key for the entry.
func (e Entry) Key() string {
	return fmt.Sprintf("%s:%d", e.Protocol, e.Port)
}

// Locker holds the set of locked port entries.
type Locker struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an initialised Locker.
func New() *Locker {
	return &Locker{entries: make(map[string]Entry)}
}

// Lock adds or replaces a port lock entry.
func (l *Locker) Lock(e Entry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries[e.Key()] = e
}

// Unlock removes a port lock entry. It is a no-op if the entry does not exist.
func (l *Locker) Unlock(protocol string, port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, fmt.Sprintf("%s:%d", protocol, port))
}

// IsLocked reports whether the given protocol/port combination is locked.
func (l *Locker) IsLocked(protocol string, port int) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	_, ok := l.entries[fmt.Sprintf("%s:%d", protocol, port)]
	return ok
}

// Get returns the Entry for a locked port and whether it was found.
func (l *Locker) Get(protocol string, port int) (Entry, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	e, ok := l.entries[fmt.Sprintf("%s:%d", protocol, port)]
	return e, ok
}

// All returns a snapshot of all locked entries.
func (l *Locker) All() []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Entry, 0, len(l.entries))
	for _, e := range l.entries {
		out = append(out, e)
	}
	return out
}

// Len returns the number of locked entries.
func (l *Locker) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.entries)
}
