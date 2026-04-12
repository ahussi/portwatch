package resolver

import (
	"testing"
)

func TestNew(t *testing.T) {
	r := New()
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
	if len(r.services) == 0 {
		t.Fatal("expected default services to be loaded")
	}
}

func TestResolve_KnownPort(t *testing.T) {
	r := New()
	svc := r.Resolve(443)
	if svc == nil {
		t.Fatal("expected service info for port 443")
	}
	if svc.Name != "https" {
		t.Errorf("expected 'https', got %q", svc.Name)
	}
	if svc.Protocol != "tcp" {
		t.Errorf("expected 'tcp', got %q", svc.Protocol)
	}
	if svc.Port != 443 {
		t.Errorf("expected port 443, got %d", svc.Port)
	}
}

func TestResolve_UnknownPort(t *testing.T) {
	r := New()
	svc := r.Resolve(9999)
	if svc != nil {
		t.Errorf("expected nil for unknown port, got %+v", svc)
	}
}

func TestName_KnownPort(t *testing.T) {
	r := New()
	if got := r.Name(22); got != "ssh" {
		t.Errorf("expected 'ssh', got %q", got)
	}
}

func TestName_UnknownPort(t *testing.T) {
	r := New()
	if got := r.Name(9999); got != "port/9999" {
		t.Errorf("expected 'port/9999', got %q", got)
	}
}

func TestRegister_AddsCustomService(t *testing.T) {
	r := New()
	r.Register(9090, "prometheus", "tcp")
	svc := r.Resolve(9090)
	if svc == nil {
		t.Fatal("expected registered service to be found")
	}
	if svc.Name != "prometheus" {
		t.Errorf("expected 'prometheus', got %q", svc.Name)
	}
}

func TestRegister_OverwritesExisting(t *testing.T) {
	r := New()
	r.Register(80, "custom-http", "tcp")
	if got := r.Name(80); got != "custom-http" {
		t.Errorf("expected 'custom-http', got %q", got)
	}
}

func TestDefaultServices(t *testing.T) {
	r := New()
	cases := map[int]string{
		21:    "ftp",
		22:    "ssh",
		53:    "dns",
		80:    "http",
		3306:  "mysql",
		5432:  "postgres",
		6379:  "redis",
		27017: "mongodb",
	}
	for port, want := range cases {
		if got := r.Name(port); got != want {
			t.Errorf("port %d: expected %q, got %q", port, want, got)
		}
	}
}
