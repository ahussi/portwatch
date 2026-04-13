// Package difflog provides a bounded, thread-safe event log for recording
// port-binding diff events detected between successive scans.
//
// Each Event captures the kind of change (added/removed), the affected port,
// protocol, owning process, and the wall-clock time of detection.
//
// The Log evicts the oldest entry when its capacity is reached, making it
// suitable for long-running daemon use without unbounded memory growth.
package difflog
