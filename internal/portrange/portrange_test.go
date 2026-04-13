package portrange_test

import (
	"testing"

	"github.com/user/portwatch/internal/portrange"
)

func TestParse_SinglePort(t *testing.T) {
	s, err := portrange.Parse("80")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Contains(80) {
		t.Error("expected set to contain port 80")
	}
	if s.Contains(81) {
		t.Error("expected set not to contain port 81")
	}
}

func TestParse_Range(t *testing.T) {
	s, err := portrange.Parse("8000-8080")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range []uint16{8000, 8040, 8080} {
		if !s.Contains(p) {
			t.Errorf("expected set to contain port %d", p)
		}
	}
	if s.Contains(7999) || s.Contains(8081) {
		t.Error("expected set not to contain ports outside range")
	}
}

func TestParse_CommaSeparated(t *testing.T) {
	s, err := portrange.Parse("22, 80, 443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range []uint16{22, 80, 443} {
		if !s.Contains(p) {
			t.Errorf("expected set to contain port %d", p)
		}
	}
	if s.Contains(8080) {
		t.Error("expected set not to contain port 8080")
	}
}

func TestParse_MixedExpr(t *testing.T) {
	s, err := portrange.Parse("22,8000-8010,443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Contains(22) || !s.Contains(8005) || !s.Contains(443) {
		t.Error("expected set to contain 22, 8005, and 443")
	}
	if s.Contains(8011) {
		t.Error("expected set not to contain port 8011")
	}
}

func TestParse_InvalidPort(t *testing.T) {
	_, err := portrange.Parse("abc")
	if err == nil {
		t.Error("expected error for non-numeric port")
	}
}

func TestParse_ReversedRange(t *testing.T) {
	_, err := portrange.Parse("9000-8000")
	if err == nil {
		t.Error("expected error for reversed range")
	}
}

func TestParse_EmptyExpression(t *testing.T) {
	_, err := portrange.Parse("")
	if err == nil {
		t.Error("expected error for empty expression")
	}
}

func TestRangeString_Single(t *testing.T) {
	r := portrange.Range{Low: 443, High: 443}
	if r.String() != "443" {
		t.Errorf("expected \"443\", got %q", r.String())
	}
}

func TestRangeString_Range(t *testing.T) {
	r := portrange.Range{Low: 8000, High: 8080}
	if r.String() != "8000-8080" {
		t.Errorf("expected \"8000-8080\", got %q", r.String())
	}
}

func TestSet_Ranges(t *testing.T) {
	s, _ := portrange.Parse("80,443")
	ranges := s.Ranges()
	if len(ranges) != 2 {
		t.Fatalf("expected 2 ranges, got %d", len(ranges))
	}
}
