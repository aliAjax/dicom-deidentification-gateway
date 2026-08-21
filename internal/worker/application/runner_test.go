package application

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

type runnerTask struct{ calls chan struct{} }

func (t runnerTask) Run(context.Context) error { t.calls <- struct{}{}; return nil }
func TestRunnerStopsTaskDispatchAfterCancel(t *testing.T) {
	calls := make(chan struct{}, 4)
	ctx, cancel := context.WithCancel(context.Background())
	r := New(slog.Default(), time.Millisecond, runnerTask{calls: calls})
	go r.Run(ctx)
	select {
	case <-calls:
		cancel()
	case <-time.After(100 * time.Millisecond):
		cancel()
		t.Fatal("task did not run")
	}
	select {
	case <-calls:
		t.Fatal("task dispatched after cancellation")
	case <-time.After(10 * time.Millisecond):
	}
}

func TestHeartbeatStopsCallbackAfterCancel(t *testing.T) {
	calls := make(chan struct{}, 4)
	ctx, cancel := context.WithCancel(context.Background())
	h := Heartbeat{Interval: time.Millisecond, Fn: func(context.Context) error { calls <- struct{}{}; return nil }}
	go h.Run(ctx)
	select {
	case <-calls:
		cancel()
	case <-time.After(100 * time.Millisecond):
		cancel()
		t.Fatal("heartbeat did not run")
	}
	time.Sleep(5 * time.Millisecond)
	select {
	case <-calls:
		t.Fatal("callback ran after cancel")
	default:
	}
}
