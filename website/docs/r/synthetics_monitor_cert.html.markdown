---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_cert_check_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-cert-check-monitor"
description: |-
Create and manage a Synthetics monitor in New Relic.
---

# Resource: newrelic\_synthetics\_cert\_check\_monitor

Use this resource to create, update, and delete the synthetics certificate check monitor in New Relic.

## Example Usage

##### Type: `CERTIFICATE CHECK`
```hcl
resource "newrelic_synthetics_cert_check_monitor" "foo" {
  name = "foo"
  domain = "example.com"
  location_public = ["AWS_US_EAST_1", "AWS_US_EAST_2"]
  certificate_expiration = "10"
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

The following are the common arguments supported for `CERTIFICATE CHECK` monitor:

* `account_id`- (Required) The account in which the Synthetics monitor will be created.
* `name` - (Required) The name for the monitor.
* `domain` - (Required) The domain of the host that will have its certificate checked.
* `location_public` - (Required) The location the monitor will run from. Valid public locations are https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/
* `location_private` - (Required) The location the monitor will run from.
* `certificate_expiration` - (Required) The desired number of remaining days until the certificate expires to trigger a monitor failure.
* `period` - (Required) The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.
* `status` - (Required) The run state of the monitor. (i.e. `ENABLED`, `DISABLED`, `MUTED`).
* `guid` - (Required) The unique identifier for the Synthetic Monitor in New Relic.


### Nested blocks

All nested `tag` blocks support the following common arguments:

* `key` - (Required) Name of the tag key.
* `values` - (Required) Values associated with the tag key.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the certificate check synthetics monitor.

## Additional Examples

##### With location_private


## Import

Synthetics certificate check monitor scripts can be imported using the `guid`, e.g.

```bash
$ terraform import newrelic_synthetics_cert_check_monitor.bar <guid>
```