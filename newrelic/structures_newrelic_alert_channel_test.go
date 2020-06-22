package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandAlertChannel(t *testing.T) {
	config := map[string]interface{}{
		"url": "https://example.com",
	}

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Expanded     *alerts.Channel
	}{
		"invalid": {
			Data: map[string]interface{}{
				"name": "testing123",
			},
			ExpectErr:    true,
			ExpectReason: "alert channel requires a config or configuration attribute",
		},
		"valid slack": {
			Data: map[string]interface{}{
				"name":   "testing123",
				"type":   "slack",
				"config": []interface{}{config},
			},
			Expanded: &alerts.Channel{
				Name: "testing123",
				Configuration: alerts.ChannelConfiguration{
					URL: "https://example.com",
				},
			},
		},
		"nil config": {
			Data: map[string]interface{}{
				"name":   "testing123",
				"type":   "slack",
				"config": []interface{}{nil},
			},
			Expanded: &alerts.Channel{
				Name: "testing123",
				Configuration: alerts.ChannelConfiguration{
					URL: "https://example.com",
				},
			},
		},
	}

	r := resourceNewRelicAlertChannel()

	for _, tc := range cases {
		d := r.TestResourceData()

		for k, v := range tc.Data {
			if err := d.Set(k, v); err != nil {
				t.Fatalf("err: %s", err)
			}
		}

		expanded, err := expandAlertChannel(d)

		if tc.ExpectErr {
			require.NotNil(t, err)
			require.Equal(t, err.Error(), tc.ExpectReason)
		} else {
			assert.Nil(t, err)
		}

		if tc.Expanded != nil {
			assert.Equal(t, tc.Expanded.Name, expanded.Name)
		}
	}

}

func TestFlattenAlertChannel(t *testing.T) {
	r := resourceNewRelicAlertChannel()

	config := map[string]interface{}{
		"url": "https://example.com",
	}

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Flattened    *alerts.Channel
	}{
		"minimal": {
			Data: map[string]interface{}{
				"name": "testing123",
			},
			Flattened: &alerts.Channel{
				Name: "testing123",
			},
		},
		"less minimal": {
			Data: map[string]interface{}{
				"name":   "testing123",
				"type":   "slack",
				"config": config,
			},
			Flattened: &alerts.Channel{
				Name: "testing123",
				Type: "slack",
				Configuration: alerts.ChannelConfiguration{
					URL: "https://example.com",
				},
			},
		},
	}

	for _, tc := range cases {
		if tc.Flattened != nil {
			d := r.TestResourceData()
			err := flattenAlertChannel(tc.Flattened, d)
			assert.NoError(t, err)

			for k, v := range tc.Data {
				var x interface{}
				var ok bool
				if x, ok = d.GetOk(k); !ok {
					t.Fatalf("err: %s", err)
				}

				if k == "config" {
					testFlattenAlertChannelConfig(t, v, tc.Flattened.Configuration)

				} else {
					assert.Equal(t, x, v)
				}
			}
		}
	}
}

func testFlattenAlertChannelConfig(t *testing.T, v interface{}, config alerts.ChannelConfiguration) {
	for ck, cv := range v.(map[string]interface{}) {
		switch ck {
		case "api_key":
			assert.Equal(t, cv, config.APIKey)
		case "auth_password":
			assert.Equal(t, cv, config.AuthPassword)

		case "auth_username":
			assert.Equal(t, cv, config.AuthUsername)

		case "base_url":
			assert.Equal(t, cv, config.BaseURL)

		case "channel":
			assert.Equal(t, cv, config.Channel)

		case "key":
			assert.Equal(t, cv, config.Key)

		case "payload":
			assert.Equal(t, cv, config.Payload)

		case "payload_type":
			assert.Equal(t, cv, config.PayloadType)

		case "recipients":
			assert.Equal(t, cv, config.Recipients)

		case "region":
			assert.Equal(t, cv, config.Region)

		case "route_key":
			assert.Equal(t, cv, config.RouteKey)

		case "service_key":
			assert.Equal(t, cv, config.ServiceKey)

		case "tags":
			assert.Equal(t, cv, config.Tags)

		case "teams":
			assert.Equal(t, cv, config.Teams)

		case "url":
			assert.Equal(t, cv, config.URL)

		case "user_id":
			assert.Equal(t, cv, config.UserID)
		}

	}
}
