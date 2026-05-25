package signal

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// ShutdownCoordinator combines a Handler with a configurable drain timeout.
// After a signal is received, it waits up to DrainTimeout for work to finish
// before forcefully cancelling a secondary context.
type ShutdownCoordinator struct {
	handler      *Handler
	drainTimeout time.Duration
	out          io.Writer
}

// NewShutdownCoordinator creates a coordinator with the given drain timeout.
// If out is nil it defaults to os.Stdout.
func NewShutdownCoordinator(out io.Writer, drainTimeout time.Duration, cancel context.CancelFunc, sigs ...os.Signal) *ShutdownCoordinator {
	if out == nil {
		out = os.Stdout
	}
	return &ShutdownCoordinator{
		handler:      New(out, cancel, sigs...),
		drainTimeout: drainTimeout,
		out:          out,
	}
}

// Start begins listening for signals. It returns a channel that is closed
// once either the work context finishes or the drain timeout expires.
func (sc *ShutdownCoordinator) Start(workCtx context.Context) <-chan struct{} {
	sc.handler.Listen()
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		select {
		case <-sc.handler.Done():
			// signal received; wait for drain or timeout
			select {
			case <-workCtx.Done():
				fmt.Fprintf(sc.out, "[signal] all work drained cleanly\n")
			case <-time.After(sc.drainTimeout):
				fmt.Fprintf(sc.out, "[signal] drain timeout exceeded — forcing shutdown\n")
			}
		case <-workCtx.Done():
			// work finished before any signal
		}
	}()
	return finished
}
