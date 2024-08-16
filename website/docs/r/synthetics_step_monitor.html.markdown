---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_step_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-step-monitor"
description: |-
    Create and manage a Synthetics Step monitor in New Relic.
---

# Resource: newrelic\_synthetics\_step\_monitor

-> **WARNING** Support for legacy Synthetics runtimes **will reach its end-of-life (EOL) on October 22, 2024**. In addition, creating **_new_** monitors using the legacy runtime **will no longer be supported after August 26, 2024**. In light of the above, kindly **upgrade your Synthetic Monitors to the new runtime** at the earliest, if they are still using the legacy runtime. Please check out [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/synthetics_legacy_runtime_eol_migration_guide) in the documentation of the Terraform Provider and [this announcement](https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm) for more details on the EOL, action needed, relevant resources, and more.

Use this resource to create, update, and delete a Synthetics Step monitor in New Relic.

## Example Usage

```hcl
resource "newrelic_synthetics_step_monitor" "foo" {
  name                                    = "Sample Step Monitor"
  enable_screenshot_on_failure_and_script = true
  locations_public                        = ["US_EAST_1", "US_EAST_2"]
  period                                  = "EVERY_6_HOURS"
  status                                  = "ENABLED"
  runtime_type                            = "CHROME_BROWSER"
  runtime_type_version                    = "100"
  steps {
    ordinal = 0
    type    = "NAVIGATE"
    values  = ["https://www.newrelic.com"]
  }
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following are the common arguments supported for `STEP` monitor:

* `account_id`- (Optional) The account in which the Synthetics monitor will be created.
* `name` - (Required) The name for the monitor.
* `locations_public` - (Required) The location the monitor will run from. Check out [this page](https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/) for a list of valid public locations. You don't need the `AWS_` prefix as the provider uses NerdGraph. At least one of either `locations_public` or `location_private` is required.
* `location_private` - (Required) The location the monitor will run from. At least one of `locations_public` or `location_private` is required. See [Nested locations_private blocks](#nested-locations-private-blocks) below for details.
* `period` - (Required) The interval at which this monitor should run. Valid values are `EVERY_MINUTE`, `EVERY_5_MINUTES`, `EVERY_10_MINUTES`, `EVERY_15_MINUTES`, `EVERY_30_MINUTES`, `EVERY_HOUR`, `EVERY_6_HOURS`, `EVERY_12_HOURS`, or `EVERY_DAY`.
* `status` - (Required) The run state of the monitor. (`ENABLED` or `DISABLED`).

-> **WARNING:** As of February 29, 2024, Synthetic Monitors no longer support the `MUTED` status. Version **3.33.0** of the New Relic Terraform Provider is released to coincide with the `MUTED` status end-of-life. Consequently, the only valid values for `status` for all types of Synthetic Monitors are mentioned above. For additional information on alternatives to the `MUTED` status of Synthetic Monitors that can be managed via Terraform, please refer to [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/upcoming_synthetics_muted_status_eol_guide).

* `runtime_type` - (Optional) The runtime that the monitor will use to run jobs.
* `runtime_type_version` - (Optional) The specific version of the runtime type selected.

-> **NOTE:** Currently, the values of `runtime_type` and `runtime_type_version` supported by this resource are `CHROME_BROWSER` and `100` respectively. In order to run the monitor in the new runtime, both `runtime_type` and `runtime_type_version` need to be specified; however, specifying neither of these attributes would set this monitor to use the legacy runtime. It may also be noted that the runtime opted for would only be effective with private locations. For public locations, all traffic has been shifted to the new runtime, irrespective of the selection made.

-> **WARNING** Support for legacy Synthetics runtimes **will reach its end-of-life (EOL) on October 22, 2024**. In addition, creating **_new_** monitors using the legacy runtime **will no longer be supported after August 26, 2024**. In light of the above, kindly **upgrade your Synthetic Monitors to the new runtime** at the earliest, if they are still using the legacy runtime. Please check out [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/synthetics_legacy_runtime_eol_migration_guide) in the documentation of the Terraform Provider and [this announcement](https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm) for more details on the EOL, action needed, relevant resources, and more.

* `steps` - (Required) The steps that make up the script the monitor will run. See [Nested steps blocks](#nested-steps-blocks) below for details.
* `tag` - (Optional) The tags that will be associated with the monitor. See [Nested tag blocks](#nested-tag-blocks) below for details.

### Nested `location private` blocks

All nested `location_private` blocks support the following common arguments:

* `guid` - (Required) The unique identifier for the Synthetics private location in New Relic.
* `vse_password` - (Optional) The location's Verified Script Execution password, only necessary if Verified Script Execution is enabled for the location.

### Nested `steps` blocks

All nested `steps` blocks support the following common arguments:

* `ordinal` - (Required) The position of the step within the script ranging from 0-100.
* `type` - (Required) Name of the tag key. Valid values are `ASSERT_ELEMENT`, `ASSERT_MODAL`, `ASSERT_TEXT`, `ASSERT_TITLE`, `CLICK_ELEMENT`, `DISMISS_MODAL`, `DOUBLE_CLICK_ELEMENT`, `HOVER_ELEMENT`, `NAVIGATE`, `SECURE_TEXT_ENTRY`, `SELECT_ELEMENT`, `TEXT_ENTRY`.
* `values` - (Optional) The metadata values related to the step.

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
  verified_script_execution = true
}

resource "newrelic_synthetics_step_monitor" "foo" {
  name   = "Sample Step Monitor"
  period = "EVERY_6_HOURS"
  status = "ENABLED"
  location_private {
    guid         = newrelic_synthetics_private_location.foo.id
    vse_password = "secret"
  }
  steps {
    ordinal = 0
    type    = "NAVIGATE"
    values  = ["https://google.com"]
  }
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```
## Attributes Reference

The following attributes are exported:

* `id` - The ID (GUID) of the synthetics step monitor.
* `period_in_minutes` - The interval in minutes at which Synthetic monitor should run.

## Import

A step monitor can be imported using its GUID, using the following command.

```bash
$ terraform import newrelic_synthetics_step_monitor.monitor <guid>
```
