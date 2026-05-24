package pipeline

import (
	"fmt"
	"io"
	"os"
)

// Summary writes a human-readable deployment summary for a slice of Results
// to the provided writer. If w is nil it defaults to os.Stdout.
func Summary(results []Result, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}

	totalApplied := 0
	totalSkipped := 0
	failures := 0

	for _, r := range results {
		totalApplied += r.Applied
		totalSkipped += r.Skipped
		if r.Err != nil {
			failures++
		}
	}

	fmt.Fprintf(w, "\n=== Pipeline Summary ===\n")
	fmt.Fprintf(w, "Hosts    : %d\n", len(results))
	fmt.Fprintf(w, "Applied  : %d\n", totalApplied)
	fmt.Fprintf(w, "Skipped  : %d\n", totalSkipped)
	fmt.Fprintf(w, "Failures : %d\n", failures)

	if failures > 0 {
		fmt.Fprintf(w, "\nFailed hosts:\n")
		for _, r := range results {
			if r.Err != nil {
				fmt.Fprintf(w, "  - %s: %v\n", r.Host, r.Err)
			}
		}
	}
}

// HasFailures returns true when at least one Result carries a non-nil error.
func HasFailures(results []Result) bool {
	for _, r := range results {
		if r.Err != nil {
			return true
		}
	}
	return false
}
