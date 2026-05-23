package reporter_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/patchwork-deploy/internal/reporter"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	r := reporter.New(nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestPrint_SuccessResult(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)

	r.Print(reporter.Result{
		Host:     "host-1",
		Patch:    "001_init.sh",
		Output:   "done",
		Err:      nil,
		Duration: 120 * time.Millisecond,
	})

	out := buf.String()
	if !strings.Contains(out, "[OK]") {
		t.Errorf("expected [OK] in output, got: %s", out)
	}
	if !strings.Contains(out, "host-1") {
		t.Errorf("expected host-1 in output, got: %s", out)
	}
}

func TestPrint_FailureResult(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)

	r.Print(reporter.Result{
		Host:     "host-2",
		Patch:    "002_migrate.sh",
		Err:      errors.New("exit status 1"),
		Duration: 50 * time.Millisecond,
	})

	out := buf.String()
	if !strings.Contains(out, "FAIL") {
		t.Errorf("expected FAIL in output, got: %s", out)
	}
	if !strings.Contains(out, "exit status 1") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}

func TestSummary_AllSuccess(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)

	results := []reporter.Result{
		{Host: "h1", Patch: "001.sh"},
		{Host: "h2", Patch: "002.sh"},
	}
	r.Summary(results)

	if !strings.Contains(buf.String(), "2/2") {
		t.Errorf("expected 2/2 in summary, got: %s", buf.String())
	}
}

func TestSummary_WithFailures(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)

	results := []reporter.Result{
		{Host: "h1", Patch: "001.sh", Err: nil},
		{Host: "h2", Patch: "002.sh", Err: errors.New("failed")},
		{Host: "h3", Patch: "003.sh", Err: errors.New("timeout")},
	}
	r.Summary(results)

	if !strings.Contains(buf.String(), "1/3") {
		t.Errorf("expected 1/3 in summary, got: %s", buf.String())
	}
}
