package baseline_test

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"

	"github.com/user/portwatch/internal/baseline"
)

func TestConcurrentSetAndHas(t *testing.T) {
	b := baseline.New(filepath.Join(t.TempDir(), "baseline.json"))
	const workers = 20
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("tcp:0.0.0.0:%d", 3000+i)
			b.Set(key, baseline.Entry{Port: 3000 + i, Protocol: "tcp"})
			_ = b.Has(key)
		}(i)
	}
	wg.Wait()
	if got := len(b.Entries()); got != workers {
		t.Fatalf("expected %d entries, got %d", workers, got)
	}
}

func TestSaveLoadRoundtrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bl.json")
	b := baseline.New(path)
	keys := []string{"tcp:0.0.0.0:80", "tcp:0.0.0.0:443", "udp:127.0.0.1:53"}
	for i, k := range keys {
		b.Set(k, baseline.Entry{Port: 80 + i, Protocol: "tcp", Process: fmt.Sprintf("proc%d", i)})
	}
	if err := b.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	b2 := baseline.New(path)
	if err := b2.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, k := range keys {
		if !b2.Has(k) {
			t.Errorf("missing key %q after reload", k)
		}
	}
	if !b2.SavedAt.IsZero() {
		t.Logf("baseline saved at %s", b2.SavedAt)
	}
}
