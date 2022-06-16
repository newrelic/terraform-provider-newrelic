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
  custom_header{
    name  =  "Name"
    value = "simpleMonitor"
  }
  treat_redirect_as_failure = true
  validation_string = "success"
  bypass_head_request = true
  verify_ssl  = true
  locations_public = ["AP_SOUTH_1"]
  name  = "%[1]s"
  period =  "EVERY_MINUTE"
  status  = "ENABLED"
  type  = "SIMPLE"
  tag{
    key = "monitor"
    values  = ["myMonitor"]
  }
  uri = "https://www.one.newrelic.com"
}
```
##### Type: `BROWSER`
```hcl
resource "newrelic_synthetics_monitor" "bar" {
  custom_headers{
    name	= "name"
    value	= "simple_browser"
  }
  enable_screenshot_on_failure_and_script = true
  validation_string = "success"
  verify_ssl  = true
  locations_public  = ["AP_SOUTH_1"]
  name  = "%s"
  period  = "EVERY_MINUTE"
  runtime_type_version  = "100"
  runtime_type  = "CHROME_BROWSER"
  script_language = "JAVASCRIPT"
  status  = "ENABLED"
  type  = "BROWSER"
  uri = "https://www.one.newrelic.com"
  tag {
    key = "name"
    values  = ["SimpleBrowserMonitor"]
  }
}

```
## Argument Reference

The following are the common arguments supported for `SIMPLE` and `BROWSER` monitors:

* `validation_string` - (Optional) TValidation text for monitor to search for at given URI.
* `verify_ssl` - (Optional) Monitor should validate SSL certificate chain.
* `period` - (Required) The interval at which the monitor runs in minutes.
* `status` - (Required) The run state of the monitor.
* `locations_public` - (Required) The locations the monitor will run from.
* `name` - (Required) The human-readable identifier for the monitor.
* `uri` - (Required) The uri the monitor runs against.
* `type` - (Required) THE monitor type. Valid values are `SIMPLE` and `BROWSER`.

Other arguments supports for `SIMPLE` monitor:

* `treat_redirect_as_failure` - (Optional) Categorize redirects during a monitor job as a failure.
* `bypass_head_request` - (Optional) Monitor should skip default HEAD request and instead use GET verb in check.

Other arguments supports for `BROWSER` monitor:

* `enable_screenshot_on_failure_and_script` - (Optional) Capture a screenshot during job execution.
* `runtime_type_version` - (Required) The runtime type that the monitor will run.
* `runtime_type` - (Required) The runtime type that the monitor will run.
* `script_language` - (Optional) The programing language that should execute the script.

### Nested blocks

All nested `custom_header` blocks support the following common arguments:

* `name` - (Required) Header name.
* `value` - (Required) Header Value.

All nested `tag` blocks support the following common arguments:

* `key` - (Required) Name of the tag key.
* `values` - (Required) Values associated with the tag key.

```
Warning: This resource will use the account ID linked to your API key. At the moment it is not possible to dynamically set the account ID.
```

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Synthetics monitor.

## Import

Synthetics monitor can be imported using the `guid`, e.g.

```bash
$ terraform import newrelic_synthetics_monitor.bar <guid>
```