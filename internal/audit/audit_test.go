package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/patchwork-deploy/internal/audit"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLog_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	err := l.Log(audit.Event{
		Type:    audit.EventDeployStart,
		Host:    "host1",
		Success: true,
		Message: "starting",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var got audit.Event
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Type != audit.EventDeployStart {
		t.Errorf("expected type %q, got %q", audit.EventDeployStart, got.Type)
	}
	if got.Host != "host1" {
		t.Errorf("expected host %q, got %q", "host1", got.Host)
	}
	if !got.Success {
		t.Error("expected success=true")
	}
	if got.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLogPatch_Success(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	if err := l.LogPatch("host2", "001_init.sh", true, "ok"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got audit.Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Type != audit.EventPatchApplied {
		t.Errorf("expected %q, got %q", audit.EventPatchApplied, got.Type)
	}
}

func TestLogPatch_Failure(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	if err := l.LogPatch("host3", "002_migrate.sh", false, "exit 1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got audit.Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Type != audit.EventPatchFailed {
		t.Errorf("expected %q, got %q", audit.EventPatchFailed, got.Type)
	}
	if got.Success {
		t.Error("expected success=false")
	}
}
