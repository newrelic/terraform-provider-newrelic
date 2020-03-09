# Examples

## [Golden Signal Alerts](modules/golden-signal-alerts)
This module encapsulates an alerting strategy based on the [Four Golden Signals](https://landing.google.com/sre/sre-book/chapters/monitoring-distributed-systems/#xref_monitoring_golden-signals) introduced in Google’s widely read book on [Site Reliability Engineering](https://landing.google.com/sre/sre-book/toc/index.html).

The signals chosen for this module are:

* *Latency*: High response time (seconds)
* *Traffic*: Low throughput (requests/minute)
* *Errors*: Error rate (errors/minute)
* *Saturation*: CPU utilization (percentage utilized)

### Requirements
Applications making use of this module need to be reporting data into both APM and Infrastructure.

### Input variables
The following input variables are accepted by the module:

* `name`: The APM application name as reported to New Relic
* `duration`: The duration to evaluate the alert conditions over, in minutes
* `cpu_threshold`: The critical threshold of the CPU utilization condition, as a percentage
* `error_percentage_threshold`: The critical threshold of the error rate condition, in errors/min
* `response_time_threshold`: The critical threshold of the response time condition, in seconds
* `throughput_threshold`: The critical threshold of the throughput condition, in requests/min

### Outputs
The following output values are provided by the module:

* `policy_id`: The ID of the created alert policy
* `cpu_condition_id`: The ID of the created high CPU alert condition
* `error_percentage_condition_id`: The ID of the created error percentage alert condition
* `response_time_condition_id`: The ID of the created response time alert condition
* `throughput_condition_id`: The ID of the created throughput alert condition


### Example usage
```terraform
data "newrelic_alert_channel" "alert_channel" {
	name = "Page Developer Toolkit Team"
}

module "webportal_alerts" {
	source = "./modules/golden-signal-alerts"
	alert_channel_ids = [data.newrelic_alert_channel.alert_channel.id]

	service = {
		name                       = "WebPortal"
		duration                   = 5
		cpu_threshold              = 90
		response_time_threshold    = 5
		error_percentage_threshold = 5
		throughput_threshold       = 5
	}
}
```