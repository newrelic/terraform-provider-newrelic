//go:build unit
// +build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
	"github.com/newrelic/newrelic-client-go/v2/pkg/notifications"

	"github.com/stretchr/testify/assert"
)

func TestExpandNotificationDestination(t *testing.T) {
	user := "test-user"
	auth := ai.AiNotificationsAuth{
		AuthType: "BASIC",
		User:     user,
	}
	basicAuth := map[string]interface{}{
		"user":     user,
		"password": "123456",
	}
	webhookProperty := map[string]interface{}{
		"key":   "url",
		"value": "https://webhook.com",
	}
	emailProperty := map[string]interface{}{
		"key":   "email",
		"value": "example@email.com",
	}

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Expanded     *notifications.AiNotificationsDestination
	}{
		"valid webhook destination  (no auth)": {
			Data: map[string]interface{}{
				"name":     "webhook-test",
				"type":     "WEBHOOK",
				"property": []interface{}{webhookProperty},
			},
			Expanded: &notifications.AiNotificationsDestination{
				Name: "webhook-test",
				Type: notifications.AiNotificationsDestinationTypeTypes.WEBHOOK,
				Properties: []notifications.AiNotificationsProperty{
					{
						Key:   "url",
						Value: "https://webhook.com",
					},
				},
			},
		},
		"valid pager duty destination (no properties)": {
			Data: map[string]interface{}{
				"name":       "pd-service-test",
				"type":       "PAGERDUTY_SERVICE_INTEGRATION",
				"auth_basic": []interface{}{basicAuth},
			},
			Expanded: &notifications.AiNotificationsDestination{
				Name: "pd-service-test",
				Type: notifications.AiNotificationsDestinationTypeTypes.PAGERDUTY_SERVICE_INTEGRATION,
				Auth: auth,
			},
		},
		"valid email destination (no auth)": {
			Data: map[string]interface{}{
				"name":     "email-test",
				"type":     "EMAIL",
				"property": []interface{}{emailProperty},
			},
			Expanded: &notifications.AiNotificationsDestination{
				Name: "email-test",
				Type: notifications.AiNotificationsDestinationTypeTypes.EMAIL,
				Properties: []notifications.AiNotificationsProperty{
					{
						Key:   "email",
						Value: "example@email.com",
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
	r := resourceNewRelicNotificationDestination()

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Flattened    *notifications.AiNotificationsDestination
	}{
		"minimal": {
			Data: map[string]interface{}{
				"name": "testing123",
				"type": "WEBHOOK",
			},
			Flattened: &notifications.AiNotificationsDestination{
				Name: "testing123",
				Type: "WEBHOOK",
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

func testFlattenNotificationDestinationAuth(t *testing.T, v interface{}, auth ai.AiNotificationsAuth) {
	for ck, cv := range v.(map[string]interface{}) {
		switch ck {
		case "type":
			assert.Equal(t, cv, string(auth.AuthType))
		case "user":
			assert.Equal(t, cv, auth.User)
		}
	}
}
