package portschedule_test

import (
	"sync"
	"testing"
	"time"

	"portwatch/internal/portschedule"
)

func TestConcurrentAddAndActive(t *testing.T) {
	sched := portschedule.New(portschedule.WithClock(func() time.Time {
		return time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	}))

	var wg sync.WaitGroup
	for port := 1; port <= 50; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			_ = sched.Add(portschedule.Rule{
				Port:    p,
				Windows: []portschedule.Window{{Start: 9 * time.Hour, End: 17 * time.Hour}},
			})
		}(port)
	}
	wg.Wait()

	for port := 1; port <= 50; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			_ = sched.Active(p)
		}(port)
	}
	wg.Wait()
}

func TestConcurrentRemoveAndActive(t *testing.T) {
	sched := portschedule.New()
	for port := 1; port <= 20; port++ {
		_ = sched.Add(portschedule.Rule{
			Port:    port,
			Windows: []portschedule.Window{{Start: 0, End: 24 * time.Hour}},
		})
	}

	var wg sync.WaitGroup
	for port := 1; port <= 20; port++ {
		wg.Add(2)
		go func(p int) {
			defer wg.Done()
			sched.Remove(p)
		}(port)
		go func(p int) {
			defer wg.Done()
			_ = sched.Active(p)
		}(port)
	}
	wg.Wait()
}
