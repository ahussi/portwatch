package metrics

import (
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("expected non-nil Metrics")
	}
	if m.startedAt.IsZero() {
		t.Error("expected startedAt to be set")
	}
}

func TestCounters_InitiallyZero(t *testing.T) {
	m := New()
	c := m.Snapshot()
	if c.Scans != 0 || c.NewBindings != 0 || c.Alerts != 0 || c.Suppressed != 0 || c.Errors != 0 {
		t.Errorf("expected all counters to be zero, got %+v", c)
	}
}

func TestIncrement(t *testing.T) {
	m := New()
	m.IncScans()
	m.IncScans()
	m.IncNewBindings()
	m.IncAlerts()
	m.IncSuppressed()
	m.IncErrors()

	c := m.Snapshot()
	if c.Scans != 2 {
		t.Errorf("Scans: want 2, got %d", c.Scans)
	}
	if c.NewBindings != 1 {
		t.Errorf("NewBindings: want 1, got %d", c.NewBindings)
	}
	if c.Alerts != 1 {
		t.Errorf("Alerts: want 1, got %d", c.Alerts)
	}
	if c.Suppressed != 1 {
		t.Errorf("Suppressed: want 1, got %d", c.Suppressed)
	}
	if c.Errors != 1 {
		t.Errorf("Errors: want 1, got %d", c.Errors)
	}
}

func TestUptime(t *testing.T) {
	m := New()
	time.Sleep(10 * time.Millisecond)
	if m.Uptime() < 10*time.Millisecond {
		t.Error("expected uptime >= 10ms")
	}
}

func TestString(t *testing.T) {
	m := New()
	m.IncScans()
	m.IncAlerts()
	s := m.String()
	for _, want := range []string{"scans=1", "alerts=1", "uptime="} {
		if !strings.Contains(s, want) {
			t.Errorf("String() missing %q, got: %s", want, s)
		}
	}
}

func TestConcurrentIncrements(t *testing.T) {
	m := New()
	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			m.IncScans()
			m.IncNewBindings()
		}()
	}
	wg.Wait()
	c := m.Snapshot()
	if c.Scans != goroutines {
		t.Errorf("Scans: want %d, got %d", goroutines, c.Scans)
	}
	if c.NewBindings != goroutines {
		t.Errorf("NewBindings: want %d, got %d", goroutines, c.NewBindings)
	}
}
