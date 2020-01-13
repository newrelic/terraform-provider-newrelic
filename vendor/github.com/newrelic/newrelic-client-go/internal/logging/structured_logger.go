package logging

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

var (
	defaultFields   map[string]string
	defaultLogLevel = "info"
)

// StructuredLogger is a logger based on logrus.
type StructuredLogger struct {
	logger *log.Logger
}

// NewStructuredLogger creates a new structured logger.
func NewStructuredLogger() StructuredLogger {
	return StructuredLogger{
		logger: log.New(),
	}
}

// SetLogLevel allows the log level to be set.
func (l StructuredLogger) SetLogLevel(logLevel string) StructuredLogger {
	if logLevel == "" {
		logLevel = defaultLogLevel
	}

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Warn(fmt.Sprintf("could not parse log level '%s', logging will proceed at %s level", logLevel, defaultLogLevel))
		level, _ = log.ParseLevel(defaultLogLevel)
	}

	l.logger.SetLevel(level)

	return l
}

// LogJSON determines whether or not to format the logs as JSON.
func (l StructuredLogger) LogJSON(value bool) StructuredLogger {
	if value {
		l.logger.SetFormatter(&log.JSONFormatter{})
	}

	return l
}

// SetDefaultFields sets fields to be logged on every use of the logger.
func (l StructuredLogger) SetDefaultFields(defaultFields map[string]string) StructuredLogger {
	l.logger.AddHook(&defaultFieldHook{})

	return l
}

// Error logs an error message.
func (l StructuredLogger) Error(msg string, fields ...interface{}) {
	l.logger.WithFields(createFieldMap(fields)).Error(msg)
}

// Info logs an info message.
func (l StructuredLogger) Info(msg string, fields ...interface{}) {
	l.logger.WithFields(createFieldMap(fields)).Info(msg)
}

// Debug logs a debug message.
func (l StructuredLogger) Debug(msg string, fields ...interface{}) {
	l.logger.WithFields(createFieldMap(fields)).Debug(msg)
}

// Warn logs an warning message.
func (l StructuredLogger) Warn(msg string, fields ...interface{}) {
	l.logger.WithFields(createFieldMap(fields)).Warn(msg)
}

func createFieldMap(fields ...interface{}) map[string]interface{} {
	m := map[string]interface{}{}

	fields = fields[0].([]interface{})

	for i := 0; i < len(fields); i += 2 {
		m[fields[i].(string)] = fields[i+1]
	}

	return m
}

type defaultFieldHook struct{}

func (h *defaultFieldHook) Levels() []log.Level {
	return log.AllLevels
}

func (h *defaultFieldHook) Fire(e *log.Entry) error {
	for k, v := range defaultFields {
		e.Data[k] = v
	}
	return nil
}
