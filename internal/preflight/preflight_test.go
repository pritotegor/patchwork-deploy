package preflight_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/preflight"
)

// fakeExec implements preflight.Executor for tests.
type fakeExec struct {
	responses map[string]string
	errors    map[string]error
}

func (f *fakeExec) Run(cmd string) (string, error) {
	if err, ok := f.errors[cmd]; ok {
		return "", err
	}
	return f.responses[cmd], nil
}

func TestNew_DefaultsToStdout(t *testing.T) {
	r := preflight.New(nil, nil)
	if r == nil {
		t.Fatal("expected non-nil runner")
	}
}

func TestNew_UsesDefaultChecksWhenNil(t *testing.T) {
	r := preflight.New(nil, nil)
	var buf strings.Builder
	exec := &fakeExec{
		responses: map[string]string{
			`df -h / | awk 'NR==2{print $5}'`: "42%",
			"bash --version | head -1":         "GNU bash, version 5.1",
			"whoami":                           "deploy",
		},
	}
	_ = r
	r2 := preflight.New(&buf, nil)
	results, err := r2.RunAll("host1", exec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != len(preflight.DefaultChecks) {
		t.Fatalf("expected %d results, got %d", len(preflight.DefaultChecks), len(results))
	}
}

func TestRunAll_AllPass(t *testing.T) {
	var buf strings.Builder
	checks := []preflight.Check{
		{Name: "echo", Command: "echo hello"},
	}
	exec := &fakeExec{
		responses: map[string]string{"echo hello": "hello"},
	}
	r := preflight.New(&buf, checks)
	results, err := r.RunAll("web-01", exec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].OK {
		t.Fatalf("expected 1 passing result, got %+v", results)
	}
	if !strings.Contains(buf.String(), "status=OK") {
		t.Errorf("expected OK in output, got: %s", buf.String())
	}
}

func TestRunAll_FailingCheck(t *testing.T) {
	var buf strings.Builder
	checks := []preflight.Check{
		{Name: "disk-check", Command: "df /"},
	}
	execErr := errors.New("command not found")
	exec := &fakeExec{
		errors: map[string]error{"df /": execErr},
	}
	r := preflight.New(&buf, checks)
	results, err := r.RunAll("db-01", exec)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, execErr) {
		t.Errorf("expected wrapped execErr, got: %v", err)
	}
	if len(results) != 1 || results[0].OK {
		t.Fatalf("expected 1 failing result, got %+v", results)
	}
	if !strings.Contains(buf.String(), "status=FAIL") {
		t.Errorf("expected FAIL in output, got: %s", buf.String())
	}
}

func TestRunAll_StopsRecordingFirstError(t *testing.T) {
	var buf strings.Builder
	checks := []preflight.Check{
		{Name: "c1", Command: "cmd1"},
		{Name: "c2", Command: "cmd2"},
	}
	e1 := errors.New("fail1")
	e2 := errors.New("fail2")
	exec := &fakeExec{
		errors: map[string]error{"cmd1": e1, "cmd2": e2},
	}
	r := preflight.New(&buf, checks)
	_, err := r.RunAll("host", exec)
	if !errors.Is(err, e1) {
		t.Errorf("expected first error e1, got: %v", err)
	}
}
