// Package filter provides patch-selection helpers for patchwork-deploy.
//
// Patches can be filtered by name prefix or by one or more tag labels
// before being handed to the executor. This allows operators to run only
// a targeted slice of the patch sequence without modifying the patch
// directory itself.
//
// Example usage:
//
//	patches := []filter.Patch{
//		{Name: "001-init.sh", Tags: []string{"db"}},
//		{Name: "002-migrate.sh", Tags: []string{"db", "schema"}},
//		{Name: "003-cache.sh", Tags: []string{"cache"}},
//	}
//
//	selected := filter.Apply(patches, filter.Options{Tags: []string{"db"}})
//	// selected contains 001-init.sh and 002-migrate.sh
package filter
