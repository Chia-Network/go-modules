package slogs

import (
	"context"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

// Logr is a custom text logger from the stdlib slog package
var Logr Logger

// Logger is a wrapper around slog that can adjust call depth for AddSource
type Logger struct {
	base *slog.Logger
	h    slog.Handler
}

// InitOptions contain options for the slog logger
type InitOptions struct {
	AddSource bool
}

// Init custom init function that accepts the log level for the application and initializes a stdout slog logger
func Init(level string, opts *InitOptions) {
	handlerOpts := slog.HandlerOptions{
		Level: parseLogLevel(level),
	}
	if opts != nil && opts.AddSource {
		handlerOpts.AddSource = true
	}

	h := slog.NewTextHandler(os.Stdout, &handlerOpts)
	l := slog.New(h)

	Logr = Logger{
		base: l,
		h:    h,
	}
}

func (l Logger) Debug(msg string, args ...any) {
	l.log(context.Background(), slog.LevelDebug, msg, args...)
}
func (l Logger) Info(msg string, args ...any) {
	l.log(context.Background(), slog.LevelInfo, msg, args...)
}
func (l Logger) Warn(msg string, args ...any) {
	l.log(context.Background(), slog.LevelWarn, msg, args...)
}
func (l Logger) Error(msg string, args ...any) {
	l.log(context.Background(), slog.LevelError, msg, args...)
}
func (l Logger) Fatal(msg string, args ...any) {
	l.log(context.Background(), slog.LevelError, msg, args...)
	os.Exit(1)
}

func (l Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelDebug, msg, args...)
}
func (l Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelInfo, msg, args...)
}
func (l Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelWarn, msg, args...)
}
func (l Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelError, msg, args...)
}
func (l Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelError, msg, args...)
	os.Exit(1)
}

// With returns a new logger with additional attributes.
func (l Logger) With(args ...any) Logger {
	// Preserve handler chain so Handle() keeps working.
	h2 := l.h.WithAttrs(slogArgsToAttrs(args))
	return Logger{
		base: slog.New(h2),
		h:    h2,
	}
}

// Handler exposes the underlying handler (optional convenience).
func (l Logger) Handler() slog.Handler { return l.h }

func (l Logger) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	// 0 runtime.Callers
	// 1 Logger.log
	// 2 Logger.Info/Debug/...
	// 3 user callsite  <-- we want this one
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])

	rec := slog.NewRecord(time.Now(), level, msg, pcs[0])
	rec.Add(args...)

	_ = l.h.Handle(ctx, rec)
}

func slogArgsToAttrs(args ...any) []slog.Attr {
	// Convert the keyvals into Attrs the same way slog does.
	// slog.Any() is fine for values; keys must be strings.
	attrs := make([]slog.Attr, 0, len(args)/2)
	for i := 0; i+1 < len(args); i += 2 {
		k, ok := args[i].(string)
		if !ok {
			// mirror slog's behavior-ish: skip bad key
			continue
		}
		attrs = append(attrs, slog.Any(k, args[i+1]))
	}
	return attrs
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
		log.Printf("unknown log level specified %q, defaulting to info level", level)
		return slog.LevelInfo
	}
}
