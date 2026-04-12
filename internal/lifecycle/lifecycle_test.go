package lifecycle_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/lifecycle"
)

func TestNew(t *testing.T) {
	lm := lifecycle.New()
	if lm == nil {
		t.Fatal("expected non-nil Manager")
	}
}

func TestRun_StartupError(t *testing.T) {
	lm := lifecycle.New()
	want := errors.New("startup failed")
	lm.OnStart(func(_ context.Context) error { return want })

	err := lm.Run(context.Background())
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestRun_CancelContext(t *testing.T) {
	lm := lifecycle.New()

	started := make(chan struct{})
	lm.OnStart(func(_ context.Context) error {
		close(started)
		return nil
	})

	stopped := make(chan struct{})
	lm.OnStop(func(_ context.Context) error {
		close(stopped)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- lm.Run(ctx) }()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("startup hook not called")
	}

	cancel()

	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("shutdown hook not called")
	}

	if err := <-done; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_ShutdownReverseOrder(t *testing.T) {
	lm := lifecycle.New()
	order := make([]int, 0, 3)

	for i := 0; i < 3; i++ {
		i := i
		lm.OnStop(func(_ context.Context) error {
			order = append(order, i)
			return nil
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_ = lm.Run(ctx)

	if len(order) != 3 || order[0] != 2 || order[1] != 1 || order[2] != 0 {
		t.Fatalf("expected reverse order [2 1 0], got %v", order)
	}
}
