# --- Function App Resources ---
resource "oci_functions_application" "logging_function_app" {
  compartment_id = var.compartment_ocid
  config = {
    "VAULT_REGION"      = var.region
    "DEBUG_ENABLED"     = var.debug_enabled
    "NEW_RELIC_REGION"  = var.new_relic_region
    "SECRET_OCID"       = var.secret_ocid
    "CLIENT_TTL"        = local.client_ttl
  }
  display_name               = local.function_app_name
  freeform_tags              = local.freeform_tags
  shape                      = local.function_app_shape
  subnet_ids                 = [var.create_vcn ? module.vcn[0].subnet_id[local.subnet] : var.function_subnet_id]
}

# --- Function Resources ---
resource "oci_functions_function" "logging_function" {
  application_id     = oci_functions_application.logging_function_app.id
  display_name       = local.function_name
  memory_in_mbs      = local.function_memory_in_mbs
  timeout_in_seconds = local.time_out_in_seconds
  freeform_tags      = local.freeform_tags
  image              = local.image_url
}

# --- Service Connector Hub - Routes logs to New Relic function ---
resource "oci_sch_service_connector" "nr_logging_service_connector" {
  for_each = var.connector_hub_details != null ? {
    for connector in jsondecode(var.connector_hub_details) : connector.display_name => connector
  } : {}

  compartment_id = var.compartment_ocid
  display_name   = each.value.display_name
  description    = each.value.description
  freeform_tags  = local.freeform_tags

  source {
    kind = "logging"
    dynamic "log_sources" {
      for_each = each.value.log_sources
      content {
        compartment_id = log_sources.value.compartment_id
        log_group_id   = log_sources.value.log_group_id
      }
    }
  }

  target {
    kind              = "functions"
    batch_size_in_kbs = var.batch_size_in_kbs
    batch_time_in_sec = var.batch_time_in_sec
    compartment_id    = var.compartment_ocid
    function_id       = oci_functions_function.logging_function.id
  }
}