package suppress_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/suppress"
)

var (
	epoch   = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	advance = func(base *time.Time) func() time.Time {
		return func() time.Time { return *base }
	}
)

func TestNew(t *testing.T) {
	m := suppress.New(nil)
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestSuppressAndIsSuppressed(t *testing.T) {
	now := epoch
	m := suppress.New(advance(&now))

	m.Suppress("tcp:8080", 5*time.Minute)

	if !m.IsSuppressed("tcp:8080") {
		t.Error("expected tcp:8080 to be suppressed")
	}
	if m.IsSuppressed("tcp:9090") {
		t.Error("expected tcp:9090 NOT to be suppressed")
	}
}

func TestIsSuppressed_Expired(t *testing.T) {
	now := epoch
	m := suppress.New(advance(&now))

	m.Suppress("tcp:8080", 1*time.Minute)
	now = epoch.Add(2 * time.Minute) // advance past expiry

	if m.IsSuppressed("tcp:8080") {
		t.Error("expected suppression to have expired")
	}
}

func TestRemove(t *testing.T) {
	now := epoch
	m := suppress.New(advance(&now))

	m.Suppress("tcp:8080", 10*time.Minute)
	m.Remove("tcp:8080")

	if m.IsSuppressed("tcp:8080") {
		t.Error("expected suppression to be removed")
	}
}

func TestPrune(t *testing.T) {
	now := epoch
	m := suppress.New(advance(&now))

	m.Suppress("tcp:8080", 1*time.Minute)
	m.Suppress("tcp:9090", 10*time.Minute)

	now = epoch.Add(2 * time.Minute)
	removed := m.Prune()

	if removed != 1 {
		t.Errorf("expected 1 pruned, got %d", removed)
	}
	if m.IsSuppressed("tcp:8080") {
		t.Error("tcp:8080 should have been pruned")
	}
	if !m.IsSuppressed("tcp:9090") {
		t.Error("tcp:9090 should still be active")
	}
}

func TestLen(t *testing.T) {
	now := epoch
	m := suppress.New(advance(&now))

	if m.Len() != 0 {
		t.Fatalf("expected 0, got %d", m.Len())
	}

	m.Suppress("tcp:8080", 5*time.Minute)
	m.Suppress("tcp:9090", 5*time.Minute)

	if m.Len() != 2 {
		t.Errorf("expected 2, got %d", m.Len())
	}

	now = epoch.Add(10 * time.Minute)
	if m.Len() != 0 {
		t.Errorf("expected 0 after expiry, got %d", m.Len())
	}
}

func TestEntry_IsExpired(t *testing.T) {
	e := suppress.Entry{
		Key:       "tcp:443",
		ExpiresAt: epoch.Add(1 * time.Minute),
	}
	if e.IsExpired(epoch) {
		t.Error("should not be expired at epoch")
	}
	if !e.IsExpired(epoch.Add(2 * time.Minute)) {
		t.Error("should be expired 2 minutes later")
	}
}
