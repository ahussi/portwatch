package snapshot_test

import (
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func TestNew(t *testing.T) {
	s := snapshot.New()
	if s == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if keys := s.Keys(); len(keys) != 0 {
		t.Fatalf("expected empty snapshot, got %d keys", len(keys))
	}
}

func TestSetAndGet(t *testing.T) {
	s := snapshot.New()
	b := snapshot.Binding{Port: 8080, Protocol: "tcp", Address: "0.0.0.0"}
	s.Set("tcp:8080", b)

	got, ok := s.Get("tcp:8080")
	if !ok {
		t.Fatal("expected binding to be present")
	}
	if got.Port != 8080 {
		t.Errorf("expected port 8080, got %d", got.Port)
	}
	if got.SeenAt.IsZero() {
		t.Error("expected SeenAt to be set")
	}
}

func TestDelete(t *testing.T) {
	s := snapshot.New()
	s.Set("tcp:9090", snapshot.Binding{Port: 9090})
	s.Delete("tcp:9090")

	_, ok := s.Get("tcp:9090")
	if ok {
		t.Error("expected binding to be deleted")
	}
}

func TestKeys(t *testing.T) {
	s := snapshot.New()
	s.Set("tcp:80", snapshot.Binding{Port: 80})
	s.Set("udp:53", snapshot.Binding{Port: 53})

	keys := s.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestDiff(t *testing.T) {
	s := snapshot.New()
	s.Set("tcp:80", snapshot.Binding{Port: 80})
	s.Set("tcp:443", snapshot.Binding{Port: 443})

	newKeys := []string{"tcp:443", "tcp:8080"}
	added, removed := s.Diff(newKeys)

	if len(added) != 1 || added[0] != "tcp:8080" {
		t.Errorf("expected added=[tcp:8080], got %v", added)
	}
	if len(removed) != 1 || removed[0] != "tcp:80" {
		t.Errorf("expected removed=[tcp:80], got %v", removed)
	}
}

func TestDiffNoChanges(t *testing.T) {
	s := snapshot.New()
	s.Set("tcp:80", snapshot.Binding{Port: 80})

	added, removed := s.Diff([]string{"tcp:80"})
	if len(added) != 0 || len(removed) != 0 {
		t.Errorf("expected no diff, got added=%v removed=%v", added, removed)
	}
}
