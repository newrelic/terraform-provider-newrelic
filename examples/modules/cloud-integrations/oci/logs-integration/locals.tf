locals {
  home_region = [
    for rs in data.oci_identity_region_subscriptions.subscriptions.region_subscriptions :
    rs.region_name if rs.region_key == data.oci_identity_tenancy.current_tenancy.home_region_key
  ][0]

  freeform_tags = {
    newrelic-terraform = "true"
  }

  terraform_suffix = "tf"

  # VCN Constants
  vcn_name         = "newrelic-${var.newrelic_logging_identifier}-${var.region}-vcn-${local.terraform_suffix}"
  nat_gateway      = "newrelic-${var.newrelic_logging_identifier}-${var.region}-natgateway-${local.terraform_suffix}"
  service_gateway  = "newrelic-${var.newrelic_logging_identifier}-${var.region}-servicegateway-${local.terraform_suffix}"
  internet_gateway = "newrelic-${var.newrelic_logging_identifier}-${var.region}-internetgateway-${local.terraform_suffix}"
  vcn_dns_label    = "nrlogging"
  vcn_cidr_block   = "10.0.0.0/16"

  # Subnet Constants
  subnet            = "newrelic-${var.newrelic_logging_identifier}-${var.region}-private-subnet-${local.terraform_suffix}"
  subnet_cidr_block = "10.0.0.0/16"
  subnet_type       = "private"

  # Route Table Constants
  internet_destination = "0.0.0.0/0"

  # Function App Constants
  function_app_name  = "newrelic-${var.newrelic_logging_identifier}-${var.region}-logs-function-app-${local.terraform_suffix}"
  function_app_shape = "GENERIC_X86"
  client_ttl         = 30

  # Function Constants
  function_name          = "newrelic-${var.newrelic_logging_identifier}-${var.region}-logs-function-${local.terraform_suffix}"
  function_memory_in_mbs = "128"
  time_out_in_seconds    = 300

  # The function runs from a copy of var.function_image in the tenancy's own Container Registry.
  ocir_host                 = "${var.region}.ocir.io"
  ocir_namespace            = oci_artifacts_container_repository.log_forwarder_repo.namespace
  function_image_repository = "newrelic-${var.newrelic_logging_identifier}-${local.terraform_suffix}/oci-log-forwarder"
  function_image_digest     = data.external.function_image.result.digest
  function_image            = "${local.ocir_host}/${local.ocir_namespace}/${local.function_image_repository}:${data.external.function_image.result.tag}"

  create_registry_token = nonsensitive(var.registry_auth_token == "")
  registry_username     = "${local.ocir_namespace}/${var.registry_username != "" ? var.registry_username : data.oci_identity_user.registry_user[0].name}"
  # Unmarked so Terraform keeps showing the copy's log output; image_mirror.py never prints it.
  registry_password = local.create_registry_token ? oci_identity_auth_token.registry_push[0].token : nonsensitive(var.registry_auth_token)

  # The tenancy's root compartment shares its OCID with the tenancy itself, but OCI's
  # Identity service can't return it from GetCompartment (data.oci_identity_compartment) --
  # only GetTenancy can. Detect that case so we read the name from the right data source.
  is_root_compartment = var.compartment_ocid == var.tenancy_ocid
  compartment_name    = local.is_root_compartment ? data.oci_identity_tenancy.current_tenancy.name : data.oci_identity_compartment.current_compartment[0].name

  user_api_key = base64decode(data.oci_secrets_secretbundle.user_api_key.secret_bundle_content[0].content)
  newrelic_graphql_endpoint = {
    US = "https://api.newrelic.com/graphql"
    EU = "https://api.eu.newrelic.com/graphql"
    JP = "https://api.jp.newrelic.com/graphql"
  }[var.new_relic_region]

  updateLinkAccount_graphql_query = <<EOF
  mutation {
    cloudUpdateAccount(
      accountId: ${var.newrelic_account_id}
      accounts: {
        oci: {
          linkedAccountId: ${var.provider_account_id}
          loggingStackOcid: "tf"
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
