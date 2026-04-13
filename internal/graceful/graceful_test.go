package graceful_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/graceful"
)

func TestNew(t *testing.T) {
	d := graceful.New()
	if d == nil {
		t.Fatal("expected non-nil Drainer")
	}
	if d.Pending() != 0 {
		t.Fatalf("expected 0 pending, got %d", d.Pending())
	}
}

func TestAddAndDone(t *testing.T) {
	d := graceful.New()
	d.Add(3)
	if d.Pending() != 3 {
		t.Fatalf("expected 3 pending, got %d", d.Pending())
	}
	d.Done()
	d.Done()
	if d.Pending() != 1 {
		t.Fatalf("expected 1 pending, got %d", d.Pending())
	}
	d.Done()
	if d.Pending() != 0 {
		t.Fatalf("expected 0 pending after all done")
	}
}

func TestAdd_PanicsOnZero(t *testing.T) {
	d := graceful.New()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for delta=0")
		}
	}()
	d.Add(0)
}

func TestWait_CleanDrain(t *testing.T) {
	d := graceful.New()
	d.Add(2)
	go func() {
		time.Sleep(10 * time.Millisecond)
		d.Done()
		d.Done()
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := d.Wait(ctx); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestWait_ContextCancelled(t *testing.T) {
	d := graceful.New()
	d.Add(1) // never Done'd
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	err := d.Wait(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestOnDrained_CalledWhenZero(t *testing.T) {
	d := graceful.New()
	d.Add(1)
	var called atomic.Bool
	d.OnDrained(func() { called.Store(true) })
	d.Done()
	time.Sleep(30 * time.Millisecond)
	if !called.Load() {
		t.Fatal("OnDrained callback was not called")
	}
}

func TestOnDrained_ImmediateWhenAlreadyZero(t *testing.T) {
	d := graceful.New()
	var called atomic.Bool
	d.OnDrained(func() { called.Store(true) })
	time.Sleep(20 * time.Millisecond)
	if !called.Load() {
		t.Fatal("expected immediate callback when pending=0")
	}
}

func TestConcurrentAddDone(t *testing.T) {
	d := graceful.New()
	const workers = 50
	d.Add(workers)
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond)
			d.Done()
		}()
	}
	wg.Wait()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := d.Wait(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Pending() != 0 {
		t.Fatalf("expected 0 pending after concurrent drain")
	}
}
