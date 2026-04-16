package porttrend

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentRecord(t *testing.T) {
	tr := New()
	const goroutines = 50
	const records = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < records; j++ {
				tr.Record("tcp:8080")
			}
		}()
	}
	wg.Wait()
	e, ok := tr.Get("tcp:8080")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Count != goroutines*records {
		t.Fatalf("want %d, got %d", goroutines*records, e.Count)
	}
}

func TestConcurrentMultiKey(t *testing.T) {
	tr := New()
	const keys = 20
	const records = 50
	var wg sync.WaitGroup
	wg.Add(keys)
	for i := 0; i < keys; i++ {
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("tcp:%d", 8000+n)
			for j := 0; j < records; j++ {
				tr.Record(key)
			}
		}(i)
	}
	wg.Wait()
	all := tr.All()
	if len(all) != keys {
		t.Fatalf("want %d entries, got %d", keys, len(all))
	}
	for _, e := range all {
		if e.Count != records {
			t.Fatalf("key %s: want %d, got %d", e.Key, records, e.Count)
		}
	}
}
