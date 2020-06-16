data "newrelic_entity" "application" {
	name = var.service.name
	type = "APPLICATION"
	domain = "APM"
}

resource "newrelic_alert_policy" "golden_signal_policy" {
	name = "Golden Signals - ${var.service.name}"
}

resource "newrelic_alert_condition" "response_time_web" {
	policy_id = newrelic_alert_policy.golden_signal_policy.id

	name            = "High Response Time (web)"
	type            = "apm_app_metric"
	entities        = [data.newrelic_application.application.application_id]
	metric          = "response_time_web"
	condition_scope = "application"

	critical {
		duration      = var.service.duration
		threshold     = var.service.response_time_threshold
		operator      = "above"
		time_function = "all"
	}
}

resource "newrelic_alert_condition" "throughput_web" {
	policy_id = newrelic_alert_policy.golden_signal_policy.id

	name            = "Low Throughput (web)"
	type            = "apm_app_metric"
	entities        = [data.newrelic_application.application.application_id]
	metric          = "throughput_web"
	condition_scope = "application"

	critical {
		duration      = var.service.duration
		threshold     = var.service.throughput_threshold
		operator      = "below"
		time_function = "all"
	}
}

resource "newrelic_alert_condition" "error_percentage" {
	policy_id = newrelic_alert_policy.golden_signal_policy.id

	name            = "High Error Percentage"
	type            = "apm_app_metric"
	entities        = [data.newrelic_application.application.application_id]
	metric          = "error_percentage"
	condition_scope = "application"

	critical {
		duration      = var.service.duration
		threshold     = var.service.error_percentage_threshold
		operator      = "above"
		time_function = "all"
	}
}

resource "newrelic_infra_alert_condition" "high_cpu" {
	policy_id = newrelic_alert_policy.golden_signal_policy.id

	name       = "High CPU usage"
	type       = "infra_metric"
	event      = "SystemSample"
	select     = "cpuPercent"
	comparison = "above"
	where      = "(`applicationId` = '${data.newrelic_application.application.application_id}')"

	critical {
		duration      = var.service.duration
		value         = var.service.cpu_threshold
		time_function = "all"
	}
}

resource "newrelic_alert_policy_channel" "alert_policy_channel" {
	policy_id  = newrelic_alert_policy.golden_signal_policy.id
	channel_ids  = var.alert_channel_ids
}
