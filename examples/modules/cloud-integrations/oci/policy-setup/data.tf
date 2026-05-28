data "oci_identity_region_subscriptions" "subscriptions" {
  tenancy_id = var.tenancy_ocid
}

data "oci_secrets_secretbundle" "user_api_key" {
  count = local.create_vault ? 0: 1
  secret_id = var.user_key_secret_ocid
  provider  = oci.home
}