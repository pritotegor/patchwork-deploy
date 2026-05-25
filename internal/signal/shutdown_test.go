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

func TestNewShutdownCoordinator_DefaultsToStdout(t *testing.T) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	sc := signal.NewShutdownCoordinator(nil, time.Second, cancel)
	if sc == nil {
		t.Fatal("expected non-nil coordinator")
	}
}

func TestStart_FinishesWhenWorkContextDone(t *testing.T) {
	var buf bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	sc := signal.NewShutdownCoordinator(&buf, 5*time.Second, cancel)
	finished := sc.Start(ctx)

	// cancel work context directly (no signal)
	cancel()

	select {
	case <-finished:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("coordinator did not finish after work context cancelled")
	}
}

func TestStart_DrainTimeoutForcesShutdown(t *testing.T) {
	var buf bytes.Buffer
	workCtx, workCancel := context.WithCancel(context.Background())
	defer workCancel()

	sigCtx, sigCancel := context.WithCancel(context.Background())
	sc := signal.NewShutdownCoordinator(&buf, 100*time.Millisecond, sigCancel, syscall.SIGUSR1)
	finished := sc.Start(workCtx)

	_ = sigCtx
	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1) //nolint:errcheck

	select {
	case <-finished:
		// drain timeout should have fired
	case <-time.After(2 * time.Second):
		t.Fatal("coordinator did not finish after drain timeout")
	}

	if !strings.Contains(buf.String(), "drain timeout") {
		t.Errorf("expected drain timeout message, got: %q", buf.String())
	}
}
