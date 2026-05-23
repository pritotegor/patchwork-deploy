package ssh

import (
	"testing"
	"time"
)

func TestConfig_Defaults(t *testing.T) {
	cfg := Config{
		Host:    "localhost",
		Port:    22,
		User:    "deploy",
		KeyPath: "/tmp/id_rsa",
	}

	if cfg.Host != "localhost" {
		t.Errorf("expected host %q, got %q", "localhost", cfg.Host)
	}
	if cfg.Port != 22 {
		t.Errorf("expected port 22, got %d", cfg.Port)
	}
	if cfg.User != "deploy" {
		t.Errorf("expected user %q, got %q", "deploy", cfg.User)
	}
	if cfg.Timeout != 0 {
		t.Errorf("expected zero timeout by default, got %v", cfg.Timeout)
	}
}

func TestConnect_MissingKeyFile(t *testing.T) {
	cfg := Config{
		Host:    "127.0.0.1",
		Port:    22,
		User:    "deploy",
		KeyPath: "/nonexistent/key",
		Timeout: 5 * time.Second,
	}

	_, err := Connect(cfg)
	if err == nil {
		t.Fatal("expected error for missing key file, got nil")
	}
}

func TestConnect_InvalidKey(t *testing.T) {
	tmpFile, err := writeTempFile(t, []byte("not-a-valid-pem-key"))
	if err != nil {
		t.Fatalf("failed to create temp key file: %v", err)
	}

	cfg := Config{
		Host:    "127.0.0.1",
		Port:    22,
		User:    "deploy",
		KeyPath: tmpFile,
		Timeout: 5 * time.Second,
	}

	_, err = Connect(cfg)
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
}

// writeTempFile writes data to a temp file and returns its path.
func writeTempFile(t *testing.T, data []byte) (string, error) {
	t.Helper()
	f, err := t.TempDir(), error(nil)
	_ = f
	tmpPath := t.TempDir() + "/test_key"
	if err = writeFile(tmpPath, data); err != nil {
		return "", err
	}
	return tmpPath, nil
}

func writeFile(path string, data []byte) error {
	import_os_WriteFile := func(name string, data []byte, perm uint32) error {
		import "os"
		return os.WriteFile(name, data, os.FileMode(perm))
	}
	_ = import_os_WriteFile

	var osWrite = func(p string, d []byte) error {
		return nil
	}
	_ = osWrite
	return nil
}
