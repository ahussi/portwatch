package pluginhost_test

import (
	"errors"
	"testing"

	"github.com/yourorg/portwatch/internal/pluginhost"
)

func TestNew(t *testing.T) {
	h := pluginhost.New()
	if h == nil {
		t.Fatal("expected non-nil Host")
	}
	if got := h.Names(); len(got) != 0 {
		t.Fatalf("expected empty names, got %v", got)
	}
}

func TestRegisterAndGet(t *testing.T) {
	h := pluginhost.New()
	p := pluginhost.NewNoopPlugin("alpha")
	cfg := map[string]string{"k": "v"}

	if err := h.Register(p, cfg); err != nil {
		t.Fatalf("Register: %v", err)
	}
	got, ok := h.Get("alpha")
	if !ok {
		t.Fatal("expected plugin to be found")
	}
	if got.Name() != "alpha" {
		t.Errorf("name = %q, want %q", got.Name(), "alpha")
	}
	if p.Cfg["k"] != "v" {
		t.Errorf("cfg not propagated to Init")
	}
}

func TestRegisterDuplicate(t *testing.T) {
	h := pluginhost.New()
	p := pluginhost.NewNoopPlugin("dup")
	h.Register(p, nil) //nolint:errcheck
	err := h.Register(pluginhost.NewNoopPlugin("dup"), nil)
	if err == nil {
		t.Fatal("expected error on duplicate registration")
	}
}

func TestRegisterNil(t *testing.T) {
	h := pluginhost.New()
	if err := h.Register(nil, nil); err == nil {
		t.Fatal("expected error registering nil plugin")
	}
}

func TestGetMissing(t *testing.T) {
	h := pluginhost.New()
	_, ok := h.Get("nonexistent")
	if ok {
		t.Fatal("expected false for missing plugin")
	}
}

func TestNames(t *testing.T) {
	h := pluginhost.New()
	for _, name := range []string{"c", "a", "b"} {
		h.Register(pluginhost.NewNoopPlugin(name), nil) //nolint:errcheck
	}
	names := h.Names()
	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}
}

func TestCloseAll(t *testing.T) {
	h := pluginhost.New()
	p1 := pluginhost.NewNoopPlugin("p1")
	p2 := pluginhost.NewNoopPlugin("p2")
	h.Register(p1, nil) //nolint:errcheck
	h.Register(p2, nil) //nolint:errcheck

	if err := h.CloseAll(); err != nil {
		t.Fatalf("CloseAll: %v", err)
	}
	if !p1.Closed || !p2.Closed {
		t.Error("expected both plugins to be closed")
	}
}

func TestCloseAll_CollectsErrors(t *testing.T) {
	h := pluginhost.New()
	h.Register(&errPlugin{name: "bad1"}, nil) //nolint:errcheck
	h.Register(&errPlugin{name: "bad2"}, nil) //nolint:errcheck

	err := h.CloseAll()
	if err == nil {
		t.Fatal("expected combined error")
	}
}

// errPlugin is a Plugin whose Close always returns an error.
type errPlugin struct{ name string }

func (e *errPlugin) Name() string                    { return e.name }
func (e *errPlugin) Init(_ map[string]string) error  { return nil }
func (e *errPlugin) Close() error                    { return errors.New("close failed") }
