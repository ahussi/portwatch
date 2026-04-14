// Package portevict tracks port eviction events for portwatch.
//
// An eviction occurs when a port binding that was present in a previous
// scan cycle is absent from the current scan — meaning the process that
// held the port has released or closed it.
//
// Tracker is safe for concurrent use. Events are stored in a bounded
// ring-buffer; when the buffer is full the oldest event is dropped to
// make room for the newest.
//
// Typical usage:
//
//	tracker := portevict.New(512)
//
//	// After diffing snapshots, record disappeared bindings:
//	tracker.Record(portevict.Event{
//		Key:   "tcp:0.0.0.0:8080",
//		Port:  8080,
//		Proto: "tcp",
//		Addr:  "0.0.0.0",
//	})
//
//	// Query events from the last minute:
//	events := tracker.Since(time.Now().Add(-time.Minute))
package portevict
