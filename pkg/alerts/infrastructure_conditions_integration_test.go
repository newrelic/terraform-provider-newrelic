// +build integration

package alerts

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func TestIntegrationListInfrastructureConditions(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	api := New(config.Config{
		APIKey: apiKey,
	})

	api.ListInfrastructureConditions(1234)
}
