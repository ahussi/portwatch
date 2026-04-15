// Package portpolicy enforces per-port access policies, deciding whether
// a newly observed binding should be permitted or flagged based on
// configurable rules (allow-list, deny-list, and wildcard patterns).
package portpolicy

import (
	"fmt"
	"sync"
)

// Action describes the outcome of a policy evaluation.
type Action int

const (
	Allow Action = iota
	Deny
)

func (a Action) String() string {
	if a == Allow {
		return "allow"
	}
	return "deny"
}

// Rule associates a port (or 0 for wildcard) and protocol with an Action.
type Rule struct {
	Port     int
	Protocol string // "tcp", "udp", or "" for any
	Action   Action
}

// Policy holds an ordered list of rules. The first matching rule wins;
// if no rule matches, the default action is returned.
type Policy struct {
	mu            sync.RWMutex
	rules         []Rule
	defaultAction Action
}

// New creates a Policy with the given default action.
func New(defaultAction Action) *Policy {
	return &Policy{defaultAction: defaultAction}
}

// Add appends a rule to the policy.
func (p *Policy) Add(r Rule) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rules = append(p.rules, r)
}

// Evaluate returns the Action that applies to the given port and protocol.
func (p *Policy) Evaluate(port int, protocol string) Action {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, r := range p.rules {
		if r.Port != 0 && r.Port != port {
			continue
		}
		if r.Protocol != "" && r.Protocol != protocol {
			continue
		}
		return r.Action
	}
	return p.defaultAction
}

// Rules returns a snapshot of the current rule list.
func (p *Policy) Rules() []Rule {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Rule, len(p.rules))
	copy(out, p.rules)
	return out
}

// String returns a human-readable summary of the policy.
func (p *Policy) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return fmt.Sprintf("portpolicy: %d rules, default=%s", len(p.rules), p.defaultAction)
}
