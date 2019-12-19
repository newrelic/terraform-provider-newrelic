package alerts

import (
	"fmt"
)

// ListAlertChannels returns all alert channels for a given account.
func (alerts *Alerts) ListAlertChannels() ([]AlertChannel, error) {
	response := alertChannelsResponse{}
	alertChannels := []AlertChannel{}
	nextURL := "/alerts_channels.json"

	for nextURL != "" {
		resp, err := alerts.client.Get(nextURL, nil, &response)

		if err != nil {
			return nil, err
		}

		alertChannels = append(alertChannels, response.Channels...)

		paging := alerts.pager.Parse(resp)
		nextURL = paging.Next
	}

	return alertChannels, nil
}

// GetAlertChannel returns a specific alert channel by ID for a given account.
func (alerts *Alerts) GetAlertChannel(id int) (*AlertChannel, error) {
	channels, err := alerts.ListAlertChannels()
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		if channel.ID == id {
			return &channel, nil
		}
	}

	return nil, fmt.Errorf("no channel found for id %d", id)
}

// CreateAlertChannel creates an alert channel within a given account.
// The configuration options different based on channel type.
// For more information on the different configurations, please
// view the New Relic API documentation for this endpoint.
// New Relic API Explorer: https://rpm.newrelic.com/api/explore/alerts_channels/create
func (alerts *Alerts) CreateAlertChannel(channel AlertChannel) (*AlertChannel, error) {
	reqBody := alertChannelRequestBody{
		Channel: channel,
	}
	resp := alertChannelsResponse{}

	_, err := alerts.client.Post("/alerts_channels.json", nil, reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Channels[0], nil
}

// DeleteAlertChannel deletes the alert channel with the specified ID.
func (alerts *Alerts) DeleteAlertChannel(id int) (*AlertChannel, error) {
	resp := alertChannelResponse{}
	url := fmt.Sprintf("/alerts_channels/%d.json", id)
	_, err := alerts.client.Delete(url, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Channel, nil
}

type alertChannelsResponse struct {
	Channels []AlertChannel `json:"channels,omitempty"`
}

type alertChannelResponse struct {
	Channel AlertChannel `json:"channel,omitempty"`
}

type alertChannelRequestBody struct {
	Channel AlertChannel `json:"channel"`
}
