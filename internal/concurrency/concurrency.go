// Package concurrency provides a worker pool for executing patch operations
// across multiple hosts in parallel, with a configurable degree of parallelism.
package concurrency

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

// Task represents a unit of work to be executed by the pool.
type Task struct {
	Host  string
	Patch string
	Fn    func(ctx context.Context) error
}

// Result holds the outcome of a single task execution.
type Result struct {
	Host  string
	Patch string
	Err   error
}

// Pool runs tasks concurrently up to a fixed worker count.
type Pool struct {
	workers int
	out     io.Writer
}

// New creates a Pool with the given worker count.
// If workers < 1 it is clamped to 1. If out is nil, os.Stdout is used.
func New(workers int, out io.Writer) *Pool {
	if workers < 1 {
		workers = 1
	}
	if out == nil {
		out = os.Stdout
	}
	return &Pool{workers: workers, out: out}
}

// Run executes all tasks using the pool's worker goroutines and returns
// one Result per task. The returned slice preserves submission order.
// Cancelling ctx causes in-flight tasks to abort; pending tasks are dropped.
func (p *Pool) Run(ctx context.Context, tasks []Task) []Result {
	results := make([]Result, len(tasks))
	type indexed struct {
		idx  int
		task Task
	}

	queue := make(chan indexed, len(tasks))
	for i, t := range tasks {
		queue <- indexed{i, t}
	}
	close(queue)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for w := 0; w < p.workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range queue {
				select {
				case <-ctx.Done():
					mu.Lock()
					results[item.idx] = Result{
						Host:  item.task.Host,
						Patch: item.task.Patch,
						Err:   fmt.Errorf("cancelled: %w", ctx.Err()),
					}
					mu.Unlock()
				default:
					err := item.task.Fn(ctx)
					mu.Lock()
					results[item.idx] = Result{
						Host:  item.task.Host,
						Patch: item.task.Patch,
						Err:   err,
					}
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	return results
}
