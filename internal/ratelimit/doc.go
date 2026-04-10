// Package ratelimit implements a key-based cooldown limiter for portwatch.
//
// It is used by the watcher to suppress duplicate alerts when the same
// port binding is detected across multiple consecutive scan intervals.
//
// Usage:
//
//	limiter := ratelimit.New(30 * time.Second)
//
//	// In alert dispatch path:
//	if limiter.Allow(binding.Key()) {
//		alertManager.Dispatch(alert)
//	}
//
// Call Prune() periodically (e.g. once per scan cycle) to release memory
// for bindings that are no longer active.
package ratelimit
