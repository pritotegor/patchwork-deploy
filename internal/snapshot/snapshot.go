// Package snapshot captures and restores the set of applied patches
// for a given host, enabling point-in-time recovery and diff reporting.
package snapshot

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Entry records a single applied patch in a snapshot.
type Entry struct {
	Patch     string    `json:"patch"`
	AppliedAt time.Time `json:"applied_at"`
}

// Snapshot holds the full state for one host at a point in time.
type Snapshot struct {
	Host      string    `json:"host"`
	CreatedAt time.Time `json:"created_at"`
	Entries   []Entry   `json:"entries"`
}

// Manager writes and reads snapshots from a directory.
type Manager struct {
	out io.Writer
	dir string
}

// New returns a Manager that stores snapshots under dir.
// If out is nil it defaults to os.Stdout.
func New(dir string, out io.Writer) *Manager {
	if out == nil {
		out = os.Stdout
	}
	return &Manager{out: out, dir: dir}
}

// Save writes a snapshot for host to disk, overwriting any previous one.
func (m *Manager) Save(host string, entries []Entry) error {
	if err := os.MkdirAll(m.dir, 0o755); err != nil {
		return fmt.Errorf("snapshot: mkdir %s: %w", m.dir, err)
	}
	s := Snapshot{
		Host:      host,
		CreatedAt: time.Now().UTC(),
		Entries:   entries,
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	path := m.filePath(host)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	fmt.Fprintf(m.out, "snapshot: saved %d entries for %s\n", len(entries), host)
	return nil
}

// Load reads the latest snapshot for host. Returns an empty Snapshot when
// no file exists yet (not an error).
func (m *Manager) Load(host string) (Snapshot, error) {
	path := m.filePath(host)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Snapshot{Host: host}, nil
	}
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: unmarshal %s: %w", path, err)
	}
	return s, nil
}

func (m *Manager) filePath(host string) string {
	return filepath.Join(m.dir, host+".json")
}
