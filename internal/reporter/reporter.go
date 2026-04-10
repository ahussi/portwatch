package reporter

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Format defines the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter renders port binding snapshots to a given writer.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a Reporter writing to out in the given format.
func New(out io.Writer, format Format) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{out: out, format: format}
}

// Render writes the current snapshot bindings to the reporter's writer.
func (r *Reporter) Render(snap *snapshot.Snapshot) error {
	keys := snap.Keys()
	if r.format == FormatJSON {
		return r.renderJSON(snap, keys)
	}
	return r.renderText(snap, keys)
}

func (r *Reporter) renderText(snap *snapshot.Snapshot, keys []string) error {
	w := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "%-25s\t%-20s\t%s\n", "BINDING", "PROCESS", "TIMESTAMP")
	fmt.Fprintf(w, "%-25s\t%-20s\t%s\n", "-------", "-------", "---------")
	for _, k := range keys {
		b, ok := snap.Get(k)
		if !ok {
			continue
		}
		fmt.Fprintf(w, "%-25s\t%-20s\t%s\n", k, b.Process, b.SeenAt.Format(time.RFC3339))
	}
	return w.Flush()
}

func (r *Reporter) renderJSON(snap *snapshot.Snapshot, keys []string) error {
	fmt.Fprintln(r.out, "[")
	for i, k := range keys {
		b, ok := snap.Get(k)
		if !ok {
			continue
		}
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Fprintf(r.out, "  {\"binding\":%q,\"process\":%q,\"seen_at\":%q}%s\n",
			k, b.Process, b.SeenAt.Format(time.RFC3339), comma)
	}
	fmt.Fprintln(r.out, "]")
	return nil
}
