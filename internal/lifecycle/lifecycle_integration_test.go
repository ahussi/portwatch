package lifecycle_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/lifecycle"
)

func TestIntegration_MultipleStartStopHooks(t *testing.T) {
	lm := lifecycle.New()

	var startCount, stopCount int32

	for i := 0; i < 5; i++ {
		lm.OnStart(func(_ context.Context) error {
			atomic.AddInt32(&startCount, 1)
			return nil
		})
		lm.OnStop(func(_ context.Context) error {
			atomic.AddInt32(&stopCount, 1)
			return nil
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = lm.Run(ctx)

	if got := atomic.LoadInt32(&startCount); got != 5 {
		t.Errorf("expected 5 start hooks, got %d", got)
	}
	if got := atomic.LoadInt32(&stopCount); got != 5 {
		t.Errorf("expected 5 stop hooks, got %d", got)
	}
}

func TestIntegration_NoHooksRunsClean(t *testing.T) {
	lm := lifecycle.New()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- lm.Run(ctx) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Run did not return after context cancellation")
	}
}
