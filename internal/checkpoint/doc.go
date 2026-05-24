// Package checkpoint provides resumable deployment state tracking.
//
// It records which patches have been successfully applied to each remote
// host and persists that state to disk. On subsequent runs the executor
// can skip already-applied patches, making deployments idempotent and
// safe to resume after a partial failure.
//
// Usage:
//
//	cp := checkpoint.New(".patchwork/checkpoints", os.Stdout)
//	if err := cp.Load(host); err != nil {
//		log.Fatal(err)
//	}
//	if !cp.Applied(patch.Name) {
//		// apply patch …
//		cp.Mark(host, patch.Name)
//	}
package checkpoint
