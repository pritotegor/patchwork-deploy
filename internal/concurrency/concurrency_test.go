package concurrency_test

import (
	"bytes"
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/patchwork-deploy/internal/concurrency"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	p := concurrency.New(2, nil)
	if p == nil {
		t.Fatal("expected non-nil pool")
	}
}

func TestNew_MinimumWorkersIsOne(t *testing.T) {
	p := concurrency.New(0, &bytes.Buffer{})
	results := p.Run(context.Background(), []concurrency.Task{
		{Host: "h1", Patch: "p1", Fn: func(_ context.Context) error { return nil }},
	})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
}

func TestRun_ReturnsOneResultPerTask(t *testing.T) {
	p := concurrency.New(3, &bytes.Buffer{})
	tasks := []concurrency.Task{
		{Host: "h1", Patch: "001", Fn: func(_ context.Context) error { return nil }},
		{Host: "h2", Patch: "002", Fn: func(_ context.Context) error { return nil }},
		{Host: "h3", Patch: "003", Fn: func(_ context.Context) error { return nil }},
	}
	results := p.Run(context.Background(), tasks)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error for %s/%s: %v", r.Host, r.Patch, r.Err)
		}
	}
}

func TestRun_PropagatesTaskError(t *testing.T) {
	p := concurrency.New(2, &bytes.Buffer{})
	want := errors.New("boom")
	tasks := []concurrency.Task{
		{Host: "h1", Patch: "001", Fn: func(_ context.Context) error { return want }},
	}
	results := p.Run(context.Background(), tasks)
	if !errors.Is(results[0].Err, want) {
		t.Fatalf("expected %v, got %v", want, results[0].Err)
	}
}

func TestRun_NeverExceedsWorkerCount(t *testing.T) {
	const maxWorkers = 3
	var active int64
	p := concurrency.New(maxWorkers, &bytes.Buffer{})

	tasks := make([]concurrency.Task, 12)
	for i := range tasks {
		tasks[i] = concurrency.Task{
			Host: "h", Patch: "p",
			Fn: func(_ context.Context) error {
				cur := atomic.AddInt64(&active, 1)
				if cur > maxWorkers {
					t.Errorf("concurrency exceeded: %d > %d", cur, maxWorkers)
				}
				time.Sleep(5 * time.Millisecond)
				atomic.AddInt64(&active, -1)
				return nil
			},
		}
	}
	p.Run(context.Background(), tasks)
}

func TestRun_CancelledContextMarksError(t *testing.T) {
	p := concurrency.New(1, &bytes.Buffer{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	tasks := []concurrency.Task{
		{Host: "h1", Patch: "001", Fn: func(_ context.Context) error { return nil }},
	}
	results := p.Run(ctx, tasks)
	if results[0].Err == nil {
		t.Fatal("expected cancellation error, got nil")
	}
}
