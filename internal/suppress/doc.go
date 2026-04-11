// Package suppress provides time-bounded suppression of port-binding alerts.
//
// A [Manager] tracks active suppressions keyed by binding key strings
// (e.g. "tcp:8080"). Suppressions expire automatically after a configured
// duration and are pruned lazily on read or explicitly via [Manager.Prune].
//
// Typical usage:
//
//	sm := suppress.New(nil)
//
//	// Suppress alerts for tcp:8080 for the next 10 minutes.
//	sm.Suppress("tcp:8080", 10*time.Minute)
//
//	// Before dispatching an alert, check suppression.
//	if !sm.IsSuppressed(key) {
//	    alertManager.Dispatch(alert)
//	}
//
// The package is safe for concurrent use.
package suppress
