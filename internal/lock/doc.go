// Package lock implements a file-based mutual exclusion mechanism for
// patchwork-deploy deployments.
//
// A Lock is scoped to a single target host. Before applying patches, the
// runner should call Acquire to create a lock file under a shared lock
// directory. On completion — successful or otherwise — Release removes the
// file so that subsequent runs may proceed.
//
// Lock files are plain text and contain a small JSON object recording the
// PID and timestamp of the process that acquired the lock, which aids
// debugging when a stale lock must be cleared manually.
//
// Example usage:
//
//	lk := lock.New("/var/run/patchwork", "web-01", os.Stdout)
//	if err := lk.Acquire(); err != nil {
//	    log.Fatal(err)
//	}
//	defer lk.Release()
package lock
