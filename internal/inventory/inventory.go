// Package inventory manages the list of remote hosts to deploy to,
// loading and validating host definitions from a JSON file.
package inventory

import (
	"encoding/json"
	"fmt"
	"os"
)

// Host represents a single remote target.
type Host struct {
	Name    string   `json:"name"`
	Address string   `json:"address"`
	User    string   `json:"user"`
	Port    int      `json:"port"`
	Tags    []string `json:"tags"`
}

// Inventory holds a collection of hosts.
type Inventory struct {
	Hosts []Host `json:"hosts"`
}

// Load reads an inventory file from the given path and returns the parsed
// Inventory. An error is returned if the file cannot be read or parsed.
func Load(path string) (*Inventory, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("inventory: read file: %w", err)
	}

	var inv Inventory
	if err := json.Unmarshal(data, &inv); err != nil {
		return nil, fmt.Errorf("inventory: parse: %w", err)
	}

	if err := validate(&inv); err != nil {
		return nil, err
	}

	return &inv, nil
}

// FilterByTag returns only the hosts that carry at least one of the given tags.
// If tags is empty, all hosts are returned.
func (inv *Inventory) FilterByTag(tags []string) []Host {
	if len(tags) == 0 {
		return inv.Hosts
	}
	want := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		want[t] = struct{}{}
	}
	var out []Host
	for _, h := range inv.Hosts {
		for _, t := range h.Tags {
			if _, ok := want[t]; ok {
				out = append(out, h)
				break
			}
		}
	}
	return out
}

func validate(inv *Inventory) error {
	for i, h := range inv.Hosts {
		if h.Name == "" {
			return fmt.Errorf("inventory: host[%d] missing name", i)
		}
		if h.Address == "" {
			return fmt.Errorf("inventory: host %q missing address", h.Name)
		}
		if h.User == "" {
			inv.Hosts[i].User = "root"
		}
		if h.Port == 0 {
			inv.Hosts[i].Port = 22
		}
	}
	return nil
}
