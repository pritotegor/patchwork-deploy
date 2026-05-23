// Package notify provides simple deployment notification support,
// allowing patchwork-deploy to emit status messages to configurable
// output channels (e.g. stdout, webhook, log file).
package notify

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Event holds the data for a single notification.
type Event struct {
	Timestamp time.Time
	Level     Level
	Host      string
	Message   string
}

// Notifier sends deployment events to an output destination.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier writing to out. If out is nil, os.Stdout is used.
func New(out io.Writer) *Notifier {
	if out == nil {
		out = os.Stdout
	}
	return &Notifier{out: out}
}

// Send emits a notification event to the configured writer.
func (n *Notifier) Send(level Level, host, message string) error {
	e := Event{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Host:      host,
		Message:   message,
	}
	return n.write(e)
}

// Info is a convenience wrapper for LevelInfo notifications.
func (n *Notifier) Info(host, message string) error {
	return n.Send(LevelInfo, host, message)
}

// Error is a convenience wrapper for LevelError notifications.
func (n *Notifier) Error(host, message string) error {
	return n.Send(LevelError, host, message)
}

func (n *Notifier) write(e Event) error {
	_, err := fmt.Fprintf(
		n.out,
		"[%s] %s host=%s msg=%s\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Host,
		e.Message,
	)
	return err
}
