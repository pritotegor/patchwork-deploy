package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// buildBinary compiles the main package into a temp binary for integration tests.
func buildBinary(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	out := filepath.Join(tmp, "patchwork")
	cmd := exec.Command("go", "build", "-o", out, ".")
	cmd.Dir = "."
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, b)
	}
	return out
}

func TestMain_Version(t *testing.T) {
	bin := buildBinary(t)
	out, err := exec.Command(bin, "-version").Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := string(out)
	if got == "" {
		t.Fatal("expected version output, got empty string")
	}
	if len(got) < 3 {
		t.Errorf("version output too short: %q", got)
	}
}

func TestMain_MissingConfig(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "-config", "/nonexistent/path.json")
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected non-zero exit for missing config, got nil")
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 1 {
			t.Errorf("expected exit code 1, got %d", exitErr.ExitCode())
		}
	}
}

func TestMain_DryRun(t *testing.T) {
	bin := buildBinary(t)
	tmp := t.TempDir()

	cfgContent := []byte(`{
	"patch_dir": "` + tmp + `",
	"hosts": [{"host": "localhost", "port": 22, "user": "deploy", "key_file": "/tmp/key"}]
}`)
	cfgPath := filepath.Join(tmp, "patchwork.json")
	if err := os.WriteFile(cfgPath, cfgContent, 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	out, err := exec.Command(bin, "-config", cfgPath, "-dry-run").Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := string(out)
	if got == "" {
		t.Error("expected dry-run output, got empty")
	}
}
