// Package lifecycle manages the startup and graceful shutdown of portwatch daemon components.
package lifecycle

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Hook is a function called during a lifecycle phase.
type Hook func(ctx context.Context) error

// Manager coordinates startup and shutdown hooks for daemon components.
type Manager struct {
	mu       sync.Mutex
	startups []Hook
	shutdowns []Hook
}

// New returns an initialised Manager.
func New() *Manager {
	return &Manager{}
}

// OnStart registers a hook to run during startup.
func (m *Manager) OnStart(h Hook) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.startups = append(m.startups, h)
}

// OnStop registers a hook to run during graceful shutdown.
func (m *Manager) OnStop(h Hook) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shutdowns = append(m.shutdowns, h)
}

// Run executes all startup hooks, blocks until the context is cancelled or an
// OS interrupt is received, then executes shutdown hooks in reverse order.
func (m *Manager) Run(ctx context.Context) error {
	m.mu.Lock()
	startups := append([]Hook(nil), m.startups...)
	shutdowns := append([]Hook(nil), m.shutdowns...)
	m.mu.Unlock()

	for _, h := range startups {
		if err := h(ctx); err != nil {
			return err
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case <-ctx.Done():
	case <-quit:
	}

	for i := len(shutdowns) - 1; i >= 0; i-- {
		_ = shutdowns[i](ctx)
	}
	return nil
}
