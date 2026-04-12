// Package healthcheck implements a lightweight liveness probe for the
// portwatch daemon.
//
// The scan loop calls Beat() after each successful iteration. An external
// caller (CLI flag, HTTP handler, or signal handler) can call Status() to
// determine whether the daemon is healthy, stale, or has not yet started.
//
// Status values:
//
//	"ok"      — a heartbeat was received within the staleness window
//	"stale"   — no heartbeat within the staleness window
//	"unknown" — Beat has never been called (daemon just started)
package healthcheck
