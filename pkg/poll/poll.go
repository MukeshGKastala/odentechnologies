package poll

import (
	"context"
	"time"
)

// NeverStop may be passed to Until to make it never stop.
var NeverStop <-chan struct{} = make(chan struct{})

// Forever calls f every period for ever.
func Forever(f func(), period time.Duration) {
	Until(f, period, NeverStop)
}

// UntilWithContext loops until context is done, running f every period.
func UntilWithContext(ctx context.Context, f func(context.Context), period time.Duration) {
	Until(func() { f(ctx) }, period, ctx.Done())
}

// Until loops until stop channel is closed, running f every period.
//
// The timer for period starts after the f completes.
//
// TODO: add ability to jitter the period - allows clients to avoid converging on periodic behavior.
// TODO: add ability to start the timer for period at the same time as the function starts
func Until(f func(), period time.Duration, stopCh <-chan struct{}) {
	for {
		select {
		case <-time.After(period):
			// TODO: handle panic
			f()
		case <-stopCh:
			return
		}
	}
}
