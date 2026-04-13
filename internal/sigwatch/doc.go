// Package sigwatch wraps OS signal handling for the portwatch daemon.
//
// It supports two signal classes:
//
//   - Shutdown signals (SIGINT, SIGTERM): trigger a graceful shutdown
//     callback and cause Run to return.
//
//   - Reload signal (SIGHUP): trigger a configuration reload callback
//     without stopping the daemon.
//
// Example usage:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	h := sigwatch.New(
//		sigwatch.WithShutdown(cancel),
//		sigwatch.WithReload(cfg.Reload),
//	)
//	go h.Run(ctx)
package sigwatch
