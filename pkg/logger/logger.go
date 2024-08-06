package logger

import (
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envTest  = "test"
)

type Level string

func (l Level) Level() slog.Level {
	switch l {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

var _ slog.Leveler = (*Level)(nil)

type Config struct {
	Level Level
	Env   string
}

func SetupLogger(cfg Config) *slog.Logger {
	var logger *slog.Logger

	switch cfg.Env {
	case envTest:
		logger = setupTestLogger(cfg.Level)
	case envLocal:
		logger = setupLocalLogger(cfg.Level)
	default:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.Level}))
	}
	return logger
}

func setupLocalLogger(level Level) *slog.Logger {
	opts := PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}

func setupTestLogger(level Level) *slog.Logger {
	opts := NullHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
	}
	handler := opts.NewNullHandler(os.Stdout)
	return slog.New(handler)
}
