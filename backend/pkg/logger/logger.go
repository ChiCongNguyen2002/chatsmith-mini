package logger

import (
	"context"
	"log/slog"
	"os"

	"chatsmith-mini/constants"
)

func New() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func WithContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, constants.ContextKeyLogger, log)
}

func FromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(constants.ContextKeyLogger).(*slog.Logger); ok {
		return log
	}

	return slog.Default()
}
