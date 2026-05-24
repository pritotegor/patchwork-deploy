package lock_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/lock"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	dir := t.TempDir()
	lk := lock.New(dir, "host1", nil)
	if lk == nil {
		t.Fatal("expected non-nil Lock")
	}
}

func TestAcquire_CreatesLockFile(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	lk := lock.New(dir, "host1", &buf)

	if err := lk.Acquire(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer lk.Release() //nolint:errcheck

	if _, err := os.Stat(filepath.Join(dir, "host1.lock")); err != nil {
		t.Errorf("expected lock file to exist: %v", err)
	}
	if !lk.Held() {
		t.Error("expected Held() to return true after Acquire")
	}
}

func TestAcquire_FailsWhenAlreadyHeld(t *testing.T) {
	dir := t.TempDir()
	lk1 := lock.New(dir, "host2", nil)
	lk2 := lock.New(dir, "host2", nil)

	if err := lk1.Acquire(); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	defer lk1.Release() //nolint:errcheck

	if err := lk2.Acquire(); err == nil {
		t.Error("expected error when lock already held")
	}
}

func TestRelease_RemovesLockFile(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	lk := lock.New(dir, "host3", &buf)

	if err := lk.Acquire(); err != nil {
		t.Fatalf("acquire failed: %v", err)
	}
	if err := lk.Release(); err != nil {
		t.Fatalf("release failed: %v", err)
	}
	if lk.Held() {
		t.Error("expected Held() to return false after Release")
	}
}

func TestRelease_IdempotentWhenNotHeld(t *testing.T) {
	dir := t.TempDir()
	lk := lock.New(dir, "host4", nil)
	if err := lk.Release(); err != nil {
		t.Errorf("expected no error releasing non-existent lock, got: %v", err)
	}
}

func TestAcquire_WritesOutputMessage(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	lk := lock.New(dir, "host5", &buf)

	if err := lk.Acquire(); err != nil {
		t.Fatalf("acquire failed: %v", err)
	}
	defer lk.Release() //nolint:errcheck

	if buf.Len() == 0 {
		t.Error("expected output message on Acquire")
	}
}
