// Package portlabel provides human-readable label assignment for monitored ports.
// Labels can be user-defined or derived from well-known service names.
package portlabel

import (
	"fmt"
	"sync"
)

// Label holds a display name and an optional category tag for a port.
type Label struct {
	Name     string
	Category string
}

// String returns a formatted representation of the label.
func (l Label) String() string {
	if l.Category != "" {
		return fmt.Sprintf("%s (%s)", l.Name, l.Category)
	}
	return l.Name
}

// Labeler assigns and retrieves labels for port/protocol pairs.
type Labeler struct {
	mu     sync.RWMutex
	labels map[string]Label
}

func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// New returns a Labeler pre-populated with a small set of well-known labels.
func New() *Labeler {
	l := &Labeler{labels: make(map[string]Label)}
	for _, e := range builtinLabels {
		l.labels[key(e.port, e.proto)] = Label{Name: e.name, Category: e.category}
	}
	return l
}

// Set assigns a label to the given port/protocol pair, overwriting any existing entry.
func (l *Labeler) Set(port int, proto, name, category string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.labels[key(port, proto)] = Label{Name: name, Category: category}
}

// Get returns the label for the given port/protocol pair and whether it was found.
func (l *Labeler) Get(port int, proto string) (Label, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	lbl, ok := l.labels[key(port, proto)]
	return lbl, ok
}

// Remove deletes the label for the given port/protocol pair.
func (l *Labeler) Remove(port int, proto string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.labels, key(port, proto))
}

// All returns a copy of all registered labels keyed by "port/proto".
func (l *Labeler) All() map[string]Label {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make(map[string]Label, len(l.labels))
	for k, v := range l.labels {
		out[k] = v
	}
	return out
}

type builtin struct {
	port     int
	proto    string
	name     string
	category string
}

var builtinLabels = []builtin{
	{22, "tcp", "SSH", "remote-access"},
	{80, "tcp", "HTTP", "web"},
	{443, "tcp", "HTTPS", "web"},
	{3306, "tcp", "MySQL", "database"},
	{5432, "tcp", "PostgreSQL", "database"},
	{6379, "tcp", "Redis", "cache"},
	{8080, "tcp", "HTTP-Alt", "web"},
}
