// Package progress provides a lightweight, concurrency-safe progress tracker
// for patchwork-deploy deployment runs.
//
// It emits structured log lines for each patch/host transition — started, ok,
// failed, or skipped — making it easy to tail a deployment in real time or
// feed the output into a log aggregator.
//
// Basic usage:
//
//	tracker := progress.New(os.Stdout)
//	tracker.Started("web-01", "001_create_user.sh")
//	// ... apply patch ...
//	tracker.Done("web-01", "001_create_user.sh", elapsed, err)
package progress
