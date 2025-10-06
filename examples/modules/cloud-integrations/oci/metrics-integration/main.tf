resource "oci_functions_application" "metrics_function_app" {
  compartment_id = var.compartment_ocid
  config = {
    "FORWARD_TO_NR"      = "True"
    "LOGGING_ENABLED"    = "False"
    "NR_METRIC_ENDPOINT" = var.newrelic_endpoint
    "TENANCY_OCID"       = var.tenancy_ocid
    "SECRET_OCID"        = var.ingest_api_secret_ocid
    "VAULT_REGION"       = local.home_region
  }
  defined_tags               = {}
  display_name               = "newrelic-${var.nr_prefix}-${var.region}-function-app-${local.terraform_suffix}"
  freeform_tags              = local.freeform_tags
  network_security_group_ids = []
  shape                      = "GENERIC_X86"
  subnet_ids = [
    data.oci_core_subnet.input_subnet.id,
  ]
}

resource "oci_functions_function" "metrics_function" {
  application_id = oci_functions_application.metrics_function_app.id
  depends_on     = [oci_functions_application.metrics_function_app]
  display_name   = "newrelic-${var.nr_prefix}-${var.region}-metrics-function-${local.terraform_suffix}"
  memory_in_mbs  = "128"
  defined_tags   = {}
  freeform_tags  = local.freeform_tags
  image          = "${var.region}.ocir.io/${var.image_bucket}/newrelic-metrics-integration/oci-metrics-forwarder:${var.image_version}"
}

resource "oci_sch_service_connector" "service_connector" {
  for_each       = { for hub in jsondecode(var.connector_hubs_data) : hub["name"] => hub }
  depends_on     = [oci_functions_function.metrics_function]
  compartment_id = var.compartment_ocid
  display_name   = "${each.value["name"]} (${local.terraform_suffix})"
  description    = each.value["description"]
  freeform_tags  = local.freeform_tags

  source {
    kind = "monitoring"

    dynamic "monitoring_sources" {
      for_each = each.value["compartments"]
      content {
        compartment_id = monitoring_sources.value["compartment_id"]
        namespace_details {
          kind = "selected"

          dynamic "namespaces" {
            for_each = monitoring_sources.value["namespaces"]
            content {
              namespace = namespaces.value
              metrics {
                kind = "all"
              }
            }
          }
        }
      }
    }
  }

  target {
    kind              = "functions"
    function_id       = oci_functions_function.metrics_function.id
    batch_size_in_kbs = 100
    batch_time_in_sec = 60
  }
}
