package circuitbreaker

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestConcurrentFailuresOpenCircuit(t *testing.T) {
	cb := New(10, time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cb.RecordFailure()
		}()
	}
	wg.Wait()
	if cb.State() != StateOpen {
		t.Fatalf("expected Open after concurrent failures")
	}
}

func TestConcurrentAllow_AllRejectedWhenOpen(t *testing.T) {
	cb := New(1, time.Hour)
	cb.RecordFailure()

	var rejected int64
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := cb.Allow(); err == ErrOpen {
				atomic.AddInt64(&rejected, 1)
			}
		}()
	}
	wg.Wait()
	if rejected != 50 {
		t.Fatalf("expected 50 rejections, got %d", rejected)
	}
}

func TestCooldownRecoveryUnderLoad(t *testing.T) {
	cb := New(3, 20*time.Millisecond)
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	time.Sleep(40 * time.Millisecond)

	// First Allow should transition to HalfOpen
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil after cooldown: %v", err)
	}
	cb.RecordSuccess()
	if cb.State() != StateClosed {
		t.Fatalf("expected Closed after probe success")
	}
}
