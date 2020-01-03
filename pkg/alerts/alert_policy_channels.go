package alerts

import (
	"strconv"
	"strings"

	"github.com/newrelic/newrelic-client-go/internal/http"
)

// UpdatePolicyChannels updates a policy by adding the specified notification channels.
func (alerts *Alerts) UpdatePolicyChannels(policyID int, channelIDs []int) (*PolicyChannels, error) {
	channelIDStrings := make([]string, len(channelIDs))

	for i, channelID := range channelIDs {
		channelIDStrings[i] = strconv.Itoa(channelID)
	}

	queryParams := []http.QueryParam{
		{Name: "policy_id", Value: strconv.Itoa(policyID)},
		{Name: "channel_ids", Value: strings.Join(channelIDStrings, ",")},
	}

	resp := updatePolicyChannelsResponse{}

	_, err := alerts.client.Put("/alerts_policy_channels.json", &queryParams, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

// DeletePolicyChannel deletes a notification channel from an alert policy.
// This method returns a response containing the AlertChannel that was deleted from the policy.
func (alerts *Alerts) DeletePolicyChannel(policyID int, channelID int) (*AlertChannel, error) {
	queryParams := []http.QueryParam{
		{Name: "policy_id", Value: strconv.Itoa(policyID)},
		{Name: "channel_id", Value: strconv.Itoa(channelID)},
	}

	resp := deletePolicyChannelResponse{}

	_, err := alerts.client.Delete("/alerts_policy_channels.json", &queryParams, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Channel, nil
}

type updatePolicyChannelsResponse struct {
	Policy PolicyChannels `json:"policy,omitempty"`
}

type deletePolicyChannelResponse struct {
	Channel AlertChannel      `json:"channel,omitempty"`
	Links   map[string]string `json:"channel.policy_ids,omitempty"`
}
