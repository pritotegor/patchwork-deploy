package progress_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/patchwork-deploy/internal/progress"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	// Should not panic when out is nil.
	tracker := progress.New(nil)
	if tracker == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestStarted_WritesLine(t *testing.T) {
	var buf bytes.Buffer
	tracker := progress.New(&buf)
	tracker.Started("db-01", "001_init.sh")

	got := buf.String()
	if !strings.Contains(got, "host=db-01") {
		t.Errorf("expected host in output, got: %s", got)
	}
	if !strings.Contains(got, "patch=001_init.sh") {
		t.Errorf("expected patch in output, got: %s", got)
	}
	if !strings.Contains(got, "status=started") {
		t.Errorf("expected status=started in output, got: %s", got)
	}
}

func TestDone_SuccessWritesOk(t *testing.T) {
	var buf bytes.Buffer
	tracker := progress.New(&buf)
	tracker.Done("web-01", "002_migrate.sh", 120*time.Millisecond, nil)

	got := buf.String()
	if !strings.Contains(got, "status=ok") {
		t.Errorf("expected status=ok, got: %s", got)
	}
	if !strings.Contains(got, "elapsed=") {
		t.Errorf("expected elapsed in output, got: %s", got)
	}
}

func TestDone_FailureWritesFailed(t *testing.T) {
	var buf bytes.Buffer
	tracker := progress.New(&buf)
	tracker.Done("web-02", "003_cleanup.sh", 50*time.Millisecond, errors.New("exit status 1"))

	got := buf.String()
	if !strings.Contains(got, "status=failed") {
		t.Errorf("expected status=failed, got: %s", got)
	}
	if !strings.Contains(got, "exit status 1") {
		t.Errorf("expected error message in output, got: %s", got)
	}
}

func TestSkipped_WritesLine(t *testing.T) {
	var buf bytes.Buffer
	tracker := progress.New(&buf)
	tracker.Skipped("cache-01", "004_index.sh")

	got := buf.String()
	if !strings.Contains(got, "status=skipped") {
		t.Errorf("expected status=skipped, got: %s", got)
	}
}

func TestRecord_UnknownStatus(t *testing.T) {
	var buf bytes.Buffer
	tracker := progress.New(&buf)
	tracker.Record(progress.Event{Host: "h", Patch: "p", Status: "pending"})

	got := buf.String()
	if !strings.Contains(got, "status=pending") {
		t.Errorf("expected status=pending in output, got: %s", got)
	}
}

func TestRecord_ConcurrentSafe(t *testing.T) {
	var buf bytes.Buffer
	tracker := progress.New(&buf)

	done := make(chan struct{})
	for i := 0; i < 20; i++ {
		go func(n int) {
			tracker.Started("host", "patch")
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 20; i++ {
		<-done
	}
}
