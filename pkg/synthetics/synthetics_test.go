// +build unit

package synthetics

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func TestDefaultEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.Config{})

	actual := a.client.Config.BaseURL
	expected := "https://synthetics.newrelic.com/synthetics/api/v3"
	if actual != expected {
		t.Errorf("expected baseURL: %s, received: %s", expected, actual)
	}
}

func TestEUEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.Config{
		Region: config.EU,
	})

	actual := a.client.Config.BaseURL
	expected := "https://synthetics.eu.newrelic.com/synthetics/api/v3"
	if actual != expected {
		t.Errorf("expected baseURL: %s, received: %s", expected, actual)
	}
}

func TestStagingEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.Config{
		Region: config.Staging,
	})

	actual := a.client.Config.BaseURL
	expected := "https://staging-synthetics.newrelic.com/synthetics/api/v3"
	if actual != expected {
		t.Errorf("expected baseURL: %s, received: %s", expected, actual)
	}
}
