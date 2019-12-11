package apm

import (
	"http"
	"testing"

	client "github.com/newrelic/newrelic-client-go/internal"
)

func TestListApplications(t *testing.T) {
	tc := client.NewTestAPIClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
		{
			"applications": [
				{
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
						"instance_count": 15
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
						]
					}
				},
				{
					"id": 204261310,
					"name": "Fulfillment Service",
					"language": "ruby",
					"health_status": "unknown",
					"reporting": true,
					"last_reported_at": "2019-12-11T19:08:40+00:00",
					"application_summary": {
						"response_time": 400,
						"throughput": 0.667,
						"error_rate": 0,
						"apdex_target": 0.5,
						"apdex_score": 0.75,
						"host_count": 1,
						"instance_count": 1
					},
					"settings": {
						"app_apdex_threshold": 0.5,
						"end_user_apdex_threshold": 7,
						"enable_real_user_monitoring": true,
						"use_server_side_config": false
					},
					"links": {
						"application_instances": [
							204577198
						],
						"servers": [],
						"application_hosts": [
							204260579
						]
					}
				}
			]
		}
		`
	))

	applicationSummary := client.ApplicationSummary{
		ResponseTime: 5.91,
		Throughput: 1,
		ErrorRate: 0,
		ApdexTarget: 0.5,
		ApdexScore: 1,
		HostCount: 1,
		InstanceCount: ,
		ConcurrentInstanceCount: 15,
	}

	applicationEndUserSummary := client.ApplicationEndUserSummary{
		ResponseTime: ,
		Throughput: ,
		ApdexTarget: ,
		ApdexScore: ,
	}

	applicationSettings := client.ApplicationSettings{
		AppApdexThreshold: 0.5,
		EndUserApdexThreshold: 7,
		EnableRealUserMonitoring: true,
		UseServerSideConfig: false,
	}

	applicationLinks := client.ApplicationLinks{
		ServerIDs: []int{},
		HostIDs: []int{204260579},
		InstanceIDs: []int{204261411},
		AlertPolicyID: false,
	}

	expected := []client.Application{
		{
			ID: 						204261410,
			Name:           "Billing Service",
			Language:       "python",
			HealthStatus:   "unknown",
			Reporting:      true,
			LastReportedAt: "2019-12-11T19:09:10+00:00",
			Summary:        applicationSummary,
			EndUserSummary: applicationEndUserSummary,
			Settings:       applicationSettings,
			Links:          applicationLinks,
		}
	}
}
