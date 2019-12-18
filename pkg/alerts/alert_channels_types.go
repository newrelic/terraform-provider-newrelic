package alerts

// AlertChannel represents a New Relic alert notification channel
type AlertChannel struct {
	ID            int                    `json:"id,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Type          string                 `json:"type,omitempty"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
	Links         AlertChannelLinks      `json:"links,omitempty"`
}

// AlertChannelLinks represent the links between policies and alert channels
type AlertChannelLinks struct {
	PolicyIDs []int `json:"policy_ids,omitempty"`
}
