package alerts

type AlertPolicyChannels struct {
	ID         int   `json:"id,omitempty"`
	ChannelIDs []int `json:"channel_ids,omitempty"`
}

type AlertPolicyChannelsParams struct {
	PolicyID   int
	ChannelIDs []int
}
