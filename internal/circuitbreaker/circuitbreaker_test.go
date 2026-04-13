package circuitbreaker

import (
	"testing"
	"time"
)

func TestNew_Defaults(t *testing.T) {
	cb := New(3, time.Second)
	if cb.State() != StateClosed {
		t.Fatalf("expected Closed, got %v", cb.State())
	}
}

func TestNew_ThresholdFloor(t *testing.T) {
	cb := New(0, time.Second)
	if cb.threshold != 1 {
		t.Fatalf("expected threshold=1, got %d", cb.threshold)
	}
}

func TestAllow_ClosedAlwaysPermits(t *testing.T) {
	cb := New(3, time.Second)
	if err := cb.Allow(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	cb := New(3, time.Second)
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	if cb.State() != StateOpen {
		t.Fatalf("expected Open after threshold failures")
	}
}

func TestAllow_RejectsWhenOpen(t *testing.T) {
	cb := New(1, 10*time.Second)
	cb.RecordFailure()
	if err := cb.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestAllow_HalfOpenAfterCooldown(t *testing.T) {
	cb := New(1, 10*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil after cooldown, got %v", err)
	}
	if cb.State() != StateHalfOpen {
		t.Fatalf("expected HalfOpen, got %v", cb.State())
	}
}

func TestRecordSuccess_ClosesFromHalfOpen(t *testing.T) {
	cb := New(1, 10*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	_ = cb.Allow() // transition to HalfOpen
	cb.RecordSuccess()
	if cb.State() != StateClosed {
		t.Fatalf("expected Closed after success in HalfOpen")
	}
}

func TestRecordFailure_ReopensFromHalfOpen(t *testing.T) {
	cb := New(1, 10*time.Millisecond)
	cb.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	_ = cb.Allow()
	cb.RecordFailure()
	if cb.State() != StateOpen {
		t.Fatalf("expected Open after failure in HalfOpen")
	}
}

func TestReset_ForcesClosed(t *testing.T) {
	cb := New(1, time.Hour)
	cb.RecordFailure()
	cb.Reset()
	if cb.State() != StateClosed {
		t.Fatalf("expected Closed after Reset")
	}
	if err := cb.Allow(); err != nil {
		t.Fatalf("unexpected error after reset: %v", err)
	}
}

func TestRecordSuccess_NoopWhenClosed(t *testing.T) {
	cb := New(3, time.Second)
	cb.RecordSuccess()
	if cb.State() != StateClosed {
		t.Fatalf("expected Closed")
	}
}
