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

-> **NOTE:** The preferred runtime is `NODE_16.10` for configuring the `SCRIPT_API` monitor. If you wish to use the old runtime, please remove the `script_language`, `runtime_type` and `runtime_type_version` attributes, or set them to empty string `""`. The old runtime will be deprecated in the future, so use the new version whenever you can.

##### Type: `SCRIPT_API`

```hcl
resource "newrelic_synthetics_script_monitor" "foo" {
  status = "ENABLED"
  name   = "monitor"
  type   = "SCRIPT_API"
  locations_public = [
    "AP_SOUTH_1",
    "AP_EAST_1"
  ]
  period = "EVERY_6_HOURS"

  script = "console.log('it works!')"

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
resource "newrelic_synthetics_script_monitor" "monitor" {
  status = "ENABLED"
  name   = "monitor"
  type   = "SCRIPT_BROWSER"
  locations_public = [
    "AP_SOUTH_1",
    "AP_EAST_1"
  ]
  period                                  = "EVERY_HOUR"
  enable_screenshot_on_failure_and_script = false

  script = "$browser.get('https://one.newrelic.com')"

  runtime_type_version = "100"
  runtime_type         = "CHROME_BROWSER"
  script_language      = "JAVASCRIPT"

  tag {
    key    = "some_key"
    values = ["some_value1", "some_value2"]
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following are the common arguments supported for `SCRIPT_API` and `SCRIPT_BROWSER` monitors:

* `account_id`- (Optional) The account in which the Synthetics monitor will be created.
* `status` - (Required) The run state of the monitor: `ENABLED` or `DISABLED`
* `name` - (Required) The name for the monitor.
* `type` - (Required) The plaintext representing the monitor script. Valid values are SCRIPT_BROWSER or SCRIPT_API
* `locations_public` - (Required) The location the monitor will run from. Valid public locations are https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/. You don't need the `AWS_` prefix as the provider uses NerdGraph. At least one of either `locations_public` or `location_private` is required.
* `location_private` - (Required) The location the monitor will run from. See [Nested location_private blocks](#nested-location-private-blocks) below for details. At least one of either `locations_public` or `location_private` is required.
* `period` - (Required) The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.
* `script` - (Required) The script that the monitor runs.
* `runtime_type` - (Optional) The runtime that the monitor will use to run jobs. Defaults to `CHROME_BROWSER`
* `runtime_type_version` - (Optional) The specific version of the runtime type selected. Defaults to `100`
* `script_language` - (Optional) The programing language that should execute the script. Defaults to `JAVASCRIPT`
* `tag` - (Optional) The tags that will be associated with the monitor. See [Nested tag blocks](#nested-tag-blocks) below for details.

The `SCRIPTED_BROWSER` monitor type supports the following additional argument:

* `enable_screenshot_on_failure_and_script` - (Optional) Capture a screenshot during job execution

#### Deprecated runtime

If you want to use the deprecated Node 10 runtime you can set the `runtime_type`, `runtime_type_version` and `script_language` to empty string `""`. The old runtime will be deprecated in the future, so use the new version whenever you can.

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
resource "newrelic_synthetics_private_location" "foo" {
  description               = "Example private location"
  name                      = "private_location"
  verified_script_execution = true
}

resource "newrelic_synthetics_script_monitor" "bar" {
  status = "ENABLED"
  name   = "Example synthetics monitor"
  type   = "SCRIPT_API"
  location_private {
    guid         = "newrelic_synthetics_private_location.private_location.id"
    vse_password = "secret"
  }
  period = "EVERY_6_HOURS"

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
resource "newrelic_synthetics_private_location" "foo" {
  description               = "Test Description"
  name                      = "private_location"
  verified_script_execution = true
}

resource "newrelic_synthetics_script_monitor" "bar" {
  status = "ENABLED"
  name   = "Example synthetics monitor"
  type   = "SCRIPT_BROWSER"
  period = "EVERY_HOUR"
  script = "$browser.get('https://one.newrelic.com')"

  enable_screenshot_on_failure_and_script = false
  location_private {
    guid         = "newrelic_synthetics_private_location.private_location.id"
    vse_password = "secret"
  }

  runtime_type_version = "100"
  runtime_type         = "CHROME_BROWSER"
  script_language      = "JAVASCRIPT"

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

