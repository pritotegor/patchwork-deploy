package throttle_test

import (
	"bytes"
	"context"
	"sync"
	"testing"
	"time"

	"github.com/patchwork-deploy/internal/throttle"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	th := throttle.New(2, nil)
	if th == nil {
		t.Fatal("expected non-nil throttle")
	}
	if th.Capacity() != 2 {
		t.Errorf("expected capacity 2, got %d", th.Capacity())
	}
}

func TestNew_MinimumCapacityIsOne(t *testing.T) {
	th := throttle.New(0, nil)
	if th.Capacity() != 1 {
		t.Errorf("expected capacity 1 for zero input, got %d", th.Capacity())
	}
}

func TestAcquireRelease_TracksActive(t *testing.T) {
	var buf bytes.Buffer
	th := throttle.New(3, &buf)

	ctx := context.Background()
	if err := th.Acquire(ctx, "host-a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th.Active() != 1 {
		t.Errorf("expected 1 active, got %d", th.Active())
	}
	th.Release("host-a")
	if th.Active() != 0 {
		t.Errorf("expected 0 active after release, got %d", th.Active())
	}
}

func TestAcquire_BlocksAtCapacity(t *testing.T) {
	var buf bytes.Buffer
	th := throttle.New(1, &buf)
	ctx := context.Background()

	if err := th.Acquire(ctx, "host-1"); err != nil {
		t.Fatal(err)
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := th.Acquire(ctxTimeout, "host-2")
	if err == nil {
		t.Error("expected error when context times out waiting for slot")
	}
	th.Release("host-1")
}

func TestConcurrent_NeverExceedsCapacity(t *testing.T) {
	var buf bytes.Buffer
	const capacity = 3
	const goroutines = 10
	th := throttle.New(capacity, &buf)

	var mu sync.Mutex
	peak := 0
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			host := "host"
			_ = th.Acquire(context.Background(), host)
			mu.Lock()
			if a := th.Active(); a > peak {
				peak = a
			}
			mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			th.Release(host)
		}(i)
	}
	wg.Wait()

	if peak > capacity {
		t.Errorf("peak active %d exceeded capacity %d", peak, capacity)
	}
}

func TestRelease_Idempotent(t *testing.T) {
	var buf bytes.Buffer
	th := throttle.New(2, &buf)
	// Release without acquire should not panic
	th.Release("ghost-host")
	if th.Active() != 0 {
		t.Errorf("expected 0 active, got %d", th.Active())
	}
}
