// +build unit

package alerts

import (
	"net/http"
	"testing"

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

	// Tests serialization of complex `headers` and `payload` fields
	testWebhookComplexHeadersAndPayloadResponseJSON = `{
		"channels": [
			{
				"id": 1,
				"name": "webhook-EMPTY-headers-and-payload",
				"type": "webhook",
				"configuration": {
					"base_url": "http://example.com",
					"headers": "",
					"payload": "",
					"payload_type": ""
				},
				"links": {
					"policy_ids": []
				}
			},
			{
				"id": 2,
				"name": "webhook-ESCAPED-STRING-headers-and-payload",
				"type": "webhook",
				"configuration": {
					"base_url": "http://example.com",
					"headers": "{\"key\":\"value\"}",
					"payload": "{\"key\":\"value\"}",
					"payload_type": "application/json"
				},
				"links": {
					"policy_ids": []
				}
			},
			{
				"id": 3,
				"name": "webhook-WEIRD-headers-and-payload",
				"type": "webhook",
				"configuration": {
					"base_url": "http://example.com",
					"headers": {
						"": ""
					},
					"payload": {
						"": ""
					},
					"payload_type": "application/json"
				},
				"links": {
					"policy_ids": []
				}
			},
			{
				"id": 4,
				"name": "webhook-COMPLEX-payload",
				"type": "webhook",
				"configuration": {
					"base_url": "http://example.com",
					"headers": {
						"key": "value",
						"invalidHeader": {
							"is": "allowed by the API"
						}
					},
					"payload": {
						"array": ["test", 1],
						"object": {
							"key": "value"
						}
					},
					"payload_type": "application/json"
				},
				"links": {
					"policy_ids": []
				}
			}
		]
	}`
)

func TestListChannels(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testListChannelsResponseJSON, http.StatusOK)

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

func TestListChannelsWebhookWithComplexHeadersAndPayload(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testWebhookComplexHeadersAndPayloadResponseJSON, http.StatusOK)

	expected := []*Channel{
		{
			ID:   1,
			Name: "webhook-EMPTY-headers-and-payload",
			Type: "webhook",
			Configuration: ChannelConfiguration{
				BaseURL:     "http://example.com",
				PayloadType: "",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		},
		{
			ID:   2,
			Name: "webhook-ESCAPED-STRING-headers-and-payload",
			Type: "webhook",
			Configuration: ChannelConfiguration{
				BaseURL:     "http://example.com",
				PayloadType: "application/json",
				Headers: MapStringInterface{
					"key": "value",
				},
				Payload: MapStringInterface{
					"key": "value",
				},
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		},
		{
			ID:   3,
			Name: "webhook-WEIRD-headers-and-payload",
			Type: "webhook",
			Configuration: ChannelConfiguration{
				BaseURL: "http://example.com",
				Headers: MapStringInterface{
					"": "",
				},
				Payload: MapStringInterface{
					"": "",
				},
				PayloadType: "application/json",
			},
			Links: ChannelLinks{
				PolicyIDs: []int{},
			},
		},
		{
			ID:   4,
			Name: "webhook-COMPLEX-payload",
			Type: "webhook",
			Configuration: ChannelConfiguration{
				BaseURL: "http://example.com",
				Headers: MapStringInterface{
					"key": "value",
					"invalidHeader": map[string]interface{}{
						"is": "allowed by the API",
					},
				},
				Payload: MapStringInterface{
					"array": []interface{}{"test", float64(1)},
					"object": map[string]interface{}{
						"key": "value",
					},
				},
				PayloadType: "application/json",
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
