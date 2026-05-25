// Package progress provides a simple progress tracker for patch deployment runs,
// reporting per-host and per-patch status to an io.Writer.
package progress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Event describes a single progress update.
type Event struct {
	Host      string
	Patch     string
	Status    string // "started", "ok", "failed", "skipped"
	Elapsed   time.Duration
	Err       error
}

// Tracker records and prints progress events.
type Tracker struct {
	mu  sync.Mutex
	out io.Writer
}

// New returns a Tracker that writes to out.
// If out is nil, os.Stdout is used.
func New(out io.Writer) *Tracker {
	if out == nil {
		out = os.Stdout
	}
	return &Tracker{out: out}
}

// Record emits a formatted progress line for the given event.
func (t *Tracker) Record(e Event) {
	t.mu.Lock()
	defer t.mu.Unlock()

	switch e.Status {
	case "started":
		fmt.Fprintf(t.out, "[progress] host=%s patch=%s status=started\n", e.Host, e.Patch)
	case "ok":
		fmt.Fprintf(t.out, "[progress] host=%s patch=%s status=ok elapsed=%s\n", e.Host, e.Patch, e.Elapsed.Round(time.Millisecond))
	case "failed":
		fmt.Fprintf(t.out, "[progress] host=%s patch=%s status=failed elapsed=%s err=%v\n", e.Host, e.Patch, e.Elapsed.Round(time.Millisecond), e.Err)
	case "skipped":
		fmt.Fprintf(t.out, "[progress] host=%s patch=%s status=skipped\n", e.Host, e.Patch)
	default:
		fmt.Fprintf(t.out, "[progress] host=%s patch=%s status=%s\n", e.Host, e.Patch, e.Status)
	}
}

// Started is a convenience wrapper for a "started" event.
func (t *Tracker) Started(host, patch string) {
	t.Record(Event{Host: host, Patch: patch, Status: "started"})
}

// Done is a convenience wrapper for an "ok" or "failed" event.
func (t *Tracker) Done(host, patch string, elapsed time.Duration, err error) {
	if err != nil {
		t.Record(Event{Host: host, Patch: patch, Status: "failed", Elapsed: elapsed, Err: err})
		return
	}
	t.Record(Event{Host: host, Patch: patch, Status: "ok", Elapsed: elapsed})
}

// Skipped is a convenience wrapper for a "skipped" event.
func (t *Tracker) Skipped(host, patch string) {
	t.Record(Event{Host: host, Patch: patch, Status: "skipped"})
}
