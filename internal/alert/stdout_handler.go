package alert

import (
	"fmt"
	"io"
	"os"
)

// StdoutHandler writes alerts to an io.Writer (defaults to os.Stdout).
type StdoutHandler struct {
	writer io.Writer
}

// NewStdoutHandler creates a StdoutHandler that writes to w.
// If w is nil, os.Stdout is used.
func NewStdoutHandler(w io.Writer) *StdoutHandler {
	if w == nil {
		w = os.Stdout
	}
	return &StdoutHandler{writer: w}
}

// Send writes the alert as a formatted line to the configured writer.
func (h *StdoutHandler) Send(a Alert) error {
	_, err := fmt.Fprintln(h.writer, a.String())
	return err
}
