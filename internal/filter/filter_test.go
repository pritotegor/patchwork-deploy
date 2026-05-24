package filter_test

import (
	"testing"

	"github.com/example/patchwork-deploy/internal/filter"
)

func makePatches() []filter.Patch {
	return []filter.Patch{
		{Name: "001-init.sh", Tags: []string{"db"}},
		{Name: "002-migrate.sh", Tags: []string{"db", "schema"}},
		{Name: "003-cache.sh", Tags: []string{"cache"}},
		{Name: "004-cleanup.sh", Tags: []string{}},
	}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	patches := makePatches()
	got := filter.Apply(patches, filter.Options{})
	if len(got) != len(patches) {
		t.Fatalf("expected %d patches, got %d", len(patches), len(got))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	got := filter.Apply(makePatches(), filter.Options{Prefix: "00"})
	if len(got) != 4 {
		t.Fatalf("expected 4, got %d", len(got))
	}

	got = filter.Apply(makePatches(), filter.Options{Prefix: "001"})
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
	if got[0].Name != "001-init.sh" {
		t.Errorf("unexpected patch name: %s", got[0].Name)
	}
}

func TestApply_TagFilter_SingleMatch(t *testing.T) {
	got := filter.Apply(makePatches(), filter.Options{Tags: []string{"cache"}})
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
	if got[0].Name != "003-cache.sh" {
		t.Errorf("unexpected patch: %s", got[0].Name)
	}
}

func TestApply_TagFilter_MultipleMatches(t *testing.T) {
	got := filter.Apply(makePatches(), filter.Options{Tags: []string{"db"}})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestApply_TagFilter_NoMatch(t *testing.T) {
	got := filter.Apply(makePatches(), filter.Options{Tags: []string{"nonexistent"}})
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestApply_PrefixAndTagCombined(t *testing.T) {
	got := filter.Apply(makePatches(), filter.Options{
		Prefix: "002",
		Tags:   []string{"schema"},
	})
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
	if got[0].Name != "002-migrate.sh" {
		t.Errorf("unexpected patch: %s", got[0].Name)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	got := filter.Apply(nil, filter.Options{Tags: []string{"db"}})
	if got != nil && len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}
