// Package metrics tracks runtime counters for portwatch daemon activity.
package metrics

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Counters holds atomic counters for key daemon events.
type Counters struct {
	Scans       uint64
	NewBindings uint64
	Alerts      uint64
	Suppressed  uint64
	Errors      uint64
}

// Metrics collects and exposes runtime statistics for the daemon.
type Metrics struct {
	mu        sync.RWMutex
	counters  Counters
	startedAt time.Time
}

// New creates a new Metrics instance, recording the current time as start.
func New() *Metrics {
	return &Metrics{startedAt: time.Now()}
}

// IncScans increments the scan counter by 1.
func (m *Metrics) IncScans() {
	atomic.AddUint64(&m.counters.Scans, 1)
}

// IncNewBindings increments the new-bindings counter by 1.
func (m *Metrics) IncNewBindings() {
	atomic.AddUint64(&m.counters.NewBindings, 1)
}

// IncAlerts increments the alerts-dispatched counter by 1.
func (m *Metrics) IncAlerts() {
	atomic.AddUint64(&m.counters.Alerts, 1)
}

// IncSuppressed increments the suppressed-alerts counter by 1.
func (m *Metrics) IncSuppressed() {
	atomic.AddUint64(&m.counters.Suppressed, 1)
}

// IncErrors increments the error counter by 1.
func (m *Metrics) IncErrors() {
	atomic.AddUint64(&m.counters.Errors, 1)
}

// Snapshot returns a point-in-time copy of the current counters.
func (m *Metrics) Snapshot() Counters {
	return Counters{
		Scans:       atomic.LoadUint64(&m.counters.Scans),
		NewBindings: atomic.LoadUint64(&m.counters.NewBindings),
		Alerts:      atomic.LoadUint64(&m.counters.Alerts),
		Suppressed:  atomic.LoadUint64(&m.counters.Suppressed),
		Errors:      atomic.LoadUint64(&m.counters.Errors),
	}
}

// Uptime returns the duration since the Metrics instance was created.
func (m *Metrics) Uptime() time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return time.Since(m.startedAt)
}

// String returns a human-readable summary of current metrics.
func (m *Metrics) String() string {
	c := m.Snapshot()
	return fmt.Sprintf(
		"uptime=%s scans=%d new_bindings=%d alerts=%d suppressed=%d errors=%d",
		m.Uptime().Round(time.Second),
		c.Scans, c.NewBindings, c.Alerts, c.Suppressed, c.Errors,
	)
}
