package alerts

type PolicyChannels struct {
	ID         int   `json:"id,omitempty"`
	ChannelIDs []int `json:"channel_ids,omitempty"`
}
