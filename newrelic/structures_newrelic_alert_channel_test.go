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
					for ck, cv := range v.(map[string]interface{}) {

						if ck == "api_key" {
							assert.Equal(t, cv, tc.Flattened.Configuration.APIKey)
						}

						if ck == "auth_password" {
							assert.Equal(t, cv, tc.Flattened.Configuration.AuthPassword)
						}

						if ck == "auth_username" {
							assert.Equal(t, cv, tc.Flattened.Configuration.AuthUsername)
						}

						if ck == "base_url" {
							assert.Equal(t, cv, tc.Flattened.Configuration.BaseURL)
						}

						if ck == "channel" {
							assert.Equal(t, cv, tc.Flattened.Configuration.Channel)
						}

						if ck == "key" {
							assert.Equal(t, cv, tc.Flattened.Configuration.Key)
						}

						if ck == "payload" {
							assert.Equal(t, cv, tc.Flattened.Configuration.Payload)
						}

						if ck == "payload_type" {
							assert.Equal(t, cv, tc.Flattened.Configuration.PayloadType)
						}

						if ck == "recipients" {
							assert.Equal(t, cv, tc.Flattened.Configuration.Recipients)
						}

						if ck == "region" {
							assert.Equal(t, cv, tc.Flattened.Configuration.Region)
						}

						if ck == "route_key" {
							assert.Equal(t, cv, tc.Flattened.Configuration.RouteKey)
						}

						if ck == "service_key" {
							assert.Equal(t, cv, tc.Flattened.Configuration.ServiceKey)
						}

						if ck == "tags" {
							assert.Equal(t, cv, tc.Flattened.Configuration.Tags)
						}

						if ck == "teams" {
							assert.Equal(t, cv, tc.Flattened.Configuration.Teams)
						}

						if ck == "url" {
							assert.Equal(t, cv, tc.Flattened.Configuration.URL)
						}

						if ck == "user_id" {
							assert.Equal(t, cv, tc.Flattened.Configuration.UserID)
						}

					}
				} else {
					assert.Equal(t, x, v)
				}
			}
		}
	}
}
