// +build integration

package alerts

import (
	"fmt"
	"testing"

	nr "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/stretchr/testify/require"
)

func TestIntegrationPolicyChannels(t *testing.T) {
	t.Parallel()

	var (
		testPolicyNameRandStr = nr.RandSeq(5)
		testIntegrationPolicy = Policy{
			IncidentPreference: "PER_POLICY",
			Name:               fmt.Sprintf("test-alert-policy-%s", testPolicyNameRandStr),
		}
		testIntegrationChannelA = Channel{
			Name: fmt.Sprintf("test-alert-channel-%s", testPolicyNameRandStr),
			Type: "slack",
			Configuration: ChannelConfiguration{
				URL:     "https://example-org.slack.com",
				Channel: testPolicyNameRandStr,
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}
		testIntegrationChannelB = Channel{
			Name: fmt.Sprintf("test-alert-channel-%s", nr.RandSeq(5)),
			Type: "slack",
			Configuration: ChannelConfiguration{
				URL:     "https://example-org.slack.com",
				Channel: nr.RandSeq(5),
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		}
	)

	client := newIntegrationTestClient(t)

	// Setup
	policyResp, err := client.CreatePolicy(testIntegrationPolicy)
	policy := *policyResp

	require.NoError(t, err)

	channelRespA, err := client.CreateChannel(testIntegrationChannelA)
	channelRespB, err := client.CreateChannel(testIntegrationChannelB)

	channelA := *channelRespA
	channelB := *channelRespB

	require.NoError(t, err)

	// Teardown
	defer func() {
		_, err = client.DeletePolicy(policy.ID)
		if err != nil {
			t.Logf("Error cleaning up alert policy %d (%s): %s", policy.ID, policy.Name, err)
		}

		_, err = client.DeleteChannel(channelA.ID)
		if err != nil {
			t.Logf("Error cleaning up alert channel %d (%s): %s", channelA.ID, channelA.Name, err)
		}

		_, err = client.DeleteChannel(channelB.ID)
		if err != nil {
			t.Logf("Error cleaning up alert channel %d (%s): %s", channelB.ID, channelB.Name, err)
		}
	}()

	// Test: Update
	updateResult, err := client.UpdatePolicyChannels(policy.ID, []int{channelA.ID, channelB.ID})

	require.NoError(t, err)
	require.NotNil(t, updateResult)
	require.Greater(t, len(updateResult.ChannelIDs), 1)

	// Test: Delete
	deleteResultA, err := client.DeletePolicyChannel(policy.ID, updateResult.ChannelIDs[0])

	require.NoError(t, err)
	require.NotNil(t, deleteResultA)

	deleteResultB, err := client.DeletePolicyChannel(policy.ID, updateResult.ChannelIDs[1])

	require.NoError(t, err)
	require.NotNil(t, deleteResultB)
}
