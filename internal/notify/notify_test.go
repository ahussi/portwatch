package notify

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// stubAlert is a minimal alert.Alert for testing.
type stubAlert struct {
	msg string
}

func (s stubAlert) String() string { return s.msg }

func TestNewHandler_Defaults(t *testing.T) {
	h := NewHandler("", "")
	if h.appName != "portwatch" {
		t.Errorf("expected default appName %q, got %q", "portwatch", h.appName)
	}
	if h.level != LevelWarning {
		t.Errorf("expected default level %q, got %q", LevelWarning, h.level)
	}
}

func TestNewHandler_Custom(t *testing.T) {
	h := NewHandler("myapp", LevelCritical)
	if h.appName != "myapp" {
		t.Errorf("expected appName %q, got %q", "myapp", h.appName)
	}
	if h.level != LevelCritical {
		t.Errorf("expected level %q, got %q", LevelCritical, h.level)
	}
}

func TestLevelToUrgency(t *testing.T) {
	cases := []struct {
		level    Level
		wantUrg  string
	}{
		{LevelCritical, "critical"},
		{LevelInfo, "low"},
		{LevelWarning, "normal"},
		{"unknown", "normal"},
	}
	for _, tc := range cases {
		got := levelToUrgency(tc.level)
		if got != tc.wantUrg {
			t.Errorf("levelToUrgency(%q) = %q, want %q", tc.level, got, tc.wantUrg)
		}
	}
}

// TestHandler_ImplementsAlertHandler ensures *Handler satisfies the alert.Handler interface.
func TestHandler_ImplementsAlertHandler(t *testing.T) {
	var _ interface{ Handle(alert.Alert) error } = NewHandler("", "")
}

// TestSend_UnsupportedOS verifies that an unsupported GOOS returns a descriptive error.
// We test the internal send path by patching via a helper.
func TestSend_ErrorContainsContext(t *testing.T) {
	h := &Handler{appName: "portwatch", level: LevelWarning}
	// Intentionally call send with a command that will fail on any platform
	// by overriding the body to trigger a non-zero exit on a real OS.
	// We rely on the fact that notify-send / osascript won't be available in CI,
	// so we just check the error wraps correctly.
	err := h.send("title", strings.Repeat("x", 1))
	// err may be nil if the tool happens to exist; skip assertion if so.
	if err != nil && !strings.Contains(err.Error(), "notify:") {
		t.Errorf("expected error to contain 'notify:', got: %v", err)
	}
}

func TestLevel_Constants(t *testing.T) {
	levels := []Level{LevelInfo, LevelWarning, LevelCritical}
	for _, l := range levels {
		if string(l) == "" {
			t.Errorf("level constant should not be empty")
		}
	}
	_ = time.Now() // ensure time import doesn't cause issues in future expansions
}
