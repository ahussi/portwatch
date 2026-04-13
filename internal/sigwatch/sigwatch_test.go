package sigwatch_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/user/portwatch/internal/sigwatch"
)

func TestNew_Defaults(t *testing.T) {
	h := sigwatch.New()
	if h == nil {
		t.Fatal("expected non-nil Handler")
	}
}

func TestRun_CancelContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	h := sigwatch.New()

	done := make(chan struct{})
	go func() {
		h.Run(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
		// ok
	case <-time.After(time.Second):
		t.Fatal("Run did not return after context cancellation")
	}
}

func TestRun_ShutdownCallback(t *testing.T) {
	ctx := context.Background()
	shutdownCalled := make(chan struct{}, 1)

	h := sigwatch.New(
		sigwatch.WithShutdown(func() { shutdownCalled <- struct{}{} }),
	)

	done := make(chan struct{})
	go func() {
		h.Run(ctx)
		close(done)
	}()

	// give goroutine time to register signal handler
	time.Sleep(20 * time.Millisecond)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	select {
	case <-shutdownCalled:
		// ok
	case <-time.After(time.Second):
		t.Fatal("shutdown callback not invoked")
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not return after SIGTERM")
	}
}

func TestRun_ReloadCallback(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reloadCalled := make(chan struct{}, 1)
	h := sigwatch.New(
		sigwatch.WithReload(func() { reloadCalled <- struct{}{} }),
	)

	go h.Run(ctx)

	time.Sleep(20 * time.Millisecond)
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGHUP)

	select {
	case <-reloadCalled:
		// ok
	case <-time.After(time.Second):
		t.Fatal("reload callback not invoked")
	}
}
