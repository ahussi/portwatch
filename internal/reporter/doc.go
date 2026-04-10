// Package reporter provides rendering of port-binding snapshots to
// human-readable (text table) or machine-readable (JSON) output.
//
// Usage:
//
//	r, err := reporter.ParseFormat(cfg.ReportFormat)
//	if err != nil { log.Fatal(err) }
//
//	rep := reporter.New(os.Stdout, r)
//	if err := rep.Render(snap); err != nil { log.Fatal(err) }
//
// Supported formats are "text" (default) and "json".
package reporter
