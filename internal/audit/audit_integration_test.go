package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
	"testing"

	"github.com/user/portwatch/internal/audit"
)

func TestConcurrentLog(t *testing.T) {
	path := tmpLog(t)
	l, err := audit.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer l.Close()

	const workers = 10
	const perWorker = 20
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < perWorker; i++ {
				if err := l.Log(audit.Entry{
					Event:    "concurrent",
					Port:     id*100 + i,
					Protocol: "tcp",
				}); err != nil {
					t.Errorf("worker %d Log: %v", id, err)
				}
			}
		}(w)
	}
	wg.Wait()
	l.Close()

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		var e audit.Entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			t.Errorf("invalid JSON on line %d: %v", count+1, err)
		}
		count++
	}
	expected := workers * perWorker
	if count != expected {
		t.Errorf("expected %d entries, got %d", expected, count)
	}
}
