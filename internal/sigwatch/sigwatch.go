// Package sigwatch provides OS signal handling for graceful shutdown
// and reload of the portwatch daemon.
package sigwatch

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Handler listens for OS signals and triggers registered callbacks.
type Handler struct {
	signals []os.Signal
	onShutdown func()
	onReload   func()
}

// Option configures a Handler.
type Option func(*Handler)

// WithShutdown registers a callback invoked on SIGINT or SIGTERM.
func WithShutdown(fn func()) Option {
	return func(h *Handler) { h.onShutdown = fn }
}

// WithReload registers a callback invoked on SIGHUP.
func WithReload(fn func()) Option {
	return func(h *Handler) { h.onReload = fn }
}

// New creates a Handler with the provided options.
func New(opts ...Option) *Handler {
	h := &Handler{
		signals:    []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP},
		onShutdown: func() {},
		onReload:   func() {},
	}
	for _, o := range opts {
		o(h)
	}
	return h
}

// Run blocks until ctx is cancelled or a handled signal is received.
// SIGINT/SIGTERM invoke the shutdown callback; SIGHUP invokes reload.
func (h *Handler) Run(ctx context.Context) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, h.signals...)
	defer signal.Stop(ch)

	for {
		select {
		case <-ctx.Done():
			return
		case sig := <-ch:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				h.onShutdown()
				return
			case syscall.SIGHUP:
				h.onReload()
			}
		}
	}
}
