// Package throttle provides concurrency limiting for parallel patch deployments.
// It ensures no more than a configured number of hosts are targeted simultaneously.
package throttle

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

// Throttle controls the maximum number of concurrent operations.
type Throttle struct {
	sem    chan struct{}
	mu     sync.Mutex
	active int
	out    io.Writer
}

// New creates a Throttle allowing up to maxConcurrent simultaneous acquisitions.
// If maxConcurrent is less than 1 it defaults to 1.
func New(maxConcurrent int, out io.Writer) *Throttle {
	if maxConcurrent < 1 {
		maxConcurrent = 1
	}
	if out == nil {
		out = os.Stdout
	}
	return &Throttle{
		sem: make(chan struct{}, maxConcurrent),
		out: out,
	}
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Returns an error if the context is done before a slot is acquired.
func (t *Throttle) Acquire(ctx context.Context, host string) error {
	select {
	case t.sem <- struct{}{}:
		t.mu.Lock()
		t.active++
		fmt.Fprintf(t.out, "[throttle] acquired slot for %s (active: %d/%d)\n", host, t.active, cap(t.sem))
		t.mu.Unlock()
		return nil
	case <-ctx.Done():
		return fmt.Errorf("throttle: context cancelled waiting for slot for %s: %w", host, ctx.Err())
	}
}

// Release frees a previously acquired slot.
func (t *Throttle) Release(host string) {
	select {
	case <-t.sem:
		t.mu.Lock()
		t.active--
		fmt.Fprintf(t.out, "[throttle] released slot for %s (active: %d/%d)\n", host, t.active, cap(t.sem))
		t.mu.Unlock()
	default:
	}
}

// Active returns the number of currently held slots.
func (t *Throttle) Active() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.active
}

// Capacity returns the maximum number of concurrent slots.
func (t *Throttle) Capacity() int {
	return cap(t.sem)
}
