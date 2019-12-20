// +build integration

package alerts

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
)

var (
	testAlertChannelEmail = AlertChannel{
		Name: "integration-test-email",
		Type: "email",
		Configuration: AlertChannelConfiguration{
			Recipients:            "devtoolkittest@newrelic.com",
			IncludeJSONAttachment: "true",
		},
		Links: AlertChannelLinks{
			PolicyIDs: []int{},
		},
	}

	testAlertChannelOpsGenie = AlertChannel{
		Name: "integration-test-opsgenie",
		Type: "opsgenie",
		Configuration: AlertChannelConfiguration{
			APIKey:     "abc123",
			Teams:      "dev-toolkit",
			Tags:       "tag1,tag2",
			Recipients: "devtoolkittest@newrelic.com",
		},
		Links: AlertChannelLinks{
			PolicyIDs: []int{},
		},
	}

	testAlertChannelSlack = AlertChannel{
		Name: "integration-test-slack",
		Type: "slack",
		Configuration: AlertChannelConfiguration{
			URL:     "https://example-org.slack.com",
			Channel: "test-channel",
		},
		Links: AlertChannelLinks{
			PolicyIDs: []int{},
		},
	}

	testAlertChannelVictorops = AlertChannel{
		Name: "integration-test-victorops",
		Type: "victorops",
		Configuration: AlertChannelConfiguration{
			Key:      "abc123",
			RouteKey: "/route-name",
		},
		Links: AlertChannelLinks{
			PolicyIDs: []int{},
		},
	}

	testAlertChannelWebhook = AlertChannel{
		Name: "integration-test-webhook",
		Type: "webhook",
		Configuration: AlertChannelConfiguration{
			BaseURL:      "https://test.com",
			AuthUsername: "devtoolkit",
			AuthPassword: "123abc",
			PayloadType:  "application/json",
			Payload: map[string]string{
				"account_id":               "$ACCOUNT_ID",
				"account_name":             "$ACCOUNT_NAME",
				"condition_id":             "$CONDITION_ID",
				"condition_name":           "$CONDITION_NAME",
				"current_state":            "$EVENT_STATE",
				"details":                  "$EVENT_DETAILS",
				"event_type":               "$EVENT_TYPE",
				"incident_acknowledge_url": "$INCIDENT_ACKNOWLEDGE_URL",
				"incident_id":              "$INCIDENT_ID",
				"incident_url":             "$INCIDENT_URL",
				"owner":                    "$EVENT_OWNER",
				"policy_name":              "$POLICY_NAME",
				"policy_url":               "$POLICY_URL",
				"runbook_url":              "$RUNBOOK_URL",
				"severity":                 "$SEVERITY",
				"targets":                  "$TARGETS",
				"timestamp":                "$TIMESTAMP",
				"violation_chart_url":      "$VIOLATION_CHART_URL",
			},
			Headers: map[string]string{
				"x-test-header": "test-header",
			},
		},
		Links: AlertChannelLinks{
			PolicyIDs: []int{},
		},
	}
)

func TestIntegrationAlertChannel(t *testing.T) {
	t.Parallel()

	channels := []AlertChannel{
		testAlertChannelEmail,
		testAlertChannelOpsGenie,
		testAlertChannelSlack,
		testAlertChannelVictorops,
		testAlertChannelWebhook,
	}

	client := newChannelsTestClient(t)

	for _, channel := range channels {
		// Test: Create
		createResult := testCreateAlertChannel(t, client, channel)

		// Test: Read
		readResult := testReadAlertChannel(t, client, createResult)

		// Test: Delete
		testDeleteAlertChannel(t, client, readResult)
	}
}

func testCreateAlertChannel(t *testing.T, client Alerts, channel AlertChannel) *AlertChannel {
	result, err := client.CreateAlertChannel(channel)

	require.NoError(t, err)

	return result
}

func testReadAlertChannel(t *testing.T, client Alerts, channel *AlertChannel) *AlertChannel {
	result, err := client.GetAlertChannel(channel.ID)

	require.NoError(t, err)

	return result
}

func testDeleteAlertChannel(t *testing.T, client Alerts, channel *AlertChannel) {
	p := *channel
	_, err := client.DeleteAlertChannel(p.ID)

	require.NoError(t, err)
}

func newChannelsTestClient(t *testing.T) Alerts {
	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	return New(config.Config{
		APIKey: apiKey,
	})
}
