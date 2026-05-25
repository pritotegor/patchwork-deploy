// Package timeout provides per-patch and per-host execution time limits.
// A Limiter wraps an operation with a context deadline and reports
// whether the operation exceeded the allowed duration.
package timeout

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// Limiter enforces a maximum duration on patch execution.
type Limiter struct {
	Duration time.Duration
	out       io.Writer
}

// New returns a Limiter with the given duration.
// If out is nil it defaults to os.Stdout.
func New(d time.Duration, out io.Writer) *Limiter {
	if out == nil {
		out = os.Stdout
	}
	return &Limiter{Duration: d, out: out}
}

// Do runs fn within the configured deadline.
// It returns an error if fn returns an error or if the deadline is exceeded.
func (l *Limiter) Do(ctx context.Context, label string, fn func(ctx context.Context) error) error {
	if l.Duration <= 0 {
		return fn(ctx)
	}

	ctx, cancel := context.WithTimeout(ctx, l.Duration)
	defer cancel()

	type result struct {
		err error
	}
	ch := make(chan result, 1)

	go func() {
		ch <- result{err: fn(ctx)}
	}()

	select {
	case r := <-ch:
		if r.err != nil {
			fmt.Fprintf(l.out, "[timeout] %s failed: %v\n", label, r.err)
		}
		return r.err
	case <-ctx.Done():
		err := fmt.Errorf("timeout: %s exceeded %s", label, l.Duration)
		fmt.Fprintf(l.out, "[timeout] %v\n", err)
		return err
	}
}
