package snapshot_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/snapshot"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	m := snapshot.New(t.TempDir(), nil)
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	buf := &bytes.Buffer{}
	m := snapshot.New(dir, buf)

	entries := []snapshot.Entry{
		{Patch: "001_init.sh", AppliedAt: time.Now().UTC()},
	}
	if err := m.Save("web-01", entries); err != nil {
		t.Fatalf("Save: %v", err)
	}

	path := filepath.Join(dir, "web-01.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file at %s: %v", path, err)
	}
	if buf.Len() == 0 {
		t.Error("expected log output after Save")
	}
}

func TestLoad_MissingFileReturnsEmpty(t *testing.T) {
	m := snapshot.New(t.TempDir(), nil)
	s, err := m.Load("unknown-host")
	if err != nil {
		t.Fatalf("Load: unexpected error: %v", err)
	}
	if s.Host != "unknown-host" {
		t.Errorf("host = %q, want %q", s.Host, "unknown-host")
	}
	if len(s.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(s.Entries))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	m := snapshot.New(dir, nil)

	now := time.Now().UTC().Truncate(time.Second)
	entries := []snapshot.Entry{
		{Patch: "001_init.sh", AppliedAt: now},
		{Patch: "002_schema.sh", AppliedAt: now.Add(time.Minute)},
	}
	if err := m.Save("db-01", entries); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s, err := m.Load("db-01")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if s.Host != "db-01" {
		t.Errorf("host = %q, want %q", s.Host, "db-01")
	}
	if len(s.Entries) != 2 {
		t.Fatalf("entries = %d, want 2", len(s.Entries))
	}
	if s.Entries[0].Patch != "001_init.sh" {
		t.Errorf("entry[0].Patch = %q", s.Entries[0].Patch)
	}
}

func TestSave_OverwritesPreviousSnapshot(t *testing.T) {
	dir := t.TempDir()
	m := snapshot.New(dir, nil)

	first := []snapshot.Entry{{Patch: "001_init.sh", AppliedAt: time.Now().UTC()}}
	if err := m.Save("app-01", first); err != nil {
		t.Fatalf("first Save: %v", err)
	}

	second := []snapshot.Entry{
		{Patch: "001_init.sh", AppliedAt: time.Now().UTC()},
		{Patch: "002_feature.sh", AppliedAt: time.Now().UTC()},
	}
	if err := m.Save("app-01", second); err != nil {
		t.Fatalf("second Save: %v", err)
	}

	s, err := m.Load("app-01")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(s.Entries) != 2 {
		t.Errorf("expected 2 entries after overwrite, got %d", len(s.Entries))
	}
}
