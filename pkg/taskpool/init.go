package taskpool

import (
	"context"
	"log/slog"

	"github.com/devchat-ai/gopool"
)

var (
	Workers gopool.GoPool
)

func InitTaskPool(ctx context.Context) {
	Workers = gopool.NewGoPool(512, gopool.WithMinWorkers(128), gopool.WithTaskQueueSize(1<<20))

	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			Workers.Wait()
			Workers.Release()
			slog.Info("task pool have been released")
			return
		}
	}(ctx)
}
