package testing

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/region"
)

const (
	AdminAPIKey    = "adminAPIKey"                                     // AdminAPIKey used in mock configs (from Environment for Integration tests)
	LogLevel       = "debug"                                           // LogLevel used in mock configs
	PersonalAPIKey = "personalAPIKey"                                  // PersonalAPIKey used in mock configs (from Environment for Integration tests)
	UserAgent      = "newrelic/newrelic-client-go (automated testing)" // UserAgent used in mock configs
)

// NewTestConfig returns a fully saturated configration with modified BaseURLs
// for all endpoints based on the test server passed in
func NewTestConfig(t *testing.T, testServer *httptest.Server) config.Config {
	cfg := config.New()

	// Set some defaults from Testing constants
	cfg.AdminAPIKey = AdminAPIKey
	cfg.LogLevel = LogLevel
	cfg.PersonalAPIKey = PersonalAPIKey
	cfg.UserAgent = UserAgent

	if testServer != nil {
		cfg.Region().SetInfrastructureBaseURL(testServer.URL)
		cfg.Region().SetNerdGraphBaseURL(testServer.URL)
		cfg.Region().SetRestBaseURL(testServer.URL)
		cfg.Region().SetSyntheticsBaseURL(testServer.URL)
	}

	return cfg
}

// NewIntegrationTestConfig grabs environment vars for required fields or skips the test.
// returns a fully saturated configuration
func NewIntegrationTestConfig(t *testing.T) config.Config {
	envPersonalAPIKey := os.Getenv("NEW_RELIC_API_KEY")
	envAdminAPIKey := os.Getenv("NEW_RELIC_ADMIN_API_KEY")
	envRegion := os.Getenv("NEW_RELIC_REGION")

	if envPersonalAPIKey == "" && envAdminAPIKey == "" {
		t.Skipf("acceptance testing requires NEW_RELIC_API_KEY and NEW_RELIC_ADMIN_API_KEY")
	}

	cfg := config.New()

	// Set some defaults
	cfg.LogLevel = LogLevel
	cfg.UserAgent = UserAgent

	cfg.PersonalAPIKey = envPersonalAPIKey
	cfg.AdminAPIKey = envAdminAPIKey

	if envRegion != "" {
		err := cfg.SetRegion(region.Parse(envRegion))
		assert.NoError(t, err)
	}

	return cfg
}
