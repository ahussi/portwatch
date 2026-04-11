// Package watchlist provides a concurrency-safe registry of ports that
// portwatch should actively monitor.
//
// Entries are keyed by protocol and port number (e.g. "tcp:8080") and
// support an optional human-readable label for reporting purposes.
//
// The watchlist is populated at startup from the loaded configuration and
// may be modified at runtime via the CLI or signal handlers.
package watchlist
