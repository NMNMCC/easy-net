package log

import (
	"log/slog"
	"os"
)

func New(component string) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	})).With("component", component)
}