---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_label"
sidebar_current: "docs-newrelic-resource-synthetics-label"
description: |-
  Create and manage a Synthetics label in New Relic.
---

# newrelic\_synthetics\_label

Use this resource to create, update, and delete a Synthetics label in New Relic.

## Example Usage

```hcl
resource "newrelic_synthetics_label" "foo" {
  monitor_id = newrelic_synthetics_monitor.foo.id
  type = "MyCategory"
  value = "MyValue"
}
```

## Argument Reference

The following arguments are supported:

  * `monitor_id` - (Required) The ID of the monitor that will be assigned the label.
  * `type` - (Required) A string representing the label key/category.
  * `value` - (Required) A string representing the label value.

## Attributes Reference

The following attributes are exported:

  * `href` - The URL of the Synthetics label.

## Import

Synthetics labels can be imported using an ID in the format `<monitor_id>:<type>:<value>`, e.g.

```
$ terraform import newrelic_synthetics_labels.foo 1a272364-f204-4cd3-ae2a-2d15a2bedadd:MyCategory:MyValue
```
