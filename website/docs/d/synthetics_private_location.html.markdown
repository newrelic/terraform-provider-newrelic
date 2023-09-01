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
  account_id = 123456
  name       = "My private location"
}

resource "newrelic_synthetics_monitor" "foo" {
  // Reference the private location data source in the monitor resource
  locations_private = [data.newrelic_synthetics_monitor_location.example.id]
}
```

```hcl
data "newrelic_synthetics_private_location" "example" {
  account_id = 123456
  name       = "My private location"
}
resource "newrelic_synthetics_step_monitor" "foo" {
  // Reference the private location data source in the monitor resource
  location_private { 
    guid = data.newrelic_synthetics_private_location.example.id 
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The New Relic account ID of the associated private location. If left empty will default to account ID specified in provider level configuration.
* `name` - (Required) The name of the Synthetics monitor private location.
