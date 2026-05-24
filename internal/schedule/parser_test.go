package schedule

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempSchedule(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "schedule.json")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp schedule: %v", err)
	}
	return p
}

func TestLoadFile_ValidSchedule(t *testing.T) {
	path := writeTempSchedule(t, `{
		"windows": [{"weekdays":["monday","tuesday"],"start":"09:00","end":"17:00"}]
	}`)
	s, err := LoadFile(path, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Windows) != 1 {
		t.Fatalf("expected 1 window, got %d", len(s.Windows))
	}
	if len(s.Windows[0].Weekdays) != 2 {
		t.Errorf("expected 2 weekdays, got %d", len(s.Windows[0].Weekdays))
	}
	if s.Windows[0].Weekdays[0] != time.Monday {
		t.Errorf("expected Monday, got %v", s.Windows[0].Weekdays[0])
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	_, err := LoadFile("/no/such/file.json", &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	path := writeTempSchedule(t, `not json`)
	_, err := LoadFile(path, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadFile_UnknownWeekday(t *testing.T) {
	path := writeTempSchedule(t, `{"windows":[{"weekdays":["funday"],"start":"09:00","end":"17:00"}]}`)
	_, err := LoadFile(path, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for unknown weekday")
	}
}

func TestLoadFile_EmptyWindowsAllowsAll(t *testing.T) {
	path := writeTempSchedule(t, `{"windows":[]}`)
	s, err := LoadFile(path, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Allowed(time.Now()) {
		t.Error("empty windows should allow any time")
	}
}
