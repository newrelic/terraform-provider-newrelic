package alerts

import (
	"strconv"
	"strings"
)

// UpdateAlertPolicyChannels updates a policy by adding the specified notification channels.
func (alerts *Alerts) UpdateAlertPolicyChannels(policyID int, channelIDs []int) (*AlertPolicyChannels, error) {
	channelIDStrings := make([]string, len(channelIDs))

	for i, channelID := range channelIDs {
		channelIDStrings[i] = strconv.Itoa(channelID)
	}

	queryParams := map[string]string{
		"policy_id":   strconv.Itoa(policyID),
		"channel_ids": strings.Join(channelIDStrings, ","),
	}

	resp := alertPolicyChannelsResponse{}

	_, err := alerts.client.Put("/alerts_policy_channels.json", &queryParams, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Policy, nil
}

// DeleteAlertPolicyChannel deletes a notification channel from an alert policy.
// This method returns a response containing the AlertChannel that was deleted from the policy.
func (alerts *Alerts) DeleteAlertPolicyChannel(policyID int, channelID int) (*AlertChannel, error) {
	queryParams := map[string]string{
		"policy_id":  strconv.Itoa(policyID),
		"channel_id": strconv.Itoa(channelID),
	}

	resp := deleteAlertPolicyChannelResponse{}

	_, err := alerts.client.Delete("/alerts_policy_channels.json", &queryParams, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Channel, nil
}

type alertPolicyChannelsResponse struct {
	Policy AlertPolicyChannels `json:"policy,omitempty"`
}

type deleteAlertPolicyChannelResponse struct {
	Channel AlertChannel      `json:"channel,omitempty"`
	Links   map[string]string `json:"channel.policy_ids,omitempty"`
}
