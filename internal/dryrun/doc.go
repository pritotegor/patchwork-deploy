// Package dryrun implements a simulation layer for patchwork-deploy.
//
// It allows operators to preview which patches would be applied to which
// hosts without making any real SSH connections or executing shell commands.
//
// Usage:
//
//	r := dryrun.New(os.Stdout)
//	results := r.Simulate(patches, hosts)
//
// Each Result indicates whether the patch would be applied or skipped and
// includes the reason when a patch is omitted (e.g. already applied via
// checkpoint, outside schedule window, or filtered by tag).
package dryrun
