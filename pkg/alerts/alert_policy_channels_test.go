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
	testUpdateAlertPolicyChannelsResponseJSON = `{
		"policy": {
			"id": 593436,
			"channel_ids": [
				2932701,
				2932702
			]
		}
	}`

	testDeleteAlertPolicyChannelResponseJSON = `{
		"channel": {
			"id": 2932701,
			"name": "devtoolkit@newrelic.com",
			"type": "email",
			"configuration": {
				"include_json_attachment": "true",
				"recipients": "devtoolkit@newrelic.com"
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

func TestUpdateAlertPolicyChannels(t *testing.T) {
	t.Parallel()
	alerts := newMocknewTestAlertPolicyChannelsClientResponse(t, testUpdateAlertPolicyChannelsResponseJSON, http.StatusOK)

	actual, err := alerts.UpdateAlertPolicyChannels(593436, []int{2932701, 2932702})

	expected := AlertPolicyChannels{
		ID:         593436,
		ChannelIDs: []int{2932701, 2932702},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, *actual)
}

func TestDeleteAlertPolicyChannel(t *testing.T) {
	t.Parallel()
	alerts := newMocknewTestAlertPolicyChannelsClientResponse(t, testDeleteAlertPolicyChannelResponseJSON, http.StatusOK)

	expected := AlertChannel{
		ID:   2932701,
		Name: "devtoolkit@newrelic.com",
		Type: "email",
		Configuration: AlertChannelConfiguration{
			IncludeJSONAttachment: "true",
			Recipients:            "devtoolkit@newrelic.com",
		},
		Links: AlertChannelLinks{
			PolicyIDs: []int{593436},
		},
	}

	actual, err := alerts.DeleteAlertPolicyChannel(593436, 2932701)

	assert.NoError(t, err)
	assert.Equal(t, expected, *actual)
}

func newTestAlertPolicyChannelsClient(handler http.Handler) Alerts {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}

func newMocknewTestAlertPolicyChannelsClientResponse(
	t *testing.T,
	mockJsonResponse string,
	statusCode int,
) Alerts {
	return newTestAlertPolicyChannelsClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		_, err := w.Write([]byte(mockJsonResponse))

		require.NoError(t, err)
	}))
}
