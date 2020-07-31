---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_monitor_location"
sidebar_current: "docs-newrelic-datasource-synthetics-monitor-location"
description: |-
  Grabs a Synthetics monitor location by label.
---

# Data Source: newrelic\_synthetics\_monitor

Use this data source to get information about a specific Synthetics monitor location in New Relic that already exist.

## Example Usage

```hcl
data "newrelic_synthetics_monitor_location" "bar" {
  label = "My private location"
}

resource "newrelic_synthetics_monitor" "foo" {
  name = "foo"
  type = "SIMPLE"
  frequency = 5
  status = "ENABLED"
  locations = [data.newrelic_synthetics_monitor_location.bar.name]

  uri                       = "https://example.com"               # Required for type "SIMPLE" and "BROWSER"
  validation_string         = "add example validation check here" # Optional for type "SIMPLE" and "BROWSER"
  verify_ssl                = true                                # Optional for type "SIMPLE" and "BROWSER"
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
