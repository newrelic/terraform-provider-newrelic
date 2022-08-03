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
  name = "broken"
  uri = "https://www.one.example.com"
  locations_public = ["AP_SOUTH_1"]
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
* `locations_public` - (Required) The location the monitor will run from. Valid public locations are https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/
* `locations_private` - (Required) The location the monitor will run from.
* `period` - (Required) The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.
* `status` - (Required) The run state of the monitor. (i.e. `ENABLED`, `DISABLED`, `MUTED`).
* `guid` - (Required) The unique identifier for the Synthetic Monitor in New Relic.

### Nested blocks

All nested `tag` blocks support the following common arguments:

* `key` - (Required) Name of the tag key.
* `values` - (Required) Values associated with the tag key.


## Additional Examples

### Create a monitor with a private location

The below example shows how you can define a private location and attach it to a monitor.

-> **NOTE:** It can take up to 10 minutes for a private location to become available.

```hcl
resource "newrelic_synthetics_private_location" "bar1" {
  description               = "Test Description"
  name                      = "private_location"
  verified_script_execution = true
}
  resource "newrelic_synthetics_broken_links_monitor" "bar" {
    name = "broken"
    uri = "https://www.one.example.com"
    locations_private = ["newrelic_synthetics_private_location.private_location.id"]
    period = "EVERY_6_HOURS"
    status = "ENABLED"
    tag {
      key = "some_key"
      values = ["some_value"]
    }
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the broken links synthetics monitor.

## Import

Synthetics broken links monitor scripts can be imported using the `guid`, e.g.

```bash
$ terraform import newrelic_synthetics_broken_links_monitor.foo <guid>
```