---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_monitor_script"
sidebar_current: "docs-newrelic-resource-synthetics-monitor-script"
description: |-
  Update and manage a Synthetics monitor script in New Relic.
---

# Resource: newrelic\_synthetics\_monitor\_script

Use this resource to update a synthetics monitor script in New Relic.

## Example Usage

```hcl
resource "newrelic_synthetics_monitor" "foo" {
  name = "foo"
  type = "SCRIPT_BROWSER"
  frequency = 5
  status = "ENABLED"
  locations = ["AWS_US_EAST_1"]
}

resource "newrelic_synthetics_monitor_script" "foo_script" {
  monitor_id = newrelic_synthetics_monitor.foo.id
  text = file("${path.module}/foo_script.js")
  location {
    name = "YWJjZAo="
    hmac = "ZmFrZWxvY2F0aW9uc2NyaXB0ZmFrZQ=="
  }
}
```

## Argument Reference

The following arguments are supported:

  * `monitor_id` - (Required) The ID of the monitor to attach the script to.
  * `text` - (Required) The plaintext representing the monitor script.
  * `location` - (Optional) A nested block that describes a monitor script location. See [Nested location blocks](#nested-`location`-blocks) below for details

### Nested `location` blocks

All nested `location` blocks support the following common arguments:

  * `name` - (Required) The monitor script location name.
  * `hmac` - (Optional) The monitor script authentication code for the location. Use one of either `hmac` or `vse_password`.
  * `vse_password` - (Optional) The password for the location used to calculate the HMAC. Use one of either `hmac` or `vse_password`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The ID of the Synthetics monitor that the script is attached to.

## Import

Synthetics monitor scripts can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_synthetics_monitor_script.main <id>
```