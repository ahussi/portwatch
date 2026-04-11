// Package notify provides desktop and system notification support for portwatch alerts.
package notify

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// Level represents the urgency level of a notification.
type Level string

const (
	LevelInfo    Level = "info"
	LevelWarning Level = "warning"
	LevelCritical Level = "critical"
)

// Handler sends desktop notifications via OS-native tooling.
type Handler struct {
	appName string
	level   Level
}

// NewHandler creates a new desktop notification Handler.
// appName is used as the notification title prefix.
func NewHandler(appName string, level Level) *Handler {
	if appName == "" {
		appName = "portwatch"
	}
	if level == "" {
		level = LevelWarning
	}
	return &Handler{appName: appName, level: level}
}

// Handle implements alert.Handler. It dispatches a system notification
// for the given alert.
func (h *Handler) Handle(a alert.Alert) error {
	title := fmt.Sprintf("%s — %s", h.appName, strings.ToUpper(string(h.level)))
	body := a.String()
	return h.send(title, body)
}

// send dispatches the notification using the appropriate OS command.
func (h *Handler) send(title, body string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		script := fmt.Sprintf(`display notification %q with title %q`, body, title)
		cmd = exec.Command("osascript", "-e", script)
	case "linux":
		urgency := levelToUrgency(h.level)
		cmd = exec.Command("notify-send", "--urgency", urgency, title, body)
	case "windows":
		// PowerShell toast notification (Windows 10+)
		ps := fmt.Sprintf(
			`[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType=WindowsRuntime] | Out-Null; `+
				`$t = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent(0); `+
				`$t.GetElementsByTagName('text')[0].AppendChild($t.CreateTextNode('%s')) | Out-Null; `+
				`[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('%s').Show($t)`,
			body, title,
		)
		cmd = exec.Command("powershell", "-Command", ps)
	default:
		return fmt.Errorf("notify: unsupported OS %q", runtime.GOOS)
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("notify: send failed: %w", err)
	}
	return nil
}

func levelToUrgency(l Level) string {
	switch l {
	case LevelCritical:
		return "critical"
	case LevelInfo:
		return "low"
	default:
		return "normal"
	}
}
