package newrelic

import (
	"github.com/newrelic/newrelic-client-go/pkg/apm"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/infrastructure"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

type NewRelic struct {
	APM            apm.APM
	Synthetics     synthetics.Synthetics
	Infrastructure infrastructure.Infrastructure
}

func New(config config.Config) NewRelic {
	return NewRelic{
		APM:            apm.New(config),
		Infrastructure: infrastructure.New(config),
		Synthetics:     synthetics.New(config),
	}
}
