package portstate_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/portstate"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew(t *testing.T) {
	tr := portstate.New(nil)
	if tr == nil {
		t.Fatal("expected non-nil Tracker")
	}
	if tr.Len() != 0 {
		t.Fatalf("expected empty tracker, got %d", tr.Len())
	}
}

func TestObserve_FirstTime(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tr := portstate.New(fixedNow(now))
	tr.Observe("tcp:8080")

	s, ok := tr.Get("tcp:8080")
	if !ok {
		t.Fatal("expected state to exist")
	}
	if s.SeenCount != 1 {
		t.Fatalf("expected SeenCount=1, got %d", s.SeenCount)
	}
	if !s.FirstSeen.Equal(now) {
		t.Fatalf("expected FirstSeen=%v, got %v", now, s.FirstSeen)
	}
	if !s.LastSeen.Equal(now) {
		t.Fatalf("expected LastSeen=%v, got %v", now, s.LastSeen)
	}
}

func TestObserve_UpdatesExisting(t *testing.T) {
	t0 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := t0.Add(5 * time.Second)
	calls := []time.Time{t0, t1}
	i := 0
	tr := portstate.New(func() time.Time {
		v := calls[i]
		if i < len(calls)-1 {
			i++
		}
		return v
	})

	tr.Observe("tcp:9090")
	tr.Observe("tcp:9090")

	s, ok := tr.Get("tcp:9090")
	if !ok {
		t.Fatal("expected state to exist")
	}
	if s.SeenCount != 2 {
		t.Fatalf("expected SeenCount=2, got %d", s.SeenCount)
	}
	if !s.FirstSeen.Equal(t0) {
		t.Fatalf("FirstSeen should remain %v, got %v", t0, s.FirstSeen)
	}
	if !s.LastSeen.Equal(t1) {
		t.Fatalf("LastSeen should be %v, got %v", t1, s.LastSeen)
	}
}

func TestGet_Missing(t *testing.T) {
	tr := portstate.New(nil)
	_, ok := tr.Get("tcp:404")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}

func TestRemove(t *testing.T) {
	tr := portstate.New(nil)
	tr.Observe("udp:53")
	tr.Remove("udp:53")
	if tr.Len() != 0 {
		t.Fatalf("expected Len=0 after Remove, got %d", tr.Len())
	}
}

func TestKeys(t *testing.T) {
	tr := portstate.New(nil)
	tr.Observe("tcp:80")
	tr.Observe("tcp:443")

	keys := tr.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	seen := map[string]bool{}
	for _, k := range keys {
		seen[k] = true
	}
	for _, want := range []string{"tcp:80", "tcp:443"} {
		if !seen[want] {
			t.Errorf("missing key %q in Keys()", want)
		}
	}
}
