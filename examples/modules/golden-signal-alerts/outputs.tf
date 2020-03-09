output "policy_id" {
  value = newrelic_alert_policy.golden_signal_policy.id
}

output "response_time_condition_id" {
  value = newrelic_alert_condition.response_time_web.id
}

output "throughput_condition_id" {
  value = newrelic_alert_condition.throughput_web.id
}

output "error_percentage_condition_id" {
  value = newrelic_alert_condition.error_percentage.id
}

output "cpu_condition_id" {
  value = newrelic_infra_alert_condition.high_cpu.id
}
