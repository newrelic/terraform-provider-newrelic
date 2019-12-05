---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-monitor"
description: |-
  Create and manage a Synthetics monitor in New Relic.
---

# newrelic\_synthetics\_monitor

Use this resource to create, update, and delete a synthetics monitor in New Relic.

## Example Usage

##### Type: `SIMPLE`
```hcl
resource "newrelic_synthetics_monitor" "foo" {
  name = "foo"
  type = "SIMPLE"
  frequency = 5
  status = "ENABLED"
  locations = ["AWS_US_EAST_1", "AWS_US_EAST_2"]

  uri                       = "https://example.com"               # Required for type "SIMPLE" and "BROWSER"
  validation_string         = "add example validation check here" # Optional for type "SIMPLE" and "BROWSER"
  verify_ssl                = true                                # Optional for type "SIMPLE" and "BROWSER"
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The title of this monitor.
  * `type` - (Required) The monitor type. Valid values are `SIMPLE`, `BROWSER`, `SCRIPT_BROWSER`, and `SCRIPT_API`.
  * `frequency` - (Required) The interval (in minutes) at which this monitor should run.
  * `status` - (Required) The monitor status (i.e. `ENABLED`, `MUTED`, `DISABLED`)
  * `locations` - (Required) The locations in which this monitor should be run.
  * `sla_threshold` - (Optional) The base threshold for the SLA report.

 The `SIMPLE` monitor type supports the following additional arguments:

  * `uri` - (Required) The URI for the monitor to hit.
  * `validation_string` - (Optional) The string to validate against in the response.
  * `verify_ssl` - (Optional) Verify SSL.
  * `bypass_head_request` - (Optional) Bypass HEAD request.
  * `treat_redirect_as_failure` - (Optional) Fail the monitor check if redirected.

The `BROWSER` monitor type supports the following additional arguments:

  * `uri` - (Required) The URI for the monitor to hit.
  * `validation_string` - (Optional) The string to validate against in the response.
  * `verify_ssl` - (Optional) Verify SSL.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the Synthetics monitor.

## Additional Examples

Type: `BROWSER`

```hcl
resource "newrelic_synthetics_monitor" "foo" {
  name = "foo"
  type = "BROWSER"
  frequency = 5
  status = "ENABLED"
  locations = ["AWS_US_EAST_1"]

  uri                       = "https://example.com"               # required for type "SIMPLE" and "BROWSER"
  validation_string         = "add example validation check here" # optional for type "SIMPLE" and "BROWSER"
  verify_ssl                = true                                # optional for type "SIMPLE" and "BROWSER"
  bypass_head_request       = true                                # Note: optional for type "BROWSER" only
  treat_redirect_as_failure = true                                # Note: optional for type "BROWSER" only
}
```

Type: `SCRIPT_BROWSER`

```hcl
resource "newrelic_synthetics_monitor" "foo" {
  name = "foo"
  type = "SCRIPT_BROWSER"
  frequency = 5
  status = "ENABLED"
  locations = ["AWS_US_EAST_1"]
}
```

Type: `SCRIPT_API`

```hcl
resource "newrelic_synthetics_monitor" "foo" {
  name = "foo"
  type = "SCRIPT_API"
  frequency = 5
  status = "ENABLED"
  locations = ["AWS_US_EAST_1"]
}
```
