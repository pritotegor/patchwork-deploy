// Package rollback provides mechanisms to track applied patches
// and revert them in reverse order when a deployment fails.
package rollback

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Entry records a single applied patch and its optional undo script path.
type Entry struct {
	PatchName string
	UndoPath  string
}

// Rollbacker tracks applied patches and executes undo steps in reverse order.
type Rollbacker struct {
	out     io.Writer
	entries []Entry
}

// New creates a Rollbacker that writes status messages to out.
// If out is nil, os.Stdout is used.
func New(out io.Writer) *Rollbacker {
	if out == nil {
		out = os.Stdout
	}
	return &Rollbacker{out: out}
}

// Record adds a successfully applied patch to the rollback stack.
func (r *Rollbacker) Record(patchName, undoPath string) {
	r.entries = append(r.entries, Entry{
		PatchName: patchName,
		UndoPath:  undoPath,
	})
}

// Entries returns a copy of the recorded entries in application order.
func (r *Rollbacker) Entries() []Entry {
	out := make([]Entry, len(r.entries))
	copy(out, r.entries)
	return out
}

// Plan returns the list of undo steps that would be executed, in reverse order.
func (r *Rollbacker) Plan() []string {
	steps := make([]string, 0, len(r.entries))
	for i := len(r.entries) - 1; i >= 0; i-- {
		e := r.entries[i]
		if strings.TrimSpace(e.UndoPath) != "" {
			steps = append(steps, fmt.Sprintf("undo %s via %s", e.PatchName, e.UndoPath))
		} else {
			steps = append(steps, fmt.Sprintf("undo %s (no undo script)", e.PatchName))
		}
	}
	return steps
}

// Execute prints the rollback plan to the writer. Actual remote execution
// is delegated to the caller via the provided exec function.
func (r *Rollbacker) Execute(exec func(undoPath string) error) error {
	for i := len(r.entries) - 1; i >= 0; i-- {
		e := r.entries[i]
		fmt.Fprintf(r.out, "[rollback] reverting %s\n", e.PatchName)
		if strings.TrimSpace(e.UndoPath) == "" {
			fmt.Fprintf(r.out, "[rollback] no undo script for %s, skipping\n", e.PatchName)
			continue
		}
		if err := exec(e.UndoPath); err != nil {
			return fmt.Errorf("rollback failed at %s: %w", e.PatchName, err)
		}
	}
	return nil
}
