---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_step_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-step-monitor"
description: |-
    Create and manage a Synthetics Step monitor in New Relic.
---

# Resource: newrelic\_synthetics\_step\_monitor

Use this resource to create, update, and delete a Synthetics Step monitor in New Relic.

## Example Usage

```hcl
resource "newrelic_synthetics_step_monitor" "monitor" {
  name                                    = "step_monitor"
  enable_screenshot_on_failure_and_script = true
  locations_public                        = ["US_EAST_1", "US_EAST_2"]
  period                                  = "EVERY_6_HOURS"
  status                                  = "ENABLED"
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
* `uri` - (Required) The uri the monitor runs against.
* `locations_public` - (Required) The location the monitor will run from. Valid public locations are https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/. You don't need the `AWS_` prefix as the provider uses NerdGraph. At least one of either `locations_public` or `location_private` is required.
* `location_private` - (Required) The location the monitor will run from. At least one of `locations_public` or `location_private` is required. See [Nested locations_private blocks](#nested-locations-private-blocks) below for details.
* `period` - (Required) The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.
* `status` - (Required) The run state of the monitor. (`ENABLED` or `DISABLED`).

-> **WARNING:** Starting with version **4.0.0** of the New Relic Terraform Provider, support for the `MUTED` status has been discontinued due to the end-of-life of the `MUTED` status for Synthetic Monitors, which took place on February 29, 2024. Consequently, `MUTED` is no longer a valid and functional value for the `status` argument of all types of Synthetic Monitors. The only valid values for `status` are mentioned above. For additional information on alternatives to the `MUTED` status of Synthetic Monitors that can be managed via Terraform, please refer to [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/upcoming_synthetics_muted_status_eol_guide).

* `steps` - (Required) The steps that make up the script the monitor will run. See [Nested steps blocks](#nested-steps-blocks) below for details.
* `tag` - (Optional) The tags that will be associated with the monitor. See [Nested tag blocks](#nested-tag-blocks) below for details.

### Nested `location private` blocks

All nested `location_private` blocks support the following common arguments:

* `guid` - (Required) The unique identifier for the Synthetics private location in New Relic.
* `vse_password` - (Optional) The location's Verified Script Execution password, only necessary if Verified Script Execution is enabled for the location.

### Nested `steps` blocks

All nested `steps` blocks support the following common arguments:

* `ordinal` - (Required) The position of the step within the script ranging from 0-100.
* `type` - (Required) Name of the tag key. Valid values are ASSERT_ELEMENT, ASSERT_MODAL, ASSERT_TEXT, ASSERT_TITLE, CLICK_ELEMENT, DISMISS_MODAL, DOUBLE_CLICK_ELEMENT, HOVER_ELEMENT, NAVIGATE, SECURE_TEXT_ENTRY, SELECT_ELEMENT, TEXT_ENTRY.
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
resource "newrelic_synthetics_private_location" "location" {
  description               = "Test Description"
  name                      = "private-location"
  verified_script_execution = true
}

resource "newrelic_synthetics_step_monitor" "bar" {
  name = "step_monitor"
  uri  = "https://www.one.example.com"
  location_private {
    guid         = newrelic_synthetics_private_location.location.id
    vse_password = "secret"
  }
  period = "EVERY_6_HOURS"
  status = "ENABLED"
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

Synthetics step monitor scripts can be imported using the `guid`, e.g.

```bash
$ terraform import newrelic_synthetics_step_monitor.monitor <guid>
```
