package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/patchwork-deploy/internal/audit"
)

func TestNewFileLogger_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "logs", "audit.jsonl")

	fl, err := audit.NewFileLogger(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer fl.Close()

	if fl.Path() != path {
		t.Errorf("expected path %q, got %q", path, fl.Path())
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("log file not created: %v", err)
	}
}

func TestNewFileLogger_WritesAndPersists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")

	fl, err := audit.NewFileLogger(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := fl.LogPatch("srv1", "001_setup.sh", true, "done"); err != nil {
		t.Fatalf("log error: %v", err)
	}
	fl.Close()

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open log: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("expected at least one line")
	}
	var e audit.Event
	if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if e.Host != "srv1" {
		t.Errorf("expected host srv1, got %q", e.Host)
	}
}

func TestNewFileLogger_InvalidPath(t *testing.T) {
	// Use a path whose parent is a file, not a directory.
	dir := t.TempDir()
	blocking := filepath.Join(dir, "block")
	if err := os.WriteFile(blocking, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := audit.NewFileLogger(filepath.Join(blocking, "audit.jsonl"))
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
