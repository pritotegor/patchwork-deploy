// Package retry provides configurable retry logic for operations
// that may fail transiently, such as SSH connections or patch execution.
package retry

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Policy defines how retries are performed.
type Policy struct {
	MaxAttempts int
	Delay       time.Duration
	Out         io.Writer
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		Delay:       2 * time.Second,
		Out:         os.Stdout,
	}
}

// Do executes fn up to MaxAttempts times, waiting Delay between attempts.
// It returns nil on the first success, or the last error if all attempts fail.
func (p Policy) Do(name string, fn func() error) error {
	if p.MaxAttempts < 1 {
		p.MaxAttempts = 1
	}
	if p.Out == nil {
		p.Out = os.Stdout
	}

	var lastErr error
	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		fmt.Fprintf(p.Out, "[retry] %s: attempt %d/%d failed: %v\n",
			name, attempt, p.MaxAttempts, lastErr)
		if attempt < p.MaxAttempts {
			time.Sleep(p.Delay)
		}
	}
	return fmt.Errorf("%s failed after %d attempts: %w", name, p.MaxAttempts, lastErr)
}
