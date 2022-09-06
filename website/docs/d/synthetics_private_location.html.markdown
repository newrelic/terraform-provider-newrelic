---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_private_location"
sidebar_current: "docs-newrelic-datasource-synthetics-private-location"
description: |-
  Grabs a Synthetics monitor location by name.
---

# Data Source: newrelic\_synthetics\_private\_location

Use this data source to get information about a specific Synthetics monitor private location in New Relic that already exists.

## Example Usage

```hcl
data "newrelic_synthetics_private_location" "example" {
  name = "My private location"
}

resource "newrelic_synthetics_monitor" "foo" {
  // Reference the private location data source in the monitor resource
  location_private = [data.newrelic_synthetics_monitor_location.example.id]
}
```

-> This data source only works for `simple`, `browser`, `cert_check` and `broken_links` monitors

```hcl
data "newrelic_synthetics_private_location" "example" {
  name = "My private location"
}
resource "newrelic_synthetics_monitor" "foo" {
  // Reference the private location data source in the monitor resource
  location_private { guid = data.newrelic_synthetics_private_location.example.id }
}
```

## Argument Reference

The following arguments are supported:


* `name` - (Required) The name of the Synthetics monitor private location.
