// +build integration

package synthetics

import (
	"os"
	"testing"
	"time"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/assert"
)

var (
	testMonitorID    = "72733a02-9701-4279-8ac3-8f6281a5a1a9"
	testTimestamp, _ = time.Parse(time.RFC3339, "2019-11-27T19:11:05.076+0000")

	testMonitorOptions = MonitorOptions{
		ValidationString:       "",
		VerifySSL:              false,
		BypassHEADRequest:      false,
		TreatRedirectAsFailure: false,
	}

	testMonitor = Monitor{
		Name:         "test-synthetics-monitor",
		Type:         MonitorTypes.Simple,
		Frequency:    15,
		URI:          "https://google.com",
		Locations:    []string{"AWS_US_EAST_1"},
		Status:       MonitorStatus.Disabled,
		SLAThreshold: 7,
		UserID:       0,
		APIVersion:   "LATEST",
		Options:      testMonitorOptions,
	}
)

func TestIntegrationMonitors(t *testing.T) {
	t.Parallel()

	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	monitors := New(config.Config{
		APIKey: apiKey,
	})

	// Test: Create
	created, err := monitors.CreateMonitor(testMonitor)
	if err != nil {
		t.Fatalf("CreateMonitors error: %s", err)
	}

	assert.NotNil(t, created)

	// Test: List
	multiple, err := monitors.ListMonitors()
	if err != nil {
		t.Fatalf("ListMonitors error: %s", err)
	}

	assert.NotNil(t, multiple)

	// Test: Get
	single, err := monitors.GetMonitor(multiple[0].ID)
	if err != nil {
		t.Fatalf("GetMonitors error: %s", err)
	}

	assert.NotNil(t, single)

	// Test: Update
	single.Name = "updated"
	updated, err := monitors.UpdateMonitor(*single)
	if err != nil {
		t.Fatalf("UpdateMonitors error: %s", err)
	}

	assert.NotNil(t, updated)

	// Test: Delete
	deleted, err := monitors.DeleteMonitor(updated.ID)
	if err != nil {
		t.Fatalf("DeleteMonitors error: %s", err)
	}

	assert.NotNil(t, deleted)
}
