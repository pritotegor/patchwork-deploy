// Package runner orchestrates the full deployment pipeline:
// loading patches, connecting via SSH, and applying each patch in order.
package runner

import (
	"fmt"
	"io"
	"os"

	"github.com/yourorg/patchwork-deploy/internal/config"
	"github.com/yourorg/patchwork-deploy/internal/patch"
	"github.com/yourorg/patchwork-deploy/internal/ssh"
)

// Runner holds the dependencies needed to execute a deployment.
type Runner struct {
	Cfg    *config.Config
	Out    io.Writer
	ErrOut io.Writer
}

// New creates a Runner with the given config, defaulting output to stdout/stderr.
func New(cfg *config.Config) *Runner {
	return &Runner{
		Cfg:    cfg,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}

// Run loads patches, connects to each host, and applies every patch in order.
// It returns the first error encountered.
func (r *Runner) Run() error {
	patches, err := patch.LoadPatches(r.Cfg.PatchDir)
	if err != nil {
		return fmt.Errorf("loading patches: %w", err)
	}

	if len(patches) == 0 {
		fmt.Fprintln(r.Out, "no patches found, nothing to do")
		return nil
	}

	fmt.Fprintf(r.Out, "found %d patch(es) to apply\n", len(patches))

	for _, host := range r.Cfg.Hosts {
		if err := r.applyToHost(host, patches); err != nil {
			return fmt.Errorf("host %s: %w", host, err)
		}
	}
	return nil
}

func (r *Runner) applyToHost(host string, patches []string) error {
	fmt.Fprintf(r.Out, "connecting to %s\n", host)

	client, err := ssh.Connect(host, r.Cfg)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer client.Close()

	exec := patch.NewExecutor(client, r.Out)

	for _, p := range patches {
		fmt.Fprintf(r.Out, "  applying %s\n", p)
		result, err := exec.Apply(p)
		if err != nil {
			return fmt.Errorf("apply %s: %w", p, err)
		}
		if result.ExitCode != 0 {
			return fmt.Errorf("patch %s exited with code %d: %s", p, result.ExitCode, result.Stderr)
		}
	}
	return nil
}
