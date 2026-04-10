package snapshot_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

// TestConcurrentAccess verifies the snapshot is safe for concurrent use.
func TestConcurrentAccess(t *testing.T) {
	s := snapshot.New()
	const goroutines = 20
	var wg sync.WaitGroup

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			key := "tcp:" + itoa(i)
			s.Set(key, snapshot.Binding{Port: i, Protocol: "tcp"})
			_, _ = s.Get(key)
			_ = s.Keys()
		}(i)
	}
	wg.Wait()

	if len(s.Keys()) != goroutines {
		t.Errorf("expected %d keys after concurrent writes, got %d", goroutines, len(s.Keys()))
	}
}

// TestDiffAfterDelete ensures removed keys appear correctly after deletion.
func TestDiffAfterDelete(t *testing.T) {
	s := snapshot.New()
	s.Set("tcp:3000", snapshot.Binding{Port: 3000})
	s.Set("tcp:3001", snapshot.Binding{Port: 3001})
	s.Delete("tcp:3001")

	// Snapshot now only has tcp:3000; new scan has neither.
	added, removed := s.Diff([]string{})
	if len(added) != 0 {
		t.Errorf("expected no added, got %v", added)
	}
	if len(removed) != 1 || removed[0] != "tcp:3000" {
		t.Errorf("expected removed=[tcp:3000], got %v", removed)
	}
}

// itoa is a minimal int-to-string helper to avoid importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
