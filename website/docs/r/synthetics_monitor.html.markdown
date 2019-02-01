---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-monitor"
description: |-
  Create and manage synthetics monitors in New Relic.
---

# newrelic\_synthetics\_monitor

## Example Usage


```hcl
resource "newrelic_synthetics_monitor" "foo" {
    name                = "foo"
    type                = "SIMPLE"
    frequency           = 10
    uri                 = "https://www.example.com"
    locations           = ["AWS_US_EAST_1","AWS_US_WEST_2"]
    status              = "ENABLED"
    sla_threshold       = 2
    validation_string   = "Some Text"
    verify_ssl          = "true"
    bypass_head         = "false"
    redirect_is_fail    = "false"
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the synthetics monitor.
  * `type` - (Optional) Type of synthetic check to execute.  New Relic allows `SIMPLE`, `BROWSER`, `SCRIPT_BROWSER`, and `SCRIPT_API`, however, currently only `SIMPLE` mode has been implemented and tested. Defaults to `SIMPLE`.
  * `frequency` - (Optional) Interval in minutes between synthetics checks.  Must be one of: `1`, `5`, `10`, `15`, `30`, `60`, `360`, `720`, `1440`).  Defaults to `15` minutes.
  * `uri` - (Required) URI to be monitored.
  * `locations` - (Required) List of New Relic location identifiers.  At least one value must be provided.  Use the [Synthetics Location API](https://docs.newrelic.com/docs/apis/synthetics-rest-api/monitor-examples/manage-synthetics-monitors-rest-api#list-locations) to retrieve a list of valid location names for your account.
  * `status` - (Optional) Must be one of (`ENABLED`, `DISABLED`, `MUTED`).  Defaults to `ENABLED`.
  * `sla_threshold` - (Required) Response time threshold in seconds for monitor check to pass.  Defaults to 2 seconds.
  * `validation_string` - (Optional) Response validation string.
  * `verify_ssl` - (Optional) Whether to verify SSL.
  * `bypass_head` - (Optional) Whether to bypass HEAD.
  * `redirect_is_fail` - (Optional) Whether to treat redirect response as a failure.

  
## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the synthetics monitor.
    
