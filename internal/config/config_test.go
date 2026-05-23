package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	return path
}

func TestLoad_ValidConfig(t *testing.T) {
	content := `{
		"patch_dir": "./patches",
		"hosts": [
			{"name": "web01", "address": "10.0.0.1", "user": "deploy", "key_file": "~/.ssh/id_rsa", "port": 22}
		]
	}`
	path := writeConfigFile(t, content)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PatchDir != "./patches" {
		t.Errorf("expected patch_dir './patches', got %q", cfg.PatchDir)
	}
	if len(cfg.Hosts) != 1 {
		t.Fatalf("expected 1 host, got %d", len(cfg.Hosts))
	}
	if cfg.Hosts[0].Name != "web01" {
		t.Errorf("expected host name 'web01', got %q", cfg.Hosts[0].Name)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := writeConfigFile(t, `{invalid json}`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestLoad_MissingPatchDir(t *testing.T) {
	content := `{"hosts": [{"name": "h", "address": "1.2.3.4", "user": "u"}]}`
	path := writeConfigFile(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing patch_dir")
	}
}

func TestLoad_NoHosts(t *testing.T) {
	content := `{"patch_dir": "./patches", "hosts": []}`
	path := writeConfigFile(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty hosts")
	}
}

func TestHost_DefaultPort(t *testing.T) {
	h := &Host{Port: 0}
	if h.DefaultPort() != 22 {
		t.Errorf("expected default port 22, got %d", h.DefaultPort())
	}
	h.Port = 2222
	if h.DefaultPort() != 2222 {
		t.Errorf("expected port 2222, got %d", h.DefaultPort())
	}
}
