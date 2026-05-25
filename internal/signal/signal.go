// Package signal provides graceful shutdown handling for patchwork-deploy.
// It listens for OS interrupt signals and coordinates a clean teardown,
// allowing in-flight patch operations to complete or be rolled back.
package signal

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// Handler listens for OS signals and cancels a context on receipt.
type Handler struct {
	out    io.Writer
	sigs   []os.Signal
	cancel context.CancelFunc
	done   chan struct{}
}

// New creates a Handler that writes log lines to out.
// If out is nil it defaults to os.Stdout.
func New(out io.Writer, cancel context.CancelFunc, sigs ...os.Signal) *Handler {
	if out == nil {
		out = os.Stdout
	}
	if len(sigs) == 0 {
		sigs = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	return &Handler{
		out:    out,
		sigs:   sigs,
		cancel: cancel,
		done:   make(chan struct{}),
	}
}

// Listen starts a goroutine that waits for the first signal.
// When received it logs the signal, calls cancel, and closes the done channel.
func (h *Handler) Listen() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, h.sigs...)
	go func() {
		defer close(h.done)
		defer signal.Stop(ch)
		select {
		case sig := <-ch:
			fmt.Fprintf(h.out, "[signal] received %s — initiating graceful shutdown\n", sig)
			h.cancel()
		}
	}()
}

// Done returns a channel that is closed once the handler has fired.
func (h *Handler) Done() <-chan struct{} {
	return h.done
}
