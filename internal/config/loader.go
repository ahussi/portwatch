package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// candidatePaths returns the ordered list of locations portwatch searches
// for a configuration file when no explicit path is provided.
func candidatePaths() []string {
	paths := []string{
		"portwatch.json",
		".portwatch.json",
	}

	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".config", "portwatch", "config.json"))
	}

	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		paths = append(paths, filepath.Join(xdg, "portwatch", "config.json"))
	}

	return paths
}

// Discover attempts to load a config from well-known locations.
// It returns Default() when no config file is found, which is not an error.
func Discover() (*Config, error) {
	for _, p := range candidatePaths() {
		if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
			continue
		}
		cfg, err := Load(p)
		if err != nil {
			return nil, fmt.Errorf("loading config from %s: %w", p, err)
		}
		return cfg, nil
	}
	return Default(), nil
}

// Validate checks a Config for logical errors and returns a combined error
// message when problems are found.
func Validate(cfg *Config) error {
	if cfg.ScanInterval.Duration <= 0 {
		return fmt.Errorf("scan_interval must be a positive duration")
	}

	seen := make(map[int]bool)
	for _, p := range cfg.WatchPorts {
		if p < 1 || p > 65535 {
			return fmt.Errorf("watch_ports: %d is not a valid port number", p)
		}
		if seen[p] {
			return fmt.Errorf("watch_ports: duplicate port %d", p)
		}
		seen[p] = true
	}

	for _, p := range cfg.AllowedPorts {
		if p < 1 || p > 65535 {
			return fmt.Errorf("allowed_ports: %d is not a valid port number", p)
		}
	}

	return nil
}
