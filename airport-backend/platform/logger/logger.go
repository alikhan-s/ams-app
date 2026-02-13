package logger

import (
	"log/slog"
	"os"
)

// New creates a new structured logger that writes JSON to stdout.
func New() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
