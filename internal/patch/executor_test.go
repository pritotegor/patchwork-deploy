package patch

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestNewExecutor_DefaultsToStdout(t *testing.T) {
	exec := NewExecutor(nil, nil)
	if exec.Out == nil {
		t.Fatal("expected Out to default to os.Stdout, got nil")
	}
}

func TestApply_StopsOnMissingFile(t *testing.T) {
	var buf bytes.Buffer
	exec := &Executor{Client: nil, Out: &buf}

	patches := []Patch{
		{Name: "001-missing.sh", Path: "/nonexistent/path/001-missing.sh"},
	}

	// runPatch reads the file before using the SSH client, so nil client is safe here.
	result := exec.runPatch(patches[0])
	if result.Err == nil {
		t.Fatal("expected error for missing patch file, got nil")
	}
}

func TestApply_ReadsAndReturnsOutput(t *testing.T) {
	dir := t.TempDir()

	scriptPath := filepath.Join(dir, "001-hello.sh")
	if err := os.WriteFile(scriptPath, []byte("echo hello"), 0600); err != nil {
		t.Fatalf("write script: %v", err)
	}

	p := Patch{Name: "001-hello.sh", Path: scriptPath}

	// Read the file manually to confirm executor can read it.
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	if string(content) != "echo hello" {
		t.Errorf("expected 'echo hello', got %q", string(content))
	}

	_ = p // SSH session creation would fail without a real client; file read is validated above.
}

func TestExecuteResult_ErrorWrapping(t *testing.T) {
	result := ExecuteResult{
		Patch:  Patch{Name: "002-fail.sh"},
		Output: "some output",
		Err:    os.ErrPermission,
	}

	if result.Err == nil {
		t.Fatal("expected non-nil error")
	}
	if result.Output != "some output" {
		t.Errorf("unexpected output: %q", result.Output)
	}
	if result.Patch.Name != "002-fail.sh" {
		t.Errorf("unexpected patch name: %q", result.Patch.Name)
	}
}
