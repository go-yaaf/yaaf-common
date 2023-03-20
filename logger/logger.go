// Package logger Logger implementation based on zap log package (https://github.com/uber-go/zap)
package logger

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

var loggerConfig zap.Config
var loggerOnce sync.Once
var loggerSingleton *zap.Logger = nil

func getLogger() (result *zap.Logger) {
	loggerOnce.Do(func() {
		if loggerSingleton == nil {
			loggerSingleton, _ = zap.NewProduction()
		}
	})
	return loggerSingleton
}

// init set default logging configuration
func init() {

	encoderConfig := zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customEncodeTime,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	loggerConfig = zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	loggerConfig.DisableCaller = true
	loggerConfig.DisableStacktrace = true
}

// region Logger configuration -----------------------------------------------------------------------------------------

// SetLevel log level DEBUG | INFO | WARN | ERROR
func SetLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "warning":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		loggerConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}
}

// EnableJsonFormat configure log output as json (true) or text line (false)
func EnableJsonFormat(value bool) {
	if value {
		loggerConfig.Encoding = "json"
	} else {
		loggerConfig.Encoding = "console"
	}
}

// EnableStacktrace configure log output to include or exclude stack trace
func EnableStacktrace(value bool) {
	loggerConfig.DisableStacktrace = !value
}

// SetTimeLayout define the log entry time layout
func SetTimeLayout(layout string) {
	if len(layout) == 0 {
		layout = "2006-01-02 15:04:05"
	}
	loggerConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(layout))
	}
}

func customEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// Init initialize logger
func Init() {
	var err error
	if loggerSingleton, err = loggerConfig.Build(); err != nil {
		loggerSingleton, _ = zap.NewProduction()
	}
}

// endregion

// region Write to log -------------------------------------------------------------------------------------------------

// Debug log level
func Debug(format string, params ...any) {
	l := getLogger()
	defer l.Sync()
	l.Debug(fmt.Sprintf(format, params...))
}

// Info log level
func Info(format string, params ...any) {
	l := getLogger()
	defer l.Sync()
	l.Info(fmt.Sprintf(format, params...))
}

// Warn log level
func Warn(format string, params ...any) {
	l := getLogger()
	defer l.Sync()
	l.Warn(fmt.Sprintf(format, params...))
}

// Error log level
func Error(format string, params ...any) {
	l := getLogger()
	defer l.Sync()
	l.Error(fmt.Sprintf(format, params...))
}

// Fatal log level
func Fatal(format string, params ...any) {
	l := getLogger()
	defer l.Sync()
	l.Fatal(fmt.Sprintf(format, params...))
}

// endregion
