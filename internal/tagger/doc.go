// Package tagger provides port-to-label mapping for portwatch.
//
// It combines built-in well-known service names (e.g. port 80 → "http")
// with user-defined tags configured at runtime. Tags are used by the
// reporter and alert subsystems to enrich output with human-readable
// context about what service is expected on a given port.
//
// Usage:
//
//	t := tagger.New()
//	t.Add(8080, "my-api")
//	tags := t.Get(8080) // [{http-alt builtin} {my-api user}]
package tagger
