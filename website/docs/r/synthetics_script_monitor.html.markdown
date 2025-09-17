---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_script_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-script-monitor"
description: |-
    Create and manage a Synthetics script monitor in New Relic.
---

# Resource: newrelic\_synthetics\_script\_monitor

Use this resource to create update, and delete a Script API or Script Browser Synthetics Monitor in New Relic.

-> **IMPORTANT:**  The **Synthetics Legacy Runtime** has reached its <b style="color:red;">end-of-life</b> on <b style="color:red;">October 22, 2024</b>. As a consequence, using the legacy runtime or blank runtime values with Synthetic monitor requests from the New Relic Terraform Provider will result in API errors. Starting with **v3.51.0** of the New Relic Terraform Provider, configurations of Synthetic monitors without runtime attributes or comprising legacy runtime values <span style="color:red;">will be deemed invalid</span>.
<br><br>
If your Synthetic monitors' configuration is not updated already with new runtime values, upgrade as soon as possible to avoid these consequences. For more details and instructions, please see the detailed warning in the [**Deprecated Runtime**](#deprecated-runtime) section.


## Example Usage

##### Type: `SCRIPT_API`

```hcl
resource "newrelic_synthetics_script_monitor" "monitor" {
  status               = "ENABLED"
  name                 = "script_monitor"
  type                 = "SCRIPT_API"
  locations_public     = ["AP_SOUTH_1", "AP_EAST_1"]
  period               = "EVERY_6_HOURS"
  
  script               = "console.log('it works!')"
  
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
  status                                  = "ENABLED"
  name                                    = "script_monitor"
  type                                    = "SCRIPT_BROWSER"
  locations_public                        = ["AP_SOUTH_1", "AP_EAST_1"]
  period                                  = "EVERY_HOUR"
  script                                  = "$browser.get('https://one.newrelic.com')"
  runtime_type_version                    = "100"
  runtime_type                            = "CHROME_BROWSER"
  script_language                         = "JAVASCRIPT"
  devices                                 = ["DESKTOP", "MOBILE_PORTRAIT", "TABLET_LANDSCAPE"]
  browsers                                = ["CHROME"]
  enable_screenshot_on_failure_and_script = false
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following are the common arguments supported for `SCRIPT_API` and `SCRIPT_BROWSER` monitors:

* `account_id`- (Optional) The account in which the Synthetics monitor will be created.
* `status` - (Required) The run state of the monitor. (`ENABLED` or `DISABLED`).
* `name` - (Required) The name for the monitor.
* `type` - (Required) The plaintext representing the monitor script. Valid values are SCRIPT_BROWSER or SCRIPT_API
* `locations_public` - (Optional) The location the monitor will run from. Check out [this page](https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/) for a list of valid public locations. The `AWS_` prefix is not needed, as the provider uses NerdGraph. **At least one of either** `locations_public` **or** `location_private` **is required**.
* `location_private` - (Optional) The location the monitor will run from. See [Nested location_private blocks](#nested-location-private-blocks) below for details. **At least one of either** `locations_public` **or** `location_private` **is required**.
* `period` - (Required) The interval at which this monitor should run. Valid values are `EVERY_MINUTE`, `EVERY_5_MINUTES`, `EVERY_10_MINUTES`, `EVERY_15_MINUTES`, `EVERY_30_MINUTES`, `EVERY_HOUR`, `EVERY_6_HOURS`, `EVERY_12_HOURS`, or `EVERY_DAY`.
* `script` - (Required) The script that the monitor runs.
* `runtime_type` - (Optional) The runtime that the monitor will use to run jobs. For the `SCRIPT_API` monitor type, a valid value is `NODE_API`. For the `SCRIPT_BROWSER` monitor type, a valid value is `CHROME_BROWSER`.
* `runtime_type_version` - (Optional) The specific version of the runtime type selected. For the `SCRIPT_API` monitor type, a valid value is `16.10`, which corresponds to the version of Node.js. For the `SCRIPT_BROWSER` monitor type, a valid value is `100`, which corresponds to the version of the Chrome browser.
* `script_language` - (Optional) The programing language that should execute the script.
* `tag` - (Optional) The tags that will be associated with the monitor. See [Nested tag blocks](#nested-tag-blocks) below for details.

The `SCRIPTED_BROWSER` monitor type supports the following additional arguments:

* `enable_screenshot_on_failure_and_script` - (Optional) Capture a screenshot during job execution.
* `browsers` - (Optional) The multiple browsers list on which synthetic monitors will run. Valid values are `CHROME` and `FIREFOX`.
* `devices` - (Optional) The multiple devices list on which synthetic monitors will run. Valid values are `DESKTOP`, `MOBILE_LANDSCAPE`, `MOBILE_PORTRAIT`, `TABLET_LANDSCAPE` and `TABLET_PORTRAIT`.
* `device_orientation` - (Optional) Device emulation orientation field. Valid values are `LANDSCAPE` and `PORTRAIT`. We recommend you to use `devices` field instead of `device_type`,`device_orientation` fields, as it allows you to select multiple combinations of device types and orientations.
* `device_type` - (Optional) Device emulation type field. Valid values are `MOBILE` and `TABLET`. We recommend you to use `devices` field instead of `device_type`,`device_orientation` fields, as it allows you to select multiple combinations of device types and orientations.

#### Deprecated Runtime

~~If you want to use a legacy runtime (Node 10 or Chrome 72) you can set the `runtime_type`, `runtime_type_version` and `script_language` to empty string `""`.~~ 

-> **WARNING:**  The <b style="color:red;">end-of-life</b> of the **Synthetics Legacy Runtime** took effect on <b style="color:red;">October 22, 2024</b>, implying that support for using the deprecated Synthetics Legacy Runtime with **new and existing** Synthetic monitors <b style="color:maroon;">officially ended as of October 22, 2024</b>. As a consequence of this API change, all requests associated with Synthetic Monitors (except Ping Monitors) going out of the New Relic Terraform Provider <span style="color:maroon;">will be blocked by an API error</span> if they include values corresponding to the legacy runtime or blank runtime values.
<br><br>
Following these changes, starting with <b style="color:red;">v3.51.0</b> of the New Relic Terraform Provider, configuration of **new and existing** Synthetic monitors without runtime attributes (or) comprising runtime attributes signifying the legacy runtime <span style="color:red;">will be deemed invalid</span> (this applies to all Synthetic monitor resources, except `newrelic_synthetics_monitor` with type `SIMPLE`). If your monitors' configuration <span style="color:red;">is not updated with new runtime values</span>, you will see the consequences stated here. New Synthetic monitors created after August 26, 2024 already adhere to these restrictions, as part of the first phase of the EOL.
<br><br>
We kindly recommend that you upgrade your Synthetic Monitors to the new runtime as soon as possible <span style="color:red;">if they are still using the legacy runtime</span>, to avoid seeing the aforementioned consequences. Please check out [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/synthetics_legacy_runtime_eol_migration_guide) in the documentation of the Terraform Provider (specifically, the table at the bottom of the guide, if you're looking for updates to be made to the configuration of Synthetic monitors) and [this announcement](https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm) for more details on the EOL, actions needed, relevant resources, and more.
<br><br>
You would not be affected by the EOL if your Synthetic monitors' Terraform configuration comprises new runtime values.

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
resource "newrelic_synthetics_private_location" "location" {
  description               = "Example private location"
  name                      = "private_location"
  verified_script_execution = true
}

resource "newrelic_synthetics_script_monitor" "monitor" {
  status = "ENABLED"
  name   = "script_monitor"
  type   = "SCRIPT_API"
  location_private {
    guid         = newrelic_synthetics_private_location.location.id
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
resource "newrelic_synthetics_private_location" "location" {
  description               = "Test Description"
  name                      = "private_location"
  verified_script_execution = true
}

resource "newrelic_synthetics_script_monitor" "monitor" {
  status                                  = "ENABLED"
  name                                    = "script_monitor"
  type                                    = "SCRIPT_BROWSER"
  period                                  = "EVERY_HOUR"
  script                                  = "$browser.get('https://one.newrelic.com')"
  runtime_type_version                    = "100"
  runtime_type                            = "CHROME_BROWSER"
  script_language                         = "JAVASCRIPT"
  devices                                 = ["DESKTOP", "MOBILE_PORTRAIT", "TABLET_LANDSCAPE"]
  browsers                                = ["CHROME"]
  enable_screenshot_on_failure_and_script = false
  location_private {
    guid         = newrelic_synthetics_private_location.location.id
    vse_password = "secret"
  }
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```

