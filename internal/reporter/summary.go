package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// RunSummary holds aggregate results from a full deployment run.
type RunSummary struct {
	Host        string
	Total       int
	Succeeded   int
	Failed      int
	Skipped     int
	Elapsed     time.Duration
	PatchErrors []PatchError
}

// PatchError pairs a patch name with the error it produced.
type PatchError struct {
	Patch string
	Err   error
}

// SummaryReporter writes a human-readable deployment summary.
type SummaryReporter struct {
	out io.Writer
}

// NewSummaryReporter returns a SummaryReporter writing to out.
// If out is nil, os.Stdout is used.
func NewSummaryReporter(out io.Writer) *SummaryReporter {
	if out == nil {
		out = os.Stdout
	}
	return &SummaryReporter{out: out}
}

// Write prints the full RunSummary to the configured writer.
func (s *SummaryReporter) Write(rs RunSummary) {
	fmt.Fprintf(s.out, "\n%s\n", strings.Repeat("─", 48))
	fmt.Fprintf(s.out, "Deployment summary for host: %s\n", rs.Host)
	fmt.Fprintf(s.out, "%s\n", strings.Repeat("─", 48))
	fmt.Fprintf(s.out, "  Total patches : %d\n", rs.Total)
	fmt.Fprintf(s.out, "  Succeeded     : %d\n", rs.Succeeded)
	fmt.Fprintf(s.out, "  Failed        : %d\n", rs.Failed)
	fmt.Fprintf(s.out, "  Skipped       : %d\n", rs.Skipped)
	fmt.Fprintf(s.out, "  Elapsed       : %s\n", rs.Elapsed.Round(time.Millisecond))

	if len(rs.PatchErrors) > 0 {
		fmt.Fprintf(s.out, "\nFailed patches:\n")
		for _, pe := range rs.PatchErrors {
			fmt.Fprintf(s.out, "  ✗ %s: %v\n", pe.Patch, pe.Err)
		}
	}

	status := "SUCCESS"
	if rs.Failed > 0 {
		status = "FAILURE"
	}
	fmt.Fprintf(s.out, "\nResult: %s\n", status)
	fmt.Fprintf(s.out, "%s\n", strings.Repeat("─", 48))
}
