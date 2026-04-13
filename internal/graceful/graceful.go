// Package graceful provides utilities for draining in-flight work before
// a process shuts down. A Drainer accepts units of work via Add, and
// Wait blocks until all work completes or the supplied context expires.
package graceful

import (
	"context"
	"sync"
	"sync/atomic"
)

// Drainer tracks in-flight work and allows callers to wait for
// completion with a context-aware deadline.
type Drainer struct {
	wg      sync.WaitGroup
	count   atomic.Int64
	mu      sync.Mutex
	drained []func() // callbacks invoked when count reaches zero
}

// New returns an initialised Drainer.
func New() *Drainer {
	return &Drainer{}
}

// Add increments the in-flight counter by delta (must be > 0).
// Panics if delta <= 0.
func (d *Drainer) Add(delta int) {
	if delta <= 0 {
		panic("graceful: Add called with non-positive delta")
	}
	d.count.Add(int64(delta))
	d.wg.Add(delta)
}

// Done decrements the in-flight counter by one.
func (d *Drainer) Done() {
	d.count.Add(-1)
	d.wg.Done()
	if d.count.Load() == 0 {
		d.runDrained()
	}
}

// Pending returns the number of in-flight units of work.
func (d *Drainer) Pending() int {
	return int(d.count.Load())
}

// OnDrained registers fn to be called when the pending count reaches zero.
// If the count is already zero fn is called immediately.
func (d *Drainer) OnDrained(fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.count.Load() == 0 {
		go fn()
		return
	}
	d.drained = append(d.drained, fn)
}

// Wait blocks until all in-flight work completes or ctx is cancelled.
// Returns ctx.Err() on timeout/cancellation, nil on clean drain.
func (d *Drainer) Wait(ctx context.Context) error {
	doneCh := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(doneCh)
	}()
	select {
	case <-doneCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (d *Drainer) runDrained() {
	d.mu.Lock()
	cbs := d.drained
	d.drained = nil
	d.mu.Unlock()
	for _, fn := range cbs {
		go fn()
	}
}
