package ratelimit_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/ratelimit"
)

// TestConcurrentAllow verifies the limiter is safe for concurrent use.
func TestConcurrentAllow(t *testing.T) {
	l := ratelimit.New(50 * time.Millisecond)
	const goroutines = 20
	var allowed atomic.Int64
	var wg sync.WaitGroup

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("tcp:%d", 8000+(id%5))
			if l.Allow(key) {
				allowed.Add(1)
			}
		}(i)
	}
	wg.Wait()

	// Only 5 unique keys exist; each should be allowed exactly once.
	if got := allowed.Load(); got != 5 {
		t.Errorf("expected 5 allowed events (one per unique key), got %d", got)
	}
}

// TestPruneConcurrent verifies Prune does not race with Allow.
func TestPruneConcurrent(t *testing.T) {
	l := ratelimit.New(5 * time.Millisecond)
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			l.Allow(fmt.Sprintf("tcp:%d", i))
			time.Sleep(time.Millisecond)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			l.Prune()
			time.Sleep(5 * time.Millisecond)
		}
	}()

	wg.Wait() // no race detector errors expected
}

// TestAllowRespectsCooldown verifies that the same key is denied within the
// cooldown window and allowed again only after the window has elapsed.
func TestAllowRespectsCooldown(t *testing.T) {
	cooldown := 50 * time.Millisecond
	l := ratelimit.New(cooldown)
	key := "tcp:9000"

	if !l.Allow(key) {
		t.Fatal("expected first Allow to return true")
	}
	if l.Allow(key) {
		t.Error("expected second Allow within cooldown to return false")
	}

	time.Sleep(cooldown + 10*time.Millisecond)

	if !l.Allow(key) {
		t.Error("expected Allow after cooldown expiry to return true")
	}
}
