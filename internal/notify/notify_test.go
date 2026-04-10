package notify

import (
	"testing"
)

func TestLevelString(t *testing.T) {
	cases := []struct {
		level Level
		want  string
	}{
		{LevelInfo, "info"},
		{LevelWarn, "warn"},
		{LevelCritical, "critical"},
		{Level(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("Level(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestNew(t *testing.T) {
	n := New(true)
	if n == nil {
		t.Fatal("New returned nil")
	}
	if !n.enabled {
		t.Error("expected enabled=true")
	}

	n2 := New(false)
	if n2.enabled {
		t.Error("expected enabled=false")
	}
}

// TestSendDisabledIsNoOp ensures Send never errors when the notifier
// is disabled, regardless of content.
func TestSendDisabledIsNoOp(t *testing.T) {
	n := New(false)
	notif := Notification{
		Title:   "Port conflict",
		Message: "Process foo bound :8080",
		Level:   LevelCritical,
	}
	if err := n.Send(notif); err != nil {
		t.Errorf("Send on disabled notifier returned error: %v", err)
	}
}

// TestNotificationFields verifies struct fields are accessible.
func TestNotificationFields(t *testing.T) {
	notif := Notification{
		Title:   "New binding",
		Message: "nginx bound :443",
		Level:   LevelWarn,
	}
	if notif.Title != "New binding" {
		t.Errorf("unexpected Title: %q", notif.Title)
	}
	if notif.Level != LevelWarn {
		t.Errorf("unexpected Level: %v", notif.Level)
	}
}
