package checkpoint_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/checkpoint"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	cp := checkpoint.New(t.TempDir(), nil)
	if cp == nil {
		t.Fatal("expected non-nil Checkpoint")
	}
}

func TestApplied_FalseBeforeMark(t *testing.T) {
	cp := checkpoint.New(t.TempDir(), &bytes.Buffer{})
	if cp.Applied("001-init.sh") {
		t.Error("expected patch to not be applied before Mark")
	}
}

func TestMark_PersistsToFile(t *testing.T) {
	dir := t.TempDir()
	cp := checkpoint.New(dir, &bytes.Buffer{})

	if err := cp.Mark("host1", "001-init.sh"); err != nil {
		t.Fatalf("Mark: %v", err)
	}

	expectedFile := filepath.Join(dir, "host1.checkpoint.json")
	if _, err := os.Stat(expectedFile); err != nil {
		t.Fatalf("checkpoint file not created: %v", err)
	}
}

func TestLoad_RestoresAppliedState(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer

	cp1 := checkpoint.New(dir, &buf)
	if err := cp1.Mark("host2", "002-setup.sh"); err != nil {
		t.Fatalf("Mark: %v", err)
	}

	cp2 := checkpoint.New(dir, &buf)
	if err := cp2.Load("host2"); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cp2.Applied("002-setup.sh") {
		t.Error("expected patch to be applied after Load")
	}
}

func TestLoad_MissingFileIsNoError(t *testing.T) {
	cp := checkpoint.New(t.TempDir(), &bytes.Buffer{})
	if err := cp.Load("unknown-host"); err != nil {
		t.Errorf("expected no error for missing checkpoint, got %v", err)
	}
}

func TestReset_ClearsState(t *testing.T) {
	dir := t.TempDir()
	cp := checkpoint.New(dir, &bytes.Buffer{})

	_ = cp.Mark("host3", "001-init.sh")
	if err := cp.Reset("host3"); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if cp.Applied("001-init.sh") {
		t.Error("expected patch to be cleared after Reset")
	}

	file := filepath.Join(dir, "host3.checkpoint.json")
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		t.Error("expected checkpoint file to be removed after Reset")
	}
}

func TestReset_IdempotentWhenNoFile(t *testing.T) {
	cp := checkpoint.New(t.TempDir(), &bytes.Buffer{})
	if err := cp.Reset("ghost-host"); err != nil {
		t.Errorf("Reset on missing file should not error, got %v", err)
	}
}

func TestLoad_WritesLogMessage(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer

	cp1 := checkpoint.New(dir, &buf)
	_ = cp1.Mark("host4", "003-migrate.sh")

	buf.Reset()
	cp2 := checkpoint.New(dir, &buf)
	_ = cp2.Load("host4")

	if buf.Len() == 0 {
		t.Error("expected log output after Load")
	}
}
