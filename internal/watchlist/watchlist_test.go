package watchlist

import (
	"fmt"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	wl := New()
	if wl == nil {
		t.Fatal("expected non-nil Watchlist")
	}
	if wl.Len() != 0 {
		t.Fatalf("expected empty watchlist, got %d entries", wl.Len())
	}
}

func TestAddAndHas(t *testing.T) {
	wl := New()
	wl.Add(Entry{Port: 8080, Protocol: "tcp", Label: "api"})

	if !wl.Has("tcp", 8080) {
		t.Error("expected tcp:8080 to be present")
	}
	if wl.Has("udp", 8080) {
		t.Error("udp:8080 should not be present")
	}
}

func TestRemove(t *testing.T) {
	wl := New()
	wl.Add(Entry{Port: 9000, Protocol: "tcp"})
	wl.Remove("tcp", 9000)

	if wl.Has("tcp", 9000) {
		t.Error("expected tcp:9000 to be removed")
	}
	if wl.Len() != 0 {
		t.Errorf("expected length 0, got %d", wl.Len())
	}
}

func TestAll(t *testing.T) {
	wl := New()
	wl.Add(Entry{Port: 80, Protocol: "tcp", Label: "http"})
	wl.Add(Entry{Port: 443, Protocol: "tcp", Label: "https"})
	wl.Add(Entry{Port: 53, Protocol: "udp", Label: "dns"})

	all := wl.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(all))
	}
}

func TestEntryKey(t *testing.T) {
	e := Entry{Port: 8443, Protocol: "tcp"}
	if e.Key() != "tcp:8443" {
		t.Errorf("unexpected key: %s", e.Key())
	}
}

func TestConcurrentAdd(t *testing.T) {
	wl := New()
	var wg sync.WaitGroup
	for i := 1000; i < 1050; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			wl.Add(Entry{Port: port, Protocol: "tcp", Label: fmt.Sprintf("svc-%d", port)})
		}(i)
	}
	wg.Wait()
	if wl.Len() != 50 {
		t.Errorf("expected 50 entries, got %d", wl.Len())
	}
}

func TestRemoveNonExistent(t *testing.T) {
	wl := New()
	// should not panic
	wl.Remove("tcp", 9999)
	if wl.Len() != 0 {
		t.Errorf("expected empty watchlist, got %d", wl.Len())
	}
}
