// Package ratelimit provides a token-bucket style rate limiter for controlling
// how frequently patch operations are dispatched to remote hosts.
package ratelimit

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Limiter controls the rate at which operations are allowed to proceed.
type Limiter struct {
	out      io.Writer
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
}

// New creates a Limiter that allows up to burst operations, refilling at
// ratePerSec tokens per second. If out is nil, os.Stdout is used.
func New(ratePerSec float64, burst int, out io.Writer) (*Limiter, error) {
	if ratePerSec <= 0 {
		return nil, fmt.Errorf("ratelimit: ratePerSec must be positive, got %v", ratePerSec)
	}
	if burst < 1 {
		burst = 1
	}
	if out == nil {
		out = os.Stdout
	}
	return &Limiter{
		out:      out,
		tokens:   float64(burst),
		max:      float64(burst),
		rate:     ratePerSec,
		lastTick: time.Now(),
	}, nil
}

// Wait blocks until a token is available, then consumes one.
func (l *Limiter) Wait() {
	for {
		l.mu.Lock()
		l.refill()
		if l.tokens >= 1 {
			l.tokens--
			fmt.Fprintf(l.out, "[ratelimit] token acquired (%.1f remaining)\n", l.tokens)
			l.mu.Unlock()
			return
		}
		// Calculate how long until next token is available.
		wait := time.Duration((1-l.tokens)/l.rate*1e9) * time.Nanosecond
		l.mu.Unlock()
		time.Sleep(wait)
	}
}

// TryAcquire attempts to consume a token without blocking.
// Returns true if a token was available and consumed.
func (l *Limiter) TryAcquire() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refill()
	if l.tokens >= 1 {
		l.tokens--
		fmt.Fprintf(l.out, "[ratelimit] token acquired (%.1f remaining)\n", l.tokens)
		return true
	}
	fmt.Fprintf(l.out, "[ratelimit] no tokens available\n")
	return false
}

// refill adds tokens based on elapsed time. Must be called with l.mu held.
func (l *Limiter) refill() {
	now := time.Now()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}
	l.lastTick = now
}
