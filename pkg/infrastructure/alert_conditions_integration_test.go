// +build integration

package infrastructure

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func TestIntegrationListAlertConditions(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	api := New(config.Config{
		APIKey: apiKey,
	})

	api.ListAlertConditions(1234)
}
