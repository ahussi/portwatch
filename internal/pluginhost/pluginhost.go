// Package pluginhost provides a lightweight plugin registry that allows
// external alert handlers and scanner extensions to be registered and
// invoked at runtime without modifying core portwatch code.
package pluginhost

import (
	"errors"
	"fmt"
	"sync"
)

// Plugin is the interface every portwatch plugin must satisfy.
type Plugin interface {
	// Name returns the unique identifier for the plugin.
	Name() string
	// Init is called once when the plugin is registered.
	Init(cfg map[string]string) error
	// Close releases any resources held by the plugin.
	Close() error
}

// Host manages the lifecycle of registered plugins.
type Host struct {
	mu      sync.RWMutex
	plugins map[string]Plugin
}

// New returns an initialised, empty Host.
func New() *Host {
	return &Host{
		plugins: make(map[string]Plugin),
	}
}

// Register adds a plugin to the host and calls Init with the supplied cfg.
// Returns an error if a plugin with the same name is already registered or
// if Init fails.
func (h *Host) Register(p Plugin, cfg map[string]string) error {
	if p == nil {
		return errors.New("pluginhost: plugin must not be nil")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	name := p.Name()
	if _, exists := h.plugins[name]; exists {
		return fmt.Errorf("pluginhost: plugin %q already registered", name)
	}
	if err := p.Init(cfg); err != nil {
		return fmt.Errorf("pluginhost: init %q: %w", name, err)
	}
	h.plugins[name] = p
	return nil
}

// Get returns the plugin registered under name, or false if absent.
func (h *Host) Get(name string) (Plugin, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	p, ok := h.plugins[name]
	return p, ok
}

// Names returns the sorted list of registered plugin names.
func (h *Host) Names() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	names := make([]string, 0, len(h.plugins))
	for n := range h.plugins {
		names = append(names, n)
	}
	return names
}

// CloseAll calls Close on every registered plugin in arbitrary order.
// All errors are collected and returned as a single combined error.
func (h *Host) CloseAll() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	var errs []error
	for name, p := range h.plugins {
		if err := p.Close(); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}
