---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_script_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-script-monitor"
description: |-
Create and manage a Synthetics script monitor in New Relic.
---

# Resource: newrelic\_synthetics\_script\_monitor

Use this resource to create update, and delete a Script API or Script Browser Synthetics Monitor in New Relic.

## Example Usage

##### Type: `SCRIPT_API`

-> **NOTE:** The preferred runtime is `NODE_16.10.0` while configuring the `SCRIPT_API` monitor. Other runtime may be deprecated in the future and receive fewer product updates. 

```hcl
resource "newrelic_synthetics_script_monitor" "monitor" {
  name                 = "monitor"
  type                 = "SCRIPT_API"
  locations_public      = ["AP_SOUTH_1", "AP_EAST_1"]
  period               = "EVERY_6_HOURS"
  status               = "ENABLED"
  script               = "console.log('terraform integration test updated')"
  script_language      = "JAVASCRIPT"
  runtime_type         = "NODE_API"
  runtime_type_version = "16.10"
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```
##### Type: `SCRIPT_BROWSER`

-> **NOTE:** The preferred runtime is `CHROME_BROWSER_100` while configuring the `SCRIPT_BROWSER` monitor. Other runtime may be deprecated in the future and receive fewer product updates.

```hcl
resource "newrelic_synthetics_script_monitor" "monitor" {
  enable_screenshot_on_failure_and_script = false
  locations_public                         = ["AP_SOUTH_1", "AP_EAST_1"]
  name                                    = "monitor"
  period                                  = "EVERY_HOUR"
  runtime_type_version                    = "100"
  runtime_type                            = "CHROME_BROWSER"
  script_language                         = "JAVASCRIPT"
  status                                  = "DISABLED"
  type                                    = "SCRIPT_BROWSER"
  script                                  = "$browser.get('https://one.newrelic.com')"
  tag {
    key    = "Name"
    values = ["scriptedMonitor", "hello"]
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following are the common arguments supported for `SCRIPT_API` and `SCRIPT_BROWSER` monitors:

* `account_id`- (Optional) The account in which the Synthetics monitor will be created.
* `name` - (Required) The name for the monitor.
* `type` - (Required) The plaintext representing the monitor script. Valid values are SCRIPT_BROWSER or SCRIPT_API
* `locations_public` - (Required) The location the monitor will run from. Valid public locations are https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/. At least one of either `locations_public` or `location_private` is required.
* `location_private` - (Required) The location the monitor will run from. See [Nested location_private blocks](#nested-location-private-blocks) below for details. At least one of either `locations_public` or `location_private` is required.
* `period` - (Required) The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.
* `runtime_type` - (Required) The runtime that the monitor will use to run jobs.
* `runtime_type_version` - (Required) The specific version of the runtime type selected.
* `script_language` - (Optional) The programing language that should execute the script.
* `status` - (Required) The run state of the monitor.
* `script` - (Required) The script that the monitor runs.
* `tag` - (Optional) The tags that will be associated with the monitor. See [Nested tag blocks](#nested-tag-blocks) below for details.

The `SCRIPTED_BROWSER` monitor type supports the following additional argument:

* `enable_screenshot_on_failure_and_script` - (Optional) Capture a screenshot during job execution

### Nested `tag` blocks

All nested `tag` blocks support the following common arguments:

* `key` - (Required) Name of the tag key.
* `values` - (Required) Values associated with the tag key.

### Nested `location private` blocks

All nested `location_private` blocks support the following common arguments:

* `guid` - (Required) The unique identifier for the Synthetics private location in New Relic.
* `vse_password` - (Optional) The location's Verified Script Execution password, Only necessary if Verified Script Execution is enabled for the location.

## Additional Examples

### Create a monitor with a private location

The below example shows how you can define a private location and attach it to a monitor.

-> **NOTE:** It can take up to 10 minutes for a private location to become available.

##### Type: `SCRIPT_API`

```hcl
resource "newrelic_synthetics_private_location" "private_location" {
  description               = "Test Description"
  name                      = "private_location"
  verified_script_execution = true
}

resource "newrelic_synthetics_script_monitor" "monitor" {
  name = "monitor"
  type = "SCRIPT_API"
  location_private {
    guid         = "newrelic_synthetics_private_location.private_location.id"
    vse_password = "secret"
  }
  period               = "EVERY_6_HOURS"
  status               = "ENABLED"
  script               = "console.log('terraform integration test updated')"
  script_language      = "JAVASCRIPT"
  runtime_type         = "NODE_API"
  runtime_type_version = "16.10"
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}   
```
##### Type: `SCRIPT_BROWSER`

```hcl
resource "newrelic_synthetics_private_location" "private_location" {
  description               = "Test Description"
  name                      = "private_location"
  verified_script_execution = true
}

resource "newrelic_synthetics_script_monitor" "monitor" {
  enable_screenshot_on_failure_and_script = false
  location_private {
    guid         = "newrelic_synthetics_private_location.private_location.id"
    vse_password = "secret"
  }
  name                 = "monitor"
  period               = "EVERY_HOUR"
  runtime_type_version = "100"
  runtime_type         = "CHROME_BROWSER"
  script_language      = "JAVASCRIPT"
  status               = "DISABLED"
  type                 = "SCRIPT_BROWSER"
  script               = "$browser.get('https://one.newrelic.com')"
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Synthetics monitor that the script is attached to.

## Import

Synthetics monitor scripts can be imported using the `guid`, e.g.

```bash
$ terraform import newrelic_synthetics_monitor_script.bar <guid>
```

