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

data "template_file" "foo_script" {
  template = file("${path.module}/foo_script.tpl")
}

resource "newrelic_synthetics_monitor_script" "foo_script" {
  monitor_id = newrelic_synthetics_monitor.foo.id
  text = data.template_file.foo_script.rendered
}
```

## Argument Reference

The following arguments are supported:

  * `monitor_id` - (Required) The ID of the monitor to attach the script to.
  * `text` - (Required) The plaintext representing the monitor script.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The ID of the Synthetics monitor that the script is attached to.
