package schedule

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func monday1000() time.Time {
	// 2024-01-08 is a Monday
	return time.Date(2024, 1, 8, 10, 0, 0, 0, time.UTC)
}

func TestNew_DefaultsToStdout(t *testing.T) {
	s := New(nil, nil)
	if s.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestAllowed_EmptyWindowsAlwaysTrue(t *testing.T) {
	s := New(nil, &bytes.Buffer{})
	if !s.Allowed(monday1000()) {
		t.Error("empty window list should always allow")
	}
}

func TestAllowed_WithinWindow(t *testing.T) {
	win := Window{Start: "09:00", End: "17:00"}
	s := New([]Window{win}, &bytes.Buffer{})
	if !s.Allowed(monday1000()) {
		t.Error("10:00 should be inside 09:00-17:00")
	}
}

func TestAllowed_OutsideWindow(t *testing.T) {
	win := Window{Start: "09:00", End: "10:00"}
	s := New([]Window{win}, &bytes.Buffer{})
	// 10:00 is NOT before end (exclusive)
	if s.Allowed(monday1000()) {
		t.Error("10:00 should be outside 09:00-10:00 (exclusive end)")
	}
}

func TestAllowed_WeekdayMatch(t *testing.T) {
	win := Window{Weekdays: []time.Weekday{time.Monday}, Start: "09:00", End: "17:00"}
	s := New([]Window{win}, &bytes.Buffer{})
	if !s.Allowed(monday1000()) {
		t.Error("Monday 10:00 should be allowed")
	}
}

func TestAllowed_WeekdayMismatch(t *testing.T) {
	win := Window{Weekdays: []time.Weekday{time.Saturday}, Start: "09:00", End: "17:00"}
	s := New([]Window{win}, &bytes.Buffer{})
	if s.Allowed(monday1000()) {
		t.Error("Monday should not be allowed when only Saturday is configured")
	}
}

func TestCheck_ReturnsNilWhenAllowed(t *testing.T) {
	s := New(nil, &bytes.Buffer{})
	if err := s.Check(monday1000()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_ReturnsErrorWhenBlocked(t *testing.T) {
	win := Window{Start: "22:00", End: "23:00"}
	s := New([]Window{win}, &bytes.Buffer{})
	err := s.Check(monday1000())
	if err == nil {
		t.Fatal("expected error for blocked time")
	}
	if !strings.Contains(err.Error(), "deployment blocked") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAllowed_LogsMatchingWindow(t *testing.T) {
	var buf bytes.Buffer
	win := Window{Start: "09:00", End: "17:00"}
	s := New([]Window{win}, &buf)
	s.Allowed(monday1000())
	if !strings.Contains(buf.String(), "[schedule]") {
		t.Error("expected log line when window matches")
	}
}
