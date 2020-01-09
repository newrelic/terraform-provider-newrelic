package logging

import (
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

var (
	internalLogger config.Logger = StructuredLogger{}
)

// SetLogger allows changing the underlying logger implementation.
func SetLogger(logger config.Logger) {
	internalLogger = logger
}

// Debug logs a message at level Debug on the standard logger.
func Debug(msg string, args ...interface{}) {
	internalLogger.Debug(msg, args...)
}

// Info logs a message at level Info on the standard logger.
func Info(msg string, args ...interface{}) {
	internalLogger.Info(msg, args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(msg string, args ...interface{}) {
	internalLogger.Warn(msg, args...)
}

// Error logs a message at level Error on the standard logger.
func Error(msg string, args ...interface{}) {
	internalLogger.Error(msg, args...)
}
