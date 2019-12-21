// +build integration

package alerts

import (
	"fmt"
	"os"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/require"
)

var (
	testRandStr                = nr.RandSeq(5)
	testIntegrationAlertPolicy = AlertPolicy{
		IncidentPreference: "PER_POLICY",
		Name:               fmt.Sprintf("test-alert-policy-%s", testRandStr),
	}
	testIntegrationAlertChannel = AlertChannel{
		Name: fmt.Sprintf("test-alert-channel-%s", testRandStr),
		Type: "slack",
		Configuration: AlertChannelConfiguration{
			URL:     "https://example-org.slack.com",
			Channel: testRandStr,
		},
		Links: AlertChannelLinks{
			PolicyIDs: []int{},
		},
	}
)

func TestIntegrationAlertPolicyChannels(t *testing.T) {
	t.Parallel()

	client := newAlertPolicyChannelsTestClient(t)

	// Setup
	policyResp, err := client.CreateAlertPolicy(testIntegrationAlertPolicy)
	policy := *policyResp

	require.NoError(t, err)

	channelResp, err := client.CreateAlertChannel(testIntegrationAlertChannel)
	channel := *channelResp

	require.NoError(t, err)

	// Teardown
	defer func() {
		_, err = client.DeleteAlertPolicy(policy.ID)
		if err != nil {
			t.Logf("Error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}

		_, err = client.DeleteAlertChannel(channel.ID)
		if err != nil {
			t.Logf("Error cleaning up alert channel %d (%s): %s", channel.ID, channel.Name, err)
		}
	}()

	// Test: Update
	updateResult, err := client.UpdateAlertPolicyChannels(policy.ID, []int{channel.ID})

	require.NoError(t, err)
	require.NotNil(t, updateResult)

	// Test: Delete
	deleteResult, err := client.DeleteAlertPolicyChannel(policy.ID, updateResult.ChannelIDs[0])

	require.NoError(t, err)
	require.NotNil(t, deleteResult)
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
