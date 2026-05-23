package runner

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewHookRunner_DefaultsToStdout(t *testing.T) {
	hr := NewHookRunner(nil)
	if hr.out == nil {
		t.Fatal("expected non-nil writer when nil passed to NewHookRunner")
	}
}

func TestRun_EmptyCommandReturnsError(t *testing.T) {
	var buf bytes.Buffer
	hr := NewHookRunner(&buf)

	err := hr.Run(Hook{Type: HookPrePatch, Command: "   "})
	if err == nil {
		t.Fatal("expected error for empty command, got nil")
	}
	if !strings.Contains(err.Error(), "empty command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRun_SuccessfulCommand(t *testing.T) {
	var buf bytes.Buffer
	hr := NewHookRunner(&buf)

	err := hr.Run(Hook{Type: HookPostPatch, Command: "echo hello"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected output to contain 'hello', got: %q", buf.String())
	}
}

func TestRun_FailingCommandReturnsError(t *testing.T) {
	var buf bytes.Buffer
	hr := NewHookRunner(&buf)

	err := hr.Run(Hook{Type: HookPrePatch, Command: "exit 1"})
	if err == nil {
		t.Fatal("expected error for failing command, got nil")
	}
	if !strings.Contains(err.Error(), string(HookPrePatch)) {
		t.Errorf("expected error to mention hook type, got: %v", err)
	}
}

func TestRunAll_StopsOnFirstFailure(t *testing.T) {
	var buf bytes.Buffer
	hr := NewHookRunner(&buf)

	hooks := []Hook{
		{Type: HookPrePatch, Command: "echo first"},
		{Type: HookPrePatch, Command: "exit 2"},
		{Type: HookPrePatch, Command: "echo third"},
	}

	err := hr.RunAll(hooks)
	if err == nil {
		t.Fatal("expected error from RunAll, got nil")
	}
	if strings.Contains(buf.String(), "third") {
		t.Error("expected execution to stop before third hook")
	}
}

func TestRunAll_AllSucceed(t *testing.T) {
	var buf bytes.Buffer
	hr := NewHookRunner(&buf)

	hooks := []Hook{
		{Type: HookPrePatch, Command: "echo a"},
		{Type: HookPostPatch, Command: "echo b"},
	}

	if err := hr.RunAll(hooks); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "a") || !strings.Contains(out, "b") {
		t.Errorf("expected both hooks to run, output: %q", out)
	}
}
