---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_cert_check_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-cert-check-monitor"
description: |-
Create and manage a Synthetics Cert Check monitor in New Relic.
---

# Resource: newrelic\_synthetics\_cert\_check\_monitor

Use this resource to create, update, and delete a synthetics certificate check monitor in New Relic.

## Example Usage

```hcl
resource "newrelic_synthetics_cert_check_monitor" "cert-check-monitor" {
  name                   = "cert-check-monitor"
  domain                 = "www.example.com"
  locations_public       = ["AP_SOUTH_1"]
  certificate_expiration = "10"
  period                 = "EVERY_6_HOURS"
  status                 = "ENABLED"
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following are the common arguments supported for `CERTIFICATE CHECK` monitor:

* `account_id` - (Optional) The account in which the Synthetics monitor will be created.
* `name` - (Required) The name for the monitor.
* `domain` - (Required) The domain of the host that will have its certificate checked.
* `locations_public` - (Required) The location the monitor will run from. Valid public locations are https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/. You don't need the `AWS_` prefix as the provider uses NerdGraph. At least one of either `locations_public` or `location_private` is required.
* `locations_private` - (Required) The location the monitor will run from. Accepts a list of private location GUIDs. At least one of either `locations_public` or `locations_private` is required.
* `certificate_expiration` - (Required) The desired number of remaining days until the certificate expires to trigger a monitor failure.
* `period` - (Required) The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.
* `status` - (Required) The run state of the monitor. (i.e. `ENABLED`, `DISABLED`, `MUTED`).
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
resource "newrelic_synthetics_private_location" "private_location" {
  description               = "Test Description"
  name                      = "private_location"
  verified_script_execution = false
}

resource "newrelic_synthetics_cert_check_monitor" "monitor" {
  name              = "cert_check"
  uri               = "https://www.one.example.com"
  locations_private = ["newrelic_synthetics_private_location.private_location.id"]
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

* `id` - The ID of the certificate check synthetics monitor.

## Import

Synthetics certificate check monitor scripts can be imported using the `guid`, e.g.

```bash
$ terraform import newrelic_synthetics_cert_check_monitor.monitor <guid>
```
