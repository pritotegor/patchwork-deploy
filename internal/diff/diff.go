// Package diff provides utilities for comparing patch state between
// two checkpoint snapshots, identifying which patches are new, already
// applied, or removed since the last deployment run.
package diff

import (
	"fmt"
	"io"
	"os"
)

// Status describes the state of a patch relative to a prior checkpoint.
type Status int

const (
	StatusNew     Status = iota // patch not seen before
	StatusApplied               // patch already applied
	StatusRemoved               // patch was applied but is no longer present
)

// String returns a human-readable label for a Status.
func (s Status) String() string {
	switch s {
	case StatusNew:
		return "new"
	case StatusApplied:
		return "applied"
	case StatusRemoved:
		return "removed"
	default:
		return "unknown"
	}
}

// Entry holds a patch name and its computed status.
type Entry struct {
	Name   string
	Status Status
}

// Differ compares a current patch list against an applied set.
type Differ struct {
	out io.Writer
}

// New returns a Differ that writes log lines to out.
// If out is nil, os.Stdout is used.
func New(out io.Writer) *Differ {
	if out == nil {
		out = os.Stdout
	}
	return &Differ{out: out}
}

// Compute returns one Entry per unique patch name across current and applied.
// current is the ordered list of patch names on disk.
// applied is the set of patch names already recorded in the checkpoint.
func (d *Differ) Compute(current []string, applied map[string]bool) []Entry {
	seen := make(map[string]bool, len(current))
	results := make([]Entry, 0, len(current))

	for _, name := range current {
		seen[name] = true
		status := StatusNew
		if applied[name] {
			status = StatusApplied
		}
		results = append(results, Entry{Name: name, Status: status})
	}

	for name := range applied {
		if !seen[name] {
			results = append(results, Entry{Name: name, Status: StatusRemoved})
		}
	}

	fmt.Fprintf(d.out, "diff: %d current, %d applied, %d entries computed\n",
		len(current), len(applied), len(results))
	return results
}
