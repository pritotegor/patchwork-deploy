// Package filter provides utilities for selecting a subset of patches
// to apply based on tag labels or name prefix patterns.
package filter

import (
	"strings"
)

// Patch represents a minimal patch descriptor used for filtering.
type Patch struct {
	Name string
	Tags []string
}

// Options controls which patches are selected.
type Options struct {
	// Tags, when non-empty, keeps only patches that carry at least one
	// of the listed tags.
	Tags []string

	// Prefix, when non-empty, keeps only patches whose Name starts with
	// the given string.
	Prefix string
}

// Apply returns the subset of patches that match all non-empty criteria
// in opts. If opts is zero-valued every patch is returned unchanged.
func Apply(patches []Patch, opts Options) []Patch {
	var out []Patch
	for _, p := range patches {
		if opts.Prefix != "" && !strings.HasPrefix(p.Name, opts.Prefix) {
			continue
		}
		if len(opts.Tags) > 0 && !hasAnyTag(p.Tags, opts.Tags) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// hasAnyTag reports whether patch tags contain at least one of want.
func hasAnyTag(patchTags, want []string) bool {
	set := make(map[string]struct{}, len(patchTags))
	for _, t := range patchTags {
		set[t] = struct{}{}
	}
	for _, w := range want {
		if _, ok := set[w]; ok {
			return true
		}
	}
	return false
}
