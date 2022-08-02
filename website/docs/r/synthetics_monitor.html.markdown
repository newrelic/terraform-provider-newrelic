---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-monitor"
description: |-
Create and manage a Synthetics monitor in New Relic.
---

# Resource: newrelic\_synthetics\_monitor

Use this resource to create, update, and delete a Simple or Browser Synthetics Monitor in New Relic.

-> **NOTE:** The [newrelic_synthetics_private_location](newrelic_synthetics_private_location.html) resource private minion can take upto 10 minutes to be available through Terraform.

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
  location_public = ["AP_SOUTH_1"]
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
See additional [examples](#additional-examples).

```hcl
resource "newrelic_synthetics_private_location" "bar1" {
  description               = "Test Description-Updated"
  name                      = "%[1]S"
  verified_script_execution = true
}
resource "newrelic_synthetics_monitor" "foo" {
  custom_header{
    name   = "name"
    value  = "simple_browser_updTE"
  }
  location_private = ["newrelic_synthetics_private_location.bar1"]
  treat_redirect_as_failure = true
  validation_string = "success"
  bypass_head_request = true
  verify_ssl  = true
  location_public = ["AP_SOUTH_1"]
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

##### Type: `SIMPLE BROWSER`

-> **NOTE:** The preferred runtime is `CHROME_BROWSER_100` while configuring the `SIMPLE_BROWSER` monitor. Other runtime may be deprecated in the future and receive fewer product updates.

```hcl
resource "newrelic_synthetics_monitor" "bar" {
  custom_headers{
    name	= "name"
    value	= "simple_browser"
  }
  enable_screenshot_on_failure_and_script = true
  validation_string = "success"
  verify_ssl  = true
  location_public  = ["AP_SOUTH_1"]
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
See additional [examples](#additional-examples).

```hcl
resource "newrelic_synthetics_private_location" "bar1" {
  description               = "Test Description-Updated"
  name                      = "%[1]S"
  verified_script_execution = true
}
resource "newrelic_synthetics_monitor" "bar" {
  custom_headers {
    name  = "name"
    value = "simple_browser"
  }
  location_private                        = ["newrelic_synthetics_private_location.bar1"]
  enable_screenshot_on_failure_and_script = true
  validation_string                       = "success"
  verify_ssl                              = true
  location_public                         = ["AP_SOUTH_1"]
  name                                    = "%s"
  period                                  = "EVERY_MINUTE"
  runtime_type_version                    = "100"
  runtime_type                            = "CHROME_BROWSER"
  script_language                         = "JAVASCRIPT"
  status                                  = "ENABLED"
  type                                    = "BROWSER"
  uri                                     = "https://www.one.newrelic.com"
  tag {
    key    = "name"
    values = ["SimpleBrowserMonitor"]
  }
}
```
## Argument Reference

The following are the common arguments supported for `SIMPLE` and `BROWSER` monitors:

* `account_id`- (Required) The account in which the Synthetics monitor will be created.
* `validation_string` - (Optional) Validation text for monitor to search for at given URI.
* `verify_ssl` - (Optional) Monitor should validate SSL certificate chain.
* `period` - (Required) The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.
* `status` - (Required) The run state of the monitor.
* `location_public` - (Required) The location the monitor will run from. Valid public locations are https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/
* `location_private` - (Required) The location the monitor will run from.
* `name` - (Required) The human-readable identifier for the monitor.
* `uri` - (Required) The uri the monitor runs against.
* `type` - (Required) THE monitor type. Valid values are `SIMPLE` and `BROWSER`.
* `guid` - (Required) The unique identifier for the Synthetic Monitor in New Relic.

The `SIMPLE` monitor type supports the following additional arguments:

* `treat_redirect_as_failure` - (Optional) Categorize redirects during a monitor job as a failure.
* `bypass_head_request` - (Optional) Monitor should skip default HEAD request and instead use GET verb in check.

The `SIMPLE BROWSER` monitor type supports the following additional arguments:

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

All nested `location_private` blocks support the following common arguments:

* `guid` - (Required) The unique identifier for the Synthetics private location in New Relic.
* `vse_password` - (Optional) The location's Verified Script Execution password, Only necessary if Verified Script Execution is enabled for the location.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Synthetics monitor that the script is attached to.

## Import

Synthetics monitor can be imported using the `guid`, e.g.

```bash
$ terraform import newrelic_synthetics_monitor.bar <guid>
```