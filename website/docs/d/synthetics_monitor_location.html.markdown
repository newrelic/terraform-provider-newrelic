---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_monitor_location"
sidebar_current: "docs-newrelic-datasource-synthetics-monitor-location"
description: |-
  Grabs a Synthetics monitor location by label.
---

# Data Source: newrelic\_synthetics\_monitor\_location

Use this data source to get information about a specific Synthetics monitor location in New Relic that already exists.

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
* `high_security_mode` - Represents if high security mode is enabled for the location. A value of true means that high security mode is enabled, and a value of false means it is disabled.
* `private` - Represents if this location is a private location. A value of true means that the location is private, and a value of false means it is public.
* `description` - A description of the Synthetics monitor location.
