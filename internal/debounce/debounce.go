// Package debounce provides a mechanism to suppress rapid repeated events
// for the same key within a configurable delay window. It is useful for
// avoiding alert storms when a port binding flaps quickly.
package debounce

import (
	"sync"
	"time"
)

// Debouncer delays event processing so that only the last event within a
// window is acted upon.
type Debouncer struct {
	mu      sync.Mutex
	delay   time.Duration
	timers  map[string]*time.Timer
	callback func(key string)
}

// New creates a Debouncer with the given delay. The callback is invoked with
// the key after the delay has elapsed without another call for the same key.
func New(delay time.Duration, callback func(key string)) *Debouncer {
	if delay <= 0 {
		delay = 200 * time.Millisecond
	}
	return &Debouncer{
		delay:    delay,
		timers:   make(map[string]*time.Timer),
		callback: callback,
	}
}

// Trigger schedules the callback for key after the debounce delay. If Trigger
// is called again for the same key before the delay elapses, the timer is
// reset.
func (d *Debouncer) Trigger(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Reset(d.delay)
		return
	}

	d.timers[key] = time.AfterFunc(d.delay, func() {
		d.mu.Lock()
		delete(d.timers, key)
		d.mu.Unlock()
		d.callback(key)
	})
}

// Cancel stops a pending callback for key. It is a no-op if no timer exists.
func (d *Debouncer) Cancel(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
		delete(d.timers, key)
	}
}

// Pending returns the number of keys currently waiting to fire.
func (d *Debouncer) Pending() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.timers)
}
