// +build integration

package synthetics

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func TestAccListMonitors(t *testing.T) {
	t.Parallel()
	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	synthetics := New(config.Config{
		APIKey: apiKey,
	})

	_, err := synthetics.ListMonitors()

	if err != nil {
		t.Fatalf("ListMonitors error: %s", err)
	}
}
