// Package portschedule provides time-based scheduling rules for port monitoring.
// Rules can restrict when a port binding is considered expected or unexpected
// based on time-of-day windows (e.g. allow :8080 only during business hours).
package portschedule

import (
	"fmt"
	"sync"
	"time"
)

// Window defines a daily time range [Start, End) in which a rule is active.
type Window struct {
	Start time.Duration // offset from midnight
	End   time.Duration // offset from midnight
}

// Rule associates a port with one or more active windows.
type Rule struct {
	Port    int
	Windows []Window
}

// Scheduler holds scheduling rules and evaluates them against wall-clock time.
type Scheduler struct {
	mu    sync.RWMutex
	rules map[int][]Window
	now   func() time.Time
}

// New returns an empty Scheduler. Pass functional options to customise behaviour.
func New(opts ...Option) *Scheduler {
	s := &Scheduler{
		rules: make(map[int][]Window),
		now:   time.Now,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Add registers a scheduling rule, replacing any previous rule for that port.
func (s *Scheduler) Add(r Rule) error {
	for _, w := range r.Windows {
		if w.Start >= w.End {
			return fmt.Errorf("portschedule: window start %v must be before end %v", w.Start, w.End)
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules[r.Port] = r.Windows
	return nil
}

// Remove deletes the scheduling rule for the given port.
func (s *Scheduler) Remove(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.rules, port)
}

// Active reports whether port is within any of its scheduled windows right now.
// If no rule exists for the port, Active returns true (unscheduled = always active).
func (s *Scheduler) Active(port int) bool {
	s.mu.RLock()
	windows, ok := s.rules[port]
	s.mu.RUnlock()
	if !ok {
		return true
	}
	now := s.now()
	offset := time.Duration(now.Hour())*time.Hour +
		time.Duration(now.Minute())*time.Minute +
		time.Duration(now.Second())*time.Second
	for _, w := range windows {
		if offset >= w.Start && offset < w.End {
			return true
		}
	}
	return false
}

// Rules returns a snapshot of all registered rules.
func (s *Scheduler) Rules() []Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Rule, 0, len(s.rules))
	for port, windows := range s.rules {
		out = append(out, Rule{Port: port, Windows: append([]Window(nil), windows...)})
	}
	return out
}
