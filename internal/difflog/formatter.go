package difflog

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// FormatText renders events as human-readable lines.
// Each line follows the pattern:
//
//	2006-01-02T15:04:05Z  ADDED   tcp:8080  pid=1234  nginx
func FormatText(events []Event) string {
	if len(events) == 0 {
		return "(no events)"
	}
	var sb strings.Builder
	for _, e := range events {
		ts := e.Timestamp.UTC().Format(time.RFC3339)
		kind := strings.ToUpper(string(e.Kind))
		pid := ""
		if e.PID > 0 {
			pid = fmt.Sprintf(" pid=%d", e.PID)
		}
		proc := ""
		if e.Process != "" {
			proc = "  " + e.Process
		}
		fmt.Fprintf(&sb, "%s  %-7s %s%s%s\n", ts, kind, e.Key, pid, proc)
	}
	return sb.String()
}

// jsonEvent is the wire representation used by FormatJSON.
type jsonEvent struct {
	Kind      EventKind `json:"kind"`
	Key       string    `json:"key"`
	Port      int       `json:"port"`
	Proto     string    `json:"proto,omitempty"`
	PID       int       `json:"pid,omitempty"`
	Process   string    `json:"process,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// FormatJSON serialises events to a JSON array.
func FormatJSON(events []Event) (string, error) {
	out := make([]jsonEvent, len(events))
	for i, e := range events {
		out[i] = jsonEvent{
			Kind:      e.Kind,
			Key:       e.Key,
			Port:      e.Port,
			Proto:     e.Proto,
			PID:       e.PID,
			Process:   e.Process,
			Timestamp: e.Timestamp,
		}
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", fmt.Errorf("difflog: json marshal: %w", err)
	}
	return string(b), nil
}
