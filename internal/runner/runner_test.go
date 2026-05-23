package runner_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/config"
	"github.com/yourorg/patchwork-deploy/internal/runner"
)

func makePatchDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("write patch file: %v", err)
		}
	}
	return dir
}

func TestNew_DefaultsOutputToStdout(t *testing.T) {
	cfg := &config.Config{PatchDir: t.TempDir()}
	r := runner.New(cfg)
	if r.Out == nil {
		t.Fatal("expected Out to be non-nil")
	}
	if r.ErrOut == nil {
		t.Fatal("expected ErrOut to be non-nil")
	}
}

func TestRun_NoPatchesReturnsEarly(t *testing.T) {
	dir := makePatchDir(t, map[string]string{})
	cfg := &config.Config{PatchDir: dir, Hosts: []string{"localhost"}}

	r := runner.New(cfg)
	var out bytes.Buffer
	r.Out = &out

	if err := r.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.String() == "" {
		t.Fatal("expected some output when no patches found")
	}
}

func TestRun_InvalidPatchDirReturnsError(t *testing.T) {
	cfg := &config.Config{
		PatchDir: "/nonexistent/path/that/does/not/exist",
		Hosts:    []string{"localhost"},
	}
	r := runner.New(cfg)
	var out bytes.Buffer
	r.Out = &out

	err := r.Run()
	if err == nil {
		t.Fatal("expected error for invalid patch directory")
	}
}

func TestRun_NoHostsSkipsSSH(t *testing.T) {
	dir := makePatchDir(t, map[string]string{
		"01_init.sh": "#!/bin/sh\necho hello",
	})
	cfg := &config.Config{PatchDir: dir, Hosts: []string{}}

	r := runner.New(cfg)
	var out bytes.Buffer
	r.Out = &out

	if err := r.Run(); err != nil {
		t.Fatalf("unexpected error with no hosts: %v", err)
	}
}
