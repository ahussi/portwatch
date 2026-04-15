package portpolicy

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/config"
)

// FromConfig builds a Policy from the application Config.
// Ports listed in config.AllowedPorts receive an Allow rule;
// ports listed in config.DeniedPorts receive a Deny rule.
// The default action is Allow when a watch-list is configured,
// otherwise Deny.
func FromConfig(cfg *config.Config) (*Policy, error) {
	if cfg == nil {
		return nil, fmt.Errorf("portpolicy: nil config")
	}

	defaultAction := Allow
	if len(cfg.WatchedPorts) > 0 {
		defaultAction = Deny
	}

	p := New(defaultAction)

	for _, spec := range cfg.AllowedPorts {
		port, proto, err := parseSpec(spec)
		if err != nil {
			return nil, fmt.Errorf("portpolicy: allowed port %q: %w", spec, err)
		}
		p.Add(Rule{Port: port, Protocol: proto, Action: Allow})
	}

	for _, spec := range cfg.DeniedPorts {
		port, proto, err := parseSpec(spec)
		if err != nil {
			return nil, fmt.Errorf("portpolicy: denied port %q: %w", spec, err)
		}
		p.Add(Rule{Port: port, Protocol: proto, Action: Deny})
	}

	return p, nil
}

// parseSpec parses "80", "80/tcp", or "80/udp".
func parseSpec(spec string) (int, string, error) {
	parts := strings.SplitN(spec, "/", 2)
	var port int
	if _, err := fmt.Sscanf(parts[0], "%d", &port); err != nil || port < 1 || port > 65535 {
		return 0, "", fmt.Errorf("invalid port %q", parts[0])
	}
	proto := ""
	if len(parts) == 2 {
		proto = strings.ToLower(parts[1])
		if proto != "tcp" && proto != "udp" {
			return 0, "", fmt.Errorf("unknown protocol %q", proto)
		}
	}
	return port, proto, nil
}
