package reporter_test

import (
	"testing"

	"github.com/user/portwatch/internal/reporter"
)

func TestParseFormat(t *testing.T) {
	cases := []struct {
		input   string
		want    reporter.Format
		wantErr bool
	}{
		{"text", reporter.FormatText, false},
		{"TEXT", reporter.FormatText, false},
		{"", reporter.FormatText, false},
		{"json", reporter.FormatJSON, false},
		{"JSON", reporter.FormatJSON, false},
		{"xml", "", true},
		{"csv", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := reporter.ParseFormat(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestFormatString(t *testing.T) {
	if reporter.FormatText.String() != "text" {
		t.Errorf("expected 'text', got %q", reporter.FormatText.String())
	}
	if reporter.FormatJSON.String() != "json" {
		t.Errorf("expected 'json', got %q", reporter.FormatJSON.String())
	}
}

func TestFormats(t *testing.T) {
	fmts := reporter.Formats()
	if len(fmts) != 2 {
		t.Errorf("expected 2 formats, got %d", len(fmts))
	}
}
