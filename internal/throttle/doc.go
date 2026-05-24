// Package throttle provides a concurrency-limiting semaphore for controlling
// how many remote hosts are targeted in parallel during a deployment run.
//
// # Usage
//
//	th := throttle.New(5, os.Stdout)
//
//	for _, host := range hosts {
//		go func(h string) {
//			if err := th.Acquire(ctx, h); err != nil {
//				return
//			}
//			defer th.Release(h)
//			// deploy to h
//		}(host)
//	}
//
// Acquire blocks until a slot is free or the context is cancelled.
// Release must always be called after a successful Acquire to avoid deadlocks.
package throttle
