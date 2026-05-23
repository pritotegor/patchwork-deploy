// Package audit provides structured, append-only audit logging for
// patchwork-deploy deployments.
//
// Events are written as newline-delimited JSON (NDJSON) so they can be
// consumed by standard log-aggregation pipelines (e.g. Loki, Splunk).
//
// Basic usage:
//
//	l := audit.New(os.Stderr)
//	l.Log(audit.Event{
//		Type:    audit.EventDeployStart,
//		Host:    "web-01",
//		Success: true,
//	})
//
// To persist events to disk use NewFileLogger:
//
//	fl, err := audit.NewFileLogger("/var/log/patchwork/audit.jsonl")
//	if err != nil { ... }
//	defer fl.Close()
//	fl.LogPatch("web-01", "001_migrate.sh", true, "")
package audit
