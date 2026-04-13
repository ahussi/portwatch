package cooldown

import (
	"testing"
	"time"
)

func TestNew_DefaultDuration(t *testing.T) {
	c := New(0)
	if c.duration != time.Second {
		t.Fatalf("expected 1s default, got %v", c.duration)
	}
}

func TestNew_CustomDuration(t *testing.T) {
	c := New(5 * time.Second)
	if c.duration != 5*time.Second {
		t.Fatalf("expected 5s, got %v", c.duration)
	}
}

func TestReady_FirstTime(t *testing.T) {
	c := New(time.Minute)
	if !c.Ready("port:8080") {
		t.Fatal("expected Ready=true for unseen key")
	}
}

func TestReady_AfterRecord_WithinCooldown(t *testing.T) {
	c := New(time.Minute)
	c.Record("port:8080")
	if c.Ready("port:8080") {
		t.Fatal("expected Ready=false within cooldown")
	}
}

func TestReady_AfterCooldownExpires(t *testing.T) {
	c := New(10 * time.Millisecond)
	c.Record("port:9090")
	time.Sleep(20 * time.Millisecond)
	if !c.Ready("port:9090") {
		t.Fatal("expected Ready=true after cooldown expired")
	}
}

func TestReadyAndRecord_FirstCall(t *testing.T) {
	c := New(time.Minute)
	if !c.ReadyAndRecord("k1") {
		t.Fatal("expected true on first call")
	}
}

func TestReadyAndRecord_SecondCallBlocked(t *testing.T) {
	c := New(time.Minute)
	c.ReadyAndRecord("k1")
	if c.ReadyAndRecord("k1") {
		t.Fatal("expected false on second call within cooldown")
	}
}

func TestReadyAndRecord_AfterExpiry(t *testing.T) {
	c := New(10 * time.Millisecond)
	c.ReadyAndRecord("k2")
	time.Sleep(20 * time.Millisecond)
	if !c.ReadyAndRecord("k2") {
		t.Fatal("expected true after cooldown expired")
	}
}

func TestReset(t *testing.T) {
	c := New(time.Minute)
	c.Record("port:443")
	c.Reset("port:443")
	if !c.Ready("port:443") {
		t.Fatal("expected Ready=true after Reset")
	}
}

func TestPrune_RemovesExpiredEntries(t *testing.T) {
	c := New(10 * time.Millisecond)
	c.Record("old")
	c.Record("fresh")
	time.Sleep(20 * time.Millisecond)
	// refresh "fresh" so it stays within cooldown
	c.Record("fresh")
	c.Prune()
	c.mu.Lock()
	_, hasOld := c.last["old"]
	_, hasFresh := c.last["fresh"]
	c.mu.Unlock()
	if hasOld {
		t.Error("expected 'old' to be pruned")
	}
	if !hasFresh {
		t.Error("expected 'fresh' to remain")
	}
}

func TestDifferentKeys_Independent(t *testing.T) {
	c := New(time.Minute)
	c.Record("a")
	if !c.Ready("b") {
		t.Fatal("key 'b' should be unaffected by key 'a'")
	}
}
