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
//
// Error handling:
//
// If a patch script exits with a non-zero status on any host, Run returns a
// [PatchError] that captures the host address, the patch filename, and the
// remote exit code, allowing callers to distinguish patch failures from
// connectivity or configuration errors.
package runner
