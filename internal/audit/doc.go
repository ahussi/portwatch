// Package audit provides a persistent, append-only audit log for portwatch.
//
// Each binding event (new binding, removed binding, conflict) is serialised as
// a newline-delimited JSON record and appended to the configured log file.
//
// Usage:
//
//	l, err := audit.New("/var/log/portwatch/audit.log")
//	if err != nil { ... }
//	defer l.Close()
//
//	l.Log(audit.Entry{
//		Event:    "new_binding",
//		Port:     8080,
//		Protocol: "tcp",
//		Process:  "nginx",
//		PID:      1234,
//	})
package audit
