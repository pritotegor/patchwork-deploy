// Package dryrun provides a dry-run executor that simulates patch application
// without executing commands on remote hosts.
package dryrun

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/patchwork-deploy/internal/patch"
)

// Result holds the outcome of a simulated patch application.
type Result struct {
	PatchName string
	Host      string
	Simulated bool
	Skipped   bool
	Reason    string
}

// Runner simulates applying patches without executing them.
type Runner struct {
	out io.Writer
}

// New returns a new dry-run Runner. If out is nil, os.Stdout is used.
func New(out io.Writer) *Runner {
	if out == nil {
		out = os.Stdout
	}
	return &Runner{out: out}
}

// Simulate prints what would be applied for each patch and host pair
// without making any real changes. It returns one Result per combination.
func (r *Runner) Simulate(patches []patch.Patch, hosts []string) []Result {
	var results []Result

	for _, p := range patches {
		for _, host := range hosts {
			res := Result{
				PatchName: p.Name,
				Host:      host,
				Simulated: true,
			}
			fmt.Fprintf(r.out, "[dry-run] %s  patch=%s  host=%s\n",
				time.Now().UTC().Format(time.RFC3339), p.Name, host)
			results = append(results, res)
		}
	}

	if len(results) == 0 {
		fmt.Fprintf(r.out, "[dry-run] no patches to simulate\n")
	}

	return results
}

// Skip records a patch that would be skipped (e.g. already applied) and
// writes a log line to the runner's output.
func (r *Runner) Skip(patchName, host, reason string) Result {
	fmt.Fprintf(r.out, "[dry-run] SKIP  patch=%s  host=%s  reason=%s\n", patchName, host, reason)
	return Result{
		PatchName: patchName,
		Host:      host,
		Skipped:   true,
		Reason:    reason,
	}
}
