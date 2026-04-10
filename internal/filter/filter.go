// Package filter provides port filtering utilities for portwatch.
// It allows callers to determine whether a given port or binding
// should be included or excluded based on configured rules.
package filter

import "github.com/yourorg/portwatch/internal/config"

// Result describes why a port was accepted or rejected.
type Result int

const (
	Allowed Result = iota
	DeniedNotWatched
	DeniedExplicit
)

// String returns a human-readable label for the Result.
func (r Result) String() string {
	switch r {
	case Allowed:
		return "allowed"
	case DeniedNotWatched:
		return "denied:not-watched"
	case DeniedExplicit:
		return "denied:explicit"
	default:
		return "unknown"
	}
}

// Filter evaluates ports against a Config and returns filtering decisions.
type Filter struct {
	cfg *config.Config
}

// New creates a Filter backed by the provided Config.
func New(cfg *config.Config) *Filter {
	return &Filter{cfg: cfg}
}

// Check returns the Result for a given port number.
// A port is Allowed when it is in the watched list and not in the
// allowed (ignored) list. Ports in the allowed list are silently
// accepted without raising alerts (DeniedExplicit). Ports outside
// the watched list are DeniedNotWatched.
func (f *Filter) Check(port int) Result {
	if !f.cfg.IsWatched(port) {
		return DeniedNotWatched
	}
	if f.cfg.IsAllowed(port) {
		return DeniedExplicit
	}
	return Allowed
}

// ShouldAlert returns true when the port should trigger an alert.
func (f *Filter) ShouldAlert(port int) bool {
	return f.Check(port) == Allowed
}
