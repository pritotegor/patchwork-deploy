package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/patchwork-deploy/internal/config"
	"github.com/yourorg/patchwork-deploy/internal/runner"
)

const version = "0.1.0"

func main() {
	var (
		configPath  = flag.String("config", "patchwork.json", "path to config file")
		showVersion = flag.Bool("version", false, "print version and exit")
		dryRun      = flag.Bool("dry-run", false, "load and validate config without applying patches")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("patchwork-deploy v%s\n", version)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	if *dryRun {
		fmt.Printf("config loaded successfully\n")
		fmt.Printf("  patch dir : %s\n", cfg.PatchDir)
		fmt.Printf("  hosts     : %d\n", len(cfg.Hosts))
		for _, h := range cfg.Hosts {
			fmt.Printf("    - %s@%s:%d\n", h.User, h.Host, h.Port)
		}
		os.Exit(0)
	}

	r := runner.New(cfg, os.Stdout)
	if err := r.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "deployment failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("deployment complete")
}
