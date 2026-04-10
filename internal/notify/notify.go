// Package notify provides desktop and system notification support
// for portwatch alerts, dispatching messages via OS-native mechanisms.
package notify

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Level represents the severity of a notification.
type Level int

const (
	LevelInfo Level = iota
	LevelWarn
	LevelCritical
)

// String returns a human-readable label for the level.
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// Notification holds the data for a single desktop notification.
type Notification struct {
	Title   string
	Message string
	Level   Level
}

// Notifier dispatches system notifications.
type Notifier struct {
	enabled bool
}

// New creates a Notifier. If enabled is false, Send is a no-op.
func New(enabled bool) *Notifier {
	return &Notifier{enabled: enabled}
}

// Send dispatches n via the OS-native notification mechanism.
// Returns an error if the underlying command fails; unsupported
// platforms silently succeed.
func (n *Notifier) Send(notif Notification) error {
	if !n.enabled {
		return nil
	}
	switch runtime.GOOS {
	case "darwin":
		return sendDarwin(notif)
	case "linux":
		return sendLinux(notif)
	default:
		// Windows and others: unsupported, no-op.
		return nil
	}
}

func sendDarwin(n Notification) error {
	script := fmt.Sprintf(
		`display notification %q with title %q subtitle %q`,
		n.Message, "portwatch", n.Title,
	)
	return exec.Command("osascript", "-e", script).Run()
}

func sendLinux(n Notification) error {
	urgency := "normal"
	if n.Level == LevelCritical {
		urgency = "critical"
	}
	return exec.Command(
		"notify-send",
		"--urgency", urgency,
		"--app-name", "portwatch",
		n.Title,
		n.Message,
	).Run()
}
