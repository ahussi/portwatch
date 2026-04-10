// Package baseline persists a known-good snapshot of port bindings to disk
// so that portwatch can detect deviations across restarts.
package baseline

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single persisted port binding.
type Entry struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
	PID      int    `json:"pid"`
	Process  string `json:"process"`
}

// Baseline holds the persisted set of expected bindings.
type Baseline struct {
	mu      sync.RWMutex
	entries map[string]Entry
	path    string
	SavedAt time.Time `json:"saved_at"`
}

// New creates a new Baseline backed by the given file path.
func New(path string) *Baseline {
	return &Baseline{
		entries: make(map[string]Entry),
		path:    path,
	}
}

// Set adds or replaces an entry keyed by "proto:addr:port".
func (b *Baseline) Set(key string, e Entry) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries[key] = e
}

// Has reports whether the key exists in the baseline.
func (b *Baseline) Has(key string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.entries[key]
	return ok
}

// Entries returns a shallow copy of all entries.
func (b *Baseline) Entries() map[string]Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	copy := make(map[string]Entry, len(b.entries))
	for k, v := range b.entries {
		copy[k] = v
	}
	return copy
}

// Save writes the baseline to disk as JSON.
func (b *Baseline) Save() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.SavedAt = time.Now().UTC()
	payload := struct {
		SavedAt time.Time        `json:"saved_at"`
		Entries map[string]Entry `json:"entries"`
	}{SavedAt: b.SavedAt, Entries: b.entries}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o644)
}

// Load reads a previously saved baseline from disk.
func (b *Baseline) Load() error {
	data, err := os.ReadFile(b.path)
	if err != nil {
		return err
	}
	var payload struct {
		SavedAt time.Time        `json:"saved_at"`
		Entries map[string]Entry `json:"entries"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.SavedAt = payload.SavedAt
	b.entries = payload.Entries
	return nil
}
