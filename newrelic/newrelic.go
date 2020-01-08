package newrelic

import (
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/dashboards"
	"github.com/newrelic/newrelic-client-go/pkg/plugins"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

// NewRelic is a collection of New Relic APIs.
type NewRelic struct {
	Alerts     alerts.Alerts
	APM        apm.APM
	Dashboards dashboards.Dashboards
	Plugins    plugins.Plugins
	Synthetics synthetics.Synthetics
}

// New returns a collection of New Relic APIs.
func New(config config.Config) NewRelic {
	return NewRelic{
		Alerts:     alerts.New(config),
		APM:        apm.New(config),
		Dashboards: dashboards.New(config),
		Plugins:    plugins.New(config),
		Synthetics: synthetics.New(config),
	}
}
