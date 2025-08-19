# Cross-Tenancy New Relic Read-Only Access Policy
resource "oci_identity_policy" "cross_tenancy_read_only_policy" {
  compartment_id = var.compartment_ocid
  name           = "New_Relic_Cross_Tenancy_Read_Only_Policy"
  description    = "Policy granting New Relic tenancy read-only access to connector hubs, VCNs, and log groups."
  statements     = [
    "Define tenancy NRTenancyAlias as ${var.new_relic_tenancy_ocid}",
    "Define group NRCustomerOCIAccessGroupAlias as ${var.new_relic_group_ocid}",
    "Admit group NRCustomerOCIAccessGroupAlias of tenancy NRTenancyAlias to read virtual-network-family in tenancy",
    "Admit group NRCustomerOCIAccessGroupAlias of tenancy NRTenancyAlias to read log-content in tenancy",
  ]
}

# Policies for Connector Hubs in given Compartment
resource "oci_identity_dynamic_group" "connector_hub_dg" {
  compartment_id = var.tenancy_ocid
  name           = "New_Relic_Service_Connector_Hubs_DG"
  description    = "Dynamic group for all Service Connector Hubs in the specified compartment."
  matching_rule  = "ALL {resource.type = 'serviceconnector', instance.compartment.id = '${var.compartment_ocid}'}"
}

resource "oci_identity_policy" "connector_hub_policy" {
  compartment_id = var.compartment_ocid
  name           = "New_Relic_Connector_Hub_Log_Access"
  description    = "Allows connector hubs to read logs and trigger functions."
  statements     = [
    "Allow dynamic-group ${oci_identity_dynamic_group.connector_hub_dg.name} to read log-content in tenancy",
    "Allow dynamic-group ${oci_identity_dynamic_group.connector_hub_dg.name} to use fn-function in compartment id ${var.compartment_ocid}",
  ]
}

# Cross-Regional Vault Access for Functions
resource "oci_identity_dynamic_group" "all_functions_dg" {
  compartment_id = var.tenancy_ocid
  name           = "New_Relic_All_Functions_DG"
  description    = "Dynamic group for all functions in the compartment."
  matching_rule  = "ALL {instance.compartment.id = '${var.compartment_ocid}'}"
}

resource "oci_identity_policy" "functions_vault_access_policy" {
  compartment_id = var.compartment_ocid
  name           = "New_Relic_Functions_Vault_Access_Policy"
  description    = "Policy allowing functions to read secrets from the vault."
  statements     = [
    "Allow dynamic-group ${oci_identity_dynamic_group.all_functions_dg.name} to read secret-bundles in compartment id ${var.compartment_ocid}",
  ]
}

# Vault Resources for New Relic API Key
resource "oci_kms_vault" "newrelic_vault" {
  compartment_id = var.compartment_ocid
  display_name   = var.kms_vault_name
  vault_type     = "DEFAULT"
}

resource "oci_kms_key" "newrelic_key" {
  compartment_id = var.compartment_ocid
  display_name   = "newrelic-key"
  key_shape {
    algorithm = "AES"
    length    = 32
  }
  management_endpoint = oci_kms_vault.newrelic_vault.management_endpoint
}

resource "oci_vault_secret" "api_key" {
  compartment_id = var.compartment_ocid
  vault_id       = oci_kms_vault.newrelic_vault.id
  key_id         = oci_kms_key.newrelic_key.id
  secret_name    = "NewRelicAPIKey"
  description    = "Secret containing New Relic ingest API key"
  secret_content {
    content_type = "BASE64"
    content      = base64encode(var.newrelic_ingest_api_key)
    name         = "testkey"
  }
}

# Link OCI Tenancy to New Relic Account
resource "newrelic_cloud_oci_link_account" "nr_link_account" {
  tenant_id = var.tenancy_ocid
  name = var.tenancy_ocid
  description = "Links OCI tenancy to New Relic for logging integration"
}