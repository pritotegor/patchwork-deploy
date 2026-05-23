// Package notify provides lightweight deployment notification support
// for patchwork-deploy.
//
// A Notifier writes structured, human-readable event lines to any
// io.Writer. Each line includes a UTC timestamp, severity level,
// target host, and a free-form message.
//
// Example output:
//
//	[2024-01-15T10:30:00Z] INFO  host=web-01 msg=patch applied successfully
//	[2024-01-15T10:30:01Z] ERROR host=web-02 msg=connection refused
//
// Usage:
//
//	n := notify.New(os.Stdout)
//	n.Info("web-01", "starting deployment")
//	n.Error("web-02", "patch 003 failed")
package notify
