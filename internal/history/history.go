// Package history records port binding events over time,
// enabling trend analysis and audit trails.
package history

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// EventKind describes whether a binding appeared or disappeared.
type EventKind string

const (
	EventAdded   EventKind = "added"
	EventRemoved EventKind = "removed"
)

// Event represents a single port-binding change captured at a point in time.
type Event struct {
	Timestamp time.Time
	Kind      EventKind
	Binding   scanner.Binding
}

// Record is a fixed-size circular buffer of Events.
type Record struct {
	mu     sync.RWMutex
	events []Event
	cap    int
	head   int
	size   int
}

// New creates a Record that retains at most capacity events.
func New(capacity int) *Record {
	if capacity <= 0 {
		capacity = 256
	}
	return &Record{
		events: make([]Event, capacity),
		cap:    capacity,
	}
}

// Add appends an event, overwriting the oldest entry when full.
func (r *Record) Add(kind EventKind, b scanner.Binding) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events[r.head] = Event{
		Timestamp: time.Now(),
		Kind:      kind,
		Binding:   b,
	}
	r.head = (r.head + 1) % r.cap
	if r.size < r.cap {
		r.size++
	}
}

// All returns a chronologically ordered snapshot of all stored events.
func (r *Record) All() []Event {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Event, r.size)
	start := (r.head - r.size + r.cap) % r.cap
	for i := 0; i < r.size; i++ {
		out[i] = r.events[(start+i)%r.cap]
	}
	return out
}

// Len returns the number of events currently stored.
func (r *Record) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.size
}
