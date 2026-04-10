package reporter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/snapshot"
)

func makeSnap(t *testing.T) *snapshot.Snapshot {
	t.Helper()
	snap := snapshot.New()
	snap.Set("tcp:127.0.0.1:8080", snapshot.Binding{
		Proto:   "tcp",
		Addr:    "127.0.0.1",
		Port:    8080,
		Process: "nginx",
		SeenAt:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	})
	snap.Set("tcp:0.0.0.0:443", snapshot.Binding{
		Proto:   "tcp",
		Addr:    "0.0.0.0",
		Port:    443,
		Process: "caddy",
		SeenAt:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	})
	return snap
}

func TestNew(t *testing.T) {
	r := reporter.New(nil, reporter.FormatText)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestRenderText(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Render(makeSnap(t)); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "BINDING") {
		t.Error("expected header BINDING in text output")
	}
	if !strings.Contains(out, "nginx") {
		t.Error("expected process 'nginx' in text output")
	}
	if !strings.Contains(out, "caddy") {
		t.Error("expected process 'caddy' in text output")
	}
}

func TestRenderJSON(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)
	if err := r.Render(makeSnap(t)); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[")
		|| !strings.Contains(out, "]") {
		t.Error("expected JSON array delimiters")
	}
	if !strings.Contains(out, "\"binding\"") {
		t.Error("expected 'binding' key in JSON output")
	}
}

func TestRenderEmptySnapshot(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Render(snapshot.New()); err != nil {
		t.Fatalf("Render empty: %v", err)
	}
	if !strings.Contains(buf.String(), "BINDING") {
		t.Error("expected header even for empty snapshot")
	}
}
