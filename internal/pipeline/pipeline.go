// Package pipeline coordinates the full deployment sequence for a single host:
// lock acquisition, scheduling check, patch filtering, execution, auditing,
// checkpointing and rollback registration.
package pipeline

import (
	"fmt"
	"io"
	"os"

	"github.com/yourorg/patchwork-deploy/internal/audit"
	"github.com/yourorg/patchwork-deploy/internal/checkpoint"
	"github.com/yourorg/patchwork-deploy/internal/lock"
	"github.com/yourorg/patchwork-deploy/internal/patch"
	"github.com/yourorg/patchwork-deploy/internal/rollback"
)

// Options holds the dependencies required to run a pipeline.
type Options struct {
	Host       string
	Patches    []patch.Patch
	Executor   *patch.Executor
	Auditor    *audit.Auditor
	Checkpoint *checkpoint.Checkpoint
	Rollback   *rollback.Rollback
	Lock       *lock.Lock
	Out        io.Writer
}

// Result captures the outcome of a single host pipeline run.
type Result struct {
	Host    string
	Applied int
	Skipped int
	Err     error
}

// Run executes the deployment pipeline for one host.
func Run(opts Options) Result {
	out := opts.Out
	if out == nil {
		out = os.Stdout
	}

	if err := opts.Lock.Acquire(); err != nil {
		return Result{Host: opts.Host, Err: fmt.Errorf("lock: %w", err)}
	}
	defer opts.Lock.Release() //nolint:errcheck

	var applied, skipped int

	for _, p := range opts.Patches {
		if opts.Checkpoint.Applied(p.Name) {
			fmt.Fprintf(out, "[%s] skip (already applied): %s\n", opts.Host, p.Name)
			skipped++
			continue
		}

		res := opts.Executor.Apply(p)
		opts.Auditor.LogPatch(opts.Host, p.Name, res.Output, res.Err)
		opts.Rollback.Record(p)

		if res.Err != nil {
			fmt.Fprintf(out, "[%s] FAIL: %s — %v\n", opts.Host, p.Name, res.Err)
			return Result{Host: opts.Host, Applied: applied, Skipped: skipped, Err: res.Err}
		}

		if err := opts.Checkpoint.Mark(p.Name); err != nil {
			fmt.Fprintf(out, "[%s] warn: checkpoint failed for %s: %v\n", opts.Host, p.Name, err)
		}

		fmt.Fprintf(out, "[%s] ok: %s\n", opts.Host, p.Name)
		applied++
	}

	return Result{Host: opts.Host, Applied: applied, Skipped: skipped}
}
