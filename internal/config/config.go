package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	// Ports to monitor; if empty, all ports are monitored.
	WatchPorts []int `json:"watch_ports"`

	// Ports that are always allowed and will never trigger alerts.
	AllowedPorts []int `json:"allowed_ports"`

	// ScanInterval is how often the scanner polls for port changes.
	ScanInterval Duration `json:"scan_interval"`

	// AlertOnNew triggers an alert when a new port binding is detected.
	AlertOnNew bool `json:"alert_on_new"`

	// AlertOnClose triggers an alert when a previously seen port binding closes.
	AlertOnClose bool `json:"alert_on_close"`

	// LogFile is the optional path to write alerts to (stdout if empty).
	LogFile string `json:"log_file"`
}

// Duration is a wrapper around time.Duration for JSON (un)marshalling.
type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = parsed
	return nil
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		ScanInterval: Duration{5 * time.Second},
		AlertOnNew:   true,
		AlertOnClose: false,
	}
}

// Load reads a JSON config file from the given path.
// Missing fields fall back to Default() values.
func Load(path string) (*Config, error) {
	cfg := Default()

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// IsAllowed reports whether port is in the allowed list.
func (c *Config) IsAllowed(port int) bool {
	for _, p := range c.AllowedPorts {
		if p == port {
			return true
		}
	}
	return false
}

// IsWatched reports whether port should be monitored.
// When WatchPorts is empty every port is watched.
func (c *Config) IsWatched(port int) bool {
	if len(c.WatchPorts) == 0 {
		return true
	}
	for _, p := range c.WatchPorts {
		if p == port {
			return true
		}
	}
	return false
}
