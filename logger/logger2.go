// Copyright 2022. Motty Cohen.
//
// Logger implementation based on system log package
//
package logger

import (
	"log"
	"os"
)

var (
	DebugLogger   *log.Logger
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init2() {
	DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// region Write to log =------------------------------------------------------------------------------------------------

/**
 * Debug
 */
func Debug2(format string, params ...any) {
	DebugLogger.Printf(format, params)
}

/**
 * Info
 */
func Info2(format string, params ...any) {
	InfoLogger.Printf(format, params)
}

/**
 * Warn
 */
func Warn2(format string, params ...any) {
	WarningLogger.Printf(format, params)
}

/**
 * Error
 */
func Error2(format string, params ...any) {
	ErrorLogger.Printf(format, params)
}

// endregion
