package portpolicy

import (
	"testing"
)

func TestNew_DefaultAllow(t *testing.T) {
	p := New(Allow)
	if got := p.Evaluate(9999, "tcp"); got != Allow {
		t.Fatalf("expected Allow, got %s", got)
	}
}

func TestNew_DefaultDeny(t *testing.T) {
	p := New(Deny)
	if got := p.Evaluate(9999, "tcp"); got != Deny {
		t.Fatalf("expected Deny, got %s", got)
	}
}

func TestEvaluate_ExactMatch(t *testing.T) {
	p := New(Deny)
	p.Add(Rule{Port: 80, Protocol: "tcp", Action: Allow})

	if got := p.Evaluate(80, "tcp"); got != Allow {
		t.Errorf("port 80/tcp: expected Allow, got %s", got)
	}
	if got := p.Evaluate(80, "udp"); got != Deny {
		t.Errorf("port 80/udp: expected Deny, got %s", got)
	}
}

func TestEvaluate_WildcardPort(t *testing.T) {
	p := New(Deny)
	p.Add(Rule{Port: 0, Protocol: "tcp", Action: Allow}) // all TCP

	if got := p.Evaluate(443, "tcp"); got != Allow {
		t.Errorf("expected Allow for any tcp port, got %s", got)
	}
	if got := p.Evaluate(443, "udp"); got != Deny {
		t.Errorf("expected Deny for udp, got %s", got)
	}
}

func TestEvaluate_WildcardProtocol(t *testing.T) {
	p := New(Deny)
	p.Add(Rule{Port: 53, Protocol: "", Action: Allow}) // port 53 any proto

	if got "tcp"); got != Allow {
		t.Errorf("expected Allow for 53/tcp, got %s", got)
	}
	if got := p.Evaluate(53, "udp"); got != Allow {
		t.Errorf("expected Allow for 53/udp, got %s", got)
	}
}

func TestEvaluate_FirstRuleWins(t *testing.T) {
	p := New(Allow)
	p.Add(Rule{Port: 22, Protocol: "tcp", Action: Deny})
	p.Add(Rule{Port: 22, Protocol: "tcp", Action: Allow}) // never reached

	if got := p.Evaluate(22, "tcp"); got != Deny {
		t.Errorf("expected first rule (Deny) to win, got %s", got)
	}
}

func TestRules_ReturnsCopy(t *testing.T) {
	p := New(Allow)
	p.Add(Rule{Port: 8080, Protocol: "tcp", Action: Deny})

	r := p.Rules()
	if len(r) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(r))
	}
	r[0].Port = 9999 // mutate copy
	if p.Rules()[0].Port != 8080 {
		t.Error("Rules() should return a copy, not the internal slice")
	}
}

func TestActionString(t *testing.T) {
	if Allow.String() != "allow" {
		t.Errorf("unexpected: %s", Allow.String())
	}
	if Deny.String() != "deny" {
		t.Errorf("unexpected: %s", Deny.String())
	}
}

func TestString(t *testing.T) {
	p := New(Deny)
	p.Add(Rule{Port: 80, Protocol: "tcp", Action: Allow})
	s := p.String()
	if s == "" {
		t.Error("String() should not be empty")
	}
}
