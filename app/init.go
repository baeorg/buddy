package app

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/baeorg/buddy/pkg/storage"
)

func InitBootSystem(ctx context.Context) {
	InitLog()
	storage.InitDB(ctx)
}

func Release() {
}

func InitLog() {
	opts := slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				if idx := strings.Index(source.File, "buddy/"); idx != -1 {
					source.File = source.File[idx+len("buddy/"):]
				}
				return slog.Any(slog.SourceKey, source)
			}
			return a
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &opts))
	slog.SetDefault(logger)
}
