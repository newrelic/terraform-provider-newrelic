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

  is_user_vault_key_present   = var.user_key_secret_ocid != ""
  is_ingest_vault_key_present = var.ingest_key_secret_ocid != ""
  # Vault creation is only needed for METRICS/LOGS integrations that don't supply pre-existing
  # vault secret OCIDs. COST-only integrations never need vault secrets.
  create_vault = local.newrelic_monitoring_active && (!local.is_user_vault_key_present || !local.is_ingest_vault_key_present)
  user_api_key = local.is_user_vault_key_present && local.newrelic_monitoring_active ? base64decode(data.oci_secrets_secretbundle.user_api_key[0].secret_bundle_content[0].content) : var.newrelic_user_api_key

  freeform_tags = {
    newrelic-terraform = "true"
  }

  terraform_suffix               = "tf"
  newrelic_metrics_access_policy = contains(split(",", var.instrumentation_type), "METRICS")
  newrelic_logs_access_policy    = contains(split(",", var.instrumentation_type), "LOGS")
  newrelic_cost_access_policy    = contains(split(",", var.instrumentation_type), "COST")
  newrelic_monitoring_active     = local.newrelic_metrics_access_policy || local.newrelic_logs_access_policy
  newrelic_logs_policy           = "newrelic_logs_policy_DO_NOT_REMOVE-${local.terraform_suffix}"
  newrelic_metrics_policy        = "newrelic_metrics_policy_DO_NOT_REMOVE-${local.terraform_suffix}"
  newrelic_common_policy         = "newrelic_common_policy_DO_NOT_REMOVE-${local.terraform_suffix}"
  newrelic_cost_policy           = "newrelic_cost_policy_DO_NOT_REMOVE-${local.terraform_suffix}"
  dynamic_group_name             = "newrelic_dynamic_group_DO_NOT_REMOVE-${local.terraform_suffix}"
  linked_account_name            = "${var.nr_prefix}-oci-${local.terraform_suffix}"
}
