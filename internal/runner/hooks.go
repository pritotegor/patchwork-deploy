package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// HookType represents the stage at which a hook runs.
type HookType string

const (
	HookPrePatch  HookType = "pre-patch"
	HookPostPatch HookType = "post-patch"
)

// Hook defines a shell command to run before or after patching.
type Hook struct {
	Type    HookType
	Command string
}

// HookRunner executes lifecycle hooks and writes output to a writer.
type HookRunner struct {
	out io.Writer
}

// NewHookRunner creates a HookRunner that writes to out.
// If out is nil, os.Stdout is used.
func NewHookRunner(out io.Writer) *HookRunner {
	if out == nil {
		out = os.Stdout
	}
	return &HookRunner{out: out}
}

// Run executes the hook's command via the system shell.
// Output (stdout+stderr) is forwarded to the configured writer.
// Returns an error if the command exits with a non-zero status.
func (h *HookRunner) Run(hook Hook) error {
	if strings.TrimSpace(hook.Command) == "" {
		return fmt.Errorf("hook %q has empty command", hook.Type)
	}

	cmd := exec.Command("sh", "-c", hook.Command)
	cmd.Stdout = h.out
	cmd.Stderr = h.out

	fmt.Fprintf(h.out, "[hook:%s] running: %s\n", hook.Type, hook.Command)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("hook %q failed: %w", hook.Type, err)
	}
	return nil
}

// RunAll executes a slice of hooks in order, stopping on first failure.
func (h *HookRunner) RunAll(hooks []Hook) error {
	for _, hook := range hooks {
		if err := h.Run(hook); err != nil {
			return err
		}
	}
	return nil
}
