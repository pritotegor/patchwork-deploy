// Package config loads and validates the patchwork-deploy configuration file.
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the full deployment configuration.
type Config struct {
	Hosts       []string `json:"hosts"`
	PatchDir    string   `json:"patch_dir"`
	User        string   `json:"user"`
	KeyFile     string   `json:"key_file"`
	AuditLog    string   `json:"audit_log,omitempty"`
	HooksBefore []string `json:"hooks_before,omitempty"`
	HooksAfter  []string `json:"hooks_after,omitempty"`
	DryRun      bool     `json:"dry_run,omitempty"`
	MaxParallel int      `json:"max_parallel,omitempty"`
}

// Load reads and validates a Config from the JSON file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %s: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse %s: %w", path, err)
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	applyDefaults(&cfg)
	return &cfg, nil
}

func validate(cfg *Config) error {
	if cfg.PatchDir == "" {
		return fmt.Errorf("config: missing required field: patch_dir")
	}
	if cfg.User == "" {
		return fmt.Errorf("config: missing required field: user")
	}
	if cfg.KeyFile == "" {
		return fmt.Errorf("config: missing required field: key_file")
	}
	return nil
}

func applyDefaults(cfg *Config) {
	if cfg.MaxParallel < 1 {
		cfg.MaxParallel = 1
	}
}
