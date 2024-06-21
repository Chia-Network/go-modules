package slogs

import (
	"log"
	"log/slog"
	"os"
	"strings"
)

// Logr is a custom text logger from the stdlib slog package
var Logr Logger

// Logger is a wrapper around the slog logger struct so we can have a type that is owned by this scope to create additional wrapper functions around
type Logger struct {
	*slog.Logger
}

// Init custom init function that accepts the log level for the application and initializes a stdout slog logger
func Init(level string) {
	l := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: parseLogLevel(level),
			},
		),
	)
	Logr = Logger{l}
}

// Function to convert log level string to slog.Level
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		log.Printf("unknown log level specified \"%s\", defaulting to info level", level)
		return slog.LevelInfo
	}
}

// Fatal is a wrapper around the standard slog Error function that exits 1 after it is called.
// Similar to the stdlib log.Fatal function
func (l Logger) Fatal(msg string, args ...any) {
	l.Error(msg, args...)
	os.Exit(1)
}
