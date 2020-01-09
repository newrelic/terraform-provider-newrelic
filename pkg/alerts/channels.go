package alerts

import (
	"fmt"
)

// ListChannels returns all alert channels for a given account.
func (alerts *Alerts) ListChannels() ([]*Channel, error) {
	response := alertChannelsResponse{}
	alertChannels := []*Channel{}
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

// GetChannel returns a specific alert channel by ID for a given account.
func (alerts *Alerts) GetChannel(id int) (*Channel, error) {
	channels, err := alerts.ListChannels()
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		if channel.ID == id {
			return channel, nil
		}
	}

	return nil, fmt.Errorf("no channel found for id %d", id)
}

// CreateChannel creates an alert channel within a given account.
// The configuration options different based on channel type.
// For more information on the different configurations, please
// view the New Relic API documentation for this endpoint.
// Docs: https://docs.newrelic.com/docs/alerts/rest-api-alerts/new-relic-alerts-rest-api/rest-api-calls-new-relic-alerts#channels
func (alerts *Alerts) CreateChannel(channel Channel) (*Channel, error) {
	reqBody := alertChannelRequestBody{
		Channel: channel,
	}
	resp := alertChannelsResponse{}

	_, err := alerts.client.Post("/alerts_channels.json", nil, &reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Channels[0], nil
}

// DeleteChannel deletes the alert channel with the specified ID.
func (alerts *Alerts) DeleteChannel(id int) (*Channel, error) {
	resp := alertChannelResponse{}
	url := fmt.Sprintf("/alerts_channels/%d.json", id)
	_, err := alerts.client.Delete(url, nil, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Channel, nil
}

type alertChannelsResponse struct {
	Channels []*Channel `json:"channels,omitempty"`
}

type alertChannelResponse struct {
	Channel Channel `json:"channel,omitempty"`
}

type alertChannelRequestBody struct {
	Channel Channel `json:"channel"`
}
