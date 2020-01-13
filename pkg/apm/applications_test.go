// +build unit

package apm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testApplicationSummary = ApplicationSummary{
		ResponseTime:            5.91,
		Throughput:              1,
		ErrorRate:               0,
		ApdexTarget:             0.5,
		ApdexScore:              1,
		HostCount:               1,
		InstanceCount:           15,
		ConcurrentInstanceCount: 1,
	}

	testApplicationEndUserSummary = ApplicationEndUserSummary{
		ResponseTime: 3.8,
		Throughput:   1660,
		ApdexTarget:  2.5,
		ApdexScore:   0.78,
	}

	testApplicationSettings = ApplicationSettings{
		AppApdexThreshold:        0.5,
		EndUserApdexThreshold:    7,
		EnableRealUserMonitoring: true,
		UseServerSideConfig:      false,
	}

	testApplicationLinks = ApplicationLinks{
		ServerIDs:     []int{},
		HostIDs:       []int{204260579},
		InstanceIDs:   []int{204261411},
		AlertPolicyID: 1234,
	}

	testApplication = Application{
		ID:             204261410,
		Name:           "Billing Service",
		Language:       "python",
		HealthStatus:   "unknown",
		Reporting:      true,
		LastReportedAt: "2019-12-11T19:09:10+00:00",
		Summary:        testApplicationSummary,
		EndUserSummary: testApplicationEndUserSummary,
		Settings:       testApplicationSettings,
		Links:          testApplicationLinks,
	}

	testApplicationJson = `{
		"id": 204261410,
		"name": "Billing Service",
		"language": "python",
		"health_status": "unknown",
		"reporting": true,
		"last_reported_at": "2019-12-11T19:09:10+00:00",
		"application_summary": {
			"response_time": 5.91,
			"throughput": 1,
			"error_rate": 0,
			"apdex_target": 0.5,
			"apdex_score": 1,
			"host_count": 1,
			"instance_count": 15,
			"concurrent_instance_count": 1
		},
		"end_user_summary": {
			"response_time": 3.8,
			"throughput": 1660,
			"apdex_target": 2.5,
			"apdex_score": 0.78
		},
		"settings": {
			"app_apdex_threshold": 0.5,
			"end_user_apdex_threshold": 7,
			"enable_real_user_monitoring": true,
			"use_server_side_config": false
		},
		"links": {
			"application_instances": [
				204261411
			],
			"servers": [],
			"application_hosts": [
				204260579
			],
			"alert_policy": 1234
		}
	}`
)

func TestListApplications(t *testing.T) {
	t.Parallel()
	responseJSON := fmt.Sprintf(`{ "applications": [%s] }`, testApplicationJson)
	apm := newMockResponse(t, responseJSON, http.StatusOK)

	actual, err := apm.ListApplications(nil)

	expected := []*Application{&testApplication}

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, expected, actual)
}

func TestListApplicationsWithParams(t *testing.T) {
	t.Parallel()
	expectedName := "appName"
	expectedHost := "appHost"
	expectedLanguage := "appLanguage"
	expectedIDs := "123,456"

	apm := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		name := values.Get("filter[name]")
		host := values.Get("filter[host]")
		ids := values.Get("filter[ids]")
		language := values.Get("filter[language]")

		assert.Equal(t, expectedName, name)
		assert.Equal(t, expectedHost, host)
		assert.Equal(t, expectedIDs, ids)
		assert.Equal(t, expectedLanguage, language)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"applications":[]}`))

		assert.NoError(t, err)
	}))

	params := ListApplicationsParams{
		Name:     expectedName,
		Host:     expectedHost,
		IDs:      []int{123, 456},
		Language: expectedLanguage,
	}

	_, err := apm.ListApplications(&params)

	assert.NoError(t, err)
}

func TestGetApplication(t *testing.T) {
	t.Parallel()
	responseJSON := fmt.Sprintf(`{ "application": %s}`, testApplicationJson)
	apm := newMockResponse(t, responseJSON, http.StatusOK)

	actual, err := apm.GetApplication(testApplication.ID)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, &testApplication, actual)
}

func TestUpdateApplication(t *testing.T) {
	t.Parallel()
	responseJSON := fmt.Sprintf(`{ "application": %s}`, testApplicationJson)
	apm := newMockResponse(t, responseJSON, http.StatusOK)

	params := UpdateApplicationParams{
		Name:     testApplication.Name,
		Settings: testApplication.Settings,
	}

	actual, err := apm.UpdateApplication(testApplication.ID, params)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, &testApplication, actual)
}

func TestDeleteApplication(t *testing.T) {
	t.Parallel()
	responseJSON := fmt.Sprintf(`{ "application": %s}`, testApplicationJson)
	apm := newMockResponse(t, responseJSON, http.StatusOK)

	actual, err := apm.DeleteApplication(testApplication.ID)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, &testApplication, actual)
}
