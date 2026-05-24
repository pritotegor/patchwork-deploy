package ratelimit

import (
	"bytes"
	"testing"
	"time"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	l, err := New(1.0, 5, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestNew_InvalidRateReturnsError(t *testing.T) {
	_, err := New(0, 5, nil)
	if err == nil {
		t.Fatal("expected error for zero rate")
	}
	_, err = New(-1, 5, nil)
	if err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestNew_MinimumBurstIsOne(t *testing.T) {
	l, err := New(1.0, 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.max != 1 {
		t.Fatalf("expected max=1, got %v", l.max)
	}
}

func TestTryAcquire_ConsumesToken(t *testing.T) {
	var buf bytes.Buffer
	l, _ := New(1.0, 3, &buf)

	if !l.TryAcquire() {
		t.Fatal("expected first acquire to succeed")
	}
	if l.tokens != 2 {
		t.Fatalf("expected 2 tokens remaining, got %v", l.tokens)
	}
}

func TestTryAcquire_FailsWhenEmpty(t *testing.T) {
	var buf bytes.Buffer
	l, _ := New(0.001, 1, &buf) // very slow refill

	if !l.TryAcquire() {
		t.Fatal("expected first acquire to succeed")
	}
	if l.TryAcquire() {
		t.Fatal("expected second acquire to fail with empty bucket")
	}
	if buf.Len() == 0 {
		t.Fatal("expected output to be written")
	}
}

func TestWait_EventuallyAcquires(t *testing.T) {
	var buf bytes.Buffer
	// 10 tokens/sec, burst 1 — refills quickly
	l, _ := New(10.0, 1, &buf)

	// Drain the bucket.
	l.TryAcquire()

	done := make(chan struct{})
	go func() {
		l.Wait()
		close(done)
	}()

	select {
	case <-done:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("Wait timed out")
	}
}

func TestRefill_DoesNotExceedMax(t *testing.T) {
	var buf bytes.Buffer
	l, _ := New(100.0, 5, &buf)

	// Force lastTick far in the past to simulate long idle period.
	l.mu.Lock()
	l.tokens = 0
	l.lastTick = time.Now().Add(-10 * time.Second)
	l.mu.Unlock()

	l.TryAcquire()

	l.mu.Lock()
	defer l.mu.Unlock()
	if l.tokens > l.max {
		t.Fatalf("tokens %v exceeded max %v", l.tokens, l.max)
	}
}
