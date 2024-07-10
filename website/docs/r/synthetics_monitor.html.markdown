---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-monitor"
description: |-
    Create and manage a Synthetics monitor in New Relic.
---

# Resource: newrelic\_synthetics\_monitor

-> **WARNING** Support for legacy Synthetics runtimes **will reach its end-of-life (EOL) on October 22, 2024**. In addition, creating **_new_** monitors using the legacy runtime **will no longer be supported after August 26, 2024**. In light of the above, kindly **upgrade your Synthetic Monitors to the new runtime** at the earliest, if they are still using the legacy runtime. Please check out [this page](https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm) for more details on the EOL, action needed (specific to monitors using public and private locations), relevant resources, and more.

Use this resource to create, update, and delete a Simple or Browser Synthetics Monitor in New Relic.

## Example Usage
```hcl
resource "newrelic_synthetics_monitor" "monitor" {
  status           = "ENABLED"
  name             = "monitor"
  period           = "EVERY_MINUTE"
  uri              = "https://www.one.newrelic.com"
  type             = "SIMPLE"
  locations_public = ["AP_SOUTH_1"]

  custom_header {
    name  = "some_name"
    value = "some_value"
  }

  treat_redirect_as_failure = true
  validation_string         = "success"
  bypass_head_request       = true
  verify_ssl                = true

  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```
##### Type: `SIMPLE BROWSER`

```hcl
resource "newrelic_synthetics_monitor" "monitor" {
  status           = "ENABLED"
  name             = "monitor"
  period           = "EVERY_MINUTE"
  uri              = "https://www.one.newrelic.com"
  type             = "BROWSER"
  locations_public = ["AP_SOUTH_1"]

  custom_header {
    name  = "some_name"
    value = "some_value"
  }

  enable_screenshot_on_failure_and_script = true
  validation_string                       = "success"
  verify_ssl                              = true
  runtime_type                            = "CHROME_BROWSER"
  runtime_type_version                    = "100"
  script_language                         = "JAVASCRIPT"

  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following are the common arguments supported for `SIMPLE` and `BROWSER` monitors:

* `account_id`- (Optional) The account in which the Synthetics monitor will be created.
* `status` - (Required) The run state of the monitor. (`ENABLED` or `DISABLED`).

-> **WARNING:** As of February 29, 2024, Synthetic Monitors no longer support the `MUTED` status. Version **3.33.0** of the New Relic Terraform Provider is released to coincide with the `MUTED` status end-of-life. Consequently, the only valid values for `status` for all types of Synthetic Monitors are mentioned above. For additional information on alternatives to the `MUTED` status of Synthetic Monitors that can be managed via Terraform, please refer to [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/upcoming_synthetics_muted_status_eol_guide).
* `name` - (Required) The human-readable identifier for the monitor.
* `period` - (Required) The interval at which this monitor should run. Valid values are `EVERY_MINUTE`, `EVERY_5_MINUTES`, `EVERY_10_MINUTES`, `EVERY_15_MINUTES`, `EVERY_30_MINUTES`, `EVERY_HOUR`, `EVERY_6_HOURS`, `EVERY_12_HOURS`, or `EVERY_DAY`.
* `uri` - (Required) The URI the monitor runs against.
* `type` - (Required) The monitor type. Valid values are `SIMPLE` and `BROWSER`.
* `locations_public` - (Required) The location the monitor will run from. Check out [this page](https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/) for a list of valid public locations. You don't need the `AWS_` prefix as the provider uses NerdGraph. At least one of either `locations_public` or `location_private` is required.
* `locations_private` - (Required) The location the monitor will run from. Accepts a list of private location GUIDs. At least one of either `locations_public` or `locations_private` is required.
* `custom_header`- (Optional) Custom headers to use in monitor job. See [Nested custom_header blocks](#nested-custom-header-blocks) below for details.
* `validation_string` - (Optional) Validation text for monitor to search for at given URI.
* `verify_ssl` - (Optional) Monitor should validate SSL certificate chain.
* `tag` - (Optional) The tags that will be associated with the monitor. See [Nested tag blocks](#nested-tag-blocks) below for details.

The `SIMPLE` monitor type supports the following additional arguments:

* `treat_redirect_as_failure` - (Optional) Categorize redirects during a monitor job as a failure.
* `bypass_head_request` - (Optional) Monitor should skip default HEAD request and instead use GET verb in check.

The `BROWSER` monitor type supports the following additional arguments:

* `enable_screenshot_on_failure_and_script` - (Optional) Capture a screenshot during job execution.
* `runtime_type_version` - (Optional) The runtime type that the monitor will run. Valid value is `100`.
* `runtime_type` - (Optional) The runtime type that the monitor will run. Valid value is `CHROME_BROWSER`
* `script_language` - (Optional) The programing language that should execute the script.
* `device_orientation` - (Optional) Device emulation orientation field. Valid values are `LANDSCAPE` and `PORTRAIT`.
* `device_type` - (Optional) Device emulation type field. Valid values are `MOBILE` and `TABLET`.

#### Deprecated runtime

If you want to use the legacy runtime you can set the `runtime_type`, `runtime_type_version` and `script_language` to empty string `""`. 

-> **WARNING** Support for legacy Synthetics runtimes **will reach its end-of-life (EOL) on October 22, 2024**. In addition, creating **_new_** monitors using the legacy runtime **will no longer be supported after August 26, 2024**. In light of the above, kindly **upgrade your Synthetic Monitors to the new runtime** at the earliest, if they are still using the legacy runtime. Please check out [this page](https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm) for more details on the EOL, action needed (specific to monitors using public and private locations), relevant resources, and more. 

### Example Usage

```hcl
resource "newrelic_synthetics_monitor" "synthetic-simple-browser" {
    status           = "ENABLED"
    name             = "test simple browser"
    period           = "EVERY_MINUTE"
    uri              = "https://www.one.newrelic.com"
    type             = "BROWSER"
    locations_public = ["AP_NORTHEAST_1"]

    lifecycle {
        ignore_changes = [
          script_language,
          runtime_type,
          runtime_type_version
        ]
    }
}
```

### Nested `custom header` blocks

All nested `custom_header` blocks support the following common arguments:

* `name` - (Required) Header name.
* `value` - (Required) Header Value.

### Nested `tag` blocks

All nested `tag` blocks support the following common arguments:

* `key` - (Required) Name of the tag key.
* `values` - (Required) Values associated with the tag key.

## Additional Examples

### Create a monitor with a private location

The below example shows how you can define a private location and attach it to a monitor.

-> **NOTE:** It can take up to 10 minutes for a private location to become available.

##### Type: `SIMPLE`

```hcl
resource "newrelic_synthetics_private_location" "location" {
  description               = "Example private location"
  name                      = "private_location"
  verified_script_execution = false
}

