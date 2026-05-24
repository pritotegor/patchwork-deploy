// Package retry implements a simple retry policy for use in patchwork-deploy.
//
// It is designed to wrap operations that may fail transiently, such as
// establishing SSH connections to remote hosts or executing shell patches
// where the remote environment may not be immediately ready.
//
// Usage:
//
//	p := retry.DefaultPolicy()
//	p.MaxAttempts = 5
//	p.Delay = time.Second
//
//	err := p.Do("connect to host", func() error {
//		return client.Connect()
//	})
//
All retry attempts are logged to the configured writer (defaults to os.Stdout).
package retry
