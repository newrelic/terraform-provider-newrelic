data "oci_identity_user" "current_user" {
  user_id = var.current_user_ocid
}

data "oci_identity_region_subscriptions" "subscriptions" {
  tenancy_id = var.tenancy_ocid
}

data "oci_identity_tenancy" "current_tenancy" {
  tenancy_id = var.tenancy_ocid
}

data "oci_identity_policies" "existing_policies" {
  compartment_id = var.compartment_ocid
}

data "oci_identity_dynamic_groups" "existing_dynamic_groups" {
  compartment_id = var.tenancy_ocid
}

data "oci_functions_applications" "existing_functions_apps" {
  compartment_id = var.compartment_ocid
}

locals {
  freeform_tags = {
    newrelic-terraform = "true"
  }
  # Names for the network infra
  vcn_name        = var.vcn_name
  nat_gateway     = "${local.vcn_name}-natgateway"
  service_gateway = "${local.vcn_name}-servicegateway"
  subnet          = "${local.vcn_name}-public-subnet"
}

resource "oci_kms_vault" "newrelic_vault" {
  compartment_id = var.compartment_ocid
  display_name   = var.kms_vault_name
  vault_type     = "DEFAULT"
  freeform_tags  = local.freeform_tags
}

resource "oci_kms_key" "newrelic_key" {
  compartment_id = var.compartment_ocid
  display_name   = "newrelic-key"
  key_shape {
    algorithm = "AES"
    length    = 32
  }
  management_endpoint = oci_kms_vault.newrelic_vault.management_endpoint
  freeform_tags       = local.freeform_tags
}

resource "oci_vault_secret" "api_key" {
  compartment_id = var.compartment_ocid
  vault_id       = oci_kms_vault.newrelic_vault.id
  key_id         = oci_kms_key.newrelic_key.id
  secret_name    = "NewRelicAPIKey"
  description    = "[DO NOT REMOVE] Secret containing New Relic ingest API key for metrics"
  secret_content {
    content_type = "BASE64"
    content      = base64encode(var.newrelic_ingest_api_key)
    name         = "testkey"
  }
  freeform_tags = local.freeform_tags
}

resource "newrelic_cloud_oci_link_account" "nr_link_account" {  
  tenant_id = var.tenancy_ocid
  name = var.tenancy_ocid
}

#Resource for the function application
resource "oci_functions_application" "metrics_function_app" {
  depends_on     = [oci_identity_policy.nr_metrics_policy]
  compartment_id = var.compartment_ocid
  config = {
    "FORWARD_TO_NR"                = "False"
    "LOGGING_ENABLED"              = "True"
    "NR_METRIC_ENDPOINT"           = var.newrelic_endpoint
    "TENANCY_OCID"                 = var.compartment_ocid
    "SECRET_OCID"                  = oci_vault_secret.api_key.id
    "VAULT_REGION"                 = var.region
  }
  defined_tags               = {}
  display_name               = var.newrelic_function_app
  freeform_tags              = local.freeform_tags
  network_security_group_ids = []
  shape                      = var.function_app_shape
  subnet_ids = [
    module.vcn[0].subnet_id[local.subnet], # Corrected reference
  ]
}


#Resource for the function
resource "oci_functions_function" "metrics_function" {
  depends_on = [oci_functions_application.metrics_function_app]

  application_id = oci_functions_application.metrics_function_app.id
  display_name   = "${oci_functions_application.metrics_function_app.display_name}-metrics-function"
  memory_in_mbs  = "256"

  defined_tags  = {}
  freeform_tags = local.freeform_tags
  image         = var.function_image
}

#Resource for the service connector hub
resource "oci_sch_service_connector" "nr_service_connector" {
  depends_on     = [oci_functions_function.metrics_function]
  compartment_id = var.compartment_ocid
  display_name   = var.connector_hub_name
  description   = "[DO NOT REMOVE] Service connector hub for pushing  metrics to New Relic"

  # Source Configuration with Monitoring
  source {
    kind = "monitoring"

    monitoring_sources {
      compartment_id = var.compartment_ocid
      namespace_details {
        kind = "selected"

        dynamic "namespaces" {
          for_each = var.metrics_namespaces
          content {
            namespace = namespaces.value
            metrics {
              kind = "all" // Adjust based on actual needs, possibly sum, mean, count
            }
          }
        }
      }
    }
  }

  # Target Configuration with Streaming
  target {
    #Required
    kind = "functions"

    #Optional
    batch_size_in_kbs = 100
    batch_time_in_sec = 60
    compartment_id    = var.compartment_ocid
    function_id       = oci_functions_function.metrics_function.id
  }

  # Optional tags and additional metadata
  defined_tags  = {}
  freeform_tags = {}
}

resource "newrelic_cloud_oci_integrations" "newrelic_cloud_integration_pull" {
  linked_account_id = newrelic_cloud_oci_link_account.nr_link_account.id
  oci_metadata_and_tags {}
}