// Package runner ties together the patch loader, SSH client, and patch executor
// to implement the top-level deployment pipeline.
//
// Typical usage:
//
//	cfg, err := config.Load("deploy.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	r := runner.New(cfg)
//	if err := r.Run(); err != nil {
//		log.Fatal(err)
//	}
//
// Run will:
//  1. Load and sort all *.sh files from cfg.PatchDir.
//  2. For each host in cfg.Hosts, open an SSH connection.
//  3. Apply each patch script in order, stopping on the first non-zero exit.
package runner
