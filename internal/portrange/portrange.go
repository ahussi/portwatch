// Package portrange provides utilities for parsing and evaluating
// port range expressions such as "80", "8000-8080", or "443,8443".
package portrange

import (
	"fmt"
	"strconv"
	"strings"
)

// Range represents a contiguous range of ports [Low, High].
type Range struct {
	Low  uint16
	High uint16
}

// Contains reports whether port p falls within the range.
func (r Range) Contains(p uint16) bool {
	return p >= r.Low && p <= r.High
}

// String returns a human-readable representation of the range.
func (r Range) String() string {
	if r.Low == r.High {
		return strconv.Itoa(int(r.Low))
	}
	return fmt.Sprintf("%d-%d", r.Low, r.High)
}

// Set is an ordered collection of Ranges.
type Set struct {
	ranges []Range
}

// Parse parses a comma-separated list of port specs into a Set.
// Each spec may be a single port ("80") or a range ("8000-8080").
func Parse(expr string) (*Set, error) {
	s := &Set{}
	for _, part := range strings.Split(expr, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		r, err := parseRange(part)
		if err != nil {
			return nil, err
		}
		s.ranges = append(s.ranges, r)
	}
	if len(s.ranges) == 0 {
		return nil, fmt.Errorf("portrange: empty expression")
	}
	return s, nil
}

// Contains reports whether port p is covered by any range in the set.
func (s *Set) Contains(p uint16) bool {
	for _, r := range s.ranges {
		if r.Contains(p) {
			return true
		}
	}
	return false
}

// Ranges returns a copy of the underlying range slice.
func (s *Set) Ranges() []Range {
	out := make([]Range, len(s.ranges))
	copy(out, s.ranges)
	return out
}

func parseRange(spec string) (Range, error) {
	parts := strings.SplitN(spec, "-", 2)
	low, err := parsePort(parts[0])
	if err != nil {
		return Range{}, fmt.Errorf("portrange: invalid port %q: %w", parts[0], err)
	}
	if len(parts) == 1 {
		return Range{Low: low, High: low}, nil
	}
	high, err := parsePort(parts[1])
	if err != nil {
		return Range{}, fmt.Errorf("portrange: invalid port %q: %w", parts[1], err)
	}
	if high < low {
		return Range{}, fmt.Errorf("portrange: high port %d < low port %d", high, low)
	}
	return Range{Low: low, High: high}, nil
}

func parsePort(s string) (uint16, error) {
	n, err := strconv.ParseUint(strings.TrimSpace(s), 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(n), nil
}