### Create a monitor and a secure credential

The following example shows how to use `depends_on` to create a monitor that uses a new secure credential.
The `depends_on` creates an explicit dependency between resources to ensure that the secure credential is created before the monitor that uses it.

-> **NOTE:** Use the `depends_on` when you are creating both monitor and its secure credentials together.

##### Type: `SCRIPT_BROWSER`

```hcl
resource "newrelic_synthetics_script_monitor" "example_script_monitor" {
  name             = "script_monitor"
  type             = "SCRIPT_BROWSER"
  period           = "EVERY_HOUR"
  locations_public = ["US_EAST_1"] 
  status           = "ENABLED" 

  script = <<EOT
      var assert = require('assert');
      var secureCredential = $secure.TEST_SECURE_CREDENTIAL;
    EOT

  script_language      = "JAVASCRIPT"
  runtime_type         = "CHROME_BROWSER"
  runtime_type_version = "100"

  # this is where we introduce the dependency
  depends_on = [
    newrelic_synthetics_secure_credential.example_credential
  ]
}

resource "newrelic_synthetics_secure_credential" "example_credential" {
  key   = "TEST_SECURE_CREDENTIAL"
  value = "some_value"
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The ID (GUID) of the Synthetics script monitor.
* `period_in_minutes` - The interval in minutes at which Synthetic monitor should run.
* `monitor_id` - The monitor id of the Synthetics script monitor (not to be confused with the GUID of the monitor).

## Import

Synthetics monitor scripts can be imported using the `guid`, e.g.

```bash
$ terraform import newrelic_synthetics_script_monitor.monitor <guid>
```

