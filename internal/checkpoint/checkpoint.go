// Package checkpoint tracks which patches have been successfully applied
// to a given host, allowing resumable deployments after partial failures.
package checkpoint

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Record represents a single applied patch entry.
type Record struct {
	PatchName string    `json:"patch_name"`
	Host      string    `json:"host"`
	AppliedAt time.Time `json:"applied_at"`
}

// Checkpoint manages applied patch state for a deployment run.
type Checkpoint struct {
	out     io.Writer
	dir     string
	applied map[string]bool
}

// New creates a Checkpoint that persists state under dir.
// If out is nil, os.Stdout is used for log messages.
func New(dir string, out io.Writer) *Checkpoint {
	if out == nil {
		out = os.Stdout
	}
	return &Checkpoint{
		out:     out,
		dir:     dir,
		applied: make(map[string]bool),
	}
}

// Load reads an existing checkpoint file for the given host.
func (c *Checkpoint) Load(host string) error {
	path := c.filePath(host)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("checkpoint: read %s: %w", path, err)
	}
	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("checkpoint: parse %s: %w", path, err)
	}
	for _, r := range records {
		c.applied[r.PatchName] = true
	}
	fmt.Fprintf(c.out, "checkpoint: loaded %d applied patches for %s\n", len(records), host)
	return nil
}

// Mark records a patch as successfully applied and persists the checkpoint.
func (c *Checkpoint) Mark(host, patchName string) error {
	c.applied[patchName] = true
	return c.save(host)
}

// Applied returns true if the patch has already been applied to the host.
func (c *Checkpoint) Applied(patchName string) bool {
	return c.applied[patchName]
}

// Reset removes the checkpoint file for a host.
func (c *Checkpoint) Reset(host string) error {
	path := c.filePath(host)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("checkpoint: reset %s: %w", host, err)
	}
	c.applied = make(map[string]bool)
	return nil
}

func (c *Checkpoint) save(host string) error {
	if err := os.MkdirAll(c.dir, 0o755); err != nil {
		return fmt.Errorf("checkpoint: mkdir: %w", err)
	}
	var records []Record
	for name := range c.applied {
		records = append(records, Record{PatchName: name, Host: host, AppliedAt: time.Now().UTC()})
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("checkpoint: marshal: %w", err)
	}
	return os.WriteFile(c.filePath(host), data, 0o644)
}

func (c *Checkpoint) filePath(host string) string {
	return filepath.Join(c.dir, host+".checkpoint.json")
}
