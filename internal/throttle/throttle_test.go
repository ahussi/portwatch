package throttle_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/throttle"
)

func TestNew(t *testing.T) {
	th := throttle.New(time.Second, 3)
	if th == nil {
		t.Fatal("expected non-nil Throttle")
	}
}

func TestNew_BurstFloor(t *testing.T) {
	// burst < 1 should be clamped to 1
	th := throttle.New(time.Second, 0)
	if !th.Allow("k") {
		t.Fatal("first call should always be allowed")
	}
	if th.Allow("k") {
		t.Fatal("second call should be denied when burst=1")
	}
}

func TestAllow_WithinBurst(t *testing.T) {
	th := throttle.New(time.Second, 3)
	key := "port:8080"
	for i := 0; i < 3; i++ {
		if !th.Allow(key) {
			t.Fatalf("call %d should be allowed within burst", i+1)
		}
	}
}

func TestAllow_ExceedsBurst(t *testing.T) {
	th := throttle.New(time.Second, 2)
	key := "port:9090"
	th.Allow(key)
	th.Allow(key)
	if th.Allow(key) {
		t.Fatal("third call should be denied after burst exhausted")
	}
}

func TestAllow_WindowReset(t *testing.T) {
	th := throttle.New(20*time.Millisecond, 1)
	key := "port:443"
	if !th.Allow(key) {
		t.Fatal("first call should be allowed")
	}
	if th.Allow(key) {
		t.Fatal("second call should be denied")
	}
	time.Sleep(30 * time.Millisecond)
	if !th.Allow(key) {
		t.Fatal("call after window reset should be allowed")
	}
}

func TestRemaining(t *testing.T) {
	th := throttle.New(time.Second, 3)
	key := "port:80"
	if got := th.Remaining(key); got != 3 {
		t.Fatalf("expected 3 remaining before any call, got %d", got)
	}
	th.Allow(key)
	if got := th.Remaining(key); got != 2 {
		t.Fatalf("expected 2 remaining after one allow, got %d", got)
	}
}

func TestPrune(t *testing.T) {
	th := throttle.New(20*time.Millisecond, 2)
	th.Allow("a")
	th.Allow("b")
	time.Sleep(30 * time.Millisecond)
	// After prune, both keys should be gone; next Allow refills them.
	th.Prune()
	// Remaining should reflect a fresh bucket.
	if got := th.Remaining("a"); got != 2 {
		t.Fatalf("expected 2 remaining after prune, got %d", got)
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	th := throttle.New(time.Second, 1)
	if !th.Allow("x") {
		t.Fatal("key x first call should be allowed")
	}
	if !th.Allow("y") {
		t.Fatal("key y first call should be allowed independently")
	}
	if th.Allow("x") {
		t.Fatal("key x second call should be denied")
	}
}
