locals {
  freeform_tags = {
    newrelic-terraform = "true"
  }
  image_url = "${var.region}.ocir.io/${var.tenancy_namespace}/${var.repository_name}/${var.function_name}:${var.repository_version}"
}

resource "oci_identity_dynamic_group" "nr_serviceconnector_group" {
  compartment_id = var.tenancy_ocid
  description    = "[DO NOT REMOVE] Dynamic group for service connector"
  matching_rule  = "All {resource.type = 'serviceconnector'}"
  name           = var.dynamic_group_name
  defined_tags   = {}
  freeform_tags  = local.freeform_tags
}

# Dynamic group for service connector
resource "oci_identity_policy" "log_forwarding_policy" {
  depends_on     = [oci_identity_dynamic_group.nr_serviceconnector_group]
  compartment_id = var.tenancy_ocid
  description    = "[DO NOT REMOVE] Policy to have any connector hub read from source and write to a target function"
  name           = var.newrelic_logs_policy
  statements     = [
    "Allow dynamic-group ${var.dynamic_group_name} to read logs in tenancy",
    "Allow dynamic-group ${var.dynamic_group_name} to use fn-function in tenancy",
    "Allow dynamic-group ${var.dynamic_group_name} to use fn-invocation in tenancy",
  ]
  defined_tags  = {}
  freeform_tags = local.freeform_tags
}

# New Relic API Access Key
resource "newrelic_api_access_key" "newrelic_aws_access_key" {
  account_id  = var.newrelic_account_id
  key_type    = "INGEST"
  ingest_type = "LICENSE"
  name        = "logging-integrations-ingest-key"
  notes       = "Ingest License key for OCI Logging Integrations"
}

# Function Application
resource "oci_functions_application" "logs_function_app" {
  depends_on     = [oci_identity_policy.log_forwarding_policy]
  compartment_id = var.compartment_ocid
  config = {
    "LICENSE_KEY"                  = newrelic_api_access_key.newrelic_aws_access_key.key
    "DEBUG_ENABLED"                = var.debug_enabled
    "REGION"                       = var.nr_region
    "LOG_GROUP_ID"                 = var.log_group_id
    "ACCOUNT_ID"                   = var.newrelic_account_id
  }
  defined_tags               = {}
  display_name               = var.function_app_name
  freeform_tags              = local.freeform_tags
  network_security_group_ids = []
  shape                      = var.function_app_shape
  subnet_ids                 = [var.subnet_id]
}

# Log Forwarding Function
resource "oci_functions_function" "logs_function" {
  depends_on = [oci_functions_application.logs_function_app]

  application_id = oci_functions_application.logs_function_app.id
  display_name   = "${oci_functions_application.logs_function_app.display_name}-logs-function"
  memory_in_mbs  = "256"

  defined_tags  = {}
  freeform_tags = local.freeform_tags
  image         = local.image_url
}

# Service Connector Hub
resource "oci_sch_service_connector" "nr_service_connector" {
  depends_on     = [oci_functions_function.logs_function]
  compartment_id = var.compartment_ocid
  display_name   = var.connector_hub_name

  source {
    kind = "logging"
    log_sources {
      compartment_id = var.compartment_ocid
      log_group_id   = var.log_group_id
      log_id         = var.log_id
    }
  }

  target {
    kind = "functions"
    compartment_id    = var.compartment_ocid
    function_id       = oci_functions_function.logs_function.id
  }

  description   = "Service Connector from Logging to Forwarding Function"
  freeform_tags = local.freeform_tags
}