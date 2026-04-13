package difflog

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestConcurrentAdd(t *testing.T) {
	l := New(512)
	var wg sync.WaitGroup
	workers := 8
	perWorker := 50

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < perWorker; i++ {
				l.Add(Event{
					Kind: KindAdded,
					Key:  fmt.Sprintf("worker-%d-port-%d", id, i),
					Port: id*1000 + i,
				})
			}
		}(w)
	}
	wg.Wait()

	if l.Len() != workers*perWorker {
		t.Errorf("expected %d events, got %d", workers*perWorker, l.Len())
	}
}

func TestConcurrentAddEviction(t *testing.T) {
	cap := 64
	l := New(cap)
	var wg sync.WaitGroup

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			l.Add(Event{Kind: KindRemoved, Port: n})
		}(i)
	}
	wg.Wait()

	if l.Len() > cap {
		t.Errorf("log exceeded capacity: len=%d cap=%d", l.Len(), cap)
	}
}

func TestSince_Concurrent(t *testing.T) {
	l := New(256)
	base := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			offset := time.Duration(n)*10*time.Millisecond - 200*time.Millisecond
			l.Add(Event{Kind: KindAdded, Port: n, Timestamp: base.Add(offset)})
		}(i)
	}
	wg.Wait()

	results := l.Since(base)
	for _, e := range results {
		if e.Timestamp.Before(base) {
			t.Errorf("Since returned event before cutoff: %v", e.Timestamp)
		}
	}
}
