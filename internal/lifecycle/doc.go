// Package lifecycle provides ordered startup and graceful shutdown management
// for portwatch daemon components.
//
// Usage:
//
//	lm := lifecycle.New()
//	lm.OnStart(func(ctx context.Context) error {
//		// initialise component
//		return nil
//	})
//	lm.OnStop(func(ctx context.Context) error {
//		// clean up component
//		return nil
//	})
//	if err := lm.Run(ctx); err != nil {
//		log.Fatal(err)
//	}
//
// Shutdown hooks are executed in reverse registration order so that
// dependencies are torn down cleanly.
package lifecycle
