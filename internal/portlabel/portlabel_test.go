package portlabel

import (
	"testing"
)

func TestNew_ContainsBuiltins(t *testing.T) {
	l := New()
	all := l.All()
	if len(all) == 0 {
		t.Fatal("expected built-in labels, got none")
	}
	lbl, ok := l.Get(80, "tcp")
	if !ok {
		t.Fatal("expected built-in label for 80/tcp")
	}
	if lbl.Name != "HTTP" {
		t.Errorf("expected HTTP, got %s", lbl.Name)
	}
	if lbl.Category != "web" {
		t.Errorf("expected web, got %s", lbl.Category)
	}
}

func TestSet_And_Get(t *testing.T) {
	l := New()
	l.Set(9090, "tcp", "Prometheus", "monitoring")
	lbl, ok := l.Get(9090, "tcp")
	if !ok {
		t.Fatal("expected label to be present")
	}
	if lbl.Name != "Prometheus" {
		t.Errorf("expected Prometheus, got %s", lbl.Name)
	}
	if lbl.Category != "monitoring" {
		t.Errorf("expected monitoring, got %s", lbl.Category)
	}
}

func TestGet_Missing(t *testing.T) {
	l := New()
	_, ok := l.Get(19999, "udp")
	if ok {
		t.Error("expected missing label, got one")
	}
}

func TestSet_Overwrites(t *testing.T) {
	l := New()
	l.Set(80, "tcp", "CustomHTTP", "custom")
	lbl, ok := l.Get(80, "tcp")
	if !ok {
		t.Fatal("expected label after overwrite")
	}
	if lbl.Name != "CustomHTTP" {
		t.Errorf("expected CustomHTTP, got %s", lbl.Name)
	}
}

func TestRemove(t *testing.T) {
	l := New()
	l.Set(7777, "tcp", "Test", "test")
	l.Remove(7777, "tcp")
	_, ok := l.Get(7777, "tcp")
	if ok {
		t.Error("expected label to be removed")
	}
}

func TestLabelString_WithCategory(t *testing.T) {
	lbl := Label{Name: "Redis", Category: "cache"}
	got := lbl.String()
	expected := "Redis (cache)"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestLabelString_NoCategory(t *testing.T) {
	lbl := Label{Name: "Custom"}
	if lbl.String() != "Custom" {
		t.Errorf("expected 'Custom', got %q", lbl.String())
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	l := New()
	a1 := l.All()
	a1["999/tcp"] = Label{Name: "injected"}
	_, ok := l.Get(999, "tcp")
	if ok {
		t.Error("modifying All() result should not affect internal state")
	}
}
