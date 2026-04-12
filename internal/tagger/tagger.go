// Package tagger assigns human-readable tags to port bindings based on
// well-known service names, user-defined labels, and protocol hints.
package tagger

import (
	"fmt"
	"sync"
)

// Tag represents a label attached to a port binding.
type Tag struct {
	Label  string
	Source string // "builtin", "user", or "inferred"
}

func (t Tag) String() string {
	return fmt.Sprintf("%s(%s)", t.Label, t.Source)
}

// Tagger maps ports to tags.
type Tagger struct {
	mu      sync.RWMutex
	tags    map[uint16][]Tag
	builtIn map[uint16]string
}

// New returns a Tagger pre-loaded with built-in well-known port labels.
func New() *Tagger {
	return &Tagger{
		tags: make(map[uint16][]Tag),
		builtIn: builtinPorts(),
	}
}

// Get returns all tags for the given port, including built-in ones.
func (t *Tagger) Get(port uint16) []Tag {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var result []Tag
	if label, ok := t.builtIn[port]; ok {
		result = append(result, Tag{Label: label, Source: "builtin"})
	}
	result = append(result, t.tags[port]...)
	return result
}

// Add attaches a user-defined tag to a port.
func (t *Tagger) Add(port uint16, label string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tags[port] = append(t.tags[port], Tag{Label: label, Source: "user"})
}

// Remove deletes all user-defined tags for the given port.
func (t *Tagger) Remove(port uint16) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.tags, port)
}

// HasTag reports whether the given port carries the specified label from any source.
func (t *Tagger) HasTag(port uint16, label string) bool {
	for _, tag := range t.Get(port) {
		if tag.Label == label {
			return true
		}
	}
	return false
}

// builtinPorts returns a map of well-known port numbers to service names.
func builtinPorts() map[uint16]string {
	return map[uint16]string{
		22:   "ssh",
		25:   "smtp",
		53:   "dns",
		80:   "http",
		443:  "https",
		3306: "mysql",
		5432: "postgres",
		6379: "redis",
		8080: "http-alt",
		8443: "https-alt",
		9200: "elasticsearch",
		27017: "mongodb",
	}
}
