package porttrend

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestNew(t *testing.T) {
	tr := New()
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
	if len(tr.All()) != 0 {
		t.Fatal("expected empty tracker")
	}
}

func TestRecord_FirstTime(t *testing.T) {
	tr := NewWithOptions(WithClock(fixedClock(epoch)))
	tr.Record("tcp:8080")
	e, ok := tr.Get("tcp:8080")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Count != 1 {
		t.Fatalf("want count 1, got %d", e.Count)
	}
	if !e.FirstSeen.Equal(epoch) {
		t.Fatalf("unexpected FirstSeen: %v", e.FirstSeen)
	}
}

func TestRecord_Increments(t *testing.T) {
	tr := NewWithOptions(WithClock(fixedClock(epoch)))
	tr.Record("tcp:8080")
	tr.Record("tcp:8080")
	tr.Record("tcp:8080")
	e, _ := tr.Get("tcp:8080")
	if e.Count != 3 {
		t.Fatalf("want 3, got %d", e.Count)
	}
}

func TestRecord_UpdatesLastSeen(t *testing.T) {
	later := epoch.Add(5 * time.Minute)
	calls := []time.Time{epoch, later}
	i := 0
	clock := func() time.Time { v := calls[i]; i++; return v }
	tr := NewWithOptions(WithClock(clock))
	tr.Record("udp:53")
	tr.Record("udp:53")
	e, _ := tr.Get("udp:53")
	if !e.FirstSeen.Equal(epoch) {
		t.Fatalf("unexpected FirstSeen: %v", e.FirstSeen)
	}
	if !e.LastSeen.Equal(later) {
		t.Fatalf("unexpected LastSeen: %v", e.LastSeen)
	}
}

func TestGet_Missing(t *testing.T) {
	tr := New()
	_, ok := tr.Get("tcp:9999")
	if ok {
		t.Fatal("expected missing")
	}
}

func TestAll(t *testing.T) {
	tr := New()
	tr.Record("tcp:80")
	tr.Record("tcp:443")
	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("want 2 entries, got %d", len(all))
	}
}

func TestReset(t *testing.T) {
	tr := New()
	tr.Record("tcp:80")
	tr.Reset()
	if len(tr.All()) != 0 {
		t.Fatal("expected empty after reset")
	}
}
