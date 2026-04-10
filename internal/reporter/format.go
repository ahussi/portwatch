package reporter

import (
	"fmt"
	"strings"
)

// ParseFormat converts a string to a Format value.
// Returns an error if the string is not a recognised format.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("unknown report format %q: choose 'text' or 'json'", s)
	}
}

// String implements fmt.Stringer.
func (f Format) String() string {
	return string(f)
}

// Formats returns all supported format values.
func Formats() []Format {
	return []Format{FormatText, FormatJSON}
}
