package notify_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/notify"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	n := notify.New(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}

func TestSend_WritesFormattedLine(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf)

	err := n.Send(notify.LevelInfo, "web-01", "deployment started")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := buf.String()
	if !strings.Contains(line, "INFO") {
		t.Errorf("expected INFO in output, got: %s", line)
	}
	if !strings.Contains(line, "web-01") {
		t.Errorf("expected host in output, got: %s", line)
	}
	if !strings.Contains(line, "deployment started") {
		t.Errorf("expected message in output, got: %s", line)
	}
}

func TestInfo_WritesInfoLevel(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf)

	_ = n.Info("db-01", "patch applied")

	if !strings.Contains(buf.String(), "INFO") {
		t.Errorf("expected INFO level, got: %s", buf.String())
	}
}

func TestError_WritesErrorLevel(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf)

	_ = n.Error("db-01", "connection failed")

	if !strings.Contains(buf.String(), "ERROR") {
		t.Errorf("expected ERROR level, got: %s", buf.String())
	}
}

func TestSend_MultipleEvents(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf)

	_ = n.Info("host-a", "start")
	_ = n.Send(notify.LevelWarn, "host-b", "slow response")
	_ = n.Error("host-c", "failed")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %s", len(lines), buf.String())
	}
	if !strings.Contains(lines[1], "WARN") {
		t.Errorf("expected WARN on second line, got: %s", lines[1])
	}
}
