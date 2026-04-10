package watcher

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
)

// Watcher periodically scans ports and emits alerts on changes.
type Watcher struct {
	cfg      *config.Config
	scanner  *scanner.Scanner
	alerts   *alert.Manager
	previous map[string]scanner.Binding
	interval time.Duration
}

// New creates a new Watcher with the given config, scanner, and alert manager.
func New(cfg *config.Config, s *scanner.Scanner, am *alert.Manager) *Watcher {
	return &Watcher{
		cfg:      cfg,
		scanner:  s,
		alerts:   am,
		previous: make(map[string]scanner.Binding),
		interval: time.Duration(cfg.Interval) * time.Second,
	}
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Perform an initial scan immediately.
	if err := w.tick(); err != nil {
		log.Printf("portwatch: initial scan error: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := w.tick(); err != nil {
				log.Printf("portwatch: scan error: %v", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// tick performs a single scan cycle and dispatches alerts for any changes.
func (w *Watcher) tick() error {
	current, err := w.scanner.ScanPorts()
	if err != nil {
		return err
	}

	currentMap := make(map[string]scanner.Binding, len(current))
	for _, b := range current {
		key := scanner.GetBindingKey(b)
		currentMap[key] = b

		if _, seen := w.previous[key]; !seen {
			w.alerts.Dispatch(alert.Alert{
				Kind:    alert.KindNewBinding,
				Binding: b,
				Message: "new port binding detected",
			})
		}
	}

	for key, b := range w.previous {
		if _, still := currentMap[key]; !still {
			w.alerts.Dispatch(alert.Alert{
				Kind:    alert.KindClosedBinding,
				Binding: b,
				Message: "port binding closed",
			})
		}
	}

	w.previous = currentMap
	return nil
}
