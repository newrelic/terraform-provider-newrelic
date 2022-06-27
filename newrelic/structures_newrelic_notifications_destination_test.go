//go:build unit
// +build unit

package newrelic

import (
	"github.com/newrelic/newrelic-client-go/pkg/notifications"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandNotificationDestination(t *testing.T) {
	user := "test-user"
	auth := notifications.Auth{
		AuthType: &notifications.AuthTypes.Basic,
		User:     &user,
	}
	property := map[string]interface{}{
		"key":   "url",
		"value": "https://webhook.com",
	}

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Expanded     *notifications.Destination
	}{
		"missing auth": {
			Data: map[string]interface{}{
				"name": "testing123",
				"type": "WEBHOOK",
			},
			ExpectErr:    true,
			ExpectReason: "notification destination requires an auth attribute",
		},
		"invalid auth type": {
			Data: map[string]interface{}{
				"name": "testing123",
				"type": "WEBHOOK",
				"auth": map[string]string{
					"type": "INVALID",
					"user": "test-user",
				},
			},
			ExpectErr:    true,
			ExpectReason: "auth type must be token or basic",
		},
		"missing value in token type": {
			Data: map[string]interface{}{
				"name": "testing123",
				"type": "WEBHOOK",
				"auth": map[string]string{
					"type":   "TOKEN",
					"prefix": "test-user",
				},
			},
			ExpectErr:    true,
			ExpectReason: "token and prefix are required when using token auth type",
		},
		"valid destination": {
			Data: map[string]interface{}{
				"name":       "testing123",
				"type":       "WEBHOOK",
				"properties": []interface{}{property},
				"auth": map[string]string{
					"type":     "BASIC",
					"user":     "test-user",
					"password": "1234",
				},
			},
			Expanded: &notifications.Destination{
				Name: "testing123",
				Type: notifications.DestinationTypes.Webhook,
				Auth: auth,
				Properties: []notifications.Property{
					{
						Key:   "url",
						Value: "https://webhook.site/94193c01-4a81-4782-8f1b-554d5230395b",
					},
				},
			},
		},
	}

	r := resourceNewRelicNotificationDestination()

	for _, tc := range cases {
		d := r.TestResourceData()

		for k, v := range tc.Data {
			if err := d.Set(k, v); err != nil {
				t.Fatalf("err: %s", err)
			}
		}

		expanded, err := expandNotificationDestination(d)

		if tc.ExpectErr {
			assert.NotNil(t, err)
			assert.Equal(t, err.Error(), tc.ExpectReason)
		} else {
			assert.Nil(t, err)
		}

		if tc.Expanded != nil {
			assert.Equal(t, tc.Expanded.Name, expanded.Name)
		}
	}
}

func TestFlattenNotificationDestination(t *testing.T) {
	user := "user"
	auth := notifications.Auth{
		AuthType: &notifications.AuthTypes.Basic,
		User:     &user,
	}
	r := resourceNewRelicNotificationDestination()

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Flattened    *notifications.Destination
	}{
		"minimal": {
			Data: map[string]interface{}{
				"name": "testing123",
				"type": "WEBHOOK",
				"auth": map[string]interface{}{
					"type":     "BASIC",
					"user":     "user",
					"password": "1234",
				},
			},
			Flattened: &notifications.Destination{
				Name: "testing123",
				Type: "WEBHOOK",
				Auth: auth,
			},
		},
	}

	for _, tc := range cases {
		if tc.Flattened != nil {
			d := r.TestResourceData()
			err := flattenNotificationDestination(tc.Flattened, d)
			assert.NoError(t, err)

			for k, v := range tc.Data {
				var x interface{}
				var ok bool
				if x, ok = d.GetOk(k); !ok {
					t.Fatalf("err: %s", err)
				}

				if k == "auth" {
					testFlattenNotificationDestinationAuth(t, v, tc.Flattened.Auth)
				} else {
					assert.Equal(t, x, v)
				}
			}
		}
	}
}

func testFlattenNotificationDestinationAuth(t *testing.T, v interface{}, auth notifications.Auth) {
	for ck, cv := range v.(map[string]interface{}) {
		switch ck {
		case "type":
			assert.Equal(t, cv, string(*auth.AuthType))
		case "user":
			assert.Equal(t, cv, *auth.User)
		}
	}
}
