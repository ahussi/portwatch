// Package portmap provides a registry that maps well-known port numbers to
// human-readable service names and protocol information.
package portmap

import "fmt"

// Entry describes a single port registration.
type Entry struct {
	Port     int
	Protocol string // "tcp" | "udp"
	Service  string
	Desc     string
}

// String returns a compact representation of the entry.
func (e Entry) String() string {
	return fmt.Sprintf("%d/%s (%s)", e.Port, e.Protocol, e.Service)
}

// Registry holds port-to-entry mappings.
type Registry struct {
	entries map[string]Entry // key: "<port>/<proto>"
}

// key builds the lookup key used internally.
func key(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// New returns a Registry pre-loaded with common well-known ports.
func New() *Registry {
	r := &Registry{entries: make(map[string]Entry)}
	for _, e := range builtins {
		r.entries[key(e.Port, e.Protocol)] = e
	}
	return r
}

// Lookup returns the Entry for the given port/protocol pair and whether it
// was found.
func (r *Registry) Lookup(port int, proto string) (Entry, bool) {
	e, ok := r.entries[key(port, proto)]
	return e, ok
}

// Register adds or overwrites an entry in the registry.
func (r *Registry) Register(e Entry) {
	r.entries[key(e.Port, e.Protocol)] = e
}

// All returns a copy of every entry in the registry.
func (r *Registry) All() []Entry {
	out := make([]Entry, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, e)
	}
	return out
}

// builtins is the seed data for the registry.
var builtins = []Entry{
	{22, "tcp", "ssh", "Secure Shell"},
	{25, "tcp", "smtp", "Simple Mail Transfer Protocol"},
	{53, "tcp", "dns", "Domain Name System"},
	{53, "udp", "dns", "Domain Name System"},
	{80, "tcp", "http", "Hypertext Transfer Protocol"},
	{110, "tcp", "pop3", "Post Office Protocol v3"},
	{143, "tcp", "imap", "Internet Message Access Protocol"},
	{443, "tcp", "https", "HTTP Secure"},
	{3306, "tcp", "mysql", "MySQL Database"},
	{5432, "tcp", "postgresql", "PostgreSQL Database"},
	{6379, "tcp", "redis", "Redis In-Memory Store"},
	{8080, "tcp", "http-alt", "HTTP Alternate"},
	{27017, "tcp", "mongodb", "MongoDB Database"},
}
