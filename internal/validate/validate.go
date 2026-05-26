// Package validate provides pre-apply patch validation for patchwork-deploy.
// It checks that patch files are well-formed before execution begins.
package validate

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/patchwork-deploy/internal/patch"
)

// Result holds the outcome of validating a single patch.
type Result struct {
	Patch patch.Patch
	OK    bool
	Reason string
}

// Validator checks patches before they are applied.
type Validator struct {
	out io.Writer
}

// New returns a Validator that writes diagnostic output to out.
// If out is nil, os.Stdout is used.
func New(out io.Writer) *Validator {
	if out == nil {
		out = os.Stdout
	}
	return &Validator{out: out}
}

// ValidateAll checks every patch in the slice and returns one Result per patch.
// A patch is considered valid when:
//   - its Path is non-empty
//   - the file exists and is readable
//   - the file is not empty
func (v *Validator) ValidateAll(patches []patch.Patch) []Result {
	results := make([]Result, 0, len(patches))
	for _, p := range patches {
		r := v.validate(p)
		if !r.OK {
			fmt.Fprintf(v.out, "[validate] FAIL %s: %s\n", p.Name, r.Reason)
		} else {
			fmt.Fprintf(v.out, "[validate] OK   %s\n", p.Name)
		}
		results = append(results, r)
	}
	return results
}

// AnyFailed returns true if at least one Result is not OK.
func AnyFailed(results []Result) bool {
	for _, r := range results {
		if !r.OK {
			return true
		}
	}
	return false
}

func (v *Validator) validate(p patch.Patch) Result {
	if strings.TrimSpace(p.Path) == "" {
		return Result{Patch: p, OK: false, Reason: "path is empty"}
	}
	info, err := os.Stat(p.Path)
	if err != nil {
		return Result{Patch: p, OK: false, Reason: fmt.Sprintf("file not accessible: %v", err)}
	}
	if info.Size() == 0 {
		return Result{Patch: p, OK: false, Reason: "file is empty"}
	}
	return Result{Patch: p, OK: true}
}
