// +build unit

package alerts

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/assert"
)

var (
	testListChannelsResponseJSON = `{
		"channels": [
			{
				"id": 2803426,
				"name": "unit-test-alert-channel",
				"type": "user",
				"configuration": {
					"user_id": "2680539"
				},
				"links": {
					"policy_ids": []
				}
			},
			{
				"id": 2932511,
				"name": "test@testing.com",
				"type": "email",
				"configuration": {
					"include_json_attachment": "true",
					"recipients": "test@testing.com"
				},
				"links": {
					"policy_ids": []
				}
			}
		]
	}`

	testCreateChannelResponseJSON = `{
		"channels": [
			{
				"id": 2932701,
				"name": "sblue@newrelic.com",
				"type": "email",
				"configuration": {
					"include_json_attachment": "true",
					"recipients": "sblue@newrelic.com"
				},
				"links": {
					"policy_ids": []
				}
			}
		],
		"links": {
			"channel.policy_ids": "/v2/policies/{policy_id}"
		}
	}`

	testDeleteChannelResponseJSON = `{
		"channel": {
			"id": 2932511,
			"name": "test@example.com",
			"type": "email",
			"configuration": {
				"include_json_attachment": "true",
				"recipients": "test@example.com"
			},
			"links": {
				"policy_ids": []
			}
		},
		"links": {
			"channel.policy_ids": "/v2/policies/{policy_id}"
		}
	}`
)

func TestListAlertChannels(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListChannelsResponseJSON, http.StatusOK)

	expected := []AlertChannel{
		{
			ID:   2803426,
			Name: "unit-test-alert-channel",
			Type: "user",
			Configuration: AlertChannelConfiguration{
				UserID: "2680539",
			},
			Links: AlertChannelLinks{
				PolicyIDs: []int{},
			},
		},
		{
			ID:   2932511,
			Name: "test@testing.com",
			Type: "email",
			Configuration: AlertChannelConfiguration{
				Recipients:            "test@testing.com",
				IncludeJSONAttachment: "true",
			},
			Links: AlertChannelLinks{
				PolicyIDs: []int{},
			},
		},
	}

	actual, err := alerts.ListAlertChannels()

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetAlertChannel(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListChannelsResponseJSON, http.StatusOK)

	expected := AlertChannel{
		ID:   2803426,
		Name: "unit-test-alert-channel",
		Type: "user",
		Configuration: AlertChannelConfiguration{
			UserID: "2680539",
		},
		Links: AlertChannelLinks{
			PolicyIDs: []int{},
		},
	}

	actual, err := alerts.GetAlertChannel(2803426)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, *actual)
}

func TestGetAlertChannelNotFound(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListChannelsResponseJSON, http.StatusOK)

	actual, err := alerts.GetAlertChannel(0)

	assert.Error(t, err)
	assert.Nil(t, actual)
	assert.Equal(t, "no channel found for id 0", err.Error())
}

func TestCreateAlertChannel(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testCreateChannelResponseJSON, http.StatusCreated)

	channel := AlertChannel{
		Name: "sblue@newrelic.com",
		Type: "email",
		Configuration: AlertChannelConfiguration{
			Recipients:            "sblue@newrelic.com",
			IncludeJSONAttachment: "true",
		},
	}

	expected := AlertChannel{
		ID:   2932701,
		Name: "sblue@newrelic.com",
		Type: "email",
		Configuration: AlertChannelConfiguration{
			Recipients:            "sblue@newrelic.com",
			IncludeJSONAttachment: "true",
		},
		Links: AlertChannelLinks{
			PolicyIDs: []int{},
		},
	}

	actual, err := alerts.CreateAlertChannel(channel)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, *actual)
}

func TestCreateAlertChannelInvalidChannelType(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, `{
		"error": {
			"title": "Invalid channel type"
		}
	}`, http.StatusUnprocessableEntity)

	channel := AlertChannel{
		Name:          "string",
		Type:          "string",
		Configuration: AlertChannelConfiguration{},
	}

	actual, err := alerts.CreateAlertChannel(channel)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestDeleteAlertChannel(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testDeleteChannelResponseJSON, http.StatusOK)

	_, err := alerts.DeleteAlertChannel(2932511)

	assert.NoError(t, err)
}

func newTestClient(handler http.Handler) Alerts {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}
