package ratelimit_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/ratelimit"
)

func TestNew(t *testing.T) {
	l := ratelimit.New(time.Second)
	if l == nil {
		t.Fatal("expected non-nil Limiter")
	}
	if l.Len() != 0 {
		t.Fatalf("expected empty limiter, got len=%d", l.Len())
	}
}

func TestAllowFirstEvent(t *testing.T) {
	l := ratelimit.New(time.Second)
	if !l.Allow("tcp:8080") {
		t.Error("first event should be allowed")
	}
}

func TestAllowSuppressesWithinCooldown(t *testing.T) {
	l := ratelimit.New(time.Second)
	l.Allow("tcp:8080")
	if l.Allow("tcp:8080") {
		t.Error("second event within cooldown should be suppressed")
	}
}

func TestAllowAfterCooldown(t *testing.T) {
	l := ratelimit.New(10 * time.Millisecond)
	l.Allow("tcp:9090")
	time.Sleep(20 * time.Millisecond)
	if !l.Allow("tcp:9090") {
		t.Error("event after cooldown should be allowed")
	}
}

func TestAllowDifferentKeys(t *testing.T) {
	l := ratelimit.New(time.Second)
	l.Allow("tcp:8080")
	if !l.Allow("tcp:9090") {
		t.Error("different key should be allowed independently")
	}
}

func TestReset(t *testing.T) {
	l := ratelimit.New(time.Second)
	l.Allow("tcp:8080")
	l.Reset("tcp:8080")
	if !l.Allow("tcp:8080") {
		t.Error("event should be allowed after Reset")
	}
}

func TestPrune(t *testing.T) {
	l := ratelimit.New(10 * time.Millisecond)
	l.Allow("tcp:1111")
	l.Allow("tcp:2222")
	if l.Len() != 2 {
		t.Fatalf("expected 2 entries before prune, got %d", l.Len())
	}
	time.Sleep(20 * time.Millisecond)
	l.Prune()
	if l.Len() != 0 {
		t.Fatalf("expected 0 entries after prune, got %d", l.Len())
	}
}

func TestPruneKeepsActive(t *testing.T) {
	l := ratelimit.New(time.Second)
	l.Allow("tcp:3333")
	l.Prune()
	if l.Len() != 1 {
		t.Fatalf("active entry should survive prune, got len=%d", l.Len())
	}
}
