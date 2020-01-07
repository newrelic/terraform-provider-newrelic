// +build integration

package apm

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func TestIntegrationListApplications(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	api := New(config.Config{
		APIKey:   apiKey,
		LogLevel: "debug",
	})

	_, err := api.ListApplications(nil)

	if err != nil {
		t.Fatalf("ListApplications error: %s", err)
	}
}

func TestIntegrationGetApplication(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	api := New(config.Config{
		APIKey:   apiKey,
		LogLevel: "debug",
	})

	a, err := api.ListApplications(nil)

	if err != nil {
		t.Fatal(err)
	}

	_, err = api.GetApplication(a[0].ID)

	if err != nil {
		t.Fatal(err)
	}
}

func TestIntegrationUpdateApplication(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	api := New(config.Config{
		APIKey:   apiKey,
		LogLevel: "debug",
	})

	a, err := api.ListApplications(nil)

	if err != nil {
		t.Fatal(err)
	}

	_, err = api.GetApplication(a[0].ID)

	if err != nil {
		t.Fatal(err)
	}

	params := UpdateApplicationParams{
		Name:     a[0].Name,
		Settings: a[0].Settings,
	}

	_, err = api.UpdateApplication(a[0].ID, params)

	if err != nil {
		t.Fatal(err)
	}
}

func TestIntegrationDeleteApplication(t *testing.T) {
	t.Skip()
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	api := New(config.Config{
		APIKey:   apiKey,
		LogLevel: "debug",
	})

	_, err := api.DeleteApplication(0)

	if err != nil {
		t.Fatal(err)
	}
}
