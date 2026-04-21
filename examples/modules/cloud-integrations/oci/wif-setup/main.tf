# Data sources for existing resources
# Retrieves the identity domain details based on the configured domain name
data "oci_identity_domains" "domain" {
  compartment_id = var.tenancy_ocid
  display_name   = var.identity_domain_name
}

# Retrieves the root compartment (tenancy) details
data "oci_identity_compartment" "root" {
  id = var.tenancy_ocid
}

locals {
  identity_domain_id  = data.oci_identity_domains.domain.domains[0].id
  identity_domain_url = data.oci_identity_domains.domain.domains[0].url
  suffix = "tf"

  # New Relic configuration based on region
  newrelic_config = {
    US = {
      issuer_name     = "newrelic-oci-us-production-issuer"
      subject_name    = "newrelic-oci-us-production-user"
      public_jwks_url = "https://publickeys.newrelic.com/r/oci-cmp/us/c5623ba5-1cc7-491a-8ec3-eeee809374f7/jwks.json"
    }
    EU = {
      issuer_name     = "newrelic-oci-eu-production-issuer"
      subject_name    = "newrelic-oci-eu-production-user"
      public_jwks_url = "https://publickeys.eu.newrelic.com/r/oci-cmp/eu/f923dba9-84a8-491c-b714-6c0e61b90c5b/jwks.json"
    }
    JP = {
      issuer_name     = "newrelic-oci-jp-production-issuer"
      subject_name    = "newrelic-oci-jp-production-user"
      # TODO: replace placeholder UUID once NR JP OCI WIF infrastructure is provisioned
      public_jwks_url = "https://publickeys.jp.newrelic.com/r/oci-cmp/jp/00000000-0000-0000-0000-000000000000/jwks.json"
    }
  }[var.newrelic_region]

  # Common resource naming
  resource_prefix = var.resource_prefix != "" ? var.resource_prefix : "newrelic"
}