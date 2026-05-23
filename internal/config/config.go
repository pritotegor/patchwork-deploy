package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Host represents a remote host configuration.
type Host struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	User    string `json:"user"`
	KeyFile string `json:"key_file"`
	Port    int    `json:"port"`
}

// Config holds the full deployment configuration.
type Config struct {
	PatchDir string `json:"patch_dir"`
	Hosts    []Host `json:"hosts"`
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// validate checks required fields in the config.
func (c *Config) validate() error {
	if c.PatchDir == "" {
		return fmt.Errorf("patch_dir must not be empty")
	}
	if len(c.Hosts) == 0 {
		return fmt.Errorf("at least one host must be defined")
	}
	for i, h := range c.Hosts {
		if h.Name == "" {
			return fmt.Errorf("host[%d]: name must not be empty", i)
		}
		if h.Address == "" {
			return fmt.Errorf("host[%d]: address must not be empty", i)
		}
		if h.User == "" {
			return fmt.Errorf("host[%d]: user must not be empty", i)
		}
	}
	return nil
}

// DefaultPort returns the port for a host, defaulting to 22.
func (h *Host) DefaultPort() int {
	if h.Port == 0 {
		return 22
	}
	return h.Port
}
