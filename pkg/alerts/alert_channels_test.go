// +build unit

package alerts

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/newrelic/newrelic-client-go/pkg/config"
)

func TestGetAlertChannel(t *testing.T) {
	t.Parallel()
	alerts := newMockServerClientResponse(t, `
		{
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
		}
	`)

	expected := AlertChannel{
		ID:   2803426,
		Name: "unit-test-alert-channel",
		Type: "user",
		Configuration: &AlertChannelConfiguration{
			UserID: "2680539",
		},
		Links: AlertChannelLinks{
			PolicyIDs: []int{},
		},
	}

	actual, err := alerts.GetAlertChannel(2803426)

	if err != nil {
		t.Fatalf("GetAlertChannel error: %s", err)
	}

	if actual == nil {
		t.Fatalf("GetAlertChannel result is nil")
	}

	if diff := cmp.Diff(expected, *actual); diff != "" {
		t.Fatalf("GetAlertChannel result differs from expected: %s", diff)
	}
}

func TestListAlertChannels(t *testing.T) {
	t.Parallel()
	alerts := newMockServerClientResponse(t, `
		{
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
		}
	`)

	expected := []AlertChannel{
		{
			ID:   2803426,
			Name: "unit-test-alert-channel",
			Type: "user",
			Configuration: &AlertChannelConfiguration{
				UserID: "2680539",
			},
			Links: AlertChannelLinks{
				PolicyIDs: []int{},
			},
		},
		{
			ID:   2932511,
			Name: "test@testing.com",
			Type: "email",
			Configuration: &AlertChannelConfiguration{
				Recipients:            "test@testing.com",
				IncludeJSONAttachment: "true",
			},
			Links: AlertChannelLinks{
				PolicyIDs: []int{},
			},
		},
	}

	actual, err := alerts.ListAlertChannels()

	if err != nil {
		t.Fatalf("ListAlertChannels error: %s", err)
	}

	if actual == nil {
		t.Fatalf("ListAlertChannels result is nil")
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("ListAlertChannels result differs from expected: %s", diff)
	}
}

func TestCreateAlertChannel(t *testing.T) {
	t.Parallel()
	alerts := newMockServerClientResponse(t, `
		{
			"channels": [
				{
					"id": 2932701,
					"name": "sblue@newrelic.com",
					"type": "email",
					"configuration": {
						"include_json_attachment": "true",
						"recipients": "sblue@newrelic.com"
					},
					"links": {
						"policy_ids": []
					}
				}
			],
			"links": {
				"channel.policy_ids": "/v2/policies/{policy_id}"
			}
		}
	`)

	channel := AlertChannel{
		Name: "sblue@newrelic.com",
		Type: "email",
		Configuration: &AlertChannelConfiguration{
			Recipients:            "sblue@newrelic.com",
			IncludeJSONAttachment: "true",
		},
	}

	expected := []AlertChannel{
		{
			ID:   2932701,
			Name: "sblue@newrelic.com",
			Type: "email",
			Configuration: &AlertChannelConfiguration{
				Recipients:            "sblue@newrelic.com",
				IncludeJSONAttachment: "true",
			},
			Links: AlertChannelLinks{
				PolicyIDs: []int{},
			},
		},
	}

	actual, err := alerts.CreateAlertChannel(channel)

	if err != nil {
		t.Fatalf("CreateAlertChannel error: %s", err)
	}

	if actual == nil {
		t.Fatalf("CreateAlertChannel result is nil")
	}

	if diff := cmp.Diff(expected, *actual); diff != "" {
		t.Fatalf("CreateAlertChannel result differs from expected: %s", diff)
	}
}

func TestDeleteAlertChannel(t *testing.T) {
	t.Parallel()
	alerts := newMockServerClientResponse(t, `
		{
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
		}
	`)

	err := alerts.DeleteAlertChannel(2932511)

	if err != nil {
		t.Fatalf("UpdateAlertPolicy error: %s", err)
	}
}

func newTestClient(handler http.Handler) Alerts {
	ts := httptest.NewServer(handler)

	c := New(config.Config{
		APIKey:    "abc123",
		BaseURL:   ts.URL,
		Debug:     false,
		UserAgent: "newrelic/newrelic-client-go",
	})

	return c
}

func newMockServerClientResponse(t *testing.T, mockJsonResponse string) Alerts {
	return newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(mockJsonResponse))

		if err != nil {
			t.Fatal(err)
		}
	}))
}
