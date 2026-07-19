---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_gcp_dm_integrations"
sidebar_current: "docs-newrelic-resource-cloud-gcp-dm-integrations"
description: |-
  Configure which GCP services New Relic polls for the GCP Dimensional Metrics integration.
---

# Resource: newrelic\_cloud\_gcp\_dm\_integrations

Use this resource to configure which GCP services New Relic polls as part of the **GCP Dimensional Metrics** integration. Each service is enabled by adding the corresponding block; omit a block to disable polling for that service.

## Prerequisite

This resource requires a linked GCP account created with [`newrelic_cloud_gcp_link_account`](cloud_gcp_link_account.html) using its keyless (WIF) mode — set `audience` and `service_account_email` on that resource. See the [full GCP Dimensional Metrics guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#gcp-dimensional-metrics) for complete setup instructions including the required GCP Workload Identity Federation infrastructure.

## Example Usage

```hcl
resource "newrelic_cloud_gcp_link_account" "example" {
  account_id            = var.newrelic_account_id
  name                  = "my-gcp-project"
  project_id            = "my-gcp-project-id"
  audience              = "//iam.googleapis.com/projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/POOL_ID/providers/PROVIDER_ID"
  service_account_email = "newrelic-integration@my-gcp-project-id.iam.gserviceaccount.com"
}

resource "newrelic_cloud_gcp_dm_integrations" "example" {
  account_id        = var.newrelic_account_id
  linked_account_id = newrelic_cloud_gcp_link_account.example.id

  # Standard services — 300 s minimum polling interval
  ai_platform      { metrics_polling_interval = 300 }
  api_gateway      { metrics_polling_interval = 300 }
  app_engine       { metrics_polling_interval = 300 }
  big_table        { metrics_polling_interval = 300 }
  composer         { metrics_polling_interval = 300 }
  data_store       { metrics_polling_interval = 300 }
  firebase_auth    { metrics_polling_interval = 300 }
  firebase_database { metrics_polling_interval = 300 }
  firebase_hosting { metrics_polling_interval = 300 }
  firebase_storage { metrics_polling_interval = 300 }
  firestore        { metrics_polling_interval = 300 }
  functions        { metrics_polling_interval = 300 }
  interconnect     { metrics_polling_interval = 300 }
  istio            { metrics_polling_interval = 300 }
  kubernetes       { metrics_polling_interval = 300 }
  mem_cache        { metrics_polling_interval = 300 }
  memory_store     { metrics_polling_interval = 300 }
  redis            { metrics_polling_interval = 300 }
  router           { metrics_polling_interval = 300 }
  run              { metrics_polling_interval = 300 }
  sql              { metrics_polling_interval = 300 }
  virtual_machines { metrics_polling_interval = 300 }
  vpc_access       { metrics_polling_interval = 300 }
  firebase_app_hosting { metrics_polling_interval = 300 }
  firebase_vertex_ai   { metrics_polling_interval = 300 }

  # Services where 1-minute polling is in Limited Preview (LP)
  alloy_db       { metrics_polling_interval = 60 }
  big_query      { metrics_polling_interval = 60 }
  data_flow      { metrics_polling_interval = 60 }
  data_proc      { metrics_polling_interval = 60 }
  load_balancing { metrics_polling_interval = 60 }
  managed_kafka  { metrics_polling_interval = 60 }
  pub_sub        { metrics_polling_interval = 60 }
  spanner        { metrics_polling_interval = 60 }
  storage        { metrics_polling_interval = 300 }
}
```

## Argument Reference

-> **WARNING:** Updating `linked_account_id` on an existing resource will **force a replacement** of the resource (destroy + create). Review `terraform plan` carefully before applying.

* `account_id` - (Optional) The New Relic account ID to operate on. Defaults to the `account_id` set on the provider.
* `linked_account_id` - (Required) The ID of the linked GCP account created by `newrelic_cloud_gcp_link_account` (in keyless/WIF mode).

### Standard services (300 s minimum `metrics_polling_interval`)

* `ai_platform` - (Optional) AI Platform integration. See [Integration blocks](#integration-blocks) below.
* `api_gateway` - (Optional) API Gateway integration (DM only). See [Integration blocks](#integration-blocks) below.
* `app_engine` - (Optional) App Engine integration. See [Integration blocks](#integration-blocks) below.
* `big_table` - (Optional) Bigtable integration. See [Integration blocks](#integration-blocks) below.
* `composer` - (Optional) Cloud Composer integration. See [Integration blocks](#integration-blocks) below.
* `data_store` - (Optional) Datastore integration. See [Integration blocks](#integration-blocks) below.
* `firebase_auth` - (Optional) Firebase Authentication integration (DM only). See [Integration blocks](#integration-blocks) below.
* `firebase_database` - (Optional) Firebase Realtime Database integration. See [Integration blocks](#integration-blocks) below.
* `firebase_hosting` - (Optional) Firebase Hosting integration. See [Integration blocks](#integration-blocks) below.
* `firebase_storage` - (Optional) Firebase Storage integration. See [Integration blocks](#integration-blocks) below.
* `firestore` - (Optional) Firestore integration. See [Integration blocks](#integration-blocks) below.
* `functions` - (Optional) Cloud Functions integration. See [Integration blocks](#integration-blocks) below.
* `interconnect` - (Optional) Cloud Interconnect integration. See [Integration blocks](#integration-blocks) below.
* `istio` - (Optional) Istio integration (DM only, metrics only — no entity support). See [Integration blocks](#integration-blocks) below.
* `kubernetes` - (Optional) Kubernetes Engine integration (metrics only — no entity support). See [Integration blocks](#integration-blocks) below.
* `mem_cache` - (Optional) Memcache integration. See [Integration blocks](#integration-blocks) below.
* `memory_store` - (Optional) Memorystore integration (DM only). See [Integration blocks](#integration-blocks) below.
* `redis` - (Optional) Memorystore for Redis integration. See [Integration blocks](#integration-blocks) below.
* `router` - (Optional) Cloud Router integration. See [Integration blocks](#integration-blocks) below.
* `run` - (Optional) Cloud Run integration. See [Integration blocks](#integration-blocks) below.
* `sql` - (Optional) Cloud SQL integration. See [Integration blocks](#integration-blocks) below.
* `storage` - (Optional) Cloud Storage integration. See [Integration blocks](#integration-blocks) below.
* `virtual_machines` - (Optional) Compute Engine (virtual machines) integration. See [Integration blocks](#integration-blocks) below.
* `vpc_access` - (Optional) VPC Access integration. See [Integration blocks](#integration-blocks) below.
* `firebase_app_hosting` - (Optional) Firebase App Hosting integration (DM only, metrics only — no entity support). See [Integration blocks](#integration-blocks) below.
* `firebase_vertex_ai` - (Optional) Firebase Vertex AI integration (DM only, metrics only — no entity support). See [Integration blocks](#integration-blocks) below.

### Services with Limited Preview (LP) 1-minute polling

The following services support a polling floor as low as **60 seconds**. 1-minute polling intervals for these services are in **Limited Preview**; all services are fully available at the standard 300-second interval. Set `metrics_polling_interval = 60` to enable lower-latency polling.

* `alloy_db` - (Optional) AlloyDB integration. See [Integration blocks](#integration-blocks) below.
* `big_query` - (Optional) BigQuery integration. See [Integration blocks](#integration-blocks) below.
* `data_flow` - (Optional) Dataflow integration. See [Integration blocks](#integration-blocks) below.
* `data_proc` - (Optional) Dataproc integration. See [Integration blocks](#integration-blocks) below.
* `load_balancing` - (Optional) Cloud Load Balancing integration. See [Integration blocks](#integration-blocks) below.
* `managed_kafka` - (Optional) Managed Apache Kafka integration (DM only). See [Integration blocks](#integration-blocks) below.
* `pub_sub` - (Optional) Cloud Pub/Sub integration. See [Integration blocks](#integration-blocks) below.
* `spanner` - (Optional) Cloud Spanner integration. See [Integration blocks](#integration-blocks) below.

### Integration blocks

All integration blocks support the following argument:

* `metrics_polling_interval` - (Optional) How often New Relic polls the service for metrics, **in seconds**. Minimum values: **60 s** for services where 1-minute polling is in Limited Preview (`alloy_db`, `big_query`, `data_flow`, `data_proc`, `load_balancing`, `managed_kafka`, `pub_sub`, `spanner`); **300 s** for all other services. All services accept the standard **300 s** interval.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the GCP DM integrations resource.

## Import

GCP DM integrations can be imported using the `id`:

```bash
$ terraform import newrelic_cloud_gcp_dm_integrations.example <id>
```
