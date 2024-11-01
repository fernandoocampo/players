package unittests

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	handlerOptions := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	loggerHandler := slog.NewJSONHandler(os.Stdout, handlerOptions)

	return slog.New(loggerHandler)
}
