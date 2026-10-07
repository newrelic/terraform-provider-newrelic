# --- Function App Resources ---
resource "oci_functions_application" "logging_function_app" {
  compartment_id = var.compartment_ocid
  config = {
    "VAULT_REGION"           = local.home_region
    "DEBUG_ENABLED"          = var.debug_enabled
    "NEW_RELIC_REGION"       = var.new_relic_region
    "SECRET_OCID"            = var.secret_ocid
    "CLIENT_TTL"             = local.client_ttl
    "FORWARDER_METRICS_TIER" = var.metrics_tier
    "TENANCY_NAME"           = data.oci_identity_tenancy.current_tenancy.name
    "COMPARTMENT_NAME"       = local.compartment_name
  }
  display_name  = local.function_app_name
  freeform_tags = local.freeform_tags
  shape         = local.function_app_shape
  subnet_ids    = [var.create_vcn ? module.vcn[0].subnet_id[local.subnet] : var.function_subnet_id]
}

resource "oci_artifacts_container_repository" "log_forwarder_repo" {
  compartment_id = var.compartment_ocid
  display_name   = local.function_image_repository
  is_public      = false
  freeform_tags  = local.freeform_tags
}

resource "oci_identity_auth_token" "registry_push" {
  count       = local.create_registry_token ? 1 : 0
  provider    = oci.home
  user_id     = var.user_ocid
  description = "New Relic logs module: pushes the function image to ${local.function_image_repository} in ${var.region}"
}

# Copies the image over the registry HTTP API, so the machine running Terraform needs no Docker.
# Runs again whenever var.function_image resolves to a new digest.
resource "null_resource" "mirror_function_image" {
  triggers = {
    source_digest = local.function_image_digest
    destination   = local.function_image
  }

  provisioner "local-exec" {
    command = "python3 ${path.module}/image_mirror.py copy"
    environment = {
      SOURCE_IMAGE    = var.function_image
      SOURCE_DIGEST   = local.function_image_digest
      DEST_REGISTRY   = local.ocir_host
      DEST_REPOSITORY = "${local.ocir_namespace}/${local.function_image_repository}"
      DEST_TAG        = data.external.function_image.result.tag
      DEST_USERNAME   = local.registry_username
      DEST_PASSWORD   = local.registry_password
    }
  }
}

# --- Function Resources ---
resource "oci_functions_function" "logging_function" {
  depends_on         = [null_resource.mirror_function_image]
  application_id     = oci_functions_application.logging_function_app.id
  display_name       = local.function_name
  memory_in_mbs      = local.function_memory_in_mbs
  timeout_in_seconds = local.time_out_in_seconds
  freeform_tags      = local.freeform_tags
  image              = local.function_image
  # Set explicitly: the function stays on the digest it was created with until this changes.
  image_digest = local.function_image_digest
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
resource "null_resource" "newrelic_update_account" {
  depends_on = [oci_functions_function.logging_function, oci_sch_service_connector.nr_logging_service_connector]
  count      = var.connector_hub_details != null ? 1 : 0
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
