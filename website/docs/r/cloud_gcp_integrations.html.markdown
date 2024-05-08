---
layout: "newrelic"
page_title: "New Relic: newrelic_cloud_gcp_integrations"
sidebar_current: "docs-newrelic-resource-cloud-gcp-integrations"
description: |-
Integrate GCP services with New Relic.
---

# Resource: newrelic\_cloud\_gcp\_integrations

Use this resource to integrate GCP services with New Relic.

## Prerequisite

Setup is required for this resource to work properly. This resource assumes you have [linked a GCP account](cloud_gcp_link_account.html) to New Relic and configured it to pull metrics from GCP.

New Relic doesn't automatically receive metrics from GCP services, so this resource can be used to configure integrations to those services.

## Example Usage

Leave an integration block empty to use its default configuration. You can also use the [full example, including the GCP set up, found in our guides](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/cloud_integrations_guide#gcp).

```hcl
resource "newrelic_cloud_gcp_link_account" "foo" {
  name = "example"
  project_id="<Your GCP project ID>"
}

resource "newrelic_cloud_gcp_integrations" "foo1" {
  linked_account_id = newrelic_cloud_gcp_link_account.foo.id
  app_engine {
    metrics_polling_interval = 300
  }
  big_query {
    metrics_polling_interval = 300
    fetch_tags = true
  }
  big_table {
    metrics_polling_interval = 300
  }
  composer {
    metrics_polling_interval = 300
  }
  data_flow {
    metrics_polling_interval = 300
  }
  data_proc {
    metrics_polling_interval = 300
  }
  data_store {
    metrics_polling_interval = 300
  }
  fire_base_database {
    metrics_polling_interval = 300
  }
  fire_base_hosting {
    metrics_polling_interval = 300
  }
  fire_base_storage {
    metrics_polling_interval = 300
  }
  fire_store {
    metrics_polling_interval = 300
  }
  functions {
    metrics_polling_interval = 300
  }
  interconnect {
    metrics_polling_interval = 300
  }
  kubernetes {
    metrics_polling_interval = 300
  }
  load_balancing {
    metrics_polling_interval = 300
  }
  mem_cache {
    metrics_polling_interval = 300
  }
  pub_sub {
    metrics_polling_interval = 300
    fetch_tags=true
  }
  redis {
    metrics_polling_interval = 300
  }
  router {
    metrics_polling_interval = 300
  }
  run {
    metrics_polling_interval = 300
  }
  spanner {
    metrics_polling_interval = 300
    fetch_tags=true
  }
  sql {
    metrics_polling_interval = 300
  }
  storage {
    metrics_polling_interval = 300
    fetch_tags=true
  }
  virtual_machines {
    metrics_polling_interval = 300
  }
  vpc_access {
    metrics_polling_interval = 300
  }
}
```
## Argument Reference

-> **WARNING:** Starting with [v3.27.2](https://registry.terraform.io/providers/newrelic/newrelic/3.27.2) of the New Relic Terraform Provider, updating the `linked_account_id` of a `newrelic_cloud_gcp_integrations` resource that has been applied would **force a replacement** of the resource (destruction of the resource, followed by the creation of a new resource). When such an update is performed, please carefully review the output of `terraform plan`, which would clearly indicate a replacement of this resource, before performing a `terraform apply`.

* `account_id` - (Optional) The New Relic account ID to operate on.  This allows the user to override the `account_id` attribute set on the provider. Defaults to the environment variable `NEW_RELIC_ACCOUNT_ID`.
* `linked_account_id` - (Required) The ID of the linked GCP account in New Relic.

The following supported arguments are with minimum metrics polling interval of 300 seconds

* `alloy_db` - (Optional) Alloy DB integration. See [Integration blocks](#integration-blocks) below for details.
* `app_engine` - (Optional) App Engine integration. See [Integration blocks](#integration-blocks) below for details.
* `big_query` - (Optional) Biq Query integration. See [Integration blocks](#integration-blocks) below for details.
* `big_table` - (Optional) Big Table. See [Integration blocks](#integration-blocks) below for details.
* `composer` - (Optional) Composer integration. See [Integration blocks](#integration-blocks) below for details.
* `data_flow` - (Optional) Data Flow integration. See [Integration blocks](#integration-blocks) below for details.
* `data_proc` - (Optional) Data Proc integration. See [Integration blocks](#integration-blocks) below for details.
* `data_store` - (Optional) Data Store integration. See [Integration blocks](#integration-blocks) below for details.
* `fire_base_database` - (Optional) Fire Base Database integration. See [Integration blocks](#integration-blocks) below for details.
* `fire_base_hosting` - (Optional) Fire Base Hosting integration. See [Integration blocks](#integration-blocks) below for details.
* `fire_base_storage` - (Optional) Fire Base Storage integration. See [Integration blocks](#integration-blocks) below for details.
* `fire_store` - (Optional) Fire Store integration. See [Integration blocks](#integration-blocks) below for details.
* `functions` - (Optional) Functions integration. See [Integration blocks](#integration-blocks) below for details.
* `interconnect` - (Optional) Interconnect integration. See [Integration blocks](#integration-blocks) below for details.
* `kubernetes` - (Optional) Kubernetes integration. See [Integration blocks](#integration-blocks) below for details.
* `load_balancing` - (Optional) Load Balancing integration. See [Integration blocks](#integration-blocks) below for details.
* `mem_cache` - (Optional) Mem cache integration. See [Integration blocks](#integration-blocks) below for details.
* `pub_sub` - (Optional) Pub/Sub integration. See [Integration blocks](#integration-blocks) below for details.
* `redis` - (Optional) Redis integration. See [Integration blocks](#integration-blocks) below for details.
* `router` - (Optional) Router integration. See [Integration blocks](#integration-blocks) below for details.
* `run` - (Optional) Run integration. See [Integration blocks](#integration-blocks) below for details.
* `spanner` - (Optional) Spanner integration. See [Integration blocks](#integration-blocks) below for details.
* `sql` - (Optional) SQL integration. See [Integration blocks](#integration-blocks) below for details.
* `storage` - (Optional) Storage integration. See [Integration blocks](#integration-blocks) below for details.
* `virtual_machines` - (Optional) Virtual machines integration. See [Integration blocks](#integration-blocks) below for details.
* `vpc_access` - (Optional) VPC Access integration. See [Integration blocks](#integration-blocks) below for details.

### `Integration` blocks

All `integration` blocks support the following common arguments:

* `metrics_polling_interval` - (Optional) The data polling interval **in seconds**.

-> **NOTE** For more information on the ranges of metric polling intervals of each of these integrations, head over to [this page](https://docs.newrelic.com/docs/infrastructure/google-cloud-platform-integrations/get-started/introduction-google-cloud-platform-integrations/)

Other integration supports an additional argument:

* `big_query`
* `pub_sub`
* `spanner`
* `storage`
    * `fetch_tags` - (Optional) Specify if labels and the extended inventory should be collected. May affect total data collection time and contribute to the Cloud provider API rate limit.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the GCP linked account.

## Import

Linked GCP account integrations can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_cloud_gcp_integrations.foo <id>
```
