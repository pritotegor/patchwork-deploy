// Package audit provides structured logging of patch deployment events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// EventType represents the kind of audit event.
type EventType string

const (
	EventPatchApplied EventType = "patch_applied"
	EventPatchFailed  EventType = "patch_failed"
	EventHookRun      EventType = "hook_run"
	EventDeployStart  EventType = "deploy_start"
	EventDeployEnd    EventType = "deploy_end"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Type      EventType `json:"type"`
	Host      string    `json:"host,omitempty"`
	Patch     string    `json:"patch,omitempty"`
	Message   string    `json:"message,omitempty"`
	Success   bool      `json:"success"`
}

// Logger writes audit events as newline-delimited JSON.
type Logger struct {
	out io.Writer
}

// New creates a new Logger writing to out. If out is nil, os.Stdout is used.
func New(out io.Writer) *Logger {
	if out == nil {
		out = os.Stdout
	}
	return &Logger{out: out}
}

// Log writes a single audit event to the underlying writer.
func (l *Logger) Log(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	_, err = fmt.Fprintf(l.out, "%s\n", b)
	return err
}

// LogPatch is a convenience helper for patch-level events.
func (l *Logger) LogPatch(host, patch string, success bool, msg string) error {
	et := EventPatchApplied
	if !success {
		et = EventPatchFailed
	}
	return l.Log(Event{
		Type:    et,
		Host:    host,
		Patch:   patch,
		Success: success,
		Message: msg,
	})
}
