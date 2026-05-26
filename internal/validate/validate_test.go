package validate_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/patchwork-deploy/internal/patch"
	"github.com/patchwork-deploy/internal/validate"
)

func writePatch(t *testing.T, dir, name, content string) patch.Patch {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writePatch: %v", err)
	}
	return patch.Patch{Name: name, Path: path}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	v := validate.New(nil)
	if v == nil {
		t.Fatal("expected non-nil Validator")
	}
}

func TestValidateAll_AllValid(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	v := validate.New(&buf)

	patches := []patch.Patch{
		writePatch(t, dir, "001_init.sh", "echo hello"),
		writePatch(t, dir, "002_setup.sh", "echo world"),
	}

	results := v.ValidateAll(patches)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.OK {
			t.Errorf("expected OK for %s, got reason: %s", r.Patch.Name, r.Reason)
		}
	}
	if validate.AnyFailed(results) {
		t.Error("AnyFailed should be false")
	}
}

func TestValidateAll_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	v := validate.New(&buf)

	patches := []patch.Patch{
		{Name: "missing.sh", Path: "/nonexistent/path/missing.sh"},
	}

	results := v.ValidateAll(patches)
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if results[0].OK {
		t.Error("expected validation failure for missing file")
	}
	if !validate.AnyFailed(results) {
		t.Error("AnyFailed should be true")
	}
}

func TestValidateAll_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	v := validate.New(&buf)

	emptyPatch := writePatch(t, dir, "empty.sh", "")
	results := v.ValidateAll([]patch.Patch{emptyPatch})

	if results[0].OK {
		t.Error("expected failure for empty file")
	}
	if results[0].Reason == "" {
		t.Error("expected non-empty Reason")
	}
}

func TestValidateAll_EmptyPath(t *testing.T) {
	var buf bytes.Buffer
	v := validate.New(&buf)

	results := v.ValidateAll([]patch.Patch{{Name: "no-path.sh", Path: ""}})
	if results[0].OK {
		t.Error("expected failure for empty path")
	}
}

func TestValidateAll_WritesOutputLines(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	v := validate.New(&buf)

	patches := []patch.Patch{
		writePatch(t, dir, "001.sh", "echo ok"),
		{Name: "bad.sh", Path: "/no/such/file"},
	}
	v.ValidateAll(patches)

	out := buf.String()
	if out == "" {
		t.Error("expected output written to writer")
	}
}
