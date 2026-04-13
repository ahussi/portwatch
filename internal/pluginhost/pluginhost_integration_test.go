package pluginhost_test

import (
	"sync"
	"testing"

	"github.com/yourorg/portwatch/internal/pluginhost"
)

// TestConcurrentRegisterAndGet verifies that concurrent registration and
// lookup operations do not cause data races.
func TestConcurrentRegisterAndGet(t *testing.T) {
	h := pluginhost.New()
	const workers = 20
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			name := string(rune('A' + id))
			p := pluginhost.NewNoopPlugin(name)
			// Ignore duplicate errors; multiple goroutines may share a letter.
			_ = h.Register(p, map[string]string{"id": name})
			_, _ = h.Get(name)
			_ = h.Names()
		}(i % 10) // intentionally collide on 10 names
	}
	wg.Wait()
}

// TestCloseAllIdempotent ensures CloseAll on an empty host returns nil.
func TestCloseAllIdempotent(t *testing.T) {
	h := pluginhost.New()
	if err := h.CloseAll(); err != nil {
		t.Fatalf("CloseAll on empty host: %v", err)
	}
}
