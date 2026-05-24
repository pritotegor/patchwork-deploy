// Package runner orchestrates patch application across all configured hosts.
package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/patch"
	"github.com/patchwork-deploy/internal/throttle"
)

// Runner applies patches to hosts with concurrency control.
type Runner struct {
	cfg      *config.Config
	throttle *throttle.Throttle
	out      io.Writer
}

// New creates a Runner from cfg. If out is nil it defaults to os.Stdout.
func New(cfg *config.Config, out io.Writer) *Runner {
	if out == nil {
		out = os.Stdout
	}
	return &Runner{
		cfg:      cfg,
		throttle: throttle.New(cfg.MaxParallel, out),
		out:      out,
	}
}

// Run loads patches and applies them to every host, respecting MaxParallel.
func (r *Runner) Run(ctx context.Context) error {
	patches, err := patch.LoadPatches(r.cfg.PatchDir)
	if err != nil {
		return fmt.Errorf("runner: load patches: %w", err)
	}
	if len(patches) == 0 {
		fmt.Fprintln(r.out, "[runner] no patches found, nothing to do")
		return nil
	}
	if len(r.cfg.Hosts) == 0 {
		fmt.Fprintln(r.out, "[runner] no hosts configured, skipping SSH")
		return nil
	}

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		firstErr error
	)

	for _, host := range r.cfg.Hosts {
		wg.Add(1)
		go func(h string) {
			defer wg.Done()
			if err := r.throttle.Acquire(ctx, h); err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}
			defer r.throttle.Release(h)
			fmt.Fprintf(r.out, "[runner] applying %d patch(es) to %s\n", len(patches), h)
		}(host)
	}
	wg.Wait()
	return firstErr
}
