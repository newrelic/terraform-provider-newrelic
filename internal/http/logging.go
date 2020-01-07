package http

import (
	log "github.com/sirupsen/logrus"
)

type structuredLogger struct{}

func (l structuredLogger) Error(msg string, fields ...interface{}) {
	log.WithFields(createFieldMap(fields)).Error(msg)
}

func (l structuredLogger) Info(msg string, fields ...interface{}) {
	log.WithFields(createFieldMap(fields)).Info(msg)
}

func (l structuredLogger) Debug(msg string, fields ...interface{}) {
	log.WithFields(createFieldMap(fields)).Debug(msg)
}

func (l structuredLogger) Warn(msg string, fields ...interface{}) {
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
