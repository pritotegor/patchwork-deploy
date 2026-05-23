// Package reporter provides structured result reporting for patch execution runs.
package reporter

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Result holds the outcome of applying a single patch to a single host.
type Result struct {
	Host      string
	Patch     string
	Output    string
	Err       error
	Duration  time.Duration
}

// Reporter writes a human-readable summary of patch execution results.
type Reporter struct {
	out io.Writer
}

// New returns a Reporter that writes to out. If out is nil, os.Stdout is used.
func New(out io.Writer) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{out: out}
}

// Print writes a single result line to the reporter's writer.
func (r *Reporter) Print(res Result) {
	status := "OK"
	if res.Err != nil {
		status = fmt.Sprintf("FAIL: %v", res.Err)
	}
	fmt.Fprintf(r.out, "[%s] %-30s %-20s %s (%s)\n",
		status,
		res.Host,
		res.Patch,
		res.Output,
		res.Duration.Round(time.Millisecond),
	)
}

// Summary writes an aggregate summary line given a slice of results.
func (r *Reporter) Summary(results []Result) {
	total := len(results)
	failed := 0
	for _, res := range results {
		if res.Err != nil {
			failed++
		}
	}
	fmt.Fprintf(r.out, "\nSummary: %d/%d patches succeeded.\n", total-failed, total)
}
