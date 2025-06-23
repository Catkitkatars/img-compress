package sl

import (
	"log/slog"
	"os"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type Logger struct {
	env string
	log *slog.Logger
}

func New(env string) *Logger {
	var handler slog.Handler

	switch env {
	case EnvLocal:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case EnvDev:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case EnvProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	return &Logger{
		env: env,
		log: slog.New(handler),
	}
}

func (l *Logger) Info(msg string, attrs ...slog.Attr) {
	l.log.Info(msg, convertAttrs(attrs)...)
}

func (l *Logger) Error(msg string, err error, attrs ...slog.Attr) {
	attrs = append(attrs, slog.String("error", err.Error()))
	l.log.Error(msg, convertAttrs(attrs)...)
}

func (l *Logger) Debug(msg string, attrs ...slog.Attr) {
	l.log.Debug(msg, convertAttrs(attrs)...)
}

func convertAttrs(attrs []slog.Attr) []any {
	args := make([]any, len(attrs))
	for i, a := range attrs {
		args[i] = a
	}
	return args
}
