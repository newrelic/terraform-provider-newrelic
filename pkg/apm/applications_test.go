// +build unit

package apm

import (
	"fmt"
	"net/http"
	"testing"

	"time"

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

	testMetricNames = []*MetricName{
		{"GC/System/Pauses", []string{
			"as_percentage",
			"average_time",
			"calls_per_minute",
			"max_value",
			"total_call_time_per_minute",
			"utilization",
		}},
		{"Memory/Heap/Free", []string{
			"used_bytes_by_host",
			"used_mb_by_host",
			"total_used_mb",
		}},
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

	testMetricNamesJson = `{
		"metrics": [
			{
				"name": "GC/System/Pauses",
				"values": [
					"as_percentage",
					"average_time",
					"calls_per_minute",
					"max_value",
					"total_call_time_per_minute",
					"utilization"
				]
			},
			{
				"name": "Memory/Heap/Free",
				"values": [
					"used_bytes_by_host",
					"used_mb_by_host",
					"total_used_mb"
				]
			}
		]
	}`

	testMetricDataJson = `{
		"metric_data": {
			"from": "2020-01-27T23:25:45+00:00",
			"to": "2020-01-27T23:55:45+00:00",
			"metrics_not_found": [],
			"metrics_found": [
				"GC/System/Pauses"
			],
			"metrics": [
				{
					"name": "GC/System/Pauses",
					"timeslices": [
						{
							"from": "2020-01-27T23:22:00+00:00",
							"to": "2020-01-27T23:23:00+00:00",
							"values": {
								"as_percentage": 0.0298,
								"average_time": 0.298,
								"calls_per_minute": 65.9,
								"max_value": 0.0006,
								"total_call_time_per_minute": 0.0196,
								"utilization": 0.0327
							}
						},
						{
							"from": "2020-01-27T23:23:00+00:00",
							"to": "2020-01-27T23:24:00+00:00",
							"values": {
								"as_percentage": 0.0294,
								"average_time": 0.294,
								"calls_per_minute": 67,
								"max_value": 0.0005,
								"total_call_time_per_minute": 0.0197,
								"utilization": 0.0328
							}
						}
					]
				}
			]
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

func TestGetMetricNames(t *testing.T) {
	t.Parallel()
	apm := newMockResponse(t, testMetricNamesJson, http.StatusOK)

	actual, err := apm.GetMetricNames(testApplication.ID, MetricNamesParams{})
	expected := testMetricNames

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, len(expected), len(actual))

	if len(expected) == len(actual) {
		for i := range expected {
			assert.Equal(t, expected[i], actual[i])
		}
	}
}

func TestMetricData(t *testing.T) {
	t.Parallel()
	apm := newMockResponse(t, testMetricDataJson, http.StatusOK)

	actual, err := apm.GetMetricData(testApplication.ID, MetricDataParams{})
	expectedTimeSlices := []struct {
		From   string
		To     string
		Values MetricTimesliceValues
	}{
		{
			"2020-01-27T23:22:00+00:00",
			"2020-01-27T23:23:00+00:00",
			MetricTimesliceValues{
				0.0298,
				0.298,
				65.9,
				0.0006,
				0.0196,
				0.0327,
			},
		},
		{
			"2020-01-27T23:23:00+00:00",
			"2020-01-27T23:24:00+00:00",
			MetricTimesliceValues{
				0.0294,
				0.294,
				67,
				0.0005,
				0.0197,
				0.0328,
			},
		},
	}

	assert.NoError(t, err)
	assert.NotNil(t, actual)

	for i, e := range expectedTimeSlices {
		from, err := time.Parse(time.RFC3339, e.From)
		assert.NoError(t, err)

		to, err := time.Parse(time.RFC3339, e.To)
		assert.NoError(t, err)

		assert.Equal(t, &from, actual[0].Timeslices[i].From)
		assert.Equal(t, &to, actual[0].Timeslices[i].To)

		assert.Equal(t, e.Values.AsPercentage, actual[0].Timeslices[i].Values.AsPercentage)
		assert.Equal(t, e.Values.AverageTime, actual[0].Timeslices[i].Values.AverageTime)
		assert.Equal(t, e.Values.CallsPerMinute, actual[0].Timeslices[i].Values.CallsPerMinute)
		assert.Equal(t, e.Values.MaxValue, actual[0].Timeslices[i].Values.MaxValue)
		assert.Equal(t, e.Values.TotalCallTimePerMinute, actual[0].Timeslices[i].Values.TotalCallTimePerMinute)
		assert.Equal(t, e.Values.Utilization, actual[0].Timeslices[i].Values.Utilization)
	}
}
