// Package logger Logger implementation based on slog
package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

var loggerOnce sync.Once
var loggerSingleton *slog.Logger

var handlerOptions = &slog.HandlerOptions{
	AddSource: false, // Disabled by default to match original zap config
	Level:     slog.LevelInfo,
}
var output io.Writer = os.Stderr
var isJson = true
var timeLayout = "2006-01-02 15:04:05" // Default time layout from original

func getLogger() *slog.Logger {
	loggerOnce.Do(func() {
		if loggerSingleton == nil {
			Init()
		}
	})
	return loggerSingleton
}

// init set default logging configuration
func init() {
	Init()
}

// region Logger configuration -----------------------------------------------------------------------------------------

// SetLevel log level DEBUG | INFO | WARN | ERROR
func SetLevel(level string) {
	level = strings.ToLower(level)
	lvl := slog.LevelInfo
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	}
	handlerOptions.Level = lvl
	Init()
}

// EnableJsonFormat configure log output as json (true) or text line (false)
func EnableJsonFormat(value bool) {
	isJson = value
	Init()
}

// EnableStacktrace is not supported by slog and is a no-op.
func EnableStacktrace(value bool) {
	// slog does not have built-in stacktrace support like zap.
	// For error stack traces, consider wrapping errors and logging them.
}

// EnableCaller enables/disables logging the caller source file and line number.
func EnableCaller(value bool) {
	handlerOptions.AddSource = value
	Init()
}

// SetTimeLayout define the log entry time layout.
func SetTimeLayout(layout string) {
	if len(layout) == 0 {
		layout = "2006-01-02 15:04:05"
	}
	timeLayout = layout
	Init()
}

// Init initialize logger
func Init() {
	replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
		// Customize time format
		if a.Key == slog.TimeKey && timeLayout != "" {
			if t, ok := a.Value.Any().(time.Time); ok {
				a.Value = slog.StringValue(t.Format(timeLayout))
			}
		}

		// Rename level to severity for GCP compatibility
		if a.Key == slog.LevelKey {
			a.Key = "severity"
			level := a.Value.Any().(slog.Level)
			switch level {
			case slog.LevelDebug:
				a.Value = slog.StringValue("DEBUG")
			case slog.LevelInfo:
				a.Value = slog.StringValue("INFO")
			case slog.LevelWarn:
				a.Value = slog.StringValue("WARNING")
			case slog.LevelError:
				a.Value = slog.StringValue("ERROR")
			default:
				a.Value = slog.StringValue(level.String())
			}
		}

		// Rename msg to message for GCP compatibility
		if a.Key == slog.MessageKey {
			a.Key = "message"
		}
		return a
	}

	handlerOptions.ReplaceAttr = replaceAttr

	var handler slog.Handler
	if isJson {
		handler = slog.NewJSONHandler(output, handlerOptions)
	} else {
		handler = slog.NewTextHandler(output, handlerOptions)
	}
	loggerSingleton = slog.New(handler)
	slog.SetDefault(loggerSingleton)
}

// endregion

// region Write to log -------------------------------------------------------------------------------------------------

// Debug log level
func Debug(format string, params ...any) {
	getLogger().Debug(fmt.Sprintf(format, params...))
}

// Info log level
func Info(format string, params ...any) {
	getLogger().Info(fmt.Sprintf(format, params...))
}

// Warn log level
func Warn(format string, params ...any) {
	getLogger().Warn(fmt.Sprintf(format, params...))
}

// Error log level
func Error(format string, params ...any) {
	getLogger().Error(fmt.Sprintf(format, params...))
}

// Fatal log level, followed by os.Exit(1)
func Fatal(format string, params ...any) {
	getLogger().Error(fmt.Sprintf(format, params...))
	os.Exit(1)
}

// endregion
