data "oci_secrets_secretbundle" "user_api_key" {
  secret_id = var.user_api_secret_ocid
  provider = oci.home
}
