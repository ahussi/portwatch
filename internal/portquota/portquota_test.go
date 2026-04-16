package portquota

import (
	"sync"
	"testing"
)

func TestNew_DefaultLimit(t *testing.T) {
	q := New(5)
	if q.default_ != 5 {
		t.Fatalf("expected default 5, got %d", q.default_)
	}
}

func TestNew_FloorAtOne(t *testing.T) {
	q := New(0)
	if q.default_ != 1 {
		t.Fatalf("expected floor 1, got %d", q.default_)
	}
}

func TestTrack_WithinLimit(t *testing.T) {
	q := New(3)
	for _, port := range []uint16{80, 443, 8080} {
		if err := q.Track("svc", port); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if q.Count("svc") != 3 {
		t.Fatalf("expected count 3")
	}
}

func TestTrack_ExceedsLimit(t *testing.T) {
	q := New(2)
	_ = q.Track("svc", 80)
	_ = q.Track("svc", 443)
	err := q.Track("svc", 8080)
	if err == nil {
		t.Fatal("expected ErrQuotaExceeded")
	}
	var qe *ErrQuotaExceeded
	if e, ok := err.(*ErrQuotaExceeded); !ok {
		t.Fatalf("wrong error type: %T", err)
	} else {
		qe = e
	}
	if qe.Subject != "svc" || qe.Limit != 2 {
		t.Fatalf("unexpected quota error: %v", qe)
	}
}

func TestTrack_Idempotent(t *testing.T) {
	q := New(1)
	if err := q.Track("svc", 80); err != nil {
		t.Fatal(err)
	}
	if err := q.Track("svc", 80); err != nil {
		t.Fatalf("duplicate track should be idempotent: %v", err)
	}
}

func TestSetLimit_Override(t *testing.T) {
	q := New(1)
	q.SetLimit("big", 100)
	for i := uint16(1); i <= 100; i++ {
		if err := q.Track("big", i); err != nil {
			t.Fatalf("port %d: %v", i, err)
		}
	}
}

func TestRelease_DecrementsCount(t *testing.T) {
	q := New(1)
	_ = q.Track("svc", 80)
	q.Release("svc", 80)
	if q.Count("svc") != 0 {
		t.Fatal("expected count 0 after release")
	}
	if err := q.Track("svc", 9000); err != nil {
		t.Fatalf("expected slot free after release: %v", err)
	}
}

func TestConcurrentTrack(t *testing.T) {
	q := New(50)
	var wg sync.WaitGroup
	for i := uint16(1); i <= 50; i++ {
		wg.Add(1)
		go func(p uint16) {
			defer wg.Done()
			_ = q.Track("concurrent", p)
		}(i)
	}
	wg.Wait()
	if q.Count("concurrent") > 50 {
		t.Fatal("count exceeded limit under concurrency")
	}
}
