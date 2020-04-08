---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_monitor"
sidebar_current: "docs-newrelic-datasource-synthetics-monitor"
description: |-
  Grabs a synthetics monitor by name.
---

# Data Source: newrelic\_synthetics\_monitor

Use this data source to get information about a specific synthetics monitor in New Relic that already exists. This can be used to set up a Synthetics alert condition.

## Example Usage

```hcl
data "newrelic_synthetics_monitor" "bar" {
  name = "bar"
}

resource "newrelic_synthetics_alert_condition" "baz" {
  policy_id = newrelic_alert_policy.foo.id

  name        = "baz"
  monitor_id  = data.newrelic_synthetics_monitor.bar.id
  runbook_url = "https://www.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the synthetics monitor in New Relic.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `monitor_id` - The ID of the synthetics monitor.
