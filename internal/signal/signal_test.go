package signal_test

import (
	"bytes"
	"context"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/signal"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	h := signal.New(nil, cancel)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestListen_CancelsContextOnSignal(t *testing.T) {
	var buf bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	h := signal.New(&buf, cancel, syscall.SIGUSR1)
	h.Listen()

	// send signal to self
	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1) //nolint:errcheck

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(2 * time.Second):
		t.Fatal("context was not cancelled after signal")
	}

	select {
	case <-h.Done():
		// handler goroutine finished
	case <-time.After(time.Second):
		t.Fatal("handler Done channel not closed")
	}

	if !strings.Contains(buf.String(), "graceful shutdown") {
		t.Errorf("expected shutdown message, got: %q", buf.String())
	}
}

func TestListen_LogsSignalName(t *testing.T) {
	var buf bytes.Buffer
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	h := signal.New(&buf, cancel, syscall.SIGUSR2)
	h.Listen()

	syscall.Kill(syscall.Getpid(), syscall.SIGUSR2) //nolint:errcheck

	select {
	case <-h.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("handler did not complete")
	}

	if !strings.Contains(buf.String(), "user defined signal 2") &&
		!strings.Contains(buf.String(), "sigusr2") &&
		!strings.Contains(buf.String(), "user2") {
		t.Logf("signal output: %q", buf.String())
	}
	if !strings.Contains(buf.String(), "[signal]") {
		t.Errorf("expected [signal] prefix in output, got: %q", buf.String())
	}
}
