package patch

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

// Patch represents a single shell patch file to be applied.
type Patch struct {
	Name string
	Path string
}

// ExecuteResult holds the result of applying a single patch.
type ExecuteResult struct {
	Patch  Patch
	Output string
	Err    error
}

// Executor applies patches to a remote host over an SSH client.
type Executor struct {
	Client *ssh.Client
	Out    io.Writer
}

// NewExecutor creates an Executor writing progress to out.
func NewExecutor(client *ssh.Client, out io.Writer) *Executor {
	if out == nil {
		out = os.Stdout
	}
	return &Executor{Client: client, Out: out}
}

// Apply runs each patch in order, stopping on the first error.
func (e *Executor) Apply(patches []Patch) ([]ExecuteResult, error) {
	var results []ExecuteResult

	for _, p := range patches {
		fmt.Fprintf(e.Out, "[patchwork] applying %s\n", p.Name)

		result := e.runPatch(p)
		results = append(results, result)

		if result.Err != nil {
			return results, fmt.Errorf("patch %q failed: %w", p.Name, result.Err)
		}

		fmt.Fprintf(e.Out, "[patchwork] %s ok\n%s", p.Name, result.Output)
	}

	return results, nil
}

func (e *Executor) runPatch(p Patch) ExecuteResult {
	script, err := os.ReadFile(filepath.Clean(p.Path))
	if err != nil {
		return ExecuteResult{Patch: p, Err: fmt.Errorf("read patch: %w", err)}
	}

	session, err := e.Client.NewSession()
	if err != nil {
		return ExecuteResult{Patch: p, Err: fmt.Errorf("new session: %w", err)}
	}
	defer session.Close()

	out, err := session.CombinedOutput(string(script))
	return ExecuteResult{
		Patch:  p,
		Output: string(out),
		Err:    err,
	}
}
