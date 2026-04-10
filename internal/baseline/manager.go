package baseline

import (
	"errors"
	"os"
)

// Manager wraps a Baseline and exposes higher-level operations used by the
// watcher to decide whether a newly observed binding is expected.
type Manager struct {
	b *Baseline
}

// NewManager returns a Manager. If the baseline file exists it is loaded
// automatically; a missing file is silently ignored.
func NewManager(path string) (*Manager, error) {
	b := New(path)
	if err := b.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return &Manager{b: b}, nil
}

// IsKnown reports whether the given binding key is part of the baseline.
func (m *Manager) IsKnown(key string) bool {
	return m.b.Has(key)
}

// Record adds the key/entry pair to the in-memory baseline without persisting.
func (m *Manager) Record(key string, e Entry) {
	m.b.Set(key, e)
}

// Commit persists the current in-memory baseline to disk.
func (m *Manager) Commit() error {
	return m.b.Save()
}

// Snapshot returns a copy of all baseline entries.
func (m *Manager) Snapshot() map[string]Entry {
	return m.b.Entries()
}
