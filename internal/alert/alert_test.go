package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestAlertString(t *testing.T) {
	a := Alert{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Level:     LevelWarning,
		Port:      8080,
		Process:   "nginx",
		Message:   "unexpected binding",
	}
	s := a.String()
	if !strings.Contains(s, "WARNING") {
		t.Errorf("expected WARNING in alert string, got: %s", s)
	}
	if !strings.Contains(s, "8080") {
		t.Errorf("expected port 8080 in alert string, got: %s", s)
	}
	if !strings.Contains(s, "nginx") {
		t.Errorf("expected process name in alert string, got: %s", s)
	}
}

func TestStdoutHandler(t *testing.T) {
	var buf bytes.Buffer
	h := NewStdoutHandler(&buf)
	a := Alert{
		Timestamp: time.Now(),
		Level:     LevelCritical,
		Port:      443,
		Process:   "unknown",
		Message:   "conflict detected",
	}
	if err := h.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "CRITICAL") {
		t.Errorf("expected CRITICAL in output, got: %s", buf.String())
	}
}

func TestManagerDispatch(t *testing.T) {
	var buf bytes.Buffer
	h := NewStdoutHandler(&buf)
	m := NewManager(h)

	m.Dispatch(LevelInfo, 3000, "node", "new binding observed")

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("expected INFO level in output, got: %s", output)
	}
	if !strings.Contains(output, "3000") {
		t.Errorf("expected port 3000 in output, got: %s", output)
	}
}

func TestNewStdoutHandlerDefaultsToStdout(t *testing.T) {
	h := NewStdoutHandler(nil)
	if h.writer == nil {
		t.Error("expected non-nil writer when nil is passed")
	}
}
