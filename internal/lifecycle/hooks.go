package lifecycle

import (
	"context"
	"fmt"
	"io"
	"time"
)

// LogHook returns a Hook that writes a timestamped message to w.
func LogHook(w io.Writer, msg string) Hook {
	return func(_ context.Context) error {
		_, err := fmt.Fprintf(w, "%s %s\n", time.Now().Format(time.RFC3339), msg)
		return err
	}
}

// TimeoutHook wraps h and enforces a maximum execution duration.
func TimeoutHook(h Hook, d time.Duration) Hook {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, d)
		defer cancel()

		result := make(chan error, 1)
		go func() { result <- h(ctx) }()

		select {
		case err := <-result:
			return err
		case <-ctx.Done():
			return fmt.Errorf("hook timed out after %s", d)
		}
	}
}

// NoopHook is a Hook that always succeeds without side effects.
func NoopHook(_ context.Context) error { return nil }
