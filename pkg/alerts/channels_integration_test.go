// +build integration

package alerts

import (
	"testing"

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
				BaseURL:     "https://test.com",
				PayloadType: "application/json",
				Headers: MapStringString{
					"x-test-header": "test-header",
				},
				Payload: MapStringString{
					"account_id": "123",
				},
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		testChannelWebhookEmptyHeadersAndPayload = Channel{
			Name: "integration-test-webhook-empty-headers-and-payload",
			Type: "webhook",
			Configuration: ChannelConfiguration{
				BaseURL: "https://test.com",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}

		testChannelWebhookWeirdHeadersAndPayload = Channel{
			Name: "integration-test-webhook-weird-headers-and-payload",
			Type: "webhook",
			Configuration: ChannelConfiguration{
				BaseURL: "https://test.com",
				Headers: MapStringString{
					"": "",
				},
				Payload: MapStringString{
					"": "",
				},
				PayloadType: "application/json",
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
			testChannelWebhookEmptyHeadersAndPayload,
			testChannelWebhookWeirdHeadersAndPayload,
		}
	)

	client := newIntegrationTestClient(t)

	for _, channel := range channels {
		// Test: Create
		created, err := client.CreateChannel(channel)

		require.NoError(t, err)
		require.NotNil(t, created)

		// Test: Read
		read, err := client.GetChannel(created.ID)

		require.NoError(t, err)
		require.NotNil(t, read)

		// Test: Delete
		deleted, err := client.DeleteChannel(read.ID)

		require.NoError(t, err)
		require.NotNil(t, deleted)
	}
}
