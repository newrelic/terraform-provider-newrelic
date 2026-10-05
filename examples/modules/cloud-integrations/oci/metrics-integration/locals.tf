locals {
  home_region = [
    for rs in data.oci_identity_region_subscriptions.subscriptions.region_subscriptions :
    rs.region_name if rs.region_key == data.oci_identity_tenancy.current_tenancy.home_region_key
  ][0]

  freeform_tags = {
    newrelic-terraform = "true"
  }

  terraform_suffix = "tf"

  # The function runs from a copy of var.function_image in the tenancy's own Container Registry.
  ocir_host                 = "${var.region}.ocir.io"
  ocir_namespace            = oci_artifacts_container_repository.metrics_function_repo.namespace
  function_image_repository = "newrelic-${lower(var.nr_prefix)}-${local.terraform_suffix}/oci-metrics-forwarder"
  function_image_digest     = data.external.function_image.result.digest
  function_image            = "${local.ocir_host}/${local.ocir_namespace}/${local.function_image_repository}:${data.external.function_image.result.tag}"
  create_registry_token     = nonsensitive(var.registry_auth_token == "")
  registry_username         = "${local.ocir_namespace}/${var.registry_username != "" ? var.registry_username : data.oci_identity_user.registry_user[0].name}"
  # Unmarked so Terraform keeps showing the copy's log output; image_mirror.py never prints it.
  registry_password = local.create_registry_token ? oci_identity_auth_token.registry_push[0].token : nonsensitive(var.registry_auth_token)

  # Names for the network infra
  vcn_name        = "newrelic-${var.nr_prefix}-${var.region}-vcn-${local.terraform_suffix}"
  nat_gateway     = "newrelic-${var.nr_prefix}-${var.region}-vcn-${local.terraform_suffix}"
  service_gateway = "newrelic-${var.nr_prefix}-${var.region}-vcn-${local.terraform_suffix}"
  subnet          = "newrelic-${var.nr_prefix}-${var.region}-vcn-${local.terraform_suffix}"

  user_api_key = base64decode(data.oci_secrets_secretbundle.user_api_key.secret_bundle_content[0].content)
  newrelic_graphql_endpoint = {
    US = "https://api.newrelic.com/graphql"
    EU = "https://api.eu.newrelic.com/graphql"
    JP = "https://api.jp.newrelic.com/graphql"
  }[var.newrelic_endpoint]

  updateLinkAccount_graphql_query = <<EOF
  mutation {
    cloudUpdateAccount(
      accountId: ${var.newrelic_account_id}
      accounts: {
        oci: {
          linkedAccountId: ${var.provider_account_id}
          metricStackOcid: "tf"
          ociRegion: "${var.region}"
        }
    }
  ) {
      linkedAccounts {
        id
        authLabel
        createdAt
        disabled
        externalId
        metricCollectionMode
        name
        nrAccountId
        updatedAt
      }
    }
  }
  EOF
}

