// +build integration

package alerts

import (
	"os"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestIntegrationAlertPolicyChannels(t *testing.T) {
	t.Parallel()

	client := newAlertPolicyChannelsTestClient(t)

	// Test: Update
	updateResult := testUpdateAlertPolicyChannels(t, client, 586667, []int{2932701})

	// Test: Delete
	testDeleteAlertPolicyChannel(t, client, updateResult)
}

func testUpdateAlertPolicyChannels(t *testing.T, client Alerts, id int, channelIDs []int) *AlertPolicyChannels {
	result, err := client.UpdateAlertPolicyChannels(id, channelIDs)

	require.NoError(t, err)

	return result
}

func testDeleteAlertPolicyChannel(t *testing.T, client Alerts, alertPolicyChannels *AlertPolicyChannels) {
	p := *alertPolicyChannels
	_, err := client.DeleteAlertPolicyChannel(p.ID, p.ChannelIDs[0])

	require.NoError(t, err)
}

func newAlertPolicyChannelsTestClient(t *testing.T) Alerts {
	apiKey := os.Getenv("NEWRELIC_API_KEY")

	if apiKey == "" {
		t.Skipf("acceptance testing requires an API key")
	}

	return New(config.Config{
		APIKey: apiKey,
	})
}
