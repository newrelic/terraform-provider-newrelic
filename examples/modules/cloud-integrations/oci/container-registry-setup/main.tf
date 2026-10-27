# OCIR repository names are unique tenancy+region-wide, not scoped per compartment -- a
# repository created inside the logs-integration/metrics-integration modules themselves
# collides with any other deployment in the same tenancy+region (NR-640865). Pulling repo
# creation out into its own one-time module -- the same shape as wif-setup -- means nothing
# ever tries to create a second one: every consumer takes this module's output as a required
# input instead of creating its own.
resource "oci_artifacts_container_repository" "log_forwarder_repo" {
  compartment_id = var.compartment_ocid
  display_name   = "newrelic-${var.newrelic_logging_identifier}/oci-log-forwarder"
  is_public      = false
  freeform_tags = {
    newrelic-terraform = "true"
  }
}
