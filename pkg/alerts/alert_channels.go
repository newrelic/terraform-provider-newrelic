package alerts

import (
	"fmt"
)

type createAlertChannelRequestBody struct {
	Channel AlertChannel `json:"channel"`
}

type createAlertChannelResponse struct {
	Channels []AlertChannel `json:"channels,omitempty"`
}

type listAlertChannelsResponse struct {
	Channels []AlertChannel `json:"channels,omitempty"`
}

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

func (alerts *Alerts) ListAlertChannels() ([]AlertChannel, error) {
	res := listAlertChannelsResponse{}
	err := alerts.client.Get("/alerts_channels.json", nil, &res)

	if err != nil {
		return nil, err
	}

	return res.Channels, nil
}

func (alerts *Alerts) CreateAlertChannel(channel AlertChannel) (*[]AlertChannel, error) {
	reqBody := createAlertChannelRequestBody{
		Channel: channel,
	}
	resp := createAlertChannelResponse{}

	err := alerts.client.Post("/alerts_channels.json", reqBody, &resp)

	if err != nil {
		return nil, err
	}

	return &resp.Channels, nil
}

// DeleteAlertChannel deletes the alert channel with the specified ID.
func (alerts *Alerts) DeleteAlertChannel(id int) error {
	url := fmt.Sprintf("/alerts_channels/%d.json", id)
	err := alerts.client.Delete(url)

	return err
}
