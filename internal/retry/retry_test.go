package retry_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/patchwork-deploy/internal/retry"
)

func TestDefaultPolicy_HasSensibleDefaults(t *testing.T) {
	p := retry.DefaultPolicy()
	if p.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", p.MaxAttempts)
	}
	if p.Delay != 2*time.Second {
		t.Errorf("expected Delay=2s, got %v", p.Delay)
	}
	if p.Out == nil {
		t.Error("expected non-nil Out")
	}
}

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	var buf bytes.Buffer
	p := retry.Policy{MaxAttempts: 3, Delay: 0, Out: &buf}

	calls := 0
	err := p.Do("test-op", func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output on success, got: %s", buf.String())
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	var buf bytes.Buffer
	p := retry.Policy{MaxAttempts: 3, Delay: 0, Out: &buf}

	calls := 0
	sentinel := errors.New("transient error")
	err := p.Do("test-op", func() error {
		calls++
		if calls < 3 {
			return sentinel
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ReturnsErrorAfterAllAttempts(t *testing.T) {
	var buf bytes.Buffer
	p := retry.Policy{MaxAttempts: 2, Delay: 0, Out: &buf}

	sentinel := errors.New("persistent error")
	err := p.Do("failing-op", func() error {
		return sentinel
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("expected wrapped sentinel error, got: %v", err)
	}
	if !strings.Contains(buf.String(), "failing-op") {
		t.Errorf("expected op name in output, got: %s", buf.String())
	}
}

func TestDo_ZeroMaxAttemptsRunsOnce(t *testing.T) {
	var buf bytes.Buffer
	p := retry.Policy{MaxAttempts: 0, Delay: 0, Out: &buf}

	calls := 0
	_ = p.Do("zero-op", func() error {
		calls++
		return errors.New("fail")
	})
	if calls != 1 {
		t.Errorf("expected exactly 1 call for MaxAttempts=0, got %d", calls)
	}
}
