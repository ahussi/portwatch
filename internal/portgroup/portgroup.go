// Package portgroup provides named grouping of ports for collective
// monitoring and alerting within portwatch.
package portgroup

import (
	"fmt"
	"sync"
)

// Group holds a named collection of port numbers.
type Group struct {
	Name  string
	Ports []int
}

// Registry maps group names to their port sets.
type Registry struct {
	mu     sync.RWMutex
	groups map[string]*Group
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{
		groups: make(map[string]*Group),
	}
}

// Add registers a named group with the given ports. Returns an error if the
// group name is already registered.
func (r *Registry) Add(name string, ports []int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.groups[name]; exists {
		return fmt.Errorf("portgroup: group %q already registered", name)
	}
	copied := make([]int, len(ports))
	copy(copied, ports)
	r.groups[name] = &Group{Name: name, Ports: copied}
	return nil
}

// Get returns the Group for the given name and whether it was found.
func (r *Registry) Get(name string) (*Group, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.groups[name]
	return g, ok
}

// Remove deletes a group by name. No-op if the group does not exist.
func (r *Registry) Remove(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.groups, name)
}

// Contains reports whether the given port belongs to the named group.
// Returns false if the group does not exist.
func (r *Registry) Contains(name string, port int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.groups[name]
	if !ok {
		return false
	}
	for _, p := range g.Ports {
		if p == port {
			return true
		}
	}
	return false
}

// All returns a snapshot of all registered groups.
func (r *Registry) All() []*Group {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Group, 0, len(r.groups))
	for _, g := range r.groups {
		copy := &Group{Name: g.Name, Ports: append([]int(nil), g.Ports...)}
		out = append(out, copy)
	}
	return out
}
