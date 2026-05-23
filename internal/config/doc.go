// Package config provides loading and validation of patchwork-deploy
// configuration files.
//
// Configuration is stored as JSON and describes the patch directory
// to apply as well as the list of remote hosts to target.
//
// Example usage:
//
//	cfg, err := config.Load("deploy.json")
//	if err != nil {
//		log.Fatalf("failed to load config: %v", err)
//	}
//	for _, host := range cfg.Hosts {
//		fmt.Printf("deploying to %s (%s)\n", host.Name, host.Address)
//	}
package config
