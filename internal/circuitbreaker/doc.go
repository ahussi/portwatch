// Package circuitbreaker provides a lightweight circuit breaker used by
// portwatch to prevent alert storms and protect external integrations.
//
// Usage:
//
//	cb := circuitbreaker.New(5, 30*time.Second)
//
//	if err := cb.Allow(); err != nil {
//		// circuit is open — skip the call
//		return err
//	}
//	if err := doSomething(); err != nil {
//		cb.RecordFailure()
//		return err
//	}
//	cb.RecordSuccess()
//
// States:
//
//	Closed   – normal operation; all calls are allowed.
//	Open     – too many failures; calls are rejected until the cooldown elapses.
//	HalfOpen – cooldown elapsed; one probe call is allowed to test recovery.
package circuitbreaker
