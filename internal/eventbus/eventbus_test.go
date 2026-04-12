package eventbus

import (
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	b := New()
	if b == nil {
		t.Fatal("expected non-nil Bus")
	}
	if len(b.subscribers) != 0 {
		t.Errorf("expected empty subscribers, got %d", len(b.subscribers))
	}
}

func TestSubscribeAndPublish(t *testing.T) {
	b := New()
	var received []Event

	b.Subscribe("test.topic", func(e Event) {
		received = append(received, e)
	})

	b.Publish(Event{Topic: "test.topic", Payload: "hello"})
	b.Publish(Event{Topic: "test.topic", Payload: "world"})

	if len(received) != 2 {
		t.Fatalf("expected 2 events, got %d", len(received))
	}
	if received[0].Payload != "hello" || received[1].Payload != "world" {
		t.Errorf("unexpected payloads: %v", received)
	}
}

func TestPublishNoSubscribers(t *testing.T) {
	b := New()
	// Should not panic.
	b.Publish(Event{Topic: "unknown", Payload: nil})
}

func TestUnsubscribe(t *testing.T) {
	b := New()
	count := 0

	unsub := b.Subscribe("ev", func(e Event) { count++ })
	b.Publish(Event{Topic: "ev"})
	unsub()
	b.Publish(Event{Topic: "ev"})

	if count != 1 {
		t.Errorf("expected count 1 after unsubscribe, got %d", count)
	}
}

func TestMultipleSubscribers(t *testing.T) {
	b := New()
	calls := make([]int, 3)

	for i := range calls {
		i := i
		b.Subscribe("multi", func(e Event) { calls[i]++ })
	}
	b.Publish(Event{Topic: "multi"})

	for i, c := range calls {
		if c != 1 {
			t.Errorf("handler %d: expected 1 call, got %d", i, c)
		}
	}
}

func TestTopics(t *testing.T) {
	b := New()
	b.Subscribe("a", func(Event) {})
	b.Subscribe("b", func(Event) {})

	topics := b.Topics()
	if len(topics) != 2 {
		t.Errorf("expected 2 topics, got %d: %v", len(topics), topics)
	}
}

func TestConcurrentPublish(t *testing.T) {
	b := New()
	var mu sync.Mutex
	var count int

	b.Subscribe("concurrent", func(e Event) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Publish(Event{Topic: "concurrent", Payload: 1})
		}()
	}
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if count != 50 {
		t.Errorf("expected 50 events, got %d", count)
	}
}
