locals {
  home_region = [
    for region in data.oci_identity_region_subscriptions.subscriptions.region_subscriptions : region.region_name
    if region.is_home_region
  ][0]
  home_region_key = [
    for region in data.oci_identity_region_subscriptions.subscriptions.region_subscriptions : region.region_key
    if region.is_home_region
  ][0]
  is_home_region = var.region == local.home_region || lower(var.region) == lower(local.home_region_key)

  freeform_tags = {
    newrelic-terraform = "true"
  }

  terraform_suffix               = "tf"
  newrelic_metrics_access_policy = contains(split(",", var.instrumentation_type), "METRICS")
  newrelic_logs_access_policy    = contains(split(",", var.instrumentation_type), "LOGS")
  newrelic_logs_policy           = "newrelic_logs_policy_DO_NOT_REMOVE-${local.terraform_suffix}"
  newrelic_metrics_policy        = "newrelic_metrics_policy_DO_NOT_REMOVE-${local.terraform_suffix}"
  newrelic_common_policy         = "newrelic_common_policy_DO_NOT_REMOVE-${local.terraform_suffix}"
  dynamic_group_name             = "newrelic_dynamic_group_DO_NOT_REMOVE-${local.terraform_suffix}"
  linked_account_name            = "${var.nr_prefix}-oci-${local.terraform_suffix}"
}
