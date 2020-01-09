package logging

var (
	logger = StructuredLogger{}
)

// Debug logs a message at level Debug on the standard logger.
func Debug(msg string, args ...interface{}) {
	logger.Debug(msg, args...)
}

// Info logs a message at level Info on the standard logger.
func Info(msg string, args ...interface{}) {
	logger.Info(msg, args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(msg string, args ...interface{}) {
	logger.Warn(msg, args...)
}

// Error logs a message at level Error on the standard logger.
func Error(msg string, args ...interface{}) {
	logger.Error(msg, args...)
}
