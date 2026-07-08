---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_gcp_dm_link_account"
sidebar_current: "docs-newrelic-resource-cloud-gcp-dm-link-account"
description: |-
  Link a GCP project to New Relic using Workload Identity Federation (keyless authentication).
---

# Resource: newrelic\_cloud\_gcp\_dm\_link\_account

Use this resource to link a GCP project to New Relic using the **GCP Dimensional Metrics** integration. Unlike the classic `newrelic_cloud_gcp_link_account`, this resource authenticates via [Workload Identity Federation (WIF)](https://cloud.google.com/iam/docs/workload-identity-federation) — no long-lived service account key file is created or managed.

## Prerequisites

Before applying this resource you must create the following GCP infrastructure:

1. A **Workload Identity Pool** with New Relic's OIDC issuer URI as the provider.
2. An **OIDC provider** inside the pool with `allowed_audiences = ["newrelic-gcp-integrations"]` and an `attribute_condition` restricting tokens to your New Relic account ID.
3. A **GCP service account** granted the four required roles:
   - `roles/monitoring.viewer` — metrics collection
   - `roles/serviceusage.serviceUsageConsumer` — API quota
   - `roles/cloudasset.viewer` — resource discovery
   - `roles/resourcemanager.folderViewer` — folder-level resource discovery (**must be granted at the folder level**, not the project level)
4. A `roles/iam.workloadIdentityUser` binding on the service account scoped to the WIF pool using the `attribute.nr_account_id` attribute.

See the [full GCP Dimensional Metrics guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#gcp-dimensional-metrics) for a complete Terraform configuration that creates all of the above.

## Example Usage

```hcl
resource "newrelic_cloud_gcp_dm_link_account" "example" {
  account_id = var.newrelic_account_id
  name       = "my-gcp-project"
  project_id = "my-gcp-project-id"

  # Constructed from the WIF pool provider resource:
  #   "//iam.googleapis.com/${google_iam_workload_identity_pool_provider.newrelic.name}"
  audience = "//iam.googleapis.com/projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/POOL_ID/providers/PROVIDER_ID"

  service_account_email = "newrelic-integration@my-gcp-project-id.iam.gserviceaccount.com"
}
```

## Argument Reference

The following arguments are supported:

- `account_id` - (Optional) The New Relic account ID to link the GCP project to. Defaults to the `account_id` set on the provider.
- `name` - (Required) Display name for this linked account in the New Relic UI.
- `project_id` - (Required) The GCP project ID to monitor (e.g. `my-project-123`).
- `audience` - (Required, Write-only) The WIF audience string that New Relic presents when requesting a GCP access token. Format: `//iam.googleapis.com/projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/POOL_ID/providers/PROVIDER_ID`. Typically set to `"//iam.googleapis.com/${google_iam_workload_identity_pool_provider.newrelic.name}"`.
- `service_account_email` - (Required, Write-only) Email address of the GCP service account that New Relic impersonates via WIF.

-> **NOTE:** `audience` and `service_account_email` are **write-only, ForceNew** fields. They are used internally to construct the WIF credential JSON and are never returned by the New Relic API. If you import an existing linked account, run `terraform apply` after importing to reconcile these fields — Terraform will destroy and recreate the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the linked GCP account in New Relic.

## Import

Linked GCP DM accounts can be imported using the linked account `id`:

```bash
$ terraform import newrelic_cloud_gcp_dm_link_account.example <id>
```
