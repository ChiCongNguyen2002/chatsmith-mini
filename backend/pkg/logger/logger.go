package logger

import (
	"context"
	"log/slog"
	"os"

	"chatsmith-mini/constants"
)

const (
	severityDebug   = "DEBUG"
	severityInfo    = "INFO"
	severityWarning = "WARNING"
	severityError   = "ERROR"
)

func New() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: replaceAttr}))
}

func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key != slog.LevelKey {
		return a
	}
	level, ok := a.Value.Any().(slog.Level)
	if !ok {
		return a
	}

	return slog.String(constants.LogKeySeverity, severityOf(level))
}

func severityOf(level slog.Level) string {
	switch {
	case level < slog.LevelInfo:
		return severityDebug
	case level < slog.LevelWarn:
		return severityInfo
	case level < slog.LevelError:
		return severityWarning
	default:
		return severityError
	}
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
