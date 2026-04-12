package eventbus_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/user/portwatch/internal/eventbus"
)

func TestIntegration_SubscribePublishUnsubscribe(t *testing.T) {
	bus := eventbus.New()

	var wg sync.WaitGroup
	var total atomic.Int64

	const workers = 10
	const eventsEach = 20

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			unsub := bus.Subscribe("binding.new", func(e eventbus.Event) {
				total.Add(1)
			})
			for j := 0; j < eventsEach; j++ {
				bus.Publish(eventbus.Event{Topic: "binding.new", Payload: j})
			}
			unsub()
		}())
	}
	wg.Wait()

	// After all goroutines unsubscribe, further publishes should reach no one.
	before := total.Load()
	bus.Publish(eventbus.Event{Topic: "binding.new", Payload: "after"})
	if total.Load() != before {
		t.Errorf("expected no handlers after unsubscribe, but count changed")
	}
}

func TestIntegration_MultiTopicIsolation(t *testing.T) {
	bus := eventbus.New()

	var newCount, removedCount atomic.Int64

	bus.Subscribe("binding.new", func(e eventbus.Event) { newCount.Add(1) })
	bus.Subscribe("binding.removed", func(e eventbus.Event) { removedCount.Add(1) })

	for i := 0; i < 5; i++ {
		bus.Publish(eventbus.Event{Topic: "binding.new"})
	}
	for i := 0; i < 3; i++ {
		bus.Publish(eventbus.Event{Topic: "binding.removed"})
	}

	if newCount.Load() != 5 {
		t.Errorf("binding.new: expected 5, got %d", newCount.Load())
	}
	if removedCount.Load() != 3 {
		t.Errorf("binding.removed: expected 3, got %d", removedCount.Load())
	}
}
