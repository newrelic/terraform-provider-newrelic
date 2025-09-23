resource "oci_identity_compartment" "newrelic_compartment" {
  compartment_id = var.tenancy_ocid
  name           = "newrelic-compartment-${local.terraform_suffix}"
  description    = "[DO NOT REMOVE] Compartment for New Relic integration resources"
  enable_delete  = false
  freeform_tags  = local.freeform_tags
}

#Key Vault and Secret for New Relic Ingest and User API Key
resource "oci_kms_vault" "newrelic_vault" {
  compartment_id = oci_identity_compartment.newrelic_compartment.id
  display_name   = "newrelic-vault-${local.terraform_suffix}"
  vault_type     = "DEFAULT"
  freeform_tags  = local.freeform_tags
  timeouts {
    create = "60m"
    update = "60m"
    delete = "60m"
  }
}

resource "oci_kms_key" "newrelic_key" {
  compartment_id = oci_identity_compartment.newrelic_compartment.id
  display_name   = "newrelic-key-${local.terraform_suffix}"
  key_shape {
    algorithm = "AES"
    length    = 32
  }
  management_endpoint = oci_kms_vault.newrelic_vault.management_endpoint
  freeform_tags       = local.freeform_tags
  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}

resource "oci_vault_secret" "ingest_api_key" {
  compartment_id = oci_identity_compartment.newrelic_compartment.id
  vault_id       = oci_kms_vault.newrelic_vault.id
  key_id         = oci_kms_key.newrelic_key.id
  secret_name    = "NewRelicIngestAPIKey"
  secret_content {
    content_type = "BASE64"
    content      = base64encode(var.newrelic_ingest_api_key)
  }
  freeform_tags = local.freeform_tags
  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}

resource "oci_vault_secret" "user_api_key" {
  compartment_id = oci_identity_compartment.newrelic_compartment.id
  vault_id       = oci_kms_vault.newrelic_vault.id
  key_id         = oci_kms_key.newrelic_key.id
  secret_name    = "NewRelicUserAPIKey"
  secret_content {
    content_type = "BASE64"
    content      = base64encode(var.newrelic_user_api_key)
  }
  freeform_tags = local.freeform_tags
  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}

#Resource for the dynamic group
resource "oci_identity_dynamic_group" "nr_service_connector_group" {
  compartment_id = var.tenancy_ocid
  description    = "[DO NOT REMOVE] Dynamic group for service connector"
  matching_rule  = "ANY {resource.type = 'serviceconnector', resource.type = 'fnfunc'}"
  name           = local.dynamic_group_name
  defined_tags   = {}
  freeform_tags  = local.freeform_tags
}

#Resource for the metrics policy
resource "oci_identity_policy" "nr_metrics_policy" {
  count          = local.is_home_region && local.newrelic_metrics_access_policy ? 1 : 0
  depends_on     = [oci_identity_dynamic_group.nr_service_connector_group]
  compartment_id = var.tenancy_ocid
  description    = "[DO NOT REMOVE] Policy to have read metrics for newrelic integration"
  name           = local.newrelic_metrics_policy
  statements = [
    "Allow dynamic-group ${local.dynamic_group_name} to read metrics in tenancy"
  ]
  defined_tags  = {}
  freeform_tags = local.freeform_tags
}

#Resource for the logging policy
resource "oci_identity_policy" "nr_logs_policy" {
  count          = local.is_home_region && local.newrelic_logs_access_policy ? 1 : 0
  depends_on     = [oci_identity_dynamic_group.nr_service_connector_group]
  compartment_id = var.tenancy_ocid
  description    = "[DO NOT REMOVE] Policy to have read logs for newrelic integration"
  name           = local.newrelic_logs_policy
  statements = [
    "Allow dynamic-group ${local.dynamic_group_name} to read log-content in tenancy"
  ]
  defined_tags  = {}
  freeform_tags = local.freeform_tags
}

#Resource for the metrics/Logging (Common) policies
resource "oci_identity_policy" "nr_common_policy" {
  depends_on     = [oci_identity_dynamic_group.nr_service_connector_group]
  compartment_id = var.tenancy_ocid
  description    = "[DO NOT REMOVE] Policy to have any connector hub read from monitoring source and write to a target function"
  name           = local.newrelic_common_policy
  statements = [
    "Allow dynamic-group ${local.dynamic_group_name} to use fn-function in tenancy",
    "Allow dynamic-group ${local.dynamic_group_name} to use fn-invocation in tenancy",
    "Allow dynamic-group ${local.dynamic_group_name} to read secret-bundles in tenancy",
  ]
  defined_tags  = {}
  freeform_tags = local.freeform_tags
}

resource "newrelic_cloud_oci_link_account" "linkAccount" {
  account_id           = var.newrelic_account_id
  name                 = local.linked_account_name
  compartment_ocid     = oci_identity_compartment.newrelic_compartment.id
  oci_home_region      = local.home_region
  tenant_id            = var.tenancy_ocid
  ingest_vault_ocid    = oci_vault_secret.ingest_api_key.id
  user_vault_ocid      = oci_vault_secret.user_api_key.id
  oci_client_id        = var.client_id
  oci_client_secret    = var.client_secret
  oci_domain_url       = var.oci_domain_url
  oci_svc_user_name    = var.svc_user_name
  instrumentation_type = var.instrumentation_type
}

output "compartment_ocid" {
  value = oci_identity_compartment.newrelic_compartment.id
}

output "ingest_vault_ocid" {
  value = oci_vault_secret.ingest_api_key.id
}

output "user_vault_ocid" {
  value = oci_vault_secret.user_api_key.id
}
