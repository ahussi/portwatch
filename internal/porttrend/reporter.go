package porttrend

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
)

// Report writes a human-readable trend summary to w.
func Report(w io.Writer, tr *Tracker) error {
	entries := tr.All()
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Key < entries[j].Key
	})
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tCOUNT\tFIRST SEEN\tLAST SEEN")
	fmt.Fprintln(tw, strings.Repeat("-", 60))
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\n",
			e.Key,
			e.Count,
			e.FirstSeen.Format("2006-01-02 15:04:05"),
			e.LastSeen.Format("2006-01-02 15:04:05"),
		)
	}
	return tw.Flush()
}
