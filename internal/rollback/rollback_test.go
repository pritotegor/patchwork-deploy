package rollback_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/rollback"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	rb := rollback.New(nil)
	if rb == nil {
		t.Fatal("expected non-nil Rollbacker")
	}
}

func TestRecord_AndEntries(t *testing.T) {
	rb := rollback.New(nil)
	rb.Record("001-init.sh", "001-undo.sh")
	rb.Record("002-migrate.sh", "")

	entries := rb.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].PatchName != "001-init.sh" {
		t.Errorf("unexpected first patch: %s", entries[0].PatchName)
	}
	if entries[1].UndoPath != "" {
		t.Errorf("expected empty undo path for second entry")
	}
}

func TestPlan_ReturnsReverseOrder(t *testing.T) {
	rb := rollback.New(nil)
	rb.Record("001-init.sh", "001-undo.sh")
	rb.Record("002-migrate.sh", "002-undo.sh")

	plan := rb.Plan()
	if len(plan) != 2 {
		t.Fatalf("expected 2 plan steps, got %d", len(plan))
	}
	if !strings.Contains(plan[0], "002-migrate.sh") {
		t.Errorf("expected first plan step to reference 002-migrate.sh, got: %s", plan[0])
	}
	if !strings.Contains(plan[1], "001-init.sh") {
		t.Errorf("expected second plan step to reference 001-init.sh, got: %s", plan[1])
	}
}

func TestPlan_NoUndoScript_Noted(t *testing.T) {
	rb := rollback.New(nil)
	rb.Record("001-init.sh", "")

	plan := rb.Plan()
	if !strings.Contains(plan[0], "no undo script") {
		t.Errorf("expected 'no undo script' note, got: %s", plan[0])
	}
}

func TestExecute_CallsUndoInReverseOrder(t *testing.T) {
	var buf strings.Builder
	rb := rollback.New(&buf)
	rb.Record("001-init.sh", "001-undo.sh")
	rb.Record("002-migrate.sh", "002-undo.sh")

	var called []string
	err := rb.Execute(func(undoPath string) error {
		called = append(called, undoPath)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(called) != 2 {
		t.Fatalf("expected 2 undo calls, got %d", len(called))
	}
	if called[0] != "002-undo.sh" || called[1] != "001-undo.sh" {
		t.Errorf("unexpected call order: %v", called)
	}
}

func TestExecute_StopsOnError(t *testing.T) {
	var buf strings.Builder
	rb := rollback.New(&buf)
	rb.Record("001-init.sh", "001-undo.sh")
	rb.Record("002-migrate.sh", "002-undo.sh")

	callCount := 0
	err := rb.Execute(func(undoPath string) error {
		callCount++
		return errors.New("undo failed")
	})
	if err == nil {
		t.Fatal("expected error from Execute")
	}
	if callCount != 1 {
		t.Errorf("expected exec to stop after first failure, got %d calls", callCount)
	}
}

func TestExecute_SkipsEmptyUndoPath(t *testing.T) {
	var buf strings.Builder
	rb := rollback.New(&buf)
	rb.Record("001-init.sh", "")
	rb.Record("002-migrate.sh", "002-undo.sh")

	var called []string
	_ = rb.Execute(func(undoPath string) error {
		called = append(called, undoPath)
		return nil
	})
	if len(called) != 1 || called[0] != "002-undo.sh" {
		t.Errorf("expected only 002-undo.sh to be called, got: %v", called)
	}
}
