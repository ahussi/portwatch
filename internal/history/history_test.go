package history_test

import (
	"testing"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/scanner"
)

func binding(port int) scanner.Binding {
	return scanner.Binding{Port: port, Proto: "tcp", Address: "0.0.0.0"}
}

func TestNew(t *testing.T) {
	r := history.New(10)
	if r == nil {
		t.Fatal("expected non-nil Record")
	}
	if r.Len() != 0 {
		t.Fatalf("expected empty record, got %d events", r.Len())
	}
}

func TestNewDefaultCapacity(t *testing.T) {
	r := history.New(0)
	if r == nil {
		t.Fatal("expected non-nil Record for zero capacity")
	}
}

func TestAddAndLen(t *testing.T) {
	r := history.New(5)
	r.Add(history.EventAdded, binding(8080))
	r.Add(history.EventAdded, binding(9090))
	if r.Len() != 2 {
		t.Fatalf("expected 2 events, got %d", r.Len())
	}
}

func TestAllOrder(t *testing.T) {
	r := history.New(10)
	ports := []int{80, 443, 8080}
	for _, p := range ports {
		r.Add(history.EventAdded, binding(p))
	}

	events := r.All()
	if len(events) != len(ports) {
		t.Fatalf("expected %d events, got %d", len(ports), len(events))
	}
	for i, e := range events {
		if e.Binding.Port != ports[i] {
			t.Errorf("event[%d]: expected port %d, got %d", i, ports[i], e.Binding.Port)
		}
	}
}

func TestCircularOverwrite(t *testing.T) {
	cap := 3
	r := history.New(cap)
	for i := 1; i <= 5; i++ {
		r.Add(history.EventAdded, binding(i))
	}

	if r.Len() != cap {
		t.Fatalf("expected %d events after overflow, got %d", cap, r.Len())
	}

	events := r.All()
	// oldest retained should be port 3 (events 1,2 were overwritten)
	if events[0].Binding.Port != 3 {
		t.Errorf("expected oldest port 3, got %d", events[0].Binding.Port)
	}
	if events[cap-1].Binding.Port != 5 {
		t.Errorf("expected newest port 5, got %d", events[cap-1].Binding.Port)
	}
}

func TestEventKinds(t *testing.T) {
	r := history.New(10)
	r.Add(history.EventAdded, binding(8080))
	r.Add(history.EventRemoved, binding(8080))

	events := r.All()
	if events[0].Kind != history.EventAdded {
		t.Errorf("expected EventAdded, got %s", events[0].Kind)
	}
	if events[1].Kind != history.EventRemoved {
		t.Errorf("expected EventRemoved, got %s", events[1].Kind)
	}
}
