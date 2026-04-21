locals {
  freeform_tags = {
    newrelic-terraform = "true"
  }

  terraform_suffix = "tf"

  # VCN Constants
  vcn_name          = "newrelic-${var.newrelic_logging_identifier}-${var.region}-vcn-${local.terraform_suffix}"
  nat_gateway       = "newrelic-${var.newrelic_logging_identifier}-${var.region}-natgateway-${local.terraform_suffix}"
  service_gateway   = "newrelic-${var.newrelic_logging_identifier}-${var.region}-servicegateway-${local.terraform_suffix}"
  internet_gateway  = "newrelic-${var.newrelic_logging_identifier}-${var.region}-internetgateway-${local.terraform_suffix}"
  vcn_dns_label     = "nrlogging"
  vcn_cidr_block    = "10.0.0.0/16"

  # Subnet Constants
  subnet               = "newrelic-${var.newrelic_logging_identifier}-${var.region}-private-subnet-${local.terraform_suffix}"
  subnet_cidr_block    = "10.0.0.0/16"
  subnet_type          = "private"

  # Route Table Constants
  internet_destination = "0.0.0.0/0"

  # Function App Constants
  function_app_name  = "newrelic-${var.newrelic_logging_identifier}-${var.region}-logs-function-app-${local.terraform_suffix}"
  function_app_shape = "GENERIC_X86"
  client_ttl         = 30

  # Function Constants
  function_name                 = "newrelic-${var.newrelic_logging_identifier}-${var.region}-logs-function-${local.terraform_suffix}"
  function_memory_in_mbs        = "128"
  time_out_in_seconds           = 300
  image_url                     = "${var.region}.ocir.io/idptojlonu4e/newrelic-logs-integration/oci-log-forwarder:${var.image_version}"

  user_api_key = base64decode(data.oci_secrets_secretbundle.user_api_key.secret_bundle_content[0].content)
  newrelic_graphql_endpoint = {
    US = "https://api.newrelic.com/graphql"
    EU = "https://api.eu.newrelic.com/graphql"
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
