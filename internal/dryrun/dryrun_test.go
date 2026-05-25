package dryrun_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/patchwork-deploy/internal/dryrun"
	"github.com/patchwork-deploy/internal/patch"
)

func makePatches(names ...string) []patch.Patch {
	var ps []patch.Patch
	for _, n := range names {
		ps = append(ps, patch.Patch{Name: n, Path: "/tmp/" + n + ".sh"})
	}
	return ps
}

func TestNew_DefaultsToStdout(t *testing.T) {
	r := dryrun.New(nil)
	if r == nil {
		t.Fatal("expected non-nil Runner")
	}
}

func TestSimulate_NoPatches_PrintsNotice(t *testing.T) {
	var buf bytes.Buffer
	r := dryrun.New(&buf)

	results := r.Simulate(nil, []string{"host1"})

	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
	if !strings.Contains(buf.String(), "no patches to simulate") {
		t.Errorf("expected notice in output, got: %s", buf.String())
	}
}

func TestSimulate_ReturnOneResultPerPatchAndHost(t *testing.T) {
	var buf bytes.Buffer
	r := dryrun.New(&buf)

	patches := makePatches("001-init", "002-setup")
	hosts := []string{"web1", "web2"}

	results := r.Simulate(patches, hosts)

	if len(results) != 4 {
		t.Fatalf("expected 4 results (2 patches × 2 hosts), got %d", len(results))
	}
	for _, res := range results {
		if !res.Simulated {
			t.Errorf("expected Simulated=true for %s/%s", res.PatchName, res.Host)
		}
		if res.Skipped {
			t.Errorf("expected Skipped=false for %s/%s", res.PatchName, res.Host)
		}
	}
}

func TestSimulate_LogsEachPatchAndHost(t *testing.T) {
	var buf bytes.Buffer
	r := dryrun.New(&buf)

	r.Simulate(makePatches("001-migrate"), []string{"db1"})

	out := buf.String()
	if !strings.Contains(out, "001-migrate") {
		t.Errorf("expected patch name in output: %s", out)
	}
	if !strings.Contains(out, "db1") {
		t.Errorf("expected host in output: %s", out)
	}
	if !strings.Contains(out, "[dry-run]") {
		t.Errorf("expected [dry-run] prefix in output: %s", out)
	}
}

func TestSkip_ReturnsSkippedResult(t *testing.T) {
	var buf bytes.Buffer
	r := dryrun.New(&buf)

	res := r.Skip("003-cleanup", "app1", "already applied")

	if !res.Skipped {
		t.Error("expected Skipped=true")
	}
	if res.Reason != "already applied" {
		t.Errorf("expected reason 'already applied', got %q", res.Reason)
	}
	out := buf.String()
	if !strings.Contains(out, "SKIP") {
		t.Errorf("expected SKIP in output: %s", out)
	}
	if !strings.Contains(out, "already applied") {
		t.Errorf("expected reason in output: %s", out)
	}
}
