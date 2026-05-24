// Package lock provides a simple deployment lock mechanism to prevent
// concurrent patch runs against the same target host.
package lock

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Lock represents a file-based deployment lock for a given host.
type Lock struct {
	path string
	out  io.Writer
}

// New creates a new Lock that stores lock files under dir, using host as
// the lock file name. Falls back to os.Stdout if out is nil.
func New(dir, host string, out io.Writer) *Lock {
	if out == nil {
		out = os.Stdout
	}
	return &Lock{
		path: filepath.Join(dir, host+".lock"),
		out:  out,
	}
}

// Acquire attempts to create the lock file. Returns an error if the lock
// already exists (another deployment is in progress).
func (l *Lock) Acquire() error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("lock already held for %s: %s", l.path, err)
		}
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "{\"pid\":%d,\"time\":\"%s\"}\n", os.Getpid(), time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to write lock file: %w", err)
	}
	fmt.Fprintf(l.out, "[lock] acquired: %s\n", l.path)
	return nil
}

// Release removes the lock file, allowing subsequent deployments to proceed.
func (l *Lock) Release() error {
	if err := os.Remove(l.path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to release lock: %w", err)
	}
	fmt.Fprintf(l.out, "[lock] released: %s\n", l.path)
	return nil
}

// Held reports whether the lock file currently exists on disk.
func (l *Lock) Held() bool {
	_, err := os.Stat(l.path)
	return err == nil
}
