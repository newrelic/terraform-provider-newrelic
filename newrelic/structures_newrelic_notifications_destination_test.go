//go:build unit

package newrelic

import (
	"testing"

	"github.com/newrelic/newrelic-client-go/v2/pkg/ai"
	"github.com/newrelic/newrelic-client-go/v2/pkg/notifications"

	"github.com/stretchr/testify/assert"
)

func TestExpandNotificationDestination(t *testing.T) {
	user := "test-user"
	basicAuthType := ai.AiNotificationsAuth{
		AuthType: "BASIC",
		User:     user,
	}
	basicAuth := map[string]interface{}{
		"user":     user,
		"password": "123456",
	}
	customHeadersAuthType := ai.AiNotificationsAuth{
		AuthType: "CUSTOM_HEADERS",
		CustomHeaders: []ai.AiNotificationsCustomHeaders{
			{
				Key: "testKey1",
			},
		},
	}
	customHeadersAuth := map[string]interface{}{
		"key":   "testKey1",
		"value": "testValue1",
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
				Auth: basicAuthType,
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
		"valid secureURL webhook destination": {
			Data: map[string]interface{}{
				"name":     "webhook-test",
				"type":     "WEBHOOK",
				"property": []interface{}{webhookProperty},
			},
			Expanded: &notifications.AiNotificationsDestination{
				Name:       "webhook-test",
				Type:       notifications.AiNotificationsDestinationTypeTypes.WEBHOOK,
				Properties: []notifications.AiNotificationsProperty{},
				SecureURL: notifications.AiNotificationsSecureURL{
					Prefix: "https://webhook.com",
				},
			},
		},
		"valid webhook destination with custom headers auth": {
			Data: map[string]interface{}{
				"name":               "webhook-test",
				"type":               "WEBHOOK",
				"property":           []interface{}{webhookProperty},
				"auth_custom_header": []interface{}{customHeadersAuth},
			},
			Expanded: &notifications.AiNotificationsDestination{
				Name: "webhook-test",
				Type: notifications.AiNotificationsDestinationTypeTypes.WEBHOOK,
				Properties: []notifications.AiNotificationsProperty{
					{
						Key:   "email",
						Value: "example@email.com",
					},
				},
				Auth: customHeadersAuthType,
			},
		},
		"valid webhook destination with organization scope": {
			Data: map[string]interface{}{
				"name":     "webhook-test",
				"type":     "WEBHOOK",
				"property": []interface{}{webhookProperty},
				"scope": []interface{}{
					map[string]interface{}{
						"type": "ORGANIZATION",
						"id":   "mock-organization-id",
					},
				},
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
		"valid webhook destination with account scope": {
			Data: map[string]interface{}{
				"name":     "webhook-test",
				"type":     "WEBHOOK",
				"property": []interface{}{webhookProperty},
				"scope": []interface{}{
					map[string]interface{}{
						"type": "ACCOUNT",
						"id":   "12345678",
					},
				},
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
	guid := notifications.EntityGUID("testdestinationentityguid")

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
				"guid": "testdestinationentityguid",
			},
			Flattened: &notifications.AiNotificationsDestination{
				Name: "testing123",
				Type: "WEBHOOK",
				GUID: guid,
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

func TestFlattenNotificationDestinationDataSource(t *testing.T) {
	r := dataSourceNewRelicNotificationDestination()
	guid := notifications.EntityGUID("testdestinationentityguid")

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
				"guid": "testdestinationentityguid",
			},
			Flattened: &notifications.AiNotificationsDestination{
				Name: "testing123",
				Type: "WEBHOOK",
				GUID: guid,
			},
		},
	}

	for _, tc := range cases {
		if tc.Flattened != nil {
			d := r.TestResourceData()
			err := flattenNotificationDestinationDataSource(tc.Flattened, notifications.AiNotificationsEntityScopeInput{}, d)
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

func TestFlattenNotificationDestinationWithScope_OrganizationScope(t *testing.T) {
	r := resourceNewRelicNotificationDestination()
	guid := notifications.EntityGUID("testguid")

	destination := &notifications.AiNotificationsDestination{
		Name: "org-dest",
		Type: "WEBHOOK",
		GUID: guid,
		Scope: notifications.AiNotificationsEntityScope{
			Type: notifications.AiNotificationsEntityScopeTypeTypes.ORGANIZATION,
			ID:   "org-uuid-123",
		},
		Active: true,
	}

	d := r.TestResourceData()
	err := flattenNotificationDestination(destination, d)
	assert.NoError(t, err)

	assert.Equal(t, "org-dest", d.Get("name"))
	assert.Equal(t, string("WEBHOOK"), d.Get("type"))

	scopeList := d.Get("scope").([]interface{})
	assert.Len(t, scopeList, 1)
	scopeMap := scopeList[0].(map[string]interface{})
	assert.Equal(t, "ORGANIZATION", scopeMap["type"])
	assert.Equal(t, "org-uuid-123", scopeMap["id"])
}

func TestFlattenNotificationDestinationWithScope_AccountScope(t *testing.T) {
	r := resourceNewRelicNotificationDestination()
	guid := notifications.EntityGUID("testguid")

	destination := &notifications.AiNotificationsDestination{
		Name:      "acct-dest",
		Type:      "WEBHOOK",
		GUID:      guid,
		AccountID: 12345,
		Scope: notifications.AiNotificationsEntityScope{
			Type: notifications.AiNotificationsEntityScopeTypeTypes.ACCOUNT,
			ID:   "12345",
		},
		Active: true,
	}

	d := r.TestResourceData()
	err := flattenNotificationDestination(destination, d)
	assert.NoError(t, err)

	assert.Equal(t, "acct-dest", d.Get("name"))

	scopeList := d.Get("scope").([]interface{})
	assert.Len(t, scopeList, 1)
	scopeMap := scopeList[0].(map[string]interface{})
	assert.Equal(t, "ACCOUNT", scopeMap["type"])
	assert.Equal(t, "12345", scopeMap["id"])

	assert.Equal(t, 12345, d.Get("account_id"))
}

func TestFlattenNotificationDestination_SetsAccountScope(t *testing.T) {
	r := resourceNewRelicNotificationDestination()
	guid := notifications.EntityGUID("testguid")

	destination := &notifications.AiNotificationsDestination{
		Name:      "acct-dest",
		Type:      "WEBHOOK",
		GUID:      guid,
		AccountID: 99999,
		Scope: notifications.AiNotificationsEntityScope{
			Type: notifications.AiNotificationsEntityScopeTypeTypes.ACCOUNT,
			ID:   "99999",
		},
		Active: true,
	}

	d := r.TestResourceData()
	err := flattenNotificationDestination(destination, d)
	assert.NoError(t, err)

	scopeList := d.Get("scope").([]interface{})
	assert.Len(t, scopeList, 1)
	scopeMap := scopeList[0].(map[string]interface{})
	assert.Equal(t, "ACCOUNT", scopeMap["type"])
	assert.Equal(t, "99999", scopeMap["id"])
}

func TestFlattenNotificationDestinationDataSource_OrgScope(t *testing.T) {
	r := dataSourceNewRelicNotificationDestination()
	guid := notifications.EntityGUID("testguid")

	destination := &notifications.AiNotificationsDestination{
		ID:        "dest-id-123",
		Name:      "org-dest",
		Type:      "WEBHOOK",
		GUID:      guid,
		AccountID: 12345,
		Scope: notifications.AiNotificationsEntityScope{
			Type: notifications.AiNotificationsEntityScopeTypeTypes.ORGANIZATION,
			ID:   "org-uuid-456",
		},
		Active: true,
	}

	scope := notifications.AiNotificationsEntityScopeInput{
		Type: notifications.AiNotificationsEntityScopeTypeInputTypes.ORGANIZATION,
		ID:   "org-uuid-456",
	}

	d := r.TestResourceData()
	err := flattenNotificationDestinationDataSource(destination, scope, d)
	assert.NoError(t, err)

	assert.Equal(t, "org-dest", d.Get("name"))

	scopeList := d.Get("scope").([]interface{})
	assert.Len(t, scopeList, 1)
	scopeMap := scopeList[0].(map[string]interface{})
	assert.Equal(t, "ORGANIZATION", scopeMap["type"])
	assert.Equal(t, "org-uuid-456", scopeMap["id"])
}

func TestFlattenNotificationDestinationDataSource_AccountScope(t *testing.T) {
	r := dataSourceNewRelicNotificationDestination()
	guid := notifications.EntityGUID("testguid")

	destination := &notifications.AiNotificationsDestination{
		ID:        "dest-id-789",
		Name:      "acct-dest",
		Type:      "WEBHOOK",
		GUID:      guid,
		AccountID: 54321,
		Active:    true,
	}

	scope := notifications.AiNotificationsEntityScopeInput{
		Type: notifications.AiNotificationsEntityScopeTypeInputTypes.ACCOUNT,
		ID:   "54321",
	}

	d := r.TestResourceData()
	err := flattenNotificationDestinationDataSource(destination, scope, d)
	assert.NoError(t, err)

	scopeList := d.Get("scope").([]interface{})
	assert.Len(t, scopeList, 1)
	scopeMap := scopeList[0].(map[string]interface{})
	assert.Equal(t, "ACCOUNT", scopeMap["type"])
	assert.Equal(t, "54321", scopeMap["id"])
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
