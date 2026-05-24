// Package schedule provides time-window based gating for patch deployments.
// A schedule defines allowed deployment windows; patches submitted outside
// a window are rejected with a clear error.
package schedule

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Window represents a single allowed deployment window.
type Window struct {
	// Weekdays lists allowed days (0=Sunday … 6=Saturday). Empty means all days.
	Weekdays []time.Weekday `json:"weekdays"`
	// Start and End are clock times in "15:04" format (24-hour, UTC).
	Start string `json:"start"`
	End   string `json:"end"`
}

// Schedule holds one or more deployment windows.
type Schedule struct {
	Windows []Window `json:"windows"`
	out     io.Writer
}

// New returns a Schedule that writes log lines to out.
// If out is nil, os.Stdout is used.
func New(windows []Window, out io.Writer) *Schedule {
	if out == nil {
		out = os.Stdout
	}
	return &Schedule{Windows: windows, out: out}
}

// Allowed reports whether t falls inside any of the configured windows.
// An empty window list always allows.
func (s *Schedule) Allowed(t time.Time) bool {
	if len(s.Windows) == 0 {
		return true
	}
	utc := t.UTC()
	for _, w := range s.Windows {
		if inWindow(utc, w) {
			fmt.Fprintf(s.out, "[schedule] %s is within window %s-%s\n",
				utc.Format(time.RFC3339), w.Start, w.End)
			return true
		}
	}
	return false
}

// Check returns an error when t is outside all windows.
func (s *Schedule) Check(t time.Time) error {
	if s.Allowed(t) {
		return nil
	}
	return fmt.Errorf("deployment blocked: %s is outside all configured windows",
		t.UTC().Format(time.RFC3339))
}

func inWindow(t time.Time, w Window) bool {
	if len(w.Weekdays) > 0 {
		allowed := false
		for _, d := range w.Weekdays {
			if t.Weekday() == d {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}
	start, err1 := parseTime(t, w.Start)
	end, err2 := parseTime(t, w.End)
	if err1 != nil || err2 != nil {
		return false
	}
	return !t.Before(start) && t.Before(end)
}

func parseTime(base time.Time, s string) (time.Time, error) {
	if s == "" {
		return base, nil
	}
	parsed, err := time.Parse("15:04", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time %q: %w", s, err)
	}
	return time.Date(base.Year(), base.Month(), base.Day(),
		parsed.Hour(), parsed.Minute(), 0, 0, time.UTC), nil
}
