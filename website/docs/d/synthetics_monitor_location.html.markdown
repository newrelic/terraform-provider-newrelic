---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_monitor_location"
sidebar_current: "docs-newrelic-datasource-synthetics-monitor-location"
description: |-
  Grabs a Synthetics monitor location by label.
---

# Data Source: newrelic\_synthetics\_monitor\_location

Use this data source to get information about a specific Synthetics monitor private location in New Relic that already exists.

## Example Usage

```hcl
data "newrelic_synthetics_monitor_location" "example" {
  name = "My private location"
}

resource "newrelic_synthetics_monitor" "foo" {
  // Reference the private location data source in the monitor resource
  locations = [data.newrelic_synthetics_monitor_location.example.name]
}
```

## Argument Reference

The following arguments are supported:


* `name` - (Optional) The name of the Synthetics monitor private location.
* `label` - (Optional) **DEPRECATED:** Use `name` instead. The label of the Synthetics monitor private location.


```
Warning: This data source will use the account ID linked to your API key. At the moment it is not possible to dynamically set the account ID.
```
