package alerts

import (
	"fmt"
	"net/url"
	"strconv"
)

// UpdateAlertPolicyChannels updates a policy by adding the specified notification channels.
func (alerts *Alerts) UpdateAlertPolicyChannels(policyID int, channelIDs []int) error {
	channelIDStrings := make([]string, len(channelIDs))

	for i, channelID := range channelIDs {
		channelIDStrings[i] = strconv.Itoa(channelID)
	}

	resp := alertPolicyChannelsResponse{}

	reqURL, err := url.Parse("/alerts_policy_channels.json")
	if err != nil {
		return err
	}

	qs := url.Values{
		"policy_id":   strconv.Itoa(policyID)},
		"channel_ids": channelIDStrings,
	}
	reqURL.RawQuery = qs.Encode()

	resp, err := alerts.client.Put(reqURL.String(), nil, nil, &resp)

	if err != nil {
		return nil, err
	}

	fmt.Println(resp)

	return &resp.Policy, nil
}

type alertPolicyChannelsResponse struct {
	Policy AlertPolicyChannels `json:"policy,omitempty"`
}
