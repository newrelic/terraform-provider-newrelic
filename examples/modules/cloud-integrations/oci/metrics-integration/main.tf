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

# Resource to link the New Relic account and configure the integration
resource "null_resource" "newrelic_update_account" {
  depends_on = [oci_functions_function.metrics_function, oci_sch_service_connector.service_connector]
  provisioner "local-exec" {
    command = <<EOT
      # Main execution for cloudUpdateAccount
      response=$(curl --silent --request POST \
        --url "${local.newrelic_graphql_endpoint}" \
        --header "API-Key: ${local.user_api_key}" \
        --header "Content-Type: application/json" \
        --header "User-Agent: insomnia/11.1.0" \
        --data '${jsonencode({
          query = local.updateLinkAccount_graphql_query
        })}')
        # Log the full response for debugging
        echo "Full Response: $response"
        # Extract errors from the response
        root_errors=$(echo "$response" | jq -r '.errors[]?.message // empty')
        update_account_errors=$(echo "$response" | jq -r '.data.cloudUpdateAccount.errors[]?.message // empty')
        # Check if data is null which indicates a possible error
        data_null=$(echo "$response" | jq -r 'if .data.cloudUpdateAccount == null then "true" else "false" end')
        # Combine errors
        errors="$root_errors"$'\n'"$update_account_errors"
        errors=$(echo "$errors" | grep -v '^$')
        # Check if errors exist or data is null
        if [ -n "$errors" ] || [ "$data_null" == "true" ]; then
          echo "Operation failed with the following errors:" >&2
          if [ -n "$errors" ]; then
            echo "$errors" | while IFS= read -r error; do
              echo "- $error" >&2
            done
          fi
          if [ "$data_null" == "true" ] && [ -z "$errors" ]; then
            echo "- GraphQL operation returned null data. Please verify your parameters and query." >&2
          fi
          exit 1
        fi
        echo "Successfully updated New Relic account link"
      EOT
  }
}
