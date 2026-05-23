package reporter

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewSummaryReporter_DefaultsToStdout(t *testing.T) {
	sr := NewSummaryReporter(nil)
	if sr.out == nil {
		t.Fatal("expected non-nil writer when nil passed")
	}
}

func TestWrite_SuccessfulRun(t *testing.T) {
	var buf bytes.Buffer
	sr := NewSummaryReporter(&buf)

	sr.Write(RunSummary{
		Host:      "prod-01",
		Total:     5,
		Succeeded: 5,
		Elapsed:   120 * time.Millisecond,
	})

	out := buf.String()
	assertContains(t, out, "prod-01")
	assertContains(t, out, "Total patches : 5")
	assertContains(t, out, "Succeeded     : 5")
	assertContains(t, out, "Result: SUCCESS")

	if strings.Contains(out, "Failed patches") {
		t.Error("expected no failed patches section for a successful run")
	}
}

func TestWrite_FailedRun(t *testing.T) {
	var buf bytes.Buffer
	sr := NewSummaryReporter(&buf)

	sr.Write(RunSummary{
		Host:      "staging-02",
		Total:     3,
		Succeeded: 2,
		Failed:    1,
		PatchErrors: []PatchError{
			{Patch: "003-migrate.sh", Err: errors.New("exit status 1")},
		},
		Elapsed: 300 * time.Millisecond,
	})

	out := buf.String()
	assertContains(t, out, "staging-02")
	assertContains(t, out, "Failed        : 1")
	assertContains(t, out, "Failed patches")
	assertContains(t, out, "003-migrate.sh")
	assertContains(t, out, "exit status 1")
	assertContains(t, out, "Result: FAILURE")
}

func TestWrite_SkippedPatches(t *testing.T) {
	var buf bytes.Buffer
	sr := NewSummaryReporter(&buf)

	sr.Write(RunSummary{
		Host:      "dev-03",
		Total:     4,
		Succeeded: 3,
		Skipped:   1,
		Elapsed:   50 * time.Millisecond,
	})

	out := buf.String()
	assertContains(t, out, "Skipped       : 1")
	assertContains(t, out, "Result: SUCCESS")
}

func assertContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Errorf("expected output to contain %q\ngot:\n%s", needle, haystack)
	}
}
