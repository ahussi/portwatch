// Package baseline provides persistence for a known-good set of port bindings.
//
// A Baseline is loaded from or saved to a JSON file on disk. During normal
// operation portwatch compares the live scan results against the baseline to
// distinguish expected bindings from newly appeared ones.
//
// Typical usage:
//
//	b := baseline.New("/var/lib/portwatch/baseline.json")
//	if err := b.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
//		log.Fatal(err)
//	}
//	if !b.Has("tcp:0.0.0.0:8080") {
//		// unexpected binding detected
//	}
package baseline
