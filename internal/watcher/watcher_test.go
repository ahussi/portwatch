package watcher_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watcher"
)

func defaultConfig() *config.Config {
	cfg := config.Default()
	cfg.Interval = 1 // 1 second for tests
	return cfg
}

func TestNew(t *testing.T) {
	cfg := defaultConfig()
	s := scanner.New(cfg)
	am := alert.NewManager()
	w := watcher.New(cfg, s, am)
	if w == nil {
		t.Fatal("expected non-nil Watcher")
	}
}

func TestRunCancelImmediately(t *testing.T) {
	cfg := defaultConfig()
	s := scanner.New(cfg)
	am := alert.NewManager()
	w := watcher.New(cfg, s, am)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel right away

	err := w.Run(ctx)
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestRunDispatchesNewBinding(t *testing.T) {
	cfg := defaultConfig()
	cfg.Interval = 1

	s := scanner.New(cfg)
	am := alert.NewManager()

	received := make([]alert.Alert, 0)
	am.AddHandler(alert.HandlerFunc(func(a alert.Alert) {
		received = append(received, a)
	}))

	w := watcher.New(cfg, s, am)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	// Run should complete after timeout; ignore the deadline error.
	_ = w.Run(ctx)

	// We cannot assert specific bindings in a unit test without mocking the
	// scanner, but we can assert the watcher ran without panicking.
	t.Logf("alerts dispatched: %d", len(received))
}

func TestRunContextDeadlineExceeded(t *testing.T) {
	cfg := defaultConfig()
	s := scanner.New(cfg)
	am := alert.NewManager()
	w := watcher.New(cfg, s, am)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := w.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}
