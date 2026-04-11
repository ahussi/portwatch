package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/audit"
)

func tmpLog(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.log")
}

func TestNew(t *testing.T) {
	path := tmpLog(t)
	l, err := audit.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer l.Close()
	if l.Path() != path {
		t.Errorf("Path() = %q, want %q", l.Path(), path)
	}
}

func TestNew_InvalidPath(t *testing.T) {
	_, err := audit.New("/nonexistent/dir/audit.log")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestLog_WritesJSON(t *testing.T) {
	path := tmpLog(t)
	l, _ := audit.New(path)
	defer l.Close()

	entry := audit.Entry{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Event:     "new_binding",
		Port:      9090,
		Protocol:  "tcp",
		Process:   "caddy",
		PID:       42,
	}
	if err := l.Log(entry); err != nil {
		t.Fatalf("Log: %v", err)
	}
	l.Close()

	f, _ := os.Open(path)
	defer f.Close()
	var got audit.Entry
	if err := json.NewDecoder(bufio.NewReader(f)).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Port != 9090 || got.Event != "new_binding" || got.Process != "caddy" {
		t.Errorf("unexpected entry: %+v", got)
	}
}

func TestLog_SetsTimestampWhenZero(t *testing.T) {
	path := tmpLog(t)
	l, _ := audit.New(path)
	defer l.Close()

	before := time.Now().UTC()
	if err := l.Log(audit.Entry{Event: "test", Port: 1}); err != nil {
		t.Fatalf("Log: %v", err)
	}
	after := time.Now().UTC()
	l.Close()

	f, _ := os.Open(path)
	defer f.Close()
	var got audit.Entry
	json.NewDecoder(f).Decode(&got) //nolint:errcheck
	if got.Timestamp.Before(before) || got.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", got.Timestamp, before, after)
	}
}

func TestLog_MultipleEntries(t *testing.T) {
	path := tmpLog(t)
	l, _ := audit.New(path)
	defer l.Close()

	for i := 0; i < 5; i++ {
		if err := l.Log(audit.Entry{Event: "test", Port: i + 1, Protocol: "tcp"}); err != nil {
			t.Fatalf("Log[%d]: %v", i, err)
		}
	}
	l.Close()

	f, _ := os.Open(path)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	if count != 5 {
		t.Errorf("expected 5 lines, got %d", count)
	}
}
