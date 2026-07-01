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
  identity_domain_url = trimsuffix(data.oci_identity_domains.domain.domains[0].url, ":443")
  suffix              = "tf"

  # NR-562518: trust-type derived flags + counts so module is idempotent for both flows.
  is_upst = var.trust_type == "UPST"
  is_rpst = var.trust_type == "RPST"

  # New Relic configuration based on region. UPST and RPST share the same JWKS endpoint
  # (NR reuses the existing UPST signing key for RPST) but use different issuer/subject strings
  # because OCI enforces one Identity Propagation Trust per issuer per Identity Domain.
  newrelic_config = {
    US = {
      issuer_name       = "newrelic-oci-us-production-issuer"
      subject_name      = "newrelic-oci-us-production-user"
      rpst_issuer_name  = "newrelic-oci-us-production-rpst-issuer"
      rpst_subject_name = "newrelic-oci-us-production-rpst-user"
      public_jwks_url   = "https://publickeys.newrelic.com/r/oci-cmp/us/c5623ba5-1cc7-491a-8ec3-eeee809374f7/jwks.json"
    }
    EU = {
      issuer_name       = "newrelic-oci-eu-production-issuer"
      subject_name      = "newrelic-oci-eu-production-user"
      rpst_issuer_name  = "newrelic-oci-eu-production-rpst-issuer"
      rpst_subject_name = "newrelic-oci-eu-production-rpst-user"
      public_jwks_url   = "https://publickeys.eu.newrelic.com/r/oci-cmp/eu/f923dba9-84a8-491c-b714-6c0e61b90c5b/jwks.json"
    }
    JP = {
      issuer_name       = "newrelic-oci-jp-production-issuer"
      subject_name      = "newrelic-oci-jp-production-user"
      rpst_issuer_name  = "newrelic-oci-jp-production-rpst-issuer"
      rpst_subject_name = "newrelic-oci-jp-production-rpst-user"
      public_jwks_url   = "https://publickeys.jp.newrelic.com/r/oci-cmp/jp/89625529-5f3e-47b8-a13a-25229f85989d/jwks.json"
    }
  }[var.newrelic_region]

  # NR-562518: static value pasted into the OCI trust as `impersonatingResource` and sent by NR
  # as `res_type` during token exchange. Same string for every NR-OCI integration.
  impersonating_resource = "newrelic-integration"

  # Common resource naming
  resource_prefix = var.resource_prefix != "" ? var.resource_prefix : "newrelic"
}
