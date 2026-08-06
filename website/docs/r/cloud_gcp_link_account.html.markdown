---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_gcp_link_account"
sidebar_current: "docs-new-relic-resource-cloud-gcp-link-account"
description: |-
Link a GCP account to New Relic.
---

# Resource: newrelic_cloud_gcp_link_account

Use this resource to link a GCP account to New Relic.

## Prerequisite

To start receiving Google Cloud Platform (GCP) data with New Relic GCP integrations, connect your Google project to New Relic infrastructure monitoring. If you don't have one already, create a New Relic account. It's free, forever.

Setup is required in GCP for this resource to work properly. The New Relic GCP integration can be done by creating a user account or a service account.

A user with Project IAM Admin role is needed to add the service account ID as a member in your GCP project.

In the GCP project IAM & admin, the service account must have the Project Viewer role and the Service Usage Consumer role or, alternatively, a custom role.

Follow the [steps outlined here](https://docs.newrelic.com/docs/infrastructure/google-cloud-platform-integrations/get-started/connect-google-cloud-platform-services-new-relic) to set up the integration.

## Example Usage

You can also use the [full example, including the GCP set up, found in our guides](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#gcp).

```hcl
resource "newrelic_cloud_gcp_link_account" "foo" {
  account_id = "account id of newrelic account"
  project_id = "id of the Project"
  name  = "account name"
}
```

### GCP Dimensional Metrics (keyless / WIF) linking

To link a GCP project for **GCP Dimensional Metrics** using keyless authentication via Workload Identity Federation (WIF) instead of a service-account key, set `use_workload_identity_federation = true` and provide `audience` and `service_account_email`. When enabled, the resource authenticates via WIF and links the project as a Dimensional Metrics account. Use this linked account with the [`newrelic_cloud_gcp_dm_integrations`](cloud_gcp_dm_integrations.html) resource.

```hcl
resource "newrelic_cloud_gcp_link_account" "dm" {
  account_id                       = "account id of newrelic account"
  name                             = "account name"
  project_id                       = "id of the Project"
  use_workload_identity_federation = true
  audience                         = "//iam.googleapis.com/projects/123456789/locations/global/workloadIdentityPools/newrelic-wif-pool/providers/newrelic-oidc-provider"
  service_account_email            = "newrelic-integration@my-project.iam.gserviceaccount.com"
}
```

## Argument Reference

The following arguments are supported:

- `account_id` - (Optional) - Account ID of the New Relic account.
- `project_id` - (Required) - Project ID of the GCP account.
- `name` - (Required) - The name of the GCP account in New Relic.
- `use_workload_identity_federation` - (Optional) - Set to `true` to link the GCP account for **GCP Dimensional Metrics** using keyless Workload Identity Federation (WIF) instead of a service-account key. When `true`, `audience` and `service_account_email` are required. Defaults to `false` (legacy service-account-key linking).
- `audience` - (Optional) - The Workload Identity Federation pool provider audience URI. Format: `//iam.googleapis.com/projects/{PROJECT_NUMBER}/locations/global/workloadIdentityPools/{POOL_ID}/providers/{PROVIDER_ID}`. Required when `use_workload_identity_federation = true`.
- `service_account_email` - (Optional) - The GCP service account email New Relic impersonates to collect metrics when linking via WIF. The service account must grant the WIF pool the `roles/iam.workloadIdentityUser` binding. Required when `use_workload_identity_federation = true`.

-> **NOTE:** `audience` and `service_account_email` are write-only, `ForceNew` fields used to construct the WIF credential internally; they are never returned by the API and are retained from state. When importing a WIF-linked account, also set `use_workload_identity_federation = true` in your configuration (it is not returned by the API and defaults to `false`), and add `audience` and `service_account_email` to `ImportStateVerifyIgnore` (or run `terraform apply` afterwards to reconcile them).

-> **WARNING:** Starting with [v3.27.2](https://registry.terraform.io/providers/newrelic/newrelic/3.27.2) of the New Relic Terraform Provider, updating any of the aforementioned attributes (except `name`) of a `newrelic_cloud_gcp_link_account` resource that has been applied would **force a replacement** of the resource (destruction of the resource, followed by the creation of a new resource). Please carefully review the output of `terraform plan`, which would clearly indicate a replacement of this resource, before performing a `terraform apply`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The id of the GCP linked account.

## Import

Linked GCP accounts can be imported using `id`, you can find the `id` of an existing GCP linked accounts in GCP dashboard under Infrastructure in Newrelic Console.

```bash

  $  terraform import newrelic_cloud_gcp_link_account.foo <id>

```
