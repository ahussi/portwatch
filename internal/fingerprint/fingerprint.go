// Package fingerprint generates stable identifiers for port bindings,
// enabling reliable deduplication across scan cycles.
package fingerprint

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// Binding represents the minimal information needed to fingerprint a port binding.
type Binding struct {
	Protocol string
	Address  string
	Port     int
	PID      int
	Process  string
}

// Fingerprinter generates and compares binding fingerprints.
type Fingerprinter struct {
	includePID bool
}

// Option configures a Fingerprinter.
type Option func(*Fingerprinter)

// WithPID includes the PID in the fingerprint, making it sensitive to process restarts.
func WithPID() Option {
	return func(f *Fingerprinter) {
		f.includePID = true
	}
}

// New returns a new Fingerprinter with the given options.
func New(opts ...Option) *Fingerprinter {
	f := &Fingerprinter{}
	for _, o := range opts {
		o(f)
	}
	return f
}

// Generate returns a hex fingerprint string for the given binding.
// By default the fingerprint is stable across PID changes (process restarts
// on the same port/address are treated as the same binding).
func (f *Fingerprinter) Generate(b Binding) string {
	var sb strings.Builder
	sb.WriteString(strings.ToLower(b.Protocol))
	sb.WriteByte('|')
	sb.WriteString(b.Address)
	sb.WriteByte('|')
	sb.WriteString(fmt.Sprintf("%d", b.Port))
	if f.includePID {
		sb.WriteByte('|')
		sb.WriteString(fmt.Sprintf("%d", b.PID))
	}
	sum := sha256.Sum256([]byte(sb.String()))
	return fmt.Sprintf("%x", sum[:8])
}

// Equal reports whether two bindings produce the same fingerprint.
func (f *Fingerprinter) Equal(a, b Binding) bool {
	return f.Generate(a) == f.Generate(b)
}
