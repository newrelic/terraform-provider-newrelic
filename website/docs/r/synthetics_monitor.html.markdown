---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-monitor"
description: |-
    Create and manage a Synthetics monitor in New Relic.
---

# Resource: newrelic\_synthetics\_monitor

Use this resource to create, update, and delete a Simple or Browser Synthetics Monitor in New Relic.

-> **IMPORTANT:**  The **Synthetics Legacy Runtime** has reached its <b style="color:red;">end-of-life</b> on <b style="color:red;">October 22, 2024</b>. As a consequence, using the legacy runtime or blank runtime values with Synthetic monitor requests from the New Relic Terraform Provider will result in API errors. Starting with **v3.51.0** of the New Relic Terraform Provider, configurations of Synthetic monitors without runtime attributes or comprising legacy runtime values <span style="color:red;">will be deemed invalid</span>.
<br><br>
If your Synthetic monitors' configuration is not updated already with new runtime values, upgrade as soon as possible to avoid these consequences. For more details and instructions, please see the detailed warning against `runtime_type` and `runtime_type_version` in the [**Argument Reference**](#runtime_type) section.


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
  status                                  = "ENABLED"
  name                                    = "monitor"
  period                                  = "EVERY_MINUTE"
  uri                                     = "https://www.one.newrelic.com"
  type                                    = "BROWSER"
  locations_public                        = ["AP_SOUTH_1"]
  enable_screenshot_on_failure_and_script = true
  validation_string                       = "success"
  verify_ssl                              = true
  runtime_type                            = "CHROME_BROWSER"
  runtime_type_version                    = "100"
  script_language                         = "JAVASCRIPT"
  devices                                 = ["DESKTOP", "TABLET_LANDSCAPE", "MOBILE_PORTRAIT"]
  browsers                                = ["CHROME"]
  custom_header {
    name  = "some_name"
    value = "some_value"
  }
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
* `script_language` - (Optional) The programing language that should execute the script.
* `browsers` - (Optional) The multiple browsers list on which synthetic monitors will run. Valid values are `CHROME` and `FIREFOX`.
* `devices` - (Optional) The multiple devices list on which synthetic monitors will run. Valid values are `DESKTOP`, `MOBILE_LANDSCAPE`, `MOBILE_PORTRAIT`, `TABLET_LANDSCAPE` and `TABLET_PORTRAIT`.
* `device_orientation` - (Optional) Device emulation orientation field. Valid values are `LANDSCAPE` and `PORTRAIT`. 
  * We recommend you to use `devices` field instead of `device_type`,`device_orientation` fields, as it allows you to select multiple combinations of device types and orientations.
* `device_type` - (Optional) Device emulation type field. Valid values are `MOBILE` and `TABLET`. 
  * We recommend you to use `devices` field instead of `device_type`,`device_orientation` fields, as it allows you to select multiple combinations of device types and orientations.
* `runtime_type` - (Optional) The runtime that the monitor will use to run jobs (`CHROME_BROWSER`).
* `runtime_type_version` - (Optional) The specific version of the runtime type selected (`100`).

#### Deprecated Runtime

-> **WARNING:**  The <b style="color:red;">end-of-life</b> of the **Synthetics Legacy Runtime** took effect on <b style="color:red;">October 22, 2024</b>, implying that support for using the deprecated Synthetics Legacy Runtime with **new and existing** Synthetic monitors <b style="color:maroon;">officially ended as of October 22, 2024</b>. As a consequence of this API change, all requests associated with Synthetic Monitors (except Ping Monitors) going out of the New Relic Terraform Provider <span style="color:maroon;">will be blocked by an API error</span> if they include values corresponding to the legacy runtime or blank runtime values.
<br><br>
Following these changes, starting with <b style="color:red;">v3.51.0</b> of the New Relic Terraform Provider, configuration of **new and existing** Synthetic monitors without runtime attributes (or) comprising runtime attributes signifying the legacy runtime <span style="color:red;">will be deemed invalid</span> (this applies to all Synthetic monitor resources, except `newrelic_synthetics_monitor` with type `SIMPLE`). If your monitors' configuration <span style="color:red;">is not updated with new runtime values</span>, you will see the consequences stated here. New Synthetic monitors created after August 26, 2024 already adhere to these restrictions, as part of the first phase of the EOL.
<br><br>
We kindly recommend that you upgrade your Synthetic Monitors to the new runtime as soon as possible <span style="color:red;">if they are still using the legacy runtime</span>, to avoid seeing the aforementioned consequences. Please check out [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/synthetics_legacy_runtime_eol_migration_guide) in the documentation of the Terraform Provider (specifically, the table at the bottom of the guide, if you're looking for updates to be made to the configuration of Synthetic monitors) and [this announcement](https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm) for more details on the EOL, actions needed, relevant resources, and more.
<br><br>
You would not be affected by the EOL if your Synthetic monitors' Terraform configuration comprises new runtime values.

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
  status                                  = "ENABLED"
  type                                    = "BROWSER"
  uri                                     = "https://www.one.newrelic.com"
  name                                    = "monitor"
  period                                  = "EVERY_MINUTE"
  locations_private                       = [newrelic_synthetics_private_location.location.id]
  enable_screenshot_on_failure_and_script = true
  validation_string                       = "success"
  verify_ssl                              = true
  runtime_type_version                    = "100"
  runtime_type                            = "CHROME_BROWSER"
  script_language                         = "JAVASCRIPT"
  devices                                 = ["DESKTOP", "TABLET_LANDSCAPE", "MOBILE_PORTRAIT"]
  browsers                                = ["CHROME"]
  custom_header {
    name  = "some_name"
    value = "some_value"
  }
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
