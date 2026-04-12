package debounce_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
)

func TestNew_DefaultDelay(t *testing.T) {
	d := debounce.New(0, func(string) {})
	if d == nil {
		t.Fatal("expected non-nil Debouncer")
	}
}

func TestTrigger_FiresAfterDelay(t *testing.T) {
	var fired atomic.Int32
	d := debounce.New(50*time.Millisecond, func(key string) {
		if key == "port:8080" {
			fired.Add(1)
		}
	})

	d.Trigger("port:8080")
	time.Sleep(100 * time.Millisecond)

	if fired.Load() != 1 {
		t.Errorf("expected callback to fire once, got %d", fired.Load())
	}
}

func TestTrigger_ResetsTimer(t *testing.T) {
	var fired atomic.Int32
	d := debounce.New(80*time.Millisecond, func(string) {
		fired.Add(1)
	})

	d.Trigger("port:9090")
	time.Sleep(40 * time.Millisecond)
	d.Trigger("port:9090") // reset
	time.Sleep(40 * time.Millisecond)
	// only 80ms since last trigger — should not have fired yet
	if fired.Load() != 0 {
		t.Errorf("callback fired too early")
	}
	time.Sleep(60 * time.Millisecond)
	if fired.Load() != 1 {
		t.Errorf("expected 1 fire after reset, got %d", fired.Load())
	}
}

func TestCancel_StopsPending(t *testing.T) {
	var fired atomic.Int32
	d := debounce.New(60*time.Millisecond, func(string) {
		fired.Add(1)
	})

	d.Trigger("port:3000")
	d.Cancel("port:3000")
	time.Sleep(100 * time.Millisecond)

	if fired.Load() != 0 {
		t.Errorf("expected no fire after cancel, got %d", fired.Load())
	}
}

func TestPending_Count(t *testing.T) {
	d := debounce.New(200*time.Millisecond, func(string) {})

	d.Trigger("a")
	d.Trigger("b")
	d.Trigger("c")

	if p := d.Pending(); p != 3 {
		t.Errorf("expected 3 pending, got %d", p)
	}

	d.Cancel("b")
	if p := d.Pending(); p != 2 {
		t.Errorf("expected 2 pending after cancel, got %d", p)
	}
}

func TestTrigger_Concurrent(t *testing.T) {
	var count atomic.Int32
	d := debounce.New(60*time.Millisecond, func(string) {
		count.Add(1)
	})

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.Trigger("shared-key")
		}()
	}
	wg.Wait()
	time.Sleep(120 * time.Millisecond)

	if count.Load() != 1 {
		t.Errorf("expected exactly 1 fire for shared key, got %d", count.Load())
	}
}
