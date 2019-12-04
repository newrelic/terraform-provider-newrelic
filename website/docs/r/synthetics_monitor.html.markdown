---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-monitor"
description: |-
  Create and manage a Synthetics monitor in New Relic.
---

# Resource: newrelic\_synthetics\_monitor

Use this resource to create, update, and delete a synthetics monitor in New Relic.

## Example Usage

```hcl
resource "newrelic_synthetics_monitor" "foo" {
  name = "foo"
  type = "SIMPLE"
  frequency = 5
  status = "ENABLED"
  locations = ["AWS_US_EAST_1"]
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The title of this monitor.
  * `type` - (Required) The monitor type (i.e. SIMPLE, BROWSER, SCRIPT_API, SCRIPT_BROWSER).
  * `frequency` - (Required) The interval (in minutes) at which this monitor should run.
  * `status` - (Required) The monitor status (i.e. ENABLED, MUTED, DISABLED)
  * `locations` - (Required) The locations in which this monitor should be run.
  * `sla_threshold` - (Optional) The base threshold for the SLA report.
  
For SIMPLE and BROWSER monitor types, the following arguments are also supported:

  * `uri` - (Required) The URI for the monitor to hit.
  * `validation_string` - (Optional) The string to validate against in the response.
  * `verify_ssl` - (Optional) Verify SSL.
  * `bypass_head_request` - (Optional) Bypass HEAD request.
  * `treat_redirect_as_failure` - (Optional) Fail the monitor check if redirected.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The ID of the Synthetics monitor.
