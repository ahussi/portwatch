package portevict_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portevict"
)

func makeEvent(key string, port int, ago time.Duration) portevict.Event {
	return portevict.Event{
		Key:       key,
		Port:      port,
		Proto:     "tcp",
		Addr:      "0.0.0.0",
		EvictedAt: time.Now().Add(-ago),
	}
}

func TestNew_DefaultCapacity(t *testing.T) {
	tr := portevict.New(0)
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestRecord_And_Len(t *testing.T) {
	tr := portevict.New(10)
	tr.Record(makeEvent("tcp:0.0.0.0:8080", 8080, 0))
	tr.Record(makeEvent("tcp:0.0.0.0:9090", 9090, 0))
	if tr.Len() != 2 {
		t.Fatalf("expected 2 events, got %d", tr.Len())
	}
}

func TestRecord_SetsTimestampWhenZero(t *testing.T) {
	tr := portevict.New(10)
	e := portevict.Event{Key: "tcp:0.0.0.0:80", Port: 80}
	tr.Record(e)
	all := tr.All()
	if all[0].EvictedAt.IsZero() {
		t.Error("expected EvictedAt to be set automatically")
	}
}

func TestRecord_EvictsOldestWhenFull(t *testing.T) {
	tr := portevict.New(3)
	tr.Record(makeEvent("tcp:0.0.0.0:1", 1, 0))
	tr.Record(makeEvent("tcp:0.0.0.0:2", 2, 0))
	tr.Record(makeEvent("tcp:0.0.0.0:3", 3, 0))
	tr.Record(makeEvent("tcp:0.0.0.0:4", 4, 0))
	all := tr.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 events after eviction, got %d", len(all))
	}
	if all[0].Port != 2 {
		t.Errorf("expected oldest surviving port=2, got %d", all[0].Port)
	}
}

func TestSince(t *testing.T) {
	tr := portevict.New(10)
	tr.Record(makeEvent("tcp:0.0.0.0:100", 100, 10*time.Minute))
	tr.Record(makeEvent("tcp:0.0.0.0:200", 200, 30*time.Second))
	tr.Record(makeEvent("tcp:0.0.0.0:300", 300, 5*time.Second))

	recent := tr.Since(time.Now().Add(-1 * time.Minute))
	if len(recent) != 2 {
		t.Fatalf("expected 2 recent events, got %d", len(recent))
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tr := portevict.New(10)
	tr.Record(makeEvent("tcp:0.0.0.0:8080", 8080, 0))
	a := tr.All()
	a[0].Port = 9999
	b := tr.All()
	if b[0].Port == 9999 {
		t.Error("All() should return an independent copy")
	}
}

func TestClear(t *testing.T) {
	tr := portevict.New(10)
	tr.Record(makeEvent("tcp:0.0.0.0:8080", 8080, 0))
	tr.Clear()
	if tr.Len() != 0 {
		t.Errorf("expected 0 events after Clear, got %d", tr.Len())
	}
}
