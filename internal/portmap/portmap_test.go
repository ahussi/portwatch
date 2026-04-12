package portmap

import (
	"strings"
	"testing"
)

func TestNew_ContainsBuiltins(t *testing.T) {
	r := New()
	if len(r.All()) == 0 {
		t.Fatal("expected built-in entries, got none")
	}
}

func TestLookup_KnownPort(t *testing.T) {
	r := New()
	e, ok := r.Lookup(443, "tcp")
	if !ok {
		t.Fatal("expected 443/tcp to be known")
	}
	if e.Service != "https" {
		t.Errorf("expected service 'https', got %q", e.Service)
	}
}

func TestLookup_UnknownPort(t *testing.T) {
	r := New()
	_, ok := r.Lookup(9999, "tcp")
	if ok {
		t.Error("expected 9999/tcp to be unknown")
	}
}

func TestLookup_ProtocolDistinct(t *testing.T) {
	r := New()
	// DNS is registered for both tcp and udp
	_, tcpOK := r.Lookup(53, "tcp")
	_, udpOK := r.Lookup(53, "udp")
	if !tcpOK || !udpOK {
		t.Errorf("expected 53/tcp and 53/udp both registered, got tcp=%v udp=%v", tcpOK, udpOK)
	}
}

func TestRegister_AddsEntry(t *testing.T) {
	r := New()
	custom := Entry{Port: 9000, Protocol: "tcp", Service: "custom-svc", Desc: "test service"}
	r.Register(custom)
	e, ok := r.Lookup(9000, "tcp")
	if !ok {
		t.Fatal("expected registered entry to be found")
	}
	if e.Service != "custom-svc" {
		t.Errorf("unexpected service name: %q", e.Service)
	}
}

func TestRegister_Overwrites(t *testing.T) {
	r := New()
	r.Register(Entry{Port: 80, Protocol: "tcp", Service: "my-http", Desc: "override"})
	e, _ := r.Lookup(80, "tcp")
	if e.Service != "my-http" {
		t.Errorf("expected overwritten service, got %q", e.Service)
	}
}

func TestEntryString(t *testing.T) {
	e := Entry{Port: 443, Protocol: "tcp", Service: "https", Desc: "HTTP Secure"}
	s := e.String()
	if !strings.Contains(s, "443/tcp") || !strings.Contains(s, "https") {
		t.Errorf("unexpected String() output: %q", s)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	r := New()
	all := r.All()
	initialLen := len(all)
	all = append(all, Entry{Port: 1, Protocol: "tcp", Service: "test"})
	if len(r.All()) != initialLen {
		t.Error("All() should return an independent copy")
	}
}
