locals {
  freeform_tags = {
    newrelic-terraform = "true"
  }
  image_url = "iad.ocir.io/idfmbxeaoavl/sanath-testing-registry/oci-function-x86:0.0.1"
  function_app_name = "newrelic-logs-function-app"
  function_app_shape = "GENERIC_X86"
  connector_hub_name = "newrelic-logs-connector-hub"
  newrelic_logs_policy = "newrelic-logs-policy"
  dynamic_group_name = "newrelic-logging-dynamic-group"
}

# ToDo: Enable once infra team completes this task https://github.com/newrelic/terraform-provider-newrelic/pull/2907
# resource "newrelic_cloud_oci_link_account" "nr_link_account" {
#   tenant_id = var.tenancy_ocid
#   name = var.tenancy_ocid
# }

resource "oci_identity_dynamic_group" "nr_logging_service_connector_dg" {
  compartment_id = var.tenancy_ocid
  description    = "[DO NOT REMOVE] Dynamic group for service connector"
  matching_rule  = "All {resource.type = 'serviceconnector'}"
  name           = local.dynamic_group_name
  defined_tags   = {}
  freeform_tags  = local.freeform_tags
}

# Dynamic group for service connector
resource "oci_identity_policy" "log_forwarding_policy" {
  depends_on     = [oci_identity_dynamic_group.nr_logging_service_connector_dg]
  compartment_id = var.tenancy_ocid
  description    = "[DO NOT REMOVE] Policy to have any connector hub read from source and write to a target function"
  name           = local.newrelic_logs_policy
  statements     = [
    "Allow dynamic-group ${local.dynamic_group_name} to read logs in tenancy",
    "Allow dynamic-group ${local.dynamic_group_name} to use fn-function in tenancy",
    "Allow dynamic-group ${local.dynamic_group_name} to use fn-invocation in tenancy",
  ]
  defined_tags  = {}
  freeform_tags = local.freeform_tags
}

# Function Application
resource "oci_functions_application" "logs_function_app" {
  depends_on     = [oci_identity_policy.log_forwarding_policy]
  compartment_id = var.compartment_ocid
  config = {
    "DEBUG_ENABLED"                = var.debug_enabled
    "REGION"                       = var.nr_region
    "LOG_GROUP_ID"                 = var.log_group_id
    "ACCOUNT_ID"                   = var.newrelic_account_id
  }
  defined_tags               = {}
  display_name               = local.function_app_name
  freeform_tags              = local.freeform_tags
  network_security_group_ids = []
  shape                      = local.function_app_shape
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
resource "oci_sch_service_connector" "nr_logging_service_connector" {
  depends_on     = [oci_functions_function.logs_function]
  compartment_id = var.compartment_ocid
  display_name   = local.connector_hub_name

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