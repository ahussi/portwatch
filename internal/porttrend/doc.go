// Package porttrend provides a thread-safe tracker that records how frequently
// each port binding is observed over time.
//
// Use Record to note each observation, Get to retrieve trend data for a
// specific key, and All to iterate over every tracked binding. The clock
// can be replaced via WithClock for deterministic testing.
package porttrend
