---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_monitor_location"
sidebar_current: "docs-newrelic-datasource-synthetics-monitor-location"
description: |-
  Grabs a Synthetics monitor location by label.
---

# Data Source: newrelic\_synthetics\_monitor

Use this data source to get information about a specific Synthetics monitor locations in New Relic that already exist. This can be used to set up a Synthetics alert condition.

## Example Usage

```hcl
data "newrelic_synthetics_monitor" "bar" {
  name = "bar"
}

resource "newrelic_synthetics_alert_condition" "baz" {
  policy_id = newrelic_alert_policy.foo.id

  name        = "baz"
  monitor_id  = data.newrelic_synthetics_monitor_location.bar.label
  runbook_url = "https://www.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The label of the Synthetics monitor location.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of the Synthetics monitor location.
* `high_security_mode` - The high security mode for the Synthetics monitor location.
* `private` - The private setting for the Synthetics monitor location.
* `description` - The description of the Synthetics monitor location.
