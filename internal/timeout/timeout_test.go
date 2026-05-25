package timeout_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/timeout"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	l := timeout.New(time.Second, nil)
	if l == nil {
		t.Fatal("expected non-nil Limiter")
	}
	if l.Duration != time.Second {
		t.Fatalf("expected 1s, got %v", l.Duration)
	}
}

func TestDo_SuccessWithinDeadline(t *testing.T) {
	var buf bytes.Buffer
	l := timeout.New(500*time.Millisecond, &buf)

	err := l.Do(context.Background(), "patch-01", func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got: %s", buf.String())
	}
}

func TestDo_ExceedsDeadlineReturnsError(t *testing.T) {
	var buf bytes.Buffer
	l := timeout.New(20*time.Millisecond, &buf)

	err := l.Do(context.Background(), "slow-patch", func(ctx context.Context) error {
		select {
		case <-time.After(200 * time.Millisecond):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("expected 'timeout' in error, got: %v", err)
	}
	if !strings.Contains(buf.String(), "slow-patch") {
		t.Fatalf("expected label in output, got: %s", buf.String())
	}
}

func TestDo_FnErrorIsReturned(t *testing.T) {
	var buf bytes.Buffer
	l := timeout.New(time.Second, &buf)
	want := errors.New("script failed")

	err := l.Do(context.Background(), "patch-02", func(ctx context.Context) error {
		return want
	})
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
	if !strings.Contains(buf.String(), "patch-02") {
		t.Fatalf("expected label in output, got: %s", buf.String())
	}
}

func TestDo_ZeroDurationSkipsDeadline(t *testing.T) {
	var buf bytes.Buffer
	l := timeout.New(0, &buf)

	called := false
	err := l.Do(context.Background(), "patch-03", func(ctx context.Context) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected fn to be called")
	}
}
