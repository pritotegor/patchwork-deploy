package inventory_test

import (
	"testing"

	"github.com/patchwork-deploy/internal/inventory"
)

// TestLoad_ExampleInventory verifies the checked-in example file is valid.
func TestLoad_ExampleInventory(t *testing.T) {
	inv, err := inventory.Load("example_inventory.json")
	if err != nil {
		t.Fatalf("example inventory failed to load: %v", err)
	}
	if len(inv.Hosts) == 0 {
		t.Fatal("example inventory has no hosts")
	}
	for _, h := range inv.Hosts {
		if h.Port == 0 {
			t.Errorf("host %q has zero port after load", h.Name)
		}
		if h.User == "" {
			t.Errorf("host %q has empty user after load", h.Name)
		}
	}
}

func TestFilterByTag_ProdOnly(t *testing.T) {
	inv, err := inventory.Load("example_inventory.json")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	prod := inv.FilterByTag([]string{"prod"})
	if len(prod) != len(inv.Hosts) {
		t.Errorf("expected all %d hosts to be prod-tagged, got %d", len(inv.Hosts), len(prod))
	}
}

func TestFilterByTag_WebOnly(t *testing.T) {
	inv, err := inventory.Load("example_inventory.json")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	web := inv.FilterByTag([]string{"web"})
	for _, h := range web {
		found := false
		for _, tag := range h.Tags {
			if tag == "web" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("host %q returned by web filter but has no web tag", h.Name)
		}
	}
}
