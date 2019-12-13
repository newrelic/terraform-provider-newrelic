package apm

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func TestDefaultEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.Config{})

	actual := a.client.Client.HostURL
	expected := "https://api.newrelic.com/v2"
	if actual != expected {
		t.Errorf("expected baseURL: %s, received: %s", expected, actual)
	}
}

func TestEUEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.Config{
		Region: config.EU,
	})

	actual := a.client.Client.HostURL
	expected := "https://api.eu.newrelic.com/v2"
	if actual != expected {
		t.Errorf("expected baseURL: %s, received: %s", expected, actual)
	}
}

func TestStagingEnvironment(t *testing.T) {
	t.Parallel()
	a := New(config.Config{
		Region: config.Staging,
	})

	actual := a.client.Client.HostURL
	expected := "https://staging-api.newrelic.com/v2"
	if actual != expected {
		t.Errorf("expected baseURL: %s, received: %s", expected, actual)
	}
}
