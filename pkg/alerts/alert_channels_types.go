package alerts

// AlertChannel represents a New Relic alert notification channel
type AlertChannel struct {
	ID            int                       `json:"id,omitempty"`
	Name          string                    `json:"name,omitempty"`
	Type          string                    `json:"type,omitempty"`
	Configuration AlertChannelConfiguration `json:"configuration,omitempty"`
	Links         AlertChannelLinks         `json:"links,omitempty"`
}

// AlertChannelLinks represent the links between policies and alert channels
type AlertChannelLinks struct {
	PolicyIDs []int `json:"policy_ids,omitempty"`
}

// AlertChannelConfiguration represents a Configuration type within AlertChannels
type AlertChannelConfiguration struct {
	Recipients            string            `json:"recipients,omitempty"`
	IncludeJSONAttachment string            `json:"include_json_attachment,omitempty"`
	AuthToken             string            `json:"auth_token,omitempty"`
	APIKey                string            `json:"api_key,omitempty"`
	Teams                 string            `json:"teams,omitempty"`
	Tags                  string            `json:"tags,omitempty"`
	URL                   string            `json:"url,omitempty"`
	Channel               string            `json:"channel,omitempty"`
	Key                   string            `json:"key,omitempty"`
	RouteKey              string            `json:"route_key,omitempty"`
	ServiceKey            string            `json:"service_key,omitempty"`
	BaseURL               string            `json:"base_url,omitempty"`
	AuthUsername          string            `json:"auth_username,omitempty"`
	AuthPassword          string            `json:"auth_password,omitempty"`
	PayloadType           string            `json:"payload_type,omitempty"`
	Region                string            `json:"region,omitempty"`
	UserID                string            `json:"user_id,omitempty"`
	Payload               map[string]string `json:"payload,omitempty"`
	Headers               map[string]string `json:"headers,omitempty"`
}
