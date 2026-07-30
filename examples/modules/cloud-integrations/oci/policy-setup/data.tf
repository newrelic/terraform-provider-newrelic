data "oci_identity_region_subscriptions" "subscriptions" {
  tenancy_id = var.tenancy_ocid
}

data "oci_secrets_secretbundle" "user_api_key" {
  # Only fetch when monitoring is active and the caller supplied an existing vault secret OCID.
  count     = local.newrelic_monitoring_active && local.is_user_vault_key_present ? 1 : 0
  secret_id = var.user_key_secret_ocid
  provider  = oci.home
}