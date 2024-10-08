---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_broken_links_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-broken-links-monitor"
description: |-
    Create and manage a Synthetics Broken Links monitor in New Relic.
---

# Resource: newrelic\_synthetics\_broken\_links\_monitor

-> **WARNING:** Support for using the Synthetics Legacy Runtime (deprecated) with **new** Synthetic monitors **has officially ended as of August 26, 2024**. As a consequence, starting with v3.43.0 of the New Relic Terraform Provider, **new** Synthetic monitors **will no longer be allowed to use the legacy runtime** (this applies to all Synthetic monitor resources). Additionally, as previously communicated by New Relic, the Synthetics Legacy Runtime **will reach its end-of-life (EOL) on October 22, 2024**. In light of the above, we kindly recommend that you upgrade your Synthetic Monitors to the new runtime at the earliest, if they are still using the legacy runtime. Please check out [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/synthetics_legacy_runtime_eol_migration_guide) in the documentation of the Terraform Provider and [this announcement](https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm) for more details on the EOL, action needed, relevant resources, and more.

Use this resource to create, update, and delete a Synthetics Broken Links monitor in New Relic.

## Example Usage

```hcl
resource "newrelic_synthetics_broken_links_monitor" "foo" {
  name                 = "Sample Broken Links Monitor"
  uri                  = "https://www.one.example.com"
  locations_public     = ["AP_SOUTH_1"]
  period               = "EVERY_6_HOURS"
  status               = "ENABLED"
  runtime_type         = "NODE_API"
  runtime_type_version = "16.10"
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following are the arguments supported by this resource.

* `account_id`- (Optional) The account in which the Synthetics monitor will be created.
* `name` - (Required) The name for the monitor.
* `uri` - (Required) The URI the monitor runs against.
* `locations_public` - (Required) The location the monitor will run from. Check out [this page](https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/) for a list of valid public locations. You don't need the `AWS_` prefix as the provider uses NerdGraph. At least one of either `locations_public` or `location_private` is required.
* `locations_private` - (Required) The location the monitor will run from. Accepts a list of private location GUIDs. At least one of either `locations_public` or `locations_private` is required.
* `period` - (Required) The interval at which this monitor should run. Valid values are `EVERY_MINUTE`, `EVERY_5_MINUTES`, `EVERY_10_MINUTES`, `EVERY_15_MINUTES`, `EVERY_30_MINUTES`, `EVERY_HOUR`, `EVERY_6_HOURS`, `EVERY_12_HOURS`, or `EVERY_DAY`.
* `status` - (Required) The run state of the monitor. (`ENABLED` or `DISABLED`). 

-> **WARNING:** As of February 29, 2024, Synthetic Monitors no longer support the `MUTED` status. Version **3.33.0** of the New Relic Terraform Provider is released to coincide with the `MUTED` status end-of-life. Consequently, the only valid values for `status` for all types of Synthetic Monitors are mentioned above. For additional information on alternatives to the `MUTED` status of Synthetic Monitors that can be managed via Terraform, please refer to [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/upcoming_synthetics_muted_status_eol_guide).
* `runtime_type` - (Optional) The runtime that the monitor will use to run jobs.
* `runtime_type_version` - (Optional) The specific version of the runtime type selected.

-> **NOTE:** Currently, the values of `runtime_type` and `runtime_type_version` supported by this resource are `NODE_API` and `16.10` respectively. In order to run the monitor in the new runtime, both `runtime_type` and `runtime_type_version` need to be specified; however, specifying neither of these attributes would set this monitor to use the legacy runtime. It may also be noted that the runtime opted for would only be effective with private locations. For public locations, all traffic has been shifted to the new runtime, irrespective of the selection made.

-> **WARNING:** Support for using the Synthetics Legacy Runtime (deprecated) with **new** Synthetic monitors **has officially ended as of August 26, 2024**. As a consequence, starting with v3.43.0 of the New Relic Terraform Provider, **new** Synthetic monitors **will no longer be allowed to use the legacy runtime** (this applies to all Synthetic monitor resources). Additionally, as previously communicated by New Relic, the Synthetics Legacy Runtime **will reach its end-of-life (EOL) on October 22, 2024**. In light of the above, we kindly recommend that you upgrade your Synthetic Monitors to the new runtime at the earliest, if they are still using the legacy runtime. Please check out [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/synthetics_legacy_runtime_eol_migration_guide) in the documentation of the Terraform Provider and [this announcement](https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm) for more details on the EOL, action needed, relevant resources, and more.

* `tag` - (Optional) The tags that will be associated with the monitor. See [Nested tag blocks](#nested-tag-blocks) below for details

### Nested `tag` blocks

All nested `tag` blocks support the following common arguments:

* `key` - (Required) Name of the tag key.
* `values` - (Required) Values associated with the tag key.

## Additional Examples

### Create a monitor with a private location

The below example shows how you can define a private location and attach it to a monitor.

-> **NOTE:** It can take up to 10 minutes for a private location to become available.

```hcl
resource "newrelic_synthetics_private_location" "foo" {
  name                      = "Sample Private Location"
  description               = "Sample Private Location Description"
  verified_script_execution = false
}

resource "newrelic_synthetics_broken_links_monitor" "foo" {
  name              = "Sample Broken Links Monitor"
  uri               = "https://www.one.example.com"
  locations_private = [newrelic_synthetics_private_location.foo.id]
  period            = "EVERY_6_HOURS"
  status            = "ENABLED"
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The ID (GUID) of the synthetics broken links monitor.
* `period_in_minutes` - The interval in minutes at which Synthetic monitor should run.

## Import

A broken links monitor can be imported using its GUID, using the following command.

```bash
$ terraform import newrelic_synthetics_broken_links_monitor.monitor <guid>
```