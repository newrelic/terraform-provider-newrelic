output "aws_config_recorder_check" {
  description = "Status of AWS Config Configuration Recorder check"
  value = {
    has_existing_recorder = local.has_existing_recorder
    should_create_new     = local.should_create_recorder
    check_log            = local.recorder_check_log
    recorder_count       = try(data.external.check_existing_recorder.result.count, "unknown")
    existing_recorders   = try(data.external.check_existing_recorder.result.recorder_names, "none")
    region              = try(data.external.check_existing_recorder.result.region, data.aws_region.current.id)
  }
}

output "newrelic_integration_details" {
  description = "New Relic integration configuration details"
  value = {
    account_id = var.newrelic_account_id
    region     = var.newrelic_account_region
    name       = var.name
  }
}
