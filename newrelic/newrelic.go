package newrelic

import (
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

type NewRelic struct {
	Alerts     alerts.Alerts
	APM        apm.APM
	Synthetics synthetics.Synthetics
}

func New(config config.Config) NewRelic {
	return NewRelic{
		Alerts:     alerts.New(config),
		APM:        apm.New(config),
		Synthetics: synthetics.New(config),
	}
}
