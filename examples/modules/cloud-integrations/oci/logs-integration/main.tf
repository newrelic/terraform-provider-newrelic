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

# Resource to link the New Relic account and configure the integration
resource "null_resource" "newrelic_link_account" {
  depends_on = [oci_functions_function.logging_function, oci_sch_service_connector.nr_logging_service_connector]
  provisioner "local-exec" {
    command = <<EOT
      # Main execution for cloudLinkAccount
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
