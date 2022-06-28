// Copyright 2022. Motty Cohen.
//
// Logger implementation based on zap log package (https://github.com/uber-go/zap)
//
package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

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

// region Write to log =------------------------------------------------------------------------------------------------

/**
 * Debug
 */
func Debug(format string, params ...any) {
	l := getLogger()
	defer l.Sync()
	l.Debug(fmt.Sprintf(format, params...))
}

/**
 * Info
 */
func Info(format string, params ...any) {
	l := getLogger()
	defer l.Sync()
	l.Info(fmt.Sprintf(format, params...))
}

/**
 * Warn
 */
func Warn(format string, params ...any) {
	l := getLogger()
	defer l.Sync()
	l.Warn(fmt.Sprintf(format, params...))
}

/**
 * Error
 */
func Error(format string, params ...any) {
	l := getLogger()
	defer l.Sync()
	l.Error(fmt.Sprintf(format, params...))
}

/**
 * Error
 */
func Fatal(format string, params ...any) {
	l := getLogger()
	defer l.Sync()
	l.Fatal(fmt.Sprintf(format, params...))
}

// endregion
