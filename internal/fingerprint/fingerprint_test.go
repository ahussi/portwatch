package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
)

func binding(proto, addr string, port, pid int, proc string) fingerprint.Binding {
	return fingerprint.Binding{Protocol: proto, Address: addr, Port: port, PID: pid, Process: proc}
}

func TestNew(t *testing.T) {
	f := fingerprint.New()
	if f == nil {
		t.Fatal("expected non-nil Fingerprinter")
	}
}

func TestGenerate_Deterministic(t *testing.T) {
	f := fingerprint.New()
	b := binding("tcp", "0.0.0.0", 8080, 1234, "nginx")
	if fingerprint.New().Generate(b) != f.Generate(b) {
		t.Error("Generate should be deterministic")
	}
}

func TestGenerate_DifferentPorts(t *testing.T) {
	f := fingerprint.New()
	a := binding("tcp", "0.0.0.0", 8080, 1, "app")
	b := binding("tcp", "0.0.0.0", 9090, 1, "app")
	if f.Generate(a) == f.Generate(b) {
		t.Error("different ports should produce different fingerprints")
	}
}

func TestGenerate_PIDIgnoredByDefault(t *testing.T) {
	f := fingerprint.New()
	a := binding("tcp", "127.0.0.1", 3000, 100, "node")
	b := binding("tcp", "127.0.0.1", 3000, 999, "node")
	if f.Generate(a) != f.Generate(b) {
		t.Error("PID should be ignored without WithPID option")
	}
}

func TestGenerate_WithPID(t *testing.T) {
	f := fingerprint.New(fingerprint.WithPID())
	a := binding("tcp", "127.0.0.1", 3000, 100, "node")
	b := binding("tcp", "127.0.0.1", 3000, 999, "node")
	if f.Generate(a) == f.Generate(b) {
		t.Error("WithPID: different PIDs should produce different fingerprints")
	}
}

func TestGenerate_ProtocolCaseInsensitive(t *testing.T) {
	f := fingerprint.New()
	a := binding("TCP", "0.0.0.0", 80, 1, "httpd")
	b := binding("tcp", "0.0.0.0", 80, 1, "httpd")
	if f.Generate(a) != f.Generate(b) {
		t.Error("protocol case should not affect fingerprint")
	}
}

func TestEqual_SameBinding(t *testing.T) {
	f := fingerprint.New()
	b := binding("udp", "0.0.0.0", 53, 42, "dns")
	if !f.Equal(b, b) {
		t.Error("same binding should be equal")
	}
}

func TestEqual_DifferentAddress(t *testing.T) {
	f := fingerprint.New()
	a := binding("tcp", "0.0.0.0", 443, 1, "nginx")
	b := binding("tcp", "127.0.0.1", 443, 1, "nginx")
	if f.Equal(a, b) {
		t.Error("different addresses should not be equal")
	}
}

func TestGenerate_LengthIsFixed(t *testing.T) {
	f := fingerprint.New()
	b := binding("tcp", "0.0.0.0", 8080, 1, "app")
	got := f.Generate(b)
	// sha256 first 8 bytes → 16 hex chars
	if len(got) != 16 {
		t.Errorf("expected fingerprint length 16, got %d", len(got))
	}
}
