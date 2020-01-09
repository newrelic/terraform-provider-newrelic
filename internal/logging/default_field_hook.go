package logging

import (
	"github.com/newrelic/newrelic-client-go/internal/version"
	log "github.com/sirupsen/logrus"
)

// DefaultFieldHook provides default fields for structured log messages.
type DefaultFieldHook struct{}

// Levels determines which levels will get default fields.
func (h *DefaultFieldHook) Levels() []log.Level {
	return log.AllLevels
}

// Fire is the hook for each log message.
func (h *DefaultFieldHook) Fire(e *log.Entry) error {
	e.Data["newrelic-client-go"] = version.Version
	return nil
}
