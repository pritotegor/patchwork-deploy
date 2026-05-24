package checkpoint

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Store lists and manages checkpoint files in a directory.
type Store struct {
	dir string
	out io.Writer
}

// NewStore returns a Store backed by dir.
func NewStore(dir string, out io.Writer) *Store {
	if out == nil {
		out = os.Stdout
	}
	return &Store{dir: dir, out: out}
}

// Hosts returns the names of all hosts that have checkpoint files.
func (s *Store) Hosts() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("store: list %s: %w", s.dir, err)
	}
	var hosts []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".checkpoint.json") {
			host := strings.TrimSuffix(name, ".checkpoint.json")
			hosts = append(hosts, host)
		}
	}
	return hosts, nil
}

// PurgeAll removes every checkpoint file in the store directory.
func (s *Store) PurgeAll() error {
	hosts, err := s.Hosts()
	if err != nil {
		return err
	}
	for _, host := range hosts {
		path := filepath.Join(s.dir, host+".checkpoint.json")
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("store: purge %s: %w", host, err)
		}
		fmt.Fprintf(s.out, "checkpoint: purged %s\n", host)
	}
	return nil
}
