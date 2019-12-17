// +build unit

package infrastructure

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func TestDefaultEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.ReplacementConfig{})

	actual := a.client.Config.BaseURL
	expected := "https://infra-api.newrelic.com/v2/alerts/conditions"
	if actual != expected {
		t.Errorf("expected baseURL: %s, received: %s", expected, actual)
	}
}

func TestEUEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.ReplacementConfig{
		Region: config.Region.EU,
	})

	actual := a.client.Config.BaseURL
	expected := "https://infra-api.eu.newrelic.com/v2/alerts/conditions"
	if actual != expected {
		t.Errorf("expected baseURL: %s, received: %s", expected, actual)
	}
}

func TestStagingEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.ReplacementConfig{
		Region: config.Region.Staging,
	})

	actual := a.client.Config.BaseURL
	expected := "https://staging-infra-api.newrelic.com/v2/alerts/conditions"
	if actual != expected {
		t.Errorf("expected baseURL: %s, received: %s", expected, actual)
	}
}
