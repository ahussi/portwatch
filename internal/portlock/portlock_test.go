package portlock_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/portlock"
)

func entry(port int, proto, owner string) portlock.Entry {
	return portlock.Entry{Port: port, Protocol: proto, Owner: owner, Reason: "test"}
}

func TestNew(t *testing.T) {
	l := portlock.New()
	if l == nil {
		t.Fatal("expected non-nil Locker")
	}
	if l.Len() != 0 {
		t.Fatalf("expected empty locker, got %d entries", l.Len())
	}
}

func TestLockAndIsLocked(t *testing.T) {
	l := portlock.New()
	l.Lock(entry(8080, "tcp", "myapp"))

	if !l.IsLocked("tcp", 8080) {
		t.Error("expected port 8080/tcp to be locked")
	}
	if l.IsLocked("udp", 8080) {
		t.Error("udp:8080 should not be locked")
	}
}

func TestGet(t *testing.T) {
	l := portlock.New()
	l.Lock(entry(443, "tcp", "nginx"))

	e, ok := l.Get("tcp", 443)
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if e.Owner != "nginx" {
		t.Fatalf("expected owner nginx, got %q", e.Owner)
	}
	if e.Key() != "tcp:443" {
		t.Fatalf("unexpected key %q", e.Key())
	}
}

func TestGet_Missing(t *testing.T) {
	l := portlock.New()
	_, ok := l.Get("tcp", 9999)
	if ok {
		t.Error("expected missing entry")
	}
}

func TestUnlock(t *testing.T) {
	l := portlock.New()
	l.Lock(entry(22, "tcp", "sshd"))
	l.Unlock("tcp", 22)

	if l.IsLocked("tcp", 22) {
		t.Error("expected port to be unlocked")
	}
	if l.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", l.Len())
	}
}

func TestUnlock_NoOp(t *testing.T) {
	l := portlock.New()
	l.Unlock("tcp", 1234) // should not panic
}

func TestAll(t *testing.T) {
	l := portlock.New()
	l.Lock(entry(80, "tcp", "httpd"))
	l.Lock(entry(53, "udp", "named"))

	all := l.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestLock_Overwrite(t *testing.T) {
	l := portlock.New()
	l.Lock(entry(8080, "tcp", "old"))
	l.Lock(entry(8080, "tcp", "new"))

	e, _ := l.Get("tcp", 8080)
	if e.Owner != "new" {
		t.Fatalf("expected owner 'new', got %q", e.Owner)
	}
	if l.Len() != 1 {
		t.Fatalf("expected 1 entry after overwrite, got %d", l.Len())
	}
}
