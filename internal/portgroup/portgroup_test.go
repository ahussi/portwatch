package portgroup_test

import (
	"testing"

	"github.com/user/portwatch/internal/portgroup"
)

func TestNew(t *testing.T) {
	r := portgroup.New()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
	if got := r.All(); len(got) != 0 {
		t.Fatalf("expected empty registry, got %d groups", len(got))
	}
}

func TestAdd_And_Get(t *testing.T) {
	r := portgroup.New()
	if err := r.Add("web", []int{80, 443, 8080}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g, ok := r.Get("web")
	if !ok {
		t.Fatal("expected group 'web' to exist")
	}
	if g.Name != "web" {
		t.Errorf("expected name 'web', got %q", g.Name)
	}
	if len(g.Ports) != 3 {
		t.Errorf("expected 3 ports, got %d", len(g.Ports))
	}
}

func TestAdd_Duplicate(t *testing.T) {
	r := portgroup.New()
	_ = r.Add("db", []int{5432})
	if err := r.Add("db", []int{3306}); err == nil {
		t.Fatal("expected error for duplicate group name")
	}
}

func TestGet_Missing(t *testing.T) {
	r := portgroup.New()
	_, ok := r.Get("nonexistent")
	if ok {
		t.Fatal("expected false for missing group")
	}
}

func TestRemove(t *testing.T) {
	r := portgroup.New()
	_ = r.Add("cache", []int{6379, 11211})
	r.Remove("cache")
	_, ok := r.Get("cache")
	if ok {
		t.Fatal("expected group to be removed")
	}
}

func TestRemove_Noop(t *testing.T) {
	r := portgroup.New()
	// should not panic
	r.Remove("ghost")
}

func TestContains(t *testing.T) {
	r := portgroup.New()
	_ = r.Add("web", []int{80, 443})
	if !r.Contains("web", 80) {
		t.Error("expected port 80 to be in 'web'")
	}
	if r.Contains("web", 8080) {
		t.Error("expected port 8080 NOT to be in 'web'")
	}
	if r.Contains("missing", 80) {
		t.Error("expected false for missing group")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	r := portgroup.New()
	_ = r.Add("a", []int{1})
	_ = r.Add("b", []int{2, 3})
	all := r.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(all))
	}
	// mutate returned slice — registry should be unaffected
	all[0].Ports[0] = 9999
	g, _ := r.Get(all[0].Name)
	if g.Ports[0] == 9999 {
		t.Error("All() should return a deep copy, not a reference")
	}
}