resource "newrelic_synthetics_monitor" "monitor" {
  status           = "ENABLED"
  name             = "monitor"
  period           = "EVERY_MINUTE"
  uri              = "https://www.one.newrelic.com"
  type             = "SIMPLE"
  locations_private = [newrelic_synthetics_private_location.location.id]

  custom_header {
    name  = "some_name"
    value = "some_value"
  }

  treat_redirect_as_failure = true
  validation_string         = "success"
  bypass_head_request       = true
  verify_ssl                = true

  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```
##### Type: `BROWSER`

```hcl
resource "newrelic_synthetics_private_location" "location" {
  description               = "Example private location"
  name                      = "private-location"
  verified_script_execution = false
}

resource "newrelic_synthetics_monitor" "monitor" {
  status            = "ENABLED"
  type              = "BROWSER"
  uri               = "https://www.one.newrelic.com"
  name              = "monitor"
  period            = "EVERY_MINUTE"
  locations_private = [newrelic_synthetics_private_location.location.id]

  custom_header {
    name  = "some_name"
    value = "some_value"
  }

  enable_screenshot_on_failure_and_script = true
  validation_string                       = "success"
  verify_ssl                              = true
  runtime_type_version                    = "100"
  runtime_type                            = "CHROME_BROWSER"
  script_language                         = "JAVASCRIPT"

  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The ID (GUID) of the Synthetics monitor that the script is attached to.
* `period_in_minutes` - The interval in minutes at which Synthetic monitor should run.

## Import

Synthetics monitor can be imported using the `guid`, e.g.

```bash
$ terraform import newrelic_synthetics_monitor.monitor <guid>
```
