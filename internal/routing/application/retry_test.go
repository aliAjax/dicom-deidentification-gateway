package application

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryWaitStopsOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	started := time.Now()
	err := Wait(ctx, time.Second)
	if !errors.Is(err, context.Canceled) || time.Since(started) > 100*time.Millisecond {
		t.Fatalf("err=%v elapsed=%s", err, time.Since(started))
	}
}
