// Package eventbus implements a lightweight publish/subscribe bus used by
// portwatch to decouple event producers (scanner, watcher) from consumers
// (alert manager, audit logger, metrics).
//
// Usage:
//
//	bus := eventbus.New()
//
//	unsub := bus.Subscribe("binding.new", func(e eventbus.Event) {
//		fmt.Println("new binding:", e.Payload)
//	})
//	defer unsub()
//
//	bus.Publish(eventbus.Event{Topic: "binding.new", Payload: someBinding})
//
// Topics used by portwatch:
//
//	"binding.new"     — a previously unseen port binding has appeared
//	"binding.removed" — a known port binding is no longer present
//	"binding.conflict"— two processes are bound to the same port
package eventbus
