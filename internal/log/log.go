package log

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	// Default to error level only
	logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
}

// EnableDebug switches to debug level logging
func EnableDebug() {
	logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

// Info logs an info message
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}
