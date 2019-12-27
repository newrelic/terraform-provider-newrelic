package alerts

// PolicyChannels represents a New Relic Alert Policy Channel
type PolicyChannels struct {
	ID         int   `json:"id,omitempty"`
	ChannelIDs []int `json:"channel_ids,omitempty"`
}
