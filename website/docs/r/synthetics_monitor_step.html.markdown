---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_step_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-step-monitor"
description: |-
Create and manage a Synthetics monitor in New Relic.
---

# Resource: newrelic\_synthetics\_step\_monitor

Use this resource to create, update, and delete the synthetics step monitor in New Relic.

## Example Usage

##### Type: `STEP MONITOR`
```hcl
resource "newrelic_synthetics_step_monitor" "foo" {
  name = "foo"
  enable_screenshot_on_failure_and_script = true
  location_public = ["AWS_US_EAST_1", "AWS_US_EAST_2"]
  period = "EVERY_6_HOURS"
  status = "ENABLED"
  steps {
    ordinal = " "
    types = " "
    values = "ASSERT_ELEMENT"
  }
  tag {
    key = "some_key"
    values = ["some_value"]
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following are the common arguments supported for `STEP` monitor:

* `account_id`- (Required) The account in which the Synthetics monitor will be created.
* `name` - (Required) The name for the monitor.
* `uri` - (Required) The uri the monitor runs against.
* `location_public` - (Required) The location the monitor will run from. Valid public locations are https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/
* `location_private` - (Required) The location the monitor will run from.
* `period` - (Required) The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.
* `status` - (Required) The run state of the monitor. (i.e. `ENABLED`, `DISABLED`, `MUTED`).

### Nested blocks

All nested `location_private` blocks support the following common arguments:

* `guid` - (Required) The unique identifier for the Synthetics private location in New Relic.
* `vse_password` - (Optional) **DEPRECATED:** The location's Verified Script Execution password, Only necessary if Verified Script Execution is enabled for the location.

All nested `steps` blocks support the following common arguments:

* `ordinal` - (Required) The position of the step within the script ranging from 1-100.
* `type` - (Required) Name of the tag key.
* `values` - (Optional) The metadata values related to the step. valid values are ASSERT_ELEMENT, ASSERT_MODAL, ASSERT_TEXT, ASSERT_TITLE, CLICK_ELEMENT, DISMISS_MODAL, DOUBLE_CLICK_ELEMENT, HOVER_ELEMENT, NAVIGATE, SECURE_TEXT_ENTRY, SELECT_ELEMENT, TEXT_ENTRY.

All nested `tag` blocks support the following common arguments:

* `key` - (Required) Name of the tag key.
* `values` - (Required) Values associated with the tag key.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the synthetics step monitor.

## Import

Synthetics step monitor scripts can be imported using the `guid`, e.g.

```bash
$ terraform import newrelic_synthetics_step_monitor.bar <guid>
```