package pipeline_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/audit"
	"github.com/yourorg/patchwork-deploy/internal/checkpoint"
	"github.com/yourorg/patchwork-deploy/internal/lock"
	"github.com/yourorg/patchwork-deploy/internal/patch"
	"github.com/yourorg/patchwork-deploy/internal/pipeline"
	"github.com/yourorg/patchwork-deploy/internal/rollback"
)

func makeOpts(t *testing.T, patches []patch.Patch, out *bytes.Buffer) pipeline.Options {
	t.Helper()
	dir := t.TempDir()

	cp, _ := checkpoint.New(filepath.Join(dir, "cp.json"), out)
	rb := rollback.New(out)
	auditor := audit.New(out)
	lk := lock.New(filepath.Join(dir, "deploy.lock"), out)
	exec := patch.NewExecutor(out)

	return pipeline.Options{
		Host:       "host1",
		Patches:    patches,
		Executor:   exec,
		Auditor:    auditor,
		Checkpoint: cp,
		Rollback:   rb,
		Lock:       lk,
		Out:        out,
	}
}

func TestRun_NoPatches_ReturnsZeroCounts(t *testing.T) {
	var buf bytes.Buffer
	opts := makeOpts(t, nil, &buf)
	res := pipeline.Run(opts)
	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if res.Applied != 0 || res.Skipped != 0 {
		t.Errorf("expected 0/0, got %d/%d", res.Applied, res.Skipped)
	}
}

func TestRun_SkipsAlreadyApplied(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "001_hello.sh")
	os.WriteFile(script, []byte("echo hello"), 0644)

	patches := []patch.Patch{{Name: "001_hello.sh", Path: script}}
	var buf bytes.Buffer
	opts := makeOpts(t, patches, &buf)
	opts.Checkpoint.Mark("001_hello.sh") //nolint:errcheck

	res := pipeline.Run(opts)
	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if res.Skipped != 1 || res.Applied != 0 {
		t.Errorf("expected 0 applied 1 skipped, got applied=%d skipped=%d", res.Applied, res.Skipped)
	}
}

func TestRun_LockAlreadyHeld_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	lockPath := filepath.Join(dir, "deploy.lock")
	os.WriteFile(lockPath, []byte("held"), 0644)

	var buf bytes.Buffer
	opts := makeOpts(t, nil, &buf)
	opts.Lock = lock.New(lockPath, &buf)
	// pre-acquire so second acquire fails
	opts.Lock.Acquire() //nolint:errcheck

	// create a second lock on same file to simulate contention
	contender := lock.New(lockPath, &buf)
	res := pipeline.Run(pipeline.Options{
		Host:       "host1",
		Patches:    nil,
		Executor:   patch.NewExecutor(&buf),
		Auditor:    audit.New(&buf),
		Checkpoint: func() *checkpoint.Checkpoint { cp, _ := checkpoint.New(filepath.Join(dir, "cp.json"), &buf); return cp }(),
		Rollback:   rollback.New(&buf),
		Lock:       contender,
		Out:        &buf,
	})
	if res.Err == nil {
		t.Fatal("expected error when lock already held")
	}
}
