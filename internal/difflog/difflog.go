// Package difflog records and retrieves port binding diff events,
// providing a time-ordered log of what changed between scans.
package difflog

import (
	"sync"
	"time"
)

// EventKind describes whether a binding appeared or disappeared.
type EventKind string

const (
	KindAdded   EventKind = "added"
	KindRemoved EventKind = "removed"
)

// Event represents a single diff entry.
type Event struct {
	Kind      EventKind
	Key       string
	Port      int
	Proto     string
	PID       int
	Process   string
	Timestamp time.Time
}

// Log is a capped, thread-safe ordered log of diff events.
type Log struct {
	mu       sync.RWMutex
	events   []Event
	capacity int
}

// New creates a Log with the given maximum capacity.
// If capacity <= 0, it defaults to 256.
func New(capacity int) *Log {
	if capacity <= 0 {
		capacity = 256
	}
	return &Log{capacity: capacity}
}

// Add appends an event to the log, evicting the oldest entry when full.
func (l *Log) Add(e Event) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.events) >= l.capacity {
		l.events = l.events[1:]
	}
	l.events = append(l.events, e)
}

// All returns a shallow copy of all events in insertion order.
func (l *Log) All() []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Event, len(l.events))
	copy(out, l.events)
	return out
}

// Since returns events that occurred at or after t.
func (l *Log) Since(t time.Time) []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var out []Event
	for _, e := range l.events {
		if !e.Timestamp.Before(t) {
			out = append(out, e)
		}
	}
	return out
}

// Len returns the current number of stored events.
func (l *Log) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.events)
}

// Clear removes all events from the log.
func (l *Log) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = l.events[:0]
}
