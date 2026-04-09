package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.ScanInterval.Duration != 5*time.Second {
		t.Errorf("expected 5s scan interval, got %v", cfg.ScanInterval.Duration)
	}
	if !cfg.AlertOnNew {
		t.Error("expected AlertOnNew to be true by default")
	}
	if cfg.AlertOnClose {
		t.Error("expected AlertOnClose to be false by default")
	}
}

func TestLoad(t *testing.T) {
	data := `{
		"watch_ports": [8080, 9090],
		"allowed_ports": [22, 443],
		"scan_interval": "10s",
		"alert_on_new": true,
		"alert_on_close": true,
		"log_file": "/tmp/portwatch.log"
	}`

	f, err := os.CreateTemp("", "portwatch-config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(data)
	f.Close()

	cfg, err := config.Load(f.Name())
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.ScanInterval.Duration != 10*time.Second {
		t.Errorf("expected 10s, got %v", cfg.ScanInterval.Duration)
	}
	if len(cfg.WatchPorts) != 2 {
		t.Errorf("expected 2 watch ports, got %d", len(cfg.WatchPorts))
	}
	if cfg.LogFile != "/tmp/portwatch.log" {
		t.Errorf("unexpected log file: %s", cfg.LogFile)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/portwatch.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestIsAllowed(t *testing.T) {
	cfg := &config.Config{AllowedPorts: []int{22, 80, 443}}
	if !cfg.IsAllowed(80) {
		t.Error("expected port 80 to be allowed")
	}
	if cfg.IsAllowed(8080) {
		t.Error("expected port 8080 to not be allowed")
	}
}

func TestIsWatched(t *testing.T) {
	cfg := &config.Config{WatchPorts: []int{8080, 9090}}
	if !cfg.IsWatched(8080) {
		t.Error("expected 8080 to be watched")
	}
	if cfg.IsWatched(3000) {
		t.Error("expected 3000 to not be watched")
	}

	// Empty WatchPorts means watch everything.
	cfgAll := &config.Config{}
	if !cfgAll.IsWatched(12345) {
		t.Error("expected all ports to be watched when WatchPorts is empty")
	}
}

func TestDurationRoundTrip(t *testing.T) {
	d := config.Duration{Duration: 30 * time.Second}
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	var d2 config.Duration
	if err := json.Unmarshal(b, &d2); err != nil {
		t.Fatal(err)
	}
	if d.Duration != d2.Duration {
		t.Errorf("round-trip mismatch: %v != %v", d.Duration, d2.Duration)
	}
}
