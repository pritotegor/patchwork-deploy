package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/diff"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	d := diff.New(nil)
	if d == nil {
		t.Fatal("expected non-nil Differ")
	}
}

func TestCompute_AllNew(t *testing.T) {
	var buf bytes.Buffer
	d := diff.New(&buf)

	entries := d.Compute([]string{"001_init.sh", "002_users.sh"}, nil)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Status != diff.StatusNew {
			t.Errorf("%s: expected StatusNew, got %s", e.Name, e.Status)
		}
	}
	if !strings.Contains(buf.String(), "2 current") {
		t.Errorf("expected log line to mention current count, got: %s", buf.String())
	}
}

func TestCompute_AlreadyApplied(t *testing.T) {
	d := diff.New(new(bytes.Buffer))
	applied := map[string]bool{"001_init.sh": true}

	entries := d.Compute([]string{"001_init.sh", "002_users.sh"}, applied)

	statuses := map[string]diff.Status{}
	for _, e := range entries {
		statuses[e.Name] = e.Status
	}

	if statuses["001_init.sh"] != diff.StatusApplied {
		t.Errorf("001_init.sh: expected StatusApplied, got %s", statuses["001_init.sh"])
	}
	if statuses["002_users.sh"] != diff.StatusNew {
		t.Errorf("002_users.sh: expected StatusNew, got %s", statuses["002_users.sh"])
	}
}

func TestCompute_RemovedPatch(t *testing.T) {
	d := diff.New(new(bytes.Buffer))
	applied := map[string]bool{"001_init.sh": true, "000_old.sh": true}

	entries := d.Compute([]string{"001_init.sh"}, applied)

	statuses := map[string]diff.Status{}
	for _, e := range entries {
		statuses[e.Name] = e.Status
	}

	if statuses["000_old.sh"] != diff.StatusRemoved {
		t.Errorf("000_old.sh: expected StatusRemoved, got %s", statuses["000_old.sh"])
	}
}

func TestCompute_EmptyInputs(t *testing.T) {
	d := diff.New(new(bytes.Buffer))
	entries := d.Compute(nil, nil)
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestStatus_String(t *testing.T) {
	cases := []struct {
		s    diff.Status
		want string
	}{
		{diff.StatusNew, "new"},
		{diff.StatusApplied, "applied"},
		{diff.StatusRemoved, "removed"},
		{diff.Status(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("Status(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}
