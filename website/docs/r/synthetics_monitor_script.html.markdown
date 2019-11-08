---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_monitor_script"
sidebar_current: "docs-newrelic-resource-synthetics-monitor-script"
description: |-
  Update and manage a Synthetics monitor script in New Relic.
---

# newrelic\_synthetics\_monitor\_script

Use this resource to update a synthetics monitor script in New Relic.

## Example Usage

```hcl
data "template_file" "foo_script" {
  template = "${file("${path.module}/foo_script.tpl")}"
}

resource "newrelic_synthetics_monitor_script" "foo_script" {
  monitor_id = "${newrelic_synthetics_monitor.foo.id}"
  text = "${data.template_file.foo_script.rendered}"
}
```

## Argument Reference

The following arguments are supported:

  * `monitor_id` - (Required) The ID of the monitor to attach the script to.
  * `text` - (Required) plaintext of the monitor script.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the Synthetics monitor that the script is attached to.
