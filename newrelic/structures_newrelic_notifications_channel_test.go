//go:build unit
// +build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/notifications"

	"github.com/stretchr/testify/assert"
)

func TestExpandNotificationChannel(t *testing.T) {
	property := map[string]interface{}{
		"key":   "payload",
		"value": "{\\n\\t\\\"id\\\": \\\"test\\\"\\n}",
		"label": "Payload Template",
	}

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Expanded     *notifications.AiNotificationsChannel
	}{
		"valid webhook channel": {
			Data: map[string]interface{}{
				"name":           "webhook-test",
				"type":           "WEBHOOK",
				"property":       []interface{}{property},
				"product":        "IINT",
				"destination_id": "b1e90a32-23b7-4028-b2c7-ffbdfe103852",
			},
			Expanded: &notifications.AiNotificationsChannel{
				Name: "webhook-test",
				Type: notifications.AiNotificationsChannelTypeTypes.WEBHOOK,
				Properties: []notifications.AiNotificationsProperty{
					{
						Key:   "payload",
						Value: "{\\n\\t\\\"id\\\": \\\"test\\\"\\n}",
						Label: "Payload Template",
					},
				},
			},
		},
		"valid email channel": {
			Data: map[string]interface{}{
				"name":           "email-test",
				"type":           "EMAIL",
				"product":        "IINT",
				"destination_id": "b1e90a32-23b7-4028-b2c7-ffbdfe103852",
			},
			Expanded: &notifications.AiNotificationsChannel{
				Name: "email-test",
				Type: notifications.AiNotificationsChannelTypeTypes.EMAIL,
				Properties: []notifications.AiNotificationsProperty{
					{
						Key:   "subject",
						Value: "{{ issueTitle }}",
					},
				},
			},
		},
	}

	r := resourceNewRelicNotificationChannel()

	for _, tc := range cases {
		d := r.TestResourceData()

		for k, v := range tc.Data {
			if err := d.Set(k, v); err != nil {
				t.Fatalf("err: %s", err)
			}
		}

		expanded := expandNotificationChannel(d)

		if tc.Expanded != nil {
			assert.Equal(t, tc.Expanded.Name, expanded.Name)
		}
	}
}

func TestFlattenNotificationChannel(t *testing.T) {
	r := resourceNewRelicNotificationChannel()

	cases := map[string]struct {
		Data         map[string]interface{}
		ExpectErr    bool
		ExpectReason string
		Flattened    *notifications.AiNotificationsChannel
	}{
		"minimal": {
			Data: map[string]interface{}{
				"name": "testing123",
				"type": "WEBHOOK",
			},
			Flattened: &notifications.AiNotificationsChannel{
				Name:          "testing123",
				Type:          "WEBHOOK",
				Product:       "IINT",
				DestinationId: "b1e90a32-23b7-4028-b2c7-ffbdfe103852",
			},
		},
	}

	for _, tc := range cases {
		if tc.Flattened != nil {
			d := r.TestResourceData()
			err := flattenNotificationChannel(tc.Flattened, d)
			assert.NoError(t, err)

			for k, v := range tc.Data {
				var x interface{}
				var ok bool
				if x, ok = d.GetOk(k); !ok {
					t.Fatalf("err: %s", err)
				}

				assert.Equal(t, x, v)
			}
		}
	}
}
