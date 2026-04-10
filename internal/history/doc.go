// Package history provides a thread-safe, fixed-capacity circular buffer
// for recording port-binding change events (additions and removals) observed
// by portwatch during its lifetime.
//
// Events are stored in chronological order. When the buffer is full the
// oldest event is silently overwritten, keeping memory usage bounded
// regardless of how long the daemon runs.
//
// Typical usage:
//
//	rec := history.New(512)
//
//	// inside the watcher loop, after diffing snapshots:
//	for _, b := range added {
//		rec.Add(history.EventAdded, b)
//	}
//	for _, b := range removed {
//		rec.Add(history.EventRemoved, b)
//	}
//
//	// later, dump all recorded events:
//	for _, e := range rec.All() {
//		fmt.Printf("%s  %-8s  %s\n", e.Timestamp.Format(time.RFC3339), e.Kind, e.Binding)
//	}
package history
