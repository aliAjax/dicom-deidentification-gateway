package application

import (
	"context"
	"log/slog"
	"time"
)

type Task interface{ Run(context.Context) error }
type Runner struct {
	tasks    []Task
	log      *slog.Logger
	interval time.Duration
}

func New(log *slog.Logger, interval time.Duration, tasks ...Task) *Runner {
	return &Runner{tasks: tasks, log: log, interval: interval}
}
func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, task := range r.tasks {
				if err := task.Run(ctx); err != nil {
					r.log.Error("worker task failed", "error", err)
				}
			}
		}
	}
}
