package portschedule

import (
	"testing"
	"time"
)

func fixedClock(h, m, s int) func() time.Time {
	return func() time.Time {
		return time.Date(2024, 1, 1, h, m, s, 0, time.UTC)
	}
}

func TestNew(t *testing.T) {
	sched := New()
	if sched == nil {
		t.Fatal("expected non-nil Scheduler")
	}
	if len(sched.Rules()) != 0 {
		t.Fatal("expected empty rules")
	}
}

func TestActive_NoRule_AlwaysTrue(t *testing.T) {
	sched := New()
	if !sched.Active(9090) {
		t.Error("unscheduled port should be active")
	}
}

func TestActive_WithinWindow(t *testing.T) {
	sched := New(WithClock(fixedClock(10, 0, 0)))
	_ = sched.Add(Rule{
		Port:    8080,
		Windows: []Window{{Start: 9 * time.Hour, End: 17 * time.Hour}},
	})
	if !sched.Active(8080) {
		t.Error("expected port to be active at 10:00")
	}
}

func TestActive_OutsideWindow(t *testing.T) {
	sched := New(WithClock(fixedClock(20, 0, 0)))
	_ = sched.Add(Rule{
		Port:    8080,
		Windows: []Window{{Start: 9 * time.Hour, End: 17 * time.Hour}},
	})
	if sched.Active(8080) {
		t.Error("expected port to be inactive at 20:00")
	}
}

func TestActive_MultipleWindows(t *testing.T) {
	sched := New(WithClock(fixedClock(22, 30, 0)))
	_ = sched.Add(Rule{
		Port: 5432,
		Windows: []Window{
			{Start: 8 * time.Hour, End: 12 * time.Hour},
			{Start: 22 * time.Hour, End: 23 * time.Hour},
		},
	})
	if !sched.Active(5432) {
		t.Error("expected port to be active in second window")
	}
}

func TestAdd_InvalidWindow(t *testing.T) {
	sched := New()
	err := sched.Add(Rule{
		Port:    1234,
		Windows: []Window{{Start: 17 * time.Hour, End: 9 * time.Hour}},
	})
	if err == nil {
		t.Error("expected error for invalid window")
	}
}

func TestRemove(t *testing.T) {
	sched := New(WithClock(fixedClock(20, 0, 0)))
	_ = sched.Add(Rule{
		Port:    8080,
		Windows: []Window{{Start: 9 * time.Hour, End: 17 * time.Hour}},
	})
	sched.Remove(8080)
	if !sched.Active(8080) {
		t.Error("removed port should be treated as always active")
	}
}

func TestRules_Snapshot(t *testing.T) {
	sched := New()
	_ = sched.Add(Rule{Port: 80, Windows: []Window{{Start: 0, End: 24 * time.Hour}}})
	_ = sched.Add(Rule{Port: 443, Windows: []Window{{Start: 0, End: 24 * time.Hour}}})
	if len(sched.Rules()) != 2 {
		t.Errorf("expected 2 rules, got %d", len(sched.Rules()))
	}
}

func TestWithClock_Option(t *testing.T) {
	clock := fixedClock(3, 0, 0)
	sched := New(WithClock(clock))
	_ = sched.Add(Rule{
		Port:    9000,
		Windows: []Window{{Start: 2 * time.Hour, End: 4 * time.Hour}},
	})
	if !sched.Active(9000) {
		t.Error("expected port active at 03:00 within 02:00-04:00 window")
	}
}
