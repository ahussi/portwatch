package portpolicy

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func makeConfig(watched, allowed, denied []string) *config.Config {
	cfg := config.Default()
	cfg.WatchedPorts = watched
	cfg.AllowedPorts = allowed
	cfg.DeniedPorts = denied
	return cfg
}

func TestFromConfig_NilConfig(t *testing.T) {
	_, err := FromConfig(nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestFromConfig_DefaultAllowWhenNoWatchList(t *testing.T) {
	cfg := makeConfig(nil, nil, nil)
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := p.Evaluate(9999, "tcp"); got != Allow {
		t.Errorf("expected Allow default, got %s", got)
	}
}

func TestFromConfig_DefaultDenyWhenWatchList(t *testing.T) {
	cfg := makeConfig([]string{"80"}, nil, nil)
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := p.Evaluate(9999, "tcp"); got != Deny {
		t.Errorf("expected Deny default, got %s", got)
	}
}

func TestFromConfig_AllowedPort(t *testing.T) {
	cfg := makeConfig([]string{"80"}, []string{"443/tcp"}, nil)
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := p.Evaluate(443, "tcp"); got != Allow {
		t.Errorf("expected Allow for 443/tcp, got %s", got)
	}
}

func TestFromConfig_DeniedPort(t *testing.T) {
	cfg := makeConfig(nil, nil, []string{"22/tcp"})
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := p.Evaluate(22, "tcp"); got != Deny {
		t.Errorf("expected Deny for 22/tcp, got %s", got)
	}
}

func TestFromConfig_InvalidAllowedPort(t *testing.T) {
	cfg := makeConfig(nil, []string{"notaport"}, nil)
	_, err := FromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid port spec")
	}
}

func TestParseSpec_ValidFormats(t *testing.T) {
	cases := []struct {
		spec  string
		port  int
		proto string
	}{
		{"80", 80, ""},
		{"80/tcp", 80, "tcp"},
		{"53/udp", 53, "udp"},
	}
	for _, c := range cases {
		port, proto, err := parseSpec(c.spec)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", c.spec, err)
			continue
		}
		if port != c.port || proto != c.proto {
			t.Errorf("%q: got (%d,%q), want (%d,%q)", c.spec, port, proto, c.port, c.proto)
		}
	}
}
