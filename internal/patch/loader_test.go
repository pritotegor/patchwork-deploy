package patch_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/patch"
)

func TestLoadPatches_OrderedAndFiltered(t *testing.T) {
	dir := t.TempDir()

	files := map[string]string{
		"002_setup.sh":  "#!/bin/sh\necho setup",
		"001_init.sh":   "#!/bin/sh\necho init",
		"003_deploy.sh": "#!/bin/sh\necho deploy",
		"notes.txt":     "not a patch",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatalf("writing temp file: %v", err)
		}
	}

	patches, err := patch.LoadPatches(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(patches) != 3 {
		t.Fatalf("expected 3 patches, got %d", len(patches))
	}

	expectedOrder := []string{"001_init.sh", "002_setup.sh", "003_deploy.sh"}
	for i, p := range patches {
		if p.Name != expectedOrder[i] {
			t.Errorf("patch[%d]: expected %q, got %q", i, expectedOrder[i], p.Name)
		}
		if p.Content == "" {
			t.Errorf("patch[%d] %q: content should not be empty", i, p.Name)
		}
	}
}

func TestLoadPatches_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	patches, err := patch.LoadPatches(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(patches) != 0 {
		t.Fatalf("expected 0 patches, got %d", len(patches))
	}
}

func TestLoadPatches_InvalidDirectory(t *testing.T) {
	_, err := patch.LoadPatches("/nonexistent/path/to/patches")
	if err == nil {
		t.Fatal("expected error for nonexistent directory, got nil")
	}
}
