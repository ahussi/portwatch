// Package portpolicy provides a lightweight, ordered-rule engine for
// deciding whether a port binding should be allowed or denied.
//
// Rules are evaluated in insertion order; the first matching rule wins.
// If no rule matches, the policy's default action is applied.
//
// Example:
//
//	p := portpolicy.New(portpolicy.Deny)   // deny-by-default
//	p.Add(portpolicy.Rule{Port: 80,  Protocol: "tcp", Action: portpolicy.Allow})
//	p.Add(portpolicy.Rule{Port: 443, Protocol: "tcp", Action: portpolicy.Allow})
//
//	action := p.Evaluate(80, "tcp")  // Allow
//	action  = p.Evaluate(22, "tcp")  // Deny  (default)
package portpolicy
