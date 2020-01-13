package alerts

import (
	"strconv"
)

// UpdatePolicyChannels updates a policy by adding the specified notification channels.
func (alerts *Alerts) UpdatePolicyChannels(policyID int, channelIDs []int) (*PolicyChannels, error) {
	channelIDStrings := make([]string, len(channelIDs))

	for i, channelID := range channelIDs {
		channelIDStrings[i] = strconv.Itoa(channelID)
	}

	queryParams := updatePolicyChannelsParams{
		PolicyID:   policyID,
		ChannelIDs: channelIDs,
	}

	resp := updatePolicyChannelsResponse{}

	_, err := alerts.client.Put("/alerts_policy_channels.json", &queryParams, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

// DeletePolicyChannel deletes a notification channel from an alert policy.
// This method returns a response containing the Channel that was deleted from the policy.
func (alerts *Alerts) DeletePolicyChannel(policyID int, channelID int) (*Channel, error) {
	queryParams := deletePolicyChannelsParams{
		PolicyID:  policyID,
		ChannelID: channelID,
	}

	resp := deletePolicyChannelResponse{}

	_, err := alerts.client.Delete("/alerts_policy_channels.json", &queryParams, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Channel, nil
}

type updatePolicyChannelsParams struct {
	PolicyID   int   `url:"policy_id,omitempty"`
	ChannelIDs []int `url:"channel_ids,omitempty"`
}

type deletePolicyChannelsParams struct {
	PolicyID  int `url:"policy_id,omitempty"`
	ChannelID int `url:"channel_id,omitempty"`
}

type updatePolicyChannelsResponse struct {
	Policy PolicyChannels `json:"policy,omitempty"`
}

type deletePolicyChannelResponse struct {
	Channel Channel           `json:"channel,omitempty"`
	Links   map[string]string `json:"channel.policy_ids,omitempty"`
}
