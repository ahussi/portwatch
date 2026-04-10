package snapshot

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Binding represents a captured port binding at a point in time.
type Binding struct {
	Port     int
	Protocol string
	Address  string
	Process  *scanner.ProcessInfo
	SeenAt   time.Time
}

// Snapshot holds the last known state of port bindings.
type Snapshot struct {
	mu       sync.RWMutex
	bindings map[string]Binding
}

// New returns an empty Snapshot.
func New() *Snapshot {
	return &Snapshot{
		bindings: make(map[string]Binding),
	}
}

// Set stores or updates a binding by its key.
func (s *Snapshot) Set(key string, b Binding) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b.SeenAt = time.Now()
	s.bindings[key] = b
}

// Get retrieves a binding by key. Returns false if not present.
func (s *Snapshot) Get(key string) (Binding, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.bindings[key]
	return b, ok
}

// Delete removes a binding by key.
func (s *Snapshot) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.bindings, key)
}

// Keys returns all current binding keys.
func (s *Snapshot) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.bindings))
	for k := range s.bindings {
		keys = append(keys, k)
	}
	return keys
}

// Diff compares the snapshot against a new set of keys and returns
// keys that are added (present in newKeys but not snapshot) and
// removed (present in snapshot but not newKeys).
func (s *Snapshot) Diff(newKeys []string) (added, removed []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	newSet := make(map[string]struct{}, len(newKeys))
	for _, k := range newKeys {
		newSet[k] = struct{}{}
		if _, exists := s.bindings[k]; !exists {
			added = append(added, k)
		}
	}
	for k := range s.bindings {
		if _, exists := newSet[k]; !exists {
			removed = append(removed, k)
		}
	}
	return added, removed
}
