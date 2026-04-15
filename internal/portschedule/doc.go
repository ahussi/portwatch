// Package portschedule allows operators to define time-of-day windows during
// which a port binding is considered expected. Bindings observed outside their
// scheduled window can be escalated as unexpected by the watcher pipeline.
//
// Usage:
//
//	sched := portschedule.New()
//	_ = sched.Add(portschedule.Rule{
//		Port: 8080,
//		Windows: []portschedule.Window{
//			{Start: 9 * time.Hour, End: 17 * time.Hour},
//		},
//	})
//	if !sched.Active(8080) {
//		// raise alert — port bound outside business hours
//	}
package portschedule
