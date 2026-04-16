// Package portquota enforces per-process or per-user port binding quotas.
package portquota

import (
	"fmt"
	"sync"
)

// ErrQuotaExceeded is returned when a subject has reached its port limit.
type ErrQuotaExceeded struct {
	Subject string
	Limit   int
}

func (e *ErrQuotaExceeded) Error() string {
	return fmt.Sprintf("quota exceeded for %q: limit is %d", e.Subject, e.Limit)
}

// Quota tracks port binding counts per subject (e.g. process name or user).
type Quota struct {
	mu      sync.Mutex
	limits  map[string]int
	counts  map[string]map[uint16]struct{}
	default_ int
}

// New creates a Quota with a default per-subject limit.
func New(defaultLimit int) *Quota {
	if defaultLimit < 1 {
		defaultLimit = 1
	}
	return &Quota{
		limits:   make(map[string]int),
		counts:   make(map[string]map[uint16]struct{}),
		default_: defaultLimit,
	}
}

// SetLimit overrides the limit for a specific subject.
func (q *Quota) SetLimit(subject string, limit int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if limit < 1 {
		limit = 1
	}
	q.limits[subject] = limit
}

// Track records a port binding for subject. Returns ErrQuotaExceeded if the
// limit would be breached. Duplicate port registrations are idempotent.
func (q *Quota) Track(subject string, port uint16) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	ports, ok := q.counts[subject]
	if !ok {
		ports = make(map[uint16]struct{})
		q.counts[subject] = ports
	}
	if _, exists := ports[port]; exists {
		return nil
	}
	limit := q.default_
	if l, ok := q.limits[subject]; ok {
		limit = l
	}
	if len(ports) >= limit {
		return &ErrQuotaExceeded{Subject: subject, Limit: limit}
	}
	ports[port] = struct{}{}
	return nil
}

// Release removes a port binding for subject.
func (q *Quota) Release(subject string, port uint16) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if ports, ok := q.counts[subject]; ok {
		delete(ports, port)
	}
}

// Count returns the current number of tracked ports for subject.
func (q *Quota) Count(subject string) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.counts[subject])
}
