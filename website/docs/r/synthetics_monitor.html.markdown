---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-monitor"
description: |-
  Create and manage a Synthetics monitor in New Relic.
---

# Resource: newrelic\_synthetics\_monitor

Use this resource to create and manage New Relic Synthetics monitor.

## Example Usage

##### Type: `SIMPLE`
```hcl
resource "newrelic_synthetics_monitor" "foo" {
  custom_headers{
    name  =  "Name"
    value = "simpleMonitor"
  }
  treat_redirect_as_failure=true
  validation_string="success"
  bypass_head_request=true
  verify_ssl=true
  locations = ["AP_SOUTH_1"]
  name      = "%[1]s"
  frequency = 5
  status    = "ENABLED"
  type      = "SIMPLE"
  tags{
    key = "monitor"
    values  = ["myMonitor"]
  }
  uri = "https://www.one.newrelic.com"
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The title of this monitor.
  * `type` - (Required) The monitor type. Valid values are `SIMPLE`, `BROWSER`.
  * `frequency` - (Required) The interval (in minutes) at which this monitor should run.
  * `status` - (Required) The monitor status (i.e. `ENABLED`, `MUTED`, `DISABLED`).
  * `locations` - (Required) The locations in which this monitor should be run.
  * `sla_threshold` - (Optional) The base threshold (in seconds) to calculate the [Apdex score](https://docs.newrelic.com/docs/apm/new-relic-apm/apdex/apdex-measure-user-satisfaction/) for use in the [SLA report](https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/pages/synthetic-monitoring-aggregate-monitor-metrics/#viewing). Default is 7 seconds.

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

```
Warning: This resource will use the account ID linked to your API key. At the moment it is not possible to dynamically set the account ID.
```

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the Synthetics monitor.

## Additional Examples

Type: `BROWSER`

```hcl
resource "newrelic_synthetics_monitor" "bar" {
  custom_headers{
    name  ="name"
    value ="simple_browser"
  }
  enable_screenshot_on_failure_and_script=false
  validation_string ="success"
  verify_ssl  =false
  locations   = ["AP_SOUTH_1","AP_EAST_1"]
  name        = "%[1]s-Updated"
  frequency   = 10
  runtime_type_version  ="100"
  runtime_type  ="CHROME_BROWSER"
  script_language ="JAVASCRIPT"
  status      = "DISABLED"
  type        = "BROWSER"
  tags{
    key = "name"
    values  = ["SimpleBrowserMonitor","my_monitor"]
  }
  uri = "https://www.one.newrelic.com"
}
```

## Import

Synthetics monitor can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_synthetics_monitor.main <id>
```
