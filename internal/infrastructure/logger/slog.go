package logger

import (
	"OnlineLeadership/internal/interfaces/http/middleware"
	"context"
	"log/slog"
	"os"
)

type SlogLogger struct {
	log *slog.Logger
}

func New(env string) *SlogLogger {
	var handler slog.Handler

	if env == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}
	return &SlogLogger{
		log: slog.New(handler),
	}
}

func (l *SlogLogger) withContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return l.log
	}

	if reqID, ok := ctx.Value(middleware.RequestIDKey).(string); ok {
		return l.log.With(
			slog.String("request_id", reqID),
		)
	}

	return l.log
}

func (l *SlogLogger) Info(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx).Info(msg, args...)
}

func (l *SlogLogger) Error(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx).Error(msg, args...)
}

func (l *SlogLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx).Warn(msg, args...)
}

func (l *SlogLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.withContext(ctx).Debug(msg, args...)
}
