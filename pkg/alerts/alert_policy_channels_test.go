// +build unit

package alerts

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	alerts := newMocknewTestPolicyChannelsClientResponse(t, testUpdatePolicyChannelsResponseJSON, http.StatusOK)

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
	alerts := newMocknewTestPolicyChannelsClientResponse(t, testDeletePolicyChannelResponseJSON, http.StatusOK)

	expected := AlertChannel{
		ID:   2932701,
		Name: "test@example.com",
		Type: "email",
		Configuration: AlertChannelConfiguration{
			IncludeJSONAttachment: "true",
			Recipients:            "test@example.com",
		},
		Links: AlertChannelLinks{
			PolicyIDs: []int{593436},
		},
	}

	actual, err := alerts.DeletePolicyChannel(593436, 2932701)

	assert.NoError(t, err)
	assert.Equal(t, expected, *actual)
}

func newTestPolicyChannelsClient(handler http.Handler) Alerts {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
	})

	return c
}

func newMocknewTestPolicyChannelsClientResponse(
	t *testing.T,
	mockJsonResponse string,
	statusCode int,
) Alerts {
	return newTestPolicyChannelsClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		_, err := w.Write([]byte(mockJsonResponse))

		require.NoError(t, err)
	}))
}
