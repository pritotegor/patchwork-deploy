package checkpoint_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/checkpoint"
)

func TestHosts_EmptyDir(t *testing.T) {
	store := checkpoint.NewStore(t.TempDir(), &bytes.Buffer{})
	hosts, err := store.Hosts()
	if err != nil {
		t.Fatalf("Hosts: %v", err)
	}
	if len(hosts) != 0 {
		t.Errorf("expected 0 hosts, got %d", len(hosts))
	}
}

func TestHosts_ReturnsCheckpointedHosts(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer

	for _, host := range []string{"web-01", "web-02", "db-01"} {
		cp := checkpoint.New(dir, &buf)
		if err := cp.Mark(host, "001-init.sh"); err != nil {
			t.Fatalf("Mark(%s): %v", host, err)
		}
	}

	store := checkpoint.NewStore(dir, &buf)
	hosts, err := store.Hosts()
	if err != nil {
		t.Fatalf("Hosts: %v", err)
	}
	if len(hosts) != 3 {
		t.Errorf("expected 3 hosts, got %d: %v", len(hosts), hosts)
	}
}

func TestHosts_NonexistentDirReturnsNil(t *testing.T) {
	store := checkpoint.NewStore("/nonexistent/path/xyz", &bytes.Buffer{})
	hosts, err := store.Hosts()
	if err != nil {
		t.Fatalf("expected no error for missing dir, got %v", err)
	}
	if hosts != nil {
		t.Errorf("expected nil hosts, got %v", hosts)
	}
}

func TestPurgeAll_RemovesFiles(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer

	for _, host := range []string{"alpha", "beta"} {
		cp := checkpoint.New(dir, &buf)
		_ = cp.Mark(host, "001-init.sh")
	}

	store := checkpoint.NewStore(dir, &buf)
	if err := store.PurgeAll(); err != nil {
		t.Fatalf("PurgeAll: %v", err)
	}

	hosts, _ := store.Hosts()
	if len(hosts) != 0 {
		t.Errorf("expected 0 hosts after purge, got %d", len(hosts))
	}
}

func TestPurgeAll_WritesLogLines(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer

	cp := checkpoint.New(dir, &buf)
	_ = cp.Mark("gamma", "001-init.sh")
	buf.Reset()

	store := checkpoint.NewStore(dir, &buf)
	_ = store.PurgeAll()

	if buf.Len() == 0 {
		t.Error("expected log output after PurgeAll")
	}
}
