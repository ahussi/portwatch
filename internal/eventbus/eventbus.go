// Package eventbus provides a simple publish/subscribe mechanism for
// broadcasting port binding events to multiple subscribers.
package eventbus

import (
	"sync"
)

// Event represents a port binding event published on the bus.
type Event struct {
	Topic   string
	Payload any
}

// Handler is a function that receives an Event.
type Handler func(Event)

// Bus is a thread-safe publish/subscribe event bus.
type Bus struct {
	mu          sync.RWMutex
	subscribers map[string][]Handler
}

// New creates and returns a new Bus.
func New() *Bus {
	return &Bus{
		subscribers: make(map[string][]Handler),
	}
}

// Subscribe registers a handler for the given topic.
// Returns an unsubscribe function that removes the handler.
func (b *Bus) Subscribe(topic string, h Handler) func() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers[topic] = append(b.subscribers[topic], h)
	idx := len(b.subscribers[topic]) - 1

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		slice := b.subscribers[topic]
		if idx < len(slice) {
			b.subscribers[topic] = append(slice[:idx], slice[idx+1:]...)
		}
	}
}

// Publish sends an Event to all handlers subscribed to the event's Topic.
// Each handler is invoked synchronously in the order it was registered.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	handlers := make([]Handler, len(b.subscribers[e.Topic]))
	copy(handlers, b.subscribers[e.Topic])
	b.mu.RUnlock()

	for _, h := range handlers {
		h(e)
	}
}

// Topics returns a snapshot of all topics that have at least one subscriber.
func (b *Bus) Topics() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	topics := make([]string, 0, len(b.subscribers))
	for t, hs := range b.subscribers {
		if len(hs) > 0 {
			topics = append(topics, t)
		}
	}
	return topics
}
