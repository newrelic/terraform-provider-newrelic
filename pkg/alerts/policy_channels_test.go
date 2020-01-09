// +build unit

package alerts

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testUpdatePolicyChannelsResponseJSON = `{
		"policy": {
			"id": 593436,
			"channel_ids": [
				2932701,
				2932702
			]
		}
	}`

	testDeletePolicyChannelResponseJSON = `{
		"channel": {
			"id": 2932701,
			"name": "test@example.com",
			"type": "email",
			"configuration": {
				"include_json_attachment": "true",
				"recipients": "test@example.com"
			},
			"links": {
				"policy_ids": [
					593436
				]
			}
		},
		"links": {
			"channel.policy_ids": "/v2/policies/{policy_id}"
		}
	}`
)

func TestUpdatePolicyChannels(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testUpdatePolicyChannelsResponseJSON, http.StatusOK)

	actual, err := alerts.UpdatePolicyChannels(593436, []int{2932701, 2932702})

	expected := PolicyChannels{
		ID:         593436,
		ChannelIDs: []int{2932701, 2932702},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, *actual)
}

func TestDeletePolicyChannel(t *testing.T) {
	t.Parallel()
	alerts := newMockResponse(t, testDeletePolicyChannelResponseJSON, http.StatusOK)

	expected := Channel{
		ID:   2932701,
		Name: "test@example.com",
		Type: "email",
		Configuration: ChannelConfiguration{
			IncludeJSONAttachment: "true",
			Recipients:            "test@example.com",
		},
		Links: ChannelLinks{
			PolicyIDs: []int{593436},
		},
	}

	actual, err := alerts.DeletePolicyChannel(593436, 2932701)

	assert.NoError(t, err)
	assert.Equal(t, expected, *actual)
}
