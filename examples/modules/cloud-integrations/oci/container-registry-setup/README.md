# OCI Container Registry Setup Module

Creates the **private Container Registry repository** the log-forwarder function's image is mirrored into, as a one-time setup step -- the same shape as [`wif-setup`](../wif-setup).

## Why this is a separate module

OCIR repository names are unique **tenancy+region-wide**, not scoped per compartment. `logs-integration` originally created this repository itself on every apply; a second deployment in the same tenancy+region (a different compartment, a repeat test run, a second New Relic account onboarding the same tenancy) collided on the same default name with a raw `409-NAMESPACE_CONFLICT` (see [NR-640865](https://new-relic.atlassian.net/browse/NR-640865)).

An automatic "check if it exists, reuse if so" fix inside `logs-integration` itself isn't safely possible in Terraform: a resource's own `precondition` can't reference that same resource ("may not refer to itself"), and no other signal (compartment, tags) can reliably tell "the repository this exact deployment already created" apart from "an unrelated repository that happens to have the same name" -- both look identical by the time a second `apply` runs.

Pulling the repository out into its own one-time module removes the problem instead of working around it: nothing in `logs-integration` (or `metrics-integration`, if adopted there too) ever tries to create this resource -- every consumer takes this module's `container_repository_id` output as a **required** input.

**Run this once per region** you deploy the logs integration into, not once per tenancy -- unlike `wif-setup` (tenancy-wide, home-region only), OCIR's uniqueness scope is per-region.

## Where this module fits in the full integration

| Order | Component | What it provisions | Required? |
|-------|-----------|--------------------|-----------|
| 1 | [`wif-setup`](../wif-setup) | OAuth apps, IAM, identity propagation trust — the auth foundation | Yes |
| 2 | [`policy-setup`](../policy-setup) | IAM policies + (optionally) Vault secrets for NR keys | Yes |
| 2b | `container-registry-setup` (this module) | Private Container Registry repository for the log-forwarder image | Yes, once per region, if deploying `logs-integration` |
| 3 | [`newrelic_cloud_oci_link_account`](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/cloud_oci_link_account) | Registers the OCI tenancy with New Relic | Yes |
| 4b | [`logs-integration`](../logs-integration) | Service Connector Hub + Function for logs shipping | Optional (need logs) |

## Inputs

| Variable | Type | Default | Required | Notes |
|----------|------|---------|----------|-------|
| `tenancy_ocid` | string | — | yes | Your OCI tenancy OCID |
| `region` | string | — | yes | OCI region to create the repository in. Run this module again for each additional region |
| `fingerprint` | string | — | yes | OCI API key fingerprint |
| `private_key` | string (sensitive) | — | yes | OCI API key PEM content |
| `compartment_ocid` | string | — | yes | Compartment to create the repository in -- the compartment already created by `policy-setup` |
| `newrelic_logging_identifier` | string | `"logs"` | no | Must match the value passed to `logs-integration` -- determines the repository's name |

## Outputs

| Output | Description |
|--------|-------------|
| `container_repository_id` | OCID of the created repository. Pass to `logs-integration`'s `container_repository_id` variable |
| `container_repository_namespace` | OCIR namespace the repository was created under |

## Example

```hcl
module "oci_container_registry_setup" {
  source = "github.com/newrelic/terraform-provider-newrelic//examples/modules/cloud-integrations/oci/container-registry-setup"

  tenancy_ocid     = "ocid1.tenancy.oc1..aaaaaaaaexampletenancy"
  region           = "us-ashburn-1"
  fingerprint      = "12:34:56:78:9a:bc:de:f0:12:34:56:78:9a:bc:de:f0"
  private_key      = "USER_PVT_KEY" # or file("/path/to/private_key.pem")
  compartment_ocid = "ocid1.compartment.oc1..bbbbbbbbexamplecmp" # or module.oci_policy_setup.compartment_ocid
}

# Pass straight into logs-integration:
# module.oci_logs_integration's container_repository_id = module.oci_container_registry_setup.container_repository_id
```

## Migrating an existing deployment

If `logs-integration` already created its own repository (any deployment before this module existed), moving to this module needs `terraform import`, not a fresh `apply`:

1. Apply this module with `terraform import oci_artifacts_container_repository.log_forwarder_repo <existing_repository_ocid>` to bring the existing repository into **this** module's state.
2. In the existing `logs-integration` deployment, run `terraform state rm oci_artifacts_container_repository.log_forwarder_repo` -- **not** `terraform destroy` -- so Terraform stops trying to manage it there.
3. Upgrade `logs-integration` to a version that takes `container_repository_id` as an input, set it to the same OCID, and apply.

Applying step 3 before steps 1-2 will make Terraform plan to destroy the repository the function is actively using.
