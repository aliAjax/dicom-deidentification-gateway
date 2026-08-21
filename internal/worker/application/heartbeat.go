package application

import (
	"context"
	"time"
)

type Heartbeat struct {
	Interval time.Duration
	Fn       func(context.Context) error
}

func (h Heartbeat) Run(ctx context.Context) {
	d := h.Interval
	if d <= 0 {
		d = 10 * time.Second
	}
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if h.Fn != nil {
				_ = h.Fn(ctx)
			}
		}
	}
}
