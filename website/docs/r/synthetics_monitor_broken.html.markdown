---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_broken_links_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-broken-links-monitor"
description: |-
Create and manage a Synthetics monitor in New Relic.
---

# Resource: newrelic\_synthetics\_broken\_links\_monitor

Use this resource to create, update, and delete the synthetics broken links monitor in New Relic.

## Example Usage

##### Type: `BROKEN LINKS`
```hcl
resource "newrelic_synthetics_broken_links_monitor" "foo" {
  name = "foo"
  uri = "https://www.one.example.com"
  locations = ["AWS_US_EAST_1", "AWS_US_EAST_2"]
  period = "EVERY_6_HOURS"
  status = "ENABLED"
  tag {
    key = "some_key"
    values = ["some_value"]
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following are the common arguments supported for `BROKEN LINKS` monitor:

* `account_id`- (Required) The account in which the Synthetics monitor will be created.
* `name` - (Required) The name for the monitor.
* `uri` - (Required) The uri the monitor runs against.
* `location_public` - (Required) The location the monitor will run from. Valid public locations are https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/
* `location_private` - (Required) The location the monitor will run from.
* `period` - (Required) The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.
* `status` - (Required) The run state of the monitor. (i.e. `ENABLED`, `DISABLED`, `MUTED`).

### Nested blocks

All nested `tag` blocks support the following common arguments:

* `key` - (Required) Name of the tag key.
* `values` - (Required) Values associated with the tag key.

All nested `location_private` blocks support the following common arguments:

* `guid` - (Required) The unique identifier for the Synthetics private location in New Relic.
* `vse_password` - (Optional) The location's Verified Script Execution password, Only necessary if Verified Script Execution is enabled for the location.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the broken links synthetics monitor.

## Import

Synthetics broken links monitor scripts can be imported using the `guid`, e.g.

```bash
$ terraform import newrelic_synthetics_broken_links_monitor.bar <guid>
```