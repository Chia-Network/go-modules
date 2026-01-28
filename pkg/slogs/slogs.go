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
	h slog.Handler

	handlerOptions slog.HandlerOptions
}

// ClientOptionFunc can be used to customize a new slogs client
type ClientOptionFunc func(*Logger)

// WithSourceContext sets the AddSource option for the slogs logger, which adds the location in the code that called the logger, to each logline
func WithSourceContext(set bool) ClientOptionFunc {
	return func(l *Logger) {
		l.setAddSource(set)
	}
}

// setAddSource sets the AddSource option on a logger
func (l *Logger) setAddSource(set bool) {
	l.handlerOptions.AddSource = set
}

// Init custom init function that accepts the log level for the application and initializes a stdout slog logger
func Init(level string, options ...ClientOptionFunc) {
	Logr.handlerOptions.Level = parseLogLevel(level)

	// Apply any given options
	for _, fn := range options {
		if fn == nil {
			continue
		}
		fn(&Logr)
	}

	h := slog.NewTextHandler(os.Stdout, &Logr.handlerOptions)

	Logr = Logger{
		h: h,
	}
}

// Debug uses the initialized logger at Debug level
func (l Logger) Debug(msg string, args ...any) {
	l.log(context.Background(), slog.LevelDebug, msg, args...)
}

// Info uses the initialized logger at Info level
func (l Logger) Info(msg string, args ...any) {
	l.log(context.Background(), slog.LevelInfo, msg, args...)
}

// Warn uses the initialized logger at Warn level
func (l Logger) Warn(msg string, args ...any) {
	l.log(context.Background(), slog.LevelWarn, msg, args...)
}

// Error uses the initialized logger at Error level
func (l Logger) Error(msg string, args ...any) {
	l.log(context.Background(), slog.LevelError, msg, args...)
}

// Fatal uses the initialized logger at Error level, and exits 1
func (l Logger) Fatal(msg string, args ...any) {
	l.log(context.Background(), slog.LevelError, msg, args...)
	os.Exit(1)
}

// DebugContext uses the initialized logger at Debug level with a given context
func (l Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelDebug, msg, args...)
}

// InfoContext uses the initialized logger at Info level with a given context
func (l Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelInfo, msg, args...)
}

// WarnContext uses the initialized logger at Warn level with a given context
func (l Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelWarn, msg, args...)
}

// ErrorContext uses the initialized logger at Error level with a given context
func (l Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelError, msg, args...)
}

// FatalContext uses the initialized logger at Error level with a given context, and exits 1
func (l Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelError, msg, args...)
	os.Exit(1)
}

// Handler exposes the underlying handler (optional convenience).
func (l Logger) Handler() slog.Handler { return l.h }

func (l Logger) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	// This retrieves the actual caller of the logger out of the call-stack.
	// When using slogs in any kind of wrapper context with the AddSource option, the logger's caller is hidden a few layers in the call stack
	// 0 runtime.Callers
	// 1 Logger.log
	// 2 Logger.[logger function]/...
	// 3 user callsite  <-- want this one
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])

	rec := slog.NewRecord(time.Now(), level, msg, pcs[0])
	rec.Add(args...)

	err := l.h.Handle(ctx, rec)
	if err != nil {
		log.Printf("slogs internal error: failed to handle record: %v", err)
	}
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
