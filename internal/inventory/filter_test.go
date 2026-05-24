package inventory

import (
	"testing"
)

func baseHosts() []Host {
	return []Host{
		{Address: "web1", Group: "web", Tags: []string{"prod", "web"}},
		{Address: "web2", Group: "web", Tags: []string{"staging", "web"}},
		{Address: "db1", Group: "db", Tags: []string{"prod", "db"}},
		{Address: "db2", Group: "db", Tags: []string{"staging", "db"}},
	}
}

func TestFilterByTag_EmptyTagsReturnsAll(t *testing.T) {
	hosts := baseHosts()
	got := FilterByTag(hosts, nil)
	if len(got) != len(hosts) {
		t.Fatalf("expected %d hosts, got %d", len(hosts), len(got))
	}
}

func TestFilterByTag_SingleTag(t *testing.T) {
	got := FilterByTag(baseHosts(), []string{"prod"})
	if len(got) != 2 {
		t.Fatalf("expected 2 prod hosts, got %d", len(got))
	}
	for _, h := range got {
		if !hasAllTags(h, []string{"prod"}) {
			t.Errorf("host %s missing tag 'prod'", h.Address)
		}
	}
}

func TestFilterByTag_MultipleTags(t *testing.T) {
	got := FilterByTag(baseHosts(), []string{"prod", "web"})
	if len(got) != 1 {
		t.Fatalf("expected 1 host, got %d", len(got))
	}
	if got[0].Address != "web1" {
		t.Errorf("expected web1, got %s", got[0].Address)
	}
}

func TestFilterByTag_NoMatch(t *testing.T) {
	got := FilterByTag(baseHosts(), []string{"nonexistent"})
	if len(got) != 0 {
		t.Fatalf("expected 0 hosts, got %d", len(got))
	}
}

func TestFilterByGroup_EmptyGroupReturnsAll(t *testing.T) {
	hosts := baseHosts()
	got := FilterByGroup(hosts, "")
	if len(got) != len(hosts) {
		t.Fatalf("expected %d hosts, got %d", len(hosts), len(got))
	}
}

func TestFilterByGroup_MatchingGroup(t *testing.T) {
	got := FilterByGroup(baseHosts(), "db")
	if len(got) != 2 {
		t.Fatalf("expected 2 db hosts, got %d", len(got))
	}
	for _, h := range got {
		if h.Group != "db" {
			t.Errorf("expected group db, got %s", h.Group)
		}
	}
}

func TestFilter_CombinesGroupAndTags(t *testing.T) {
	opts := FilterOptions{Group: "db", Tags: []string{"prod"}}
	got := Filter(baseHosts(), opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 host, got %d", len(got))
	}
	if got[0].Address != "db1" {
		t.Errorf("expected db1, got %s", got[0].Address)
	}
}

func TestFilter_EmptyOptions_ReturnsAll(t *testing.T) {
	hosts := baseHosts()
	got := Filter(hosts, FilterOptions{})
	if len(got) != len(hosts) {
		t.Fatalf("expected %d hosts, got %d", len(hosts), len(got))
	}
}
