resource "oci_identity_compartment" "newrelic_compartment" {
  compartment_id = var.tenancy_ocid
  name           = "newrelic-compartment-${local.terraform_suffix}"
  description    = "[DO NOT REMOVE] Compartment for New Relic integration resources"
  enable_delete  = false
  freeform_tags  = local.freeform_tags
}

#Key Vault and Secret for New Relic Ingest and User API Key
resource "oci_kms_vault" "newrelic_vault" {
  count = local.create_vault ? 1 : 0
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
  count = local.create_vault ? 1 : 0
  compartment_id = oci_identity_compartment.newrelic_compartment.id
  display_name   = "newrelic-key-${local.terraform_suffix}"
  key_shape {
    algorithm = "AES"
    length    = 32
  }
  management_endpoint = oci_kms_vault.newrelic_vault[count.index].management_endpoint
  freeform_tags       = local.freeform_tags
  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}

resource "oci_vault_secret" "ingest_api_key" {
  count = local.newrelic_monitoring_active && !local.is_ingest_vault_key_present ? 1 : 0
  compartment_id = oci_identity_compartment.newrelic_compartment.id
  vault_id       = oci_kms_vault.newrelic_vault[count.index].id
  key_id         = oci_kms_key.newrelic_key[count.index].id
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
  count = local.newrelic_monitoring_active && !local.is_user_vault_key_present ? 1 : 0
  compartment_id = oci_identity_compartment.newrelic_compartment.id
  vault_id       = oci_kms_vault.newrelic_vault[count.index].id
  key_id         = oci_kms_key.newrelic_key[count.index].id
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

# Dynamic group for Service Connector Hub — only needed for METRICS/LOGS integrations.
# COST-only integrations use Scheduler Primitive and do not require a connector or function.
resource "oci_identity_dynamic_group" "nr_service_connector_group" {
  count          = local.newrelic_monitoring_active ? 1 : 0
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

#Resource for the metrics/Logging (Common) policies — only needed when monitoring is active
resource "oci_identity_policy" "nr_common_policy" {
  count          = local.is_home_region && local.newrelic_monitoring_active ? 1 : 0
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

# Cross-tenancy endorsement so WIF-authenticated principals can read Oracle's shared
# cost-report bucket. Only deployed when COST is in instrumentation_type.
resource "oci_identity_policy" "nr_cost_policy" {
  count          = local.is_home_region && local.newrelic_cost_access_policy ? 1 : 0
  compartment_id = var.tenancy_ocid
  description    = "[DO NOT REMOVE] Policy to endorse New Relic principals to read OCI cost reports from the shared Oracle reporting tenancy"
  name           = local.newrelic_cost_policy
  statements = [
    "DEFINE tenancy usage-report as ocid1.tenancy.oc1..aaaaaaaaned4fkpkisbwjlr56u7cj63lf3wffbilvqknstgtvzub7vhqkggq",
    "endorse any-user to read objects in tenancy usage-report",
  ]
  defined_tags  = {}
  freeform_tags = local.freeform_tags
}

resource "newrelic_cloud_oci_link_account" "linkAccount" {
  account_id  = var.newrelic_account_id
  name        = local.linked_account_name
  tenant_id   = var.tenancy_ocid
  oci_home_region = local.home_region

  # Vault and compartment fields are only required for monitoring (METRICS/LOGS).
  # Cost-only accounts skip vault creation and leave these fields empty.
  compartment_ocid  = local.newrelic_monitoring_active ? oci_identity_compartment.newrelic_compartment.id : null
  ingest_vault_ocid = local.newrelic_monitoring_active ? (local.is_ingest_vault_key_present ? var.ingest_key_secret_ocid : try(oci_vault_secret.ingest_api_key[0].id, null)) : null
  user_vault_ocid   = local.newrelic_monitoring_active ? (local.is_user_vault_key_present ? var.user_key_secret_ocid : try(oci_vault_secret.user_api_key[0].id, null)) : null

  oci_client_id        = var.client_id
  oci_client_secret    = var.client_secret
  oci_domain_url       = var.oci_domain_url
  instrumentation_type = var.instrumentation_type
  # NR-562518: propagate the trust shape created by wif-setup to the linked account so the
  # worker mints the matching token type. Defaults preserve UPST behavior for existing callers.
  trust_type   = var.trust_type
  resource_tag = var.resource_tag
}

output "compartment_ocid" {
  value = oci_identity_compartment.newrelic_compartment.id
}

output "ingest_vault_ocid" {
  value = local.newrelic_monitoring_active ? (local.is_ingest_vault_key_present ? var.ingest_key_secret_ocid : try(oci_vault_secret.ingest_api_key[0].id, null)) : null
}

output "user_vault_ocid" {
  value = local.newrelic_monitoring_active ? (local.is_user_vault_key_present ? var.user_key_secret_ocid : try(oci_vault_secret.user_api_key[0].id, null)) : null
}

output "provider_account_id" {
  value = newrelic_cloud_oci_link_account.linkAccount.id
}
