---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_script_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-script-monitor"
description: |-
    Create and manage a Synthetics script monitor in New Relic.
---

# Resource: newrelic\_synthetics\_script\_monitor

-> **WARNING:** Support for using the Synthetics Legacy Runtime (deprecated) with **new** Synthetic monitors **has officially ended as of August 26, 2024**. As a consequence, starting with v3.x.x of the New Relic Terraform Provider, **new** Synthetic monitors **will no longer be allowed to use the legacy runtime** (this applies to all Synthetic monitor resources). Additionally, as previously communicated by New Relic, the Synthetics Legacy Runtime **will reach its end-of-life (EOL) on October 22, 2024**. In light of the above, we kindly recommend that you upgrade your Synthetic Monitors to the new runtime at the earliest, if they are still using the legacy runtime. Please check out [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/synthetics_legacy_runtime_eol_migration_guide) in the documentation of the Terraform Provider and [this announcement](https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm) for more details on the EOL, action needed, relevant resources, and more.

Use this resource to create update, and delete a Script API or Script Browser Synthetics Monitor in New Relic.

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
  status           = "ENABLED"
  name             = "script_monitor"
  type             = "SCRIPT_BROWSER"
  locations_public = ["AP_SOUTH_1", "AP_EAST_1"]
  period           = "EVERY_HOUR"

  enable_screenshot_on_failure_and_script = false

  script = "$browser.get('https://one.newrelic.com')"

  runtime_type_version = "100"
  runtime_type         = "CHROME_BROWSER"
  script_language      = "JAVASCRIPT"

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

-> **WARNING:** As of February 29, 2024, Synthetic Monitors no longer support the `MUTED` status. Version **3.33.0** of the New Relic Terraform Provider is released to coincide with the `MUTED` status end-of-life. Consequently, the only valid values for `status` for all types of Synthetic Monitors are mentioned above. For additional information on alternatives to the `MUTED` status of Synthetic Monitors that can be managed via Terraform, please refer to [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/upcoming_synthetics_muted_status_eol_guide).
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

The `SCRIPTED_BROWSER` monitor type supports the following additional argument:

* `enable_screenshot_on_failure_and_script` - (Optional) Capture a screenshot during job execution.
* `device_orientation` - (Optional) Device emulation orientation field. Valid values are `LANDSCAPE` and `PORTRAIT`.
* `device_type` - (Optional) Device emulation type field. Valid values are `MOBILE` and `TABLET`.

#### Deprecated runtime

If you want to use a legacy runtime (Node 10 or Chrome 72) you can set the `runtime_type`, `runtime_type_version` and `script_language` to empty string `""`. 

<<<<<<< HEAD
-> **WARNING** Support for legacy Synthetics runtimes **will reach its end-of-life (EOL) on October 22, 2024**. In addition, creating **_new_** monitors using the legacy runtime **will no longer be supported after August 26, 2024**. In light of the above, kindly **upgrade your Synthetic Monitors to the new runtime** at the earliest, if they are still using the legacy runtime. Please check out [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/synthetics_legacy_runtime_eol_migration_guide) in the documentation of the Terraform Provider and [this announcement](https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm) for more details on the EOL, action needed, relevant resources, and more.
=======
-> **WARNING:** Support for using the Synthetics Legacy Runtime (deprecated) with **new** Synthetic monitors **has officially ended as of August 26, 2024**. As a consequence, starting with v3.x.x of the New Relic Terraform Provider, **new** Synthetic monitors **will no longer be allowed to use the legacy runtime** (this applies to all Synthetic monitor resources). Additionally, as previously communicated by New Relic, the Synthetics Legacy Runtime **will reach its end-of-life (EOL) on October 22, 2024**. In light of the above, we kindly recommend that you upgrade your Synthetic Monitors to the new runtime at the earliest, if they are still using the legacy runtime. Please check out [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/synthetics_legacy_runtime_eol_migration_guide) in the documentation of the Terraform Provider and [this announcement](https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm) for more details on the EOL, action needed, relevant resources, and more.
>>>>>>> 13b7ec92 (docs(synthetics): fill incomplete documentation, update LREOL warnings in resource docs)

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
  status = "ENABLED"
  name   = "script_monitor"
  type   = "SCRIPT_BROWSER"
  period = "EVERY_HOUR"
  script = "$browser.get('https://one.newrelic.com')"

  enable_screenshot_on_failure_and_script = false
  location_private {
    guid         = newrelic_synthetics_private_location.location.id
    vse_password = "secret"
  }

  runtime_type_version = "100"
  runtime_type         = "CHROME_BROWSER"
  script_language      = "JAVASCRIPT"
  device_type          = "MOBILE"
  device_orientation   = "LANDSCAPE"

  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Synthetics script monitor.
* `period_in_minutes` - The interval in minutes at which Synthetic monitor should run.

## Import

Synthetics monitor scripts can be imported using the `guid`, e.g.

```bash
$ terraform import newrelic_synthetics_script_monitor.monitor <guid>
```

