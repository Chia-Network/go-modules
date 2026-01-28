package slogs

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

// Logr is a custom text logger from the stdlib slog package
var Logr Logger

// Logger is a wrapper around a slog log handler
type Logger struct {
	h                slog.Handler
	addSourceContext bool
}

type loggerOptions struct {
	// writer is any interface that implements an I/O Writer
	writer io.Writer
	// handlerOptions are the options passed to the slog handler
	handlerOptions slog.HandlerOptions
}

// ClientOptionFunc can be used to customize a new slogs client
type ClientOptionFunc func(*loggerOptions)

// WithSourceContext sets the AddSource option for the slogs logger, which adds the location in the code that called the logger, to each logline
func WithSourceContext(set bool) ClientOptionFunc {
	return func(o *loggerOptions) {
		o.handlerOptions.AddSource = set
	}
}

// WithWriter sets a given io.Writer to be the receiver for slogs logs
func WithWriter(w io.Writer) ClientOptionFunc {
	return func(o *loggerOptions) {
		o.writer = w
	}
}

// Init custom init function that accepts the log level for the application and initializes a stdout slog logger
func Init(level string, options ...ClientOptionFunc) {
	logOpts := loggerOptions{
		writer: os.Stdout,
	}
	logOpts.handlerOptions.Level = parseLogLevel(level)

	// Apply any given options
	for _, fn := range options {
		if fn == nil {
			continue
		}
		fn(&logOpts)
	}

	Logr.h = slog.NewTextHandler(logOpts.writer, &logOpts.handlerOptions)
	Logr.addSourceContext = logOpts.handlerOptions.AddSource
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

func (l Logger) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !l.h.Enabled(ctx, level) {
		return
	}

	var pc uintptr
	if l.addSourceContext {
		// skips runtime.Callers, this log function, the log level function wrapper, to the caller of that function.
		// This is for adding the correct source caller to the log's record
		var pcs [1]uintptr
		runtime.Callers(3, pcs[:])
		pc = pcs[0]
	}

	rec := slog.NewRecord(time.Now(), level, msg, pc)
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
