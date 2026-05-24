package pipeline_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/pipeline"
)

func TestSummary_AllSuccess(t *testing.T) {
	results := []pipeline.Result{
		{Host: "h1", Applied: 3, Skipped: 1},
		{Host: "h2", Applied: 2, Skipped: 0},
	}
	var buf bytes.Buffer
	pipeline.Summary(results, &buf)
	out := buf.String()

	for _, want := range []string{"Hosts    : 2", "Applied  : 5", "Skipped  : 1", "Failures : 0"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestSummary_WithFailures(t *testing.T) {
	results := []pipeline.Result{
		{Host: "h1", Applied: 1, Err: errors.New("script failed")},
		{Host: "h2", Applied: 2},
	}
	var buf bytes.Buffer
	pipeline.Summary(results, &buf)
	out := buf.String()

	if !strings.Contains(out, "Failures : 1") {
		t.Errorf("expected failure count in output:\n%s", out)
	}
	if !strings.Contains(out, "h1") {
		t.Errorf("expected failed host name in output:\n%s", out)
	}
}

func TestHasFailures_NoErrors(t *testing.T) {
	results := []pipeline.Result{
		{Host: "h1", Applied: 2},
		{Host: "h2", Applied: 1},
	}
	if pipeline.HasFailures(results) {
		t.Error("expected no failures")
	}
}

func TestHasFailures_WithError(t *testing.T) {
	results := []pipeline.Result{
		{Host: "h1"},
		{Host: "h2", Err: errors.New("boom")},
	}
	if !pipeline.HasFailures(results) {
		t.Error("expected HasFailures to return true")
	}
}

func TestSummary_DefaultsToStdout(t *testing.T) {
	// Should not panic when writer is nil.
	pipeline.Summary([]pipeline.Result{{Host: "h1", Applied: 1}}, nil)
}
