// Package pipeline provides a high-level orchestration layer that sequences
// all deployment steps for a single target host.
//
// A pipeline run proceeds as follows:
//
//  1. Acquire the distributed lock for the host to prevent concurrent deploys.
//  2. Iterate over the ordered list of patches.
//  3. Skip any patch already recorded in the checkpoint store.
//  4. Execute the patch script via the SSH executor.
//  5. Write an audit log entry for every attempt (success or failure).
//  6. Register the patch with the rollback tracker.
//  7. On success, mark the patch in the checkpoint store.
//  8. On failure, stop immediately and return the error.
//  9. Release the lock on exit.
//
// The pipeline is intentionally stateless — all dependencies are injected
// through Options so they can be replaced with test doubles in unit tests.
package pipeline
