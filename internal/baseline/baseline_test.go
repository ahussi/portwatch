package baseline

import (
	"os"
	"path/filepath"
	"testing"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestNew(t *testing.T) {
	b := New(tmpPath(t))
	if b == nil {
		t.Fatal("expected non-nil Baseline")
	}
	if len(b.Entries()) != 0 {
		t.Fatalf("expected empty entries, got %d", len(b.Entries()))
	}
}

func TestSetAndHas(t *testing.T) {
	b := New(tmpPath(t))
	key := "tcp:0.0.0.0:9000"
	if b.Has(key) {
		t.Fatal("should not have key before Set")
	}
	b.Set(key, Entry{Port: 9000, Protocol: "tcp", Address: "0.0.0.0", PID: 42, Process: "myapp"})
	if !b.Has(key) {
		t.Fatal("should have key after Set")
	}
}

func TestEntries(t *testing.T) {
	b := New(tmpPath(t))
	b.Set("tcp:127.0.0.1:80", Entry{Port: 80})
	b.Set("udp:0.0.0.0:53", Entry{Port: 53})
	e := b.Entries()
	if len(e) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(e))
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tmpPath(t)
	b := New(path)
	b.Set("tcp:0.0.0.0:8080", Entry{Port: 8080, Protocol: "tcp", Address: "0.0.0.0", PID: 1, Process: "server"})
	if err := b.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	b2 := New(path)
	if err := b2.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !b2.Has("tcp:0.0.0.0:8080") {
		t.Fatal("loaded baseline missing expected key")
	}
	e := b2.Entries()["tcp:0.0.0.0:8080"]
	if e.Port != 8080 || e.Process != "server" {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestLoadMissingFile(t *testing.T) {
	b := New("/nonexistent/path/baseline.json")
	err := b.Load()
	if !os.IsNotExist(err) {
		t.Fatalf("expected ErrNotExist, got %v", err)
	}
}

func TestSavedAt(t *testing.T) {
	path := tmpPath(t)
	b := New(path)
	if err := b.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if b.SavedAt.IsZero() {
		t.Fatal("SavedAt should be set after Save")
	}
}
