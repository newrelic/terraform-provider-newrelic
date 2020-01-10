// +build unit

package alerts

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	mock "github.com/newrelic/newrelic-client-go/internal/testing"
	"github.com/newrelic/newrelic-client-go/pkg/config"
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
				"name": "test@example.com",
				"type": "email",
				"configuration": {
					"include_json_attachment": "true",
					"recipients": "test@example.com"
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

func TestListChannels(t *testing.T) {
	t.Parallel()

	ts := mock.NewMockResponse(t, testListChannelsResponseJSON, http.StatusOK, "")

	alerts := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		UserAgent: "newrelic/newrelic-client-go",
	})

	expected := []*Channel{
		{
			ID:   2803426,
			Name: "unit-test-alert-channel",
			Type: "user",
			Configuration: ChannelConfiguration{
				UserID: "2680539",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		},
		{
			ID:   2932511,
			Name: "test@testing.com",
			Type: "email",
			Configuration: ChannelConfiguration{
				Recipients:            "test@testing.com",
				IncludeJSONAttachment: "true",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		},
	}

	actual, err := alerts.ListChannels()

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetChannel(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListChannelsResponseJSON, http.StatusOK)

	expected := &Channel{
		ID:   2803426,
		Name: "unit-test-alert-channel",
		Type: "user",
		Configuration: ChannelConfiguration{
			UserID: "2680539",
		},
		Links: ChannelLinks{
			PolicyIDs: []int{},
		},
	}

	actual, err := alerts.GetChannel(2803426)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestGetChannelNotFound(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListChannelsResponseJSON, http.StatusOK)

	actual, err := alerts.GetChannel(0)

	assert.Error(t, err)
	assert.Nil(t, actual)
	assert.Equal(t, "no channel found for id 0", err.Error())
}

func TestCreateChannel(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testCreateChannelResponseJSON, http.StatusCreated)

	channel := Channel{
		Name: "test@example.com",
		Type: "email",
		Configuration: ChannelConfiguration{
			Recipients:            "test@example.com",
			IncludeJSONAttachment: "true",
		},
	}

	expected := &Channel{
		ID:   2932701,
		Name: "test@example.com",
		Type: "email",
		Configuration: ChannelConfiguration{
			Recipients:            "test@example.com",
			IncludeJSONAttachment: "true",
		},
		Links: ChannelLinks{
			PolicyIDs: []int{},
		},
	}

	actual, err := alerts.CreateChannel(channel)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestCreateChannelInvalidChannelType(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, `{
		"error": {
			"title": "Invalid channel type"
		}
	}`, http.StatusUnprocessableEntity)

	channel := Channel{
		Name:          "string",
		Type:          "string",
		Configuration: ChannelConfiguration{},
	}

	actual, err := alerts.CreateChannel(channel)

	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestDeleteChannel(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testDeleteChannelResponseJSON, http.StatusOK)

	expected := &Channel{
		ID:   2932511,
		Name: "test@example.com",
		Type: "email",
		Configuration: ChannelConfiguration{
			Recipients:            "test@example.com",
			IncludeJSONAttachment: "true",
		},
		Links: ChannelLinks{
			PolicyIDs: []int{},
		},
	}

	actual, err := alerts.DeleteChannel(2932511)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}
