// +build integration

package alerts

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestIntegrationChannel(t *testing.T) {
	t.Parallel()

	var (
		testChannelEmail = Channel{
			Name: "integration-test-email",
			Type: "email",
			Configuration: ChannelConfiguration{
				Recipients:            "devtoolkittest@newrelic.com",
				IncludeJSONAttachment: "true",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		testChannelOpsGenie = Channel{
			Name: "integration-test-opsgenie",
			Type: "opsgenie",
			Configuration: ChannelConfiguration{
				APIKey:     "abc123",
				Teams:      "dev-toolkit",
				Tags:       "tag1,tag2",
				Recipients: "devtoolkittest@newrelic.com",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		testChannelSlack = Channel{
			Name: "integration-test-slack",
			Type: "slack",
			Configuration: ChannelConfiguration{
				URL:     "https://example-org.slack.com",
				Channel: "test-channel",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		testChannelVictorops = Channel{
			Name: "integration-test-victorops",
			Type: "victorops",
			Configuration: ChannelConfiguration{
				Key:      "abc123",
				RouteKey: "/route-name",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		testChannelWebhook = Channel{
			Name: "integration-test-webhook",
			Type: "webhook",
			Configuration: ChannelConfiguration{
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
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		channels = []Channel{
			testChannelEmail,
			testChannelOpsGenie,
			testChannelSlack,
			testChannelVictorops,
			testChannelWebhook,
		}
	)

	client := newChannelsTestClient(t)

	for _, channel := range channels {
		// Test: Create
		createResult := testCreateChannel(t, client, channel)

		// Test: Read
		readResult := testReadChannel(t, client, createResult)

		// Test: Delete
		testDeleteChannel(t, client, readResult)
	}
}

func testCreateChannel(t *testing.T, client Alerts, channel Channel) *Channel {
	result, err := client.CreateChannel(channel)

	require.NoError(t, err)

	return result
}

func testReadChannel(t *testing.T, client Alerts, channel *Channel) *Channel {
	result, err := client.GetChannel(channel.ID)

	require.NoError(t, err)

	return result
}

func testDeleteChannel(t *testing.T, client Alerts, channel *Channel) {
	p := *channel
	_, err := client.DeleteChannel(p.ID)

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
