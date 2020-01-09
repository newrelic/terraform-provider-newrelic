package logging

import (
	log "github.com/sirupsen/logrus"
)

// StructuredLogger is a logger based on logrus.
type StructuredLogger struct{}

// Error logs an error message.
func (l StructuredLogger) Error(msg string, fields ...interface{}) {
	log.WithFields(createFieldMap(fields)).Error(msg)
}

// Info logs an info message.
func (l StructuredLogger) Info(msg string, fields ...interface{}) {
	log.WithFields(createFieldMap(fields)).Info(msg)
}

// Debug logs a debug message.
func (l StructuredLogger) Debug(msg string, fields ...interface{}) {
	log.WithFields(createFieldMap(fields)).Debug(msg)
}

// Warn logs an warning message.
func (l StructuredLogger) Warn(msg string, fields ...interface{}) {
	log.WithFields(createFieldMap(fields)).Warn(msg)
}

func createFieldMap(fields ...interface{}) map[string]interface{} {
	m := map[string]interface{}{}

	fields = fields[0].([]interface{})

	for i := 0; i < len(fields); i += 2 {
		m[fields[i].(string)] = fields[i+1]
	}

	return m
}
