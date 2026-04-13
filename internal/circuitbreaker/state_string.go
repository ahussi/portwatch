package circuitbreaker

import "fmt"

// String returns a human-readable label for the circuit state.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return fmt.Sprintf("unknown(%d)", int(s))
	}
}

// IsHealthy returns true when the circuit is in the Closed state,
// meaning downstream calls are expected to succeed.
func (s State) IsHealthy() bool {
	return s == StateClosed
}
