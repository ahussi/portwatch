// Package watchlist manages the set of ports actively being monitored,
// supporting dynamic add/remove and iteration.
package watchlist

import (
	"fmt"
	"sync"
)

// Entry describes a single watched port entry.
type Entry struct {
	Port     int
	Protocol string // "tcp" or "udp"
	Label    string // optional human-readable label
}

// Key returns a unique string key for the entry.
func (e Entry) Key() string {
	return fmt.Sprintf("%s:%d", e.Protocol, e.Port)
}

// Watchlist holds the set of ports under active observation.
type Watchlist struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an empty Watchlist.
func New() *Watchlist {
	return &Watchlist{
		entries: make(map[string]Entry),
	}
}

// Add inserts or replaces an entry in the watchlist.
func (w *Watchlist) Add(e Entry) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries[e.Key()] = e
}

// Remove deletes the entry identified by protocol and port.
func (w *Watchlist) Remove(protocol string, port int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	key := fmt.Sprintf("%s:%d", protocol, port)
	delete(w.entries, key)
}

// Has reports whether the given protocol/port combination is watched.
func (w *Watchlist) Has(protocol string, port int) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	key := fmt.Sprintf("%s:%d", protocol, port)
	_, ok := w.entries[key]
	return ok
}

// All returns a snapshot of all current entries.
func (w *Watchlist) All() []Entry {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]Entry, 0, len(w.entries))
	for _, e := range w.entries {
		out = append(out, e)
	}
	return out
}

// Len returns the number of watched entries.
func (w *Watchlist) Len() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.entries)
}
