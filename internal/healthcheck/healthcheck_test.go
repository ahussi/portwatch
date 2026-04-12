package healthcheck

import (
	"testing"
	"time"
)

func TestNew_DefaultStaleness(t *testing.T) {
	c := New(0)
	if c.staleness != defaultStaleness {
		t.Fatalf("expected default staleness %v, got %v", defaultStaleness, c.staleness)
	}
}

func TestNew_CustomStaleness(t *testing.T) {
	c := New(5 * time.Second)
	if c.staleness != 5*time.Second {
		t.Fatalf("expected 5s staleness, got %v", c.staleness)
	}
}

func TestStatus_Unknown_BeforeBeat(t *testing.T) {
	c := New(time.Second)
	if got := c.Status(); got != StatusUnknown {
		t.Fatalf("expected %q before first beat, got %q", StatusUnknown, got)
	}
}

func TestStatus_OK_AfterBeat(t *testing.T) {
	c := New(time.Second)
	c.Beat()
	if got := c.Status(); got != StatusOK {
		t.Fatalf("expected %q after beat, got %q", StatusOK, got)
	}
}

func TestStatus_Stale_AfterStaleness(t *testing.T) {
	c := New(10 * time.Millisecond)
	c.Beat()
	time.Sleep(20 * time.Millisecond)
	if got := c.Status(); got != StatusStale {
		t.Fatalf("expected %q after staleness window, got %q", StatusStale, got)
	}
}

func TestBeat_UpdatesLastBeat(t *testing.T) {
	c := New(time.Second)
	before := time.Now()
	c.Beat()
	after := time.Now()
	lb := c.LastBeat()
	if lb.Before(before) || lb.After(after) {
		t.Fatalf("LastBeat %v not in expected range [%v, %v]", lb, before, after)
	}
}

func TestStatus_RecoversAfterNewBeat(t *testing.T) {
	c := New(10 * time.Millisecond)
	c.Beat()
	time.Sleep(20 * time.Millisecond)
	if c.Status() != StatusStale {
		t.Fatal("expected stale before recovery")
	}
	c.Beat()
	if got := c.Status(); got != StatusOK {
		t.Fatalf("expected %q after new beat, got %q", StatusOK, got)
	}
}

func TestString_Unknown(t *testing.T) {
	c := New(time.Second)
	s := c.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
	if !containsSubstr(s, "unknown") {
		t.Fatalf("expected 'unknown' in string, got: %s", s)
	}
}

func TestString_OK(t *testing.T) {
	c := New(time.Second)
	c.Beat()
	s := c.String()
	if !containsSubstr(s, "ok") {
		t.Fatalf("expected 'ok' in string, got: %s", s)
	}
}

func containsSubstr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && func() bool {
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	}())
}
