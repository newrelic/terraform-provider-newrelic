locals {
  home_region = [
    for rs in data.oci_identity_region_subscriptions.subscriptions.region_subscriptions :
    rs.region_name if rs.region_key == data.oci_identity_tenancy.current_tenancy.home_region_key
  ][0]

  freeform_tags = {
    newrelic-terraform = "true"
  }

  terraform_suffix = "tf"

  # Names for the network infra
  vcn_name        = "newrelic-${var.nr_prefix}-${var.region}-vcn-${local.terraform_suffix}"
  nat_gateway     = "newrelic-${var.nr_prefix}-${var.region}-vcn-${local.terraform_suffix}"
  service_gateway = "newrelic-${var.nr_prefix}-${var.region}-vcn-${local.terraform_suffix}"
  subnet          = "newrelic-${var.nr_prefix}-${var.region}-vcn-${local.terraform_suffix}"

  user_api_key = base64decode(data.oci_secrets_secretbundle.user_api_key.secret_bundle_content[0].content)
  newrelic_graphql_endpoint = {
    US = "https://api.newrelic.com/graphql"
    EU = "https://api.eu.newrelic.com/graphql"
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

