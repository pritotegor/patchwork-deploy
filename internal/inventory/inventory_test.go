package inventory_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/patchwork-deploy/internal/inventory"
)

func writeInventory(t *testing.T, v any) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	p := filepath.Join(t.TempDir(), "inventory.json")
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	return p
}

func TestLoad_ValidInventory(t *testing.T) {
	path := writeInventory(t, map[string]any{
		"hosts": []map[string]any{
			{"name": "web-01", "address": "10.0.0.1", "user": "deploy", "port": 22, "tags": []string{"web"}},
		},
	})
	inv, err := inventory.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(inv.Hosts) != 1 {
		t.Fatalf("expected 1 host, got %d", len(inv.Hosts))
	}
	if inv.Hosts[0].Name != "web-01" {
		t.Errorf("expected name web-01, got %q", inv.Hosts[0].Name)
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	path := writeInventory(t, map[string]any{
		"hosts": []map[string]any{
			{"name": "db-01", "address": "10.0.1.1"},
		},
	})
	inv, err := inventory.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	h := inv.Hosts[0]
	if h.User != "root" {
		t.Errorf("expected default user root, got %q", h.User)
	}
	if h.Port != 22 {
		t.Errorf("expected default port 22, got %d", h.Port)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := inventory.Load("/no/such/file.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(p, []byte("not-json"), 0o644)
	_, err := inventory.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoad_MissingName(t *testing.T) {
	path := writeInventory(t, map[string]any{
		"hosts": []map[string]any{
			{"address": "10.0.0.1"},
		},
	})
	_, err := inventory.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing name")
	}
}

func TestFilterByTag_NoTags_ReturnsAll(t *testing.T) {
	inv := &inventory.Inventory{
		Hosts: []inventory.Host{
			{Name: "a", Address: "1.1.1.1", Tags: []string{"web"}},
			{Name: "b", Address: "1.1.1.2", Tags: []string{"db"}},
		},
	}
	got := inv.FilterByTag(nil)
	if len(got) != 2 {
		t.Errorf("expected 2 hosts, got %d", len(got))
	}
}

func TestFilterByTag_MatchesSubset(t *testing.T) {
	inv := &inventory.Inventory{
		Hosts: []inventory.Host{
			{Name: "a", Address: "1.1.1.1", Tags: []string{"web"}},
			{Name: "b", Address: "1.1.1.2", Tags: []string{"db"}},
			{Name: "c", Address: "1.1.1.3", Tags: []string{"web", "db"}},
		},
	}
	got := inv.FilterByTag([]string{"db"})
	if len(got) != 2 {
		t.Errorf("expected 2 hosts, got %d", len(got))
	}
}
