// Package healthcheck provides a simple liveness probe for the portwatch daemon.
// It tracks whether the daemon's core scan loop is making progress and exposes
// a status that can be queried by external tools or a future HTTP endpoint.
package healthcheck

import (
	"fmt"
	"sync"
	"time"
)

// Status represents the health state of the daemon.
type Status string

const (
	StatusOK      Status = "ok"
	StatusStale   Status = "stale"
	StatusUnknown Status = "unknown"

	// defaultStaleness is how long without a heartbeat before we report stale.
	defaultStaleness = 30 * time.Second
)

// Checker tracks heartbeats from the scan loop and reports health status.
type Checker struct {
	mu        sync.RWMutex
	lastBeat  time.Time
	staleness time.Duration
	started   bool
}

// New returns a new Checker with the given staleness threshold.
// If staleness is zero, defaultStaleness is used.
func New(staleness time.Duration) *Checker {
	if staleness <= 0 {
		staleness = defaultStaleness
	}
	return &Checker{staleness: staleness}
}

// Beat records a heartbeat from the scan loop. Should be called after each
// successful scan cycle.
func (c *Checker) Beat() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastBeat = time.Now()
	c.started = true
}

// Status returns the current health status of the daemon.
func (c *Checker) Status() Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !c.started {
		return StatusUnknown
	}
	if time.Since(c.lastBeat) > c.staleness {
		return StatusStale
	}
	return StatusOK
}

// LastBeat returns the time of the most recent heartbeat.
func (c *Checker) LastBeat() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastBeat
}

// String returns a human-readable summary of the health state.
func (c *Checker) String() string {
	status := c.Status()
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !c.started {
		return fmt.Sprintf("healthcheck: %s (no heartbeat recorded)", status)
	}
	return fmt.Sprintf("healthcheck: %s (last beat %s ago)", status, time.Since(c.lastBeat).Round(time.Millisecond))
}
