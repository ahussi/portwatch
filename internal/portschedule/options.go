package portschedule

import "time"

// Option is a functional option for Scheduler.
type Option func(*Scheduler)

// WithClock replaces the wall-clock source used by Active.
// Useful for deterministic testing.
func WithClock(fn func() time.Time) Option {
	return func(s *Scheduler) {
		s.now = fn
	}
}
