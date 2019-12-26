// +build integration

package apm

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestIntegrationListComponents(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	api := New(config.Config{
		APIKey: apiKey,
	})

	_, err := api.ListComponents(nil)

	require.NoError(t, err)
}

func TestIntegrationGetComponent(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	api := New(config.Config{
		APIKey: apiKey,
	})

	a, err := api.ListComponents(nil)

	require.NoError(t, err)

	_, err = api.GetComponent(a[0].ID)

	require.NoError(t, err)
}
