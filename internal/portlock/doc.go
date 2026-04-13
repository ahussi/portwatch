// Package portlock provides a thread-safe registry of operator-reserved ports.
//
// Ports added to the Locker are considered "owned" — any process that binds
// to them without matching the expected owner should be flagged immediately.
//
// Typical usage:
//
//	l := portlock.New()
//	l.Lock(portlock.Entry{
//		Port:     443,
//		Protocol: "tcp",
//		Owner:    "nginx",
//		Reason:   "production HTTPS listener",
//	})
//
//	if l.IsLocked("tcp", 443) {
//		e, _ := l.Get("tcp", 443)
//		fmt.Println(e.Owner) // "nginx"
//	}
package portlock
