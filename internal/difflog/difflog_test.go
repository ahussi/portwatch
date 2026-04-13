package difflog

import (
	"testing"
	"time"
)

func TestNew_DefaultCapacity(t *testing.T) {
	l := New(0)
	if l.capacity != 256 {
		t.Fatalf("expected default capacity 256, got %d", l.capacity)
	}
}

func TestNew_CustomCapacity(t *testing.T) {
	l := New(10)
	if l.capacity != 10 {
		t.Fatalf("expected capacity 10, got %d", l.capacity)
	}
}

func TestAdd_SetsTimestamp(t *testing.T) {
	l := New(8)
	before := time.Now()
	l.Add(Event{Kind: KindAdded, Key: "tcp:8080"})
	after := time.Now()
	events := l.All()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	ts := events[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v outside expected range [%v, %v]", ts, before, after)
	}
}

func TestAdd_EvictsOldest(t *testing.T) {
	l := New(3)
	for i, k := range []string{"a", "b", "c", "d"} {
		l.Add(Event{Kind: KindAdded, Key: k, Port: i + 1})
	}
	events := l.All()
	if len(events) != 3 {
		t.Fatalf("expected 3 events after eviction, got %d", len(events))
	}
	if events[0].Key != "b" {
		t.Errorf("expected oldest surviving key 'b', got %q", events[0].Key)
	}
	if events[2].Key != "d" {
		t.Errorf("expected newest key 'd', got %q", events[2].Key)
	}
}

func TestSince(t *testing.T) {
	l := New(16)
	t0 := time.Now()
	l.Add(Event{Kind: KindAdded, Key: "old", Timestamp: t0.Add(-2 * time.Second)})
	l.Add(Event{Kind: KindAdded, Key: "new", Timestamp: t0.Add(1 * time.Second)})

	results := l.Since(t0)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "new" {
		t.Errorf("expected key 'new', got %q", results[0].Key)
	}
}

func TestLen(t *testing.T) {
	l := New(10)
	for i := 0; i < 5; i++ {
		l.Add(Event{Kind: KindRemoved, Port: i})
	}
	if l.Len() != 5 {
		t.Errorf("expected Len 5, got %d", l.Len())
	}
}

func TestClear(t *testing.T) {
	l := New(10)
	l.Add(Event{Kind: KindAdded, Port: 9000})
	l.Clear()
	if l.Len() != 0 {
		t.Errorf("expected empty log after Clear, got %d", l.Len())
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	l := New(10)
	l.Add(Event{Kind: KindAdded, Key: "tcp:443"})
	copy1 := l.All()
	copy1[0].Key = "mutated"
	copy2 := l.All()
	if copy2[0].Key == "mutated" {
		t.Error("All() should return an independent copy")
	}
}
