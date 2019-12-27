package newrelic

import (
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

// NewRelic is a connection to New Relic APIs
type NewRelic struct {
	Alerts     alerts.Alerts
	APM        apm.APM
	Synthetics synthetics.Synthetics
}

// New returns a NewRelic API connection struct
func New(config config.Config) NewRelic {
	return NewRelic{
		Alerts:     alerts.New(config),
		APM:        apm.New(config),
		Synthetics: synthetics.New(config),
	}
}
