package watchlist

import (
	"fmt"

	"github.com/user/portwatch/internal/config"
)

// FromConfig builds a Watchlist from the application configuration.
// Ports listed under config.Ports are added with protocol defaulting to "tcp"
// unless explicitly prefixed (e.g. "udp:5353").
func FromConfig(cfg *config.Config) (*Watchlist, error) {
	if cfg == nil {
		return nil, fmt.Errorf("watchlist: nil config")
	}
	wl := New()
	for _, p := range cfg.Ports {
		proto, port, err := parsePortSpec(p)
		if err != nil {
			return nil, fmt.Errorf("watchlist: invalid port spec %q: %w", p, err)
		}
		wl.Add(Entry{
			Port:     port,
			Protocol: proto,
			Label:    fmt.Sprintf("%s:%d", proto, port),
		})
	}
	return wl, nil
}

// parsePortSpec parses a port specification of the form "[proto:]port".
// If no protocol prefix is present, "tcp" is assumed.
func parsePortSpec(spec string) (proto string, port int, err error) {
	proto = "tcp"
	var n int
	// try "proto:port" form first
	if _, scanErr := fmt.Sscanf(spec, "%5[a-z]:%d", &proto, &n); scanErr == nil && n > 0 {
		if proto != "tcp" && proto != "udp" {
			return "", 0, fmt.Errorf("unknown protocol %q", proto)
		}
		return proto, n, nil
	}
	// plain port number
	if _, scanErr := fmt.Sscanf(spec, "%d", &n); scanErr == nil && n > 0 {
		return "tcp", n, nil
	}
	return "", 0, fmt.Errorf("cannot parse %q as port spec", spec)
}
