---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_step_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-step-monitor"
description: |-
    Create and manage a Synthetics Step monitor in New Relic.
---

# Resource: newrelic\_synthetics\_step\_monitor

Use this resource to create, update, and delete a Synthetics Step monitor in New Relic.

-> **IMPORTANT:**  The **Synthetics Legacy Runtime** reached <b style="color:red;">end-of-life</b> on <b style="color:red;">October 22, 2024</b>. Using the legacy runtime or blank runtime values with requests from the New Relic Terraform Provider will result in API errors. Starting with **v3.51.0** of the New Relic Terraform Provider, configurations without runtime attributes or comprising legacy runtime values <span style="color:red;">will be deemed invalid</span>. If your Synthetic monitors' configuration is not updated already with new runtime values, these consequences would be seen; upgrade as soon as possible to avoid consequences. For more details and instructions, please see detailed warning against `runtime_type` and `runtime_type_version` in the [**Argument Reference**](#runtime_type) section.

## Example Usage

```hcl
resource "newrelic_synthetics_step_monitor" "foo" {
  name                                    = "Sample Step Monitor"
  enable_screenshot_on_failure_and_script = true
  locations_public                        = ["US_EAST_1", "US_EAST_2"]
  period                                  = "EVERY_6_HOURS"
  status                                  = "ENABLED"
  runtime_type                            = "CHROME_BROWSER"
  runtime_type_version                    = "100"
  devices                                 = ["DESKTOP", "MOBILE_PORTRAIT", "TABLET_LANDSCAPE"]
  browsers                                = ["CHROME"]
  steps {
    ordinal = 0
    type    = "NAVIGATE"
    values  = ["https://www.newrelic.com"]
  }
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following are the common arguments supported for `STEP` monitor:

* `account_id`- (Optional) The account in which the Synthetics monitor will be created.
* `name` - (Required) The name for the monitor.
* `locations_public` - (Required) The location the monitor will run from. Check out [this page](https://docs.newrelic.com/docs/synthetics/synthetic-monitoring/administration/synthetic-public-minion-ips/) for a list of valid public locations. You don't need the `AWS_` prefix as the provider uses NerdGraph. At least one of either `locations_public` or `location_private` is required.
* `location_private` - (Required) The location the monitor will run from. At least one of `locations_public` or `location_private` is required. See [Nested locations_private blocks](#nested-locations-private-blocks) below for details.
* `period` - (Required) The interval at which this monitor should run. Valid values are `EVERY_MINUTE`, `EVERY_5_MINUTES`, `EVERY_10_MINUTES`, `EVERY_15_MINUTES`, `EVERY_30_MINUTES`, `EVERY_HOUR`, `EVERY_6_HOURS`, `EVERY_12_HOURS`, or `EVERY_DAY`.
* `status` - (Required) The run state of the monitor. (`ENABLED` or `DISABLED`).
* `browsers` - (Optional) The multiple browsers list on which synthetic monitors will run. Valid values are `CHROME` and `FIREFOX`.
* `devices` - (Optional) The multiple devices list on which synthetic monitors will run. Valid values are `DESKTOP`, `MOBILE_LANDSCAPE`, `MOBILE_PORTRAIT`, `TABLET_LANDSCAPE` and `TABLET_PORTRAIT`.
* `steps` - (Required) The steps that make up the script the monitor will run. See [Nested steps blocks](#nested-steps-blocks) below for details.
* `tag` - (Optional) The tags that will be associated with the monitor. See [Nested tag blocks](#nested-tag-blocks) below for details.
* `runtime_type` - (Optional) The runtime that the monitor will use to run jobs.
* `runtime_type_version` - (Optional) The specific version of the runtime type selected.

-> **NOTE:** Currently, the values of `runtime_type` and `runtime_type_version` supported by this resource are `CHROME_BROWSER` and `100` respectively. In order to run the monitor in the new runtime, both `runtime_type` and `runtime_type_version` need to be specified; however, specifying neither of these attributes would set this monitor to use the legacy runtime. It may also be noted that the runtime opted for would only be effective with private locations. For public locations, all traffic has been shifted to the new runtime, irrespective of the selection made.

-> **WARNING:**  The <b style="color:red;">end-of-life</b> of the **Synthetics Legacy Runtime** took effect on <b style="color:red;">October 22, 2024</b>, implying that support for using the deprecated Synthetics Legacy Runtime with **new and existing** Synthetic monitors <b style="color:maroon;">officially ended as of October 22, 2024</b>. As a consequence of this API change, all requests associated with Synthetic Monitors (except Ping Monitors) going out of the New Relic Terraform Provider <span style="color:maroon;">will be blocked by an API error</span> if they include values corresponding to the legacy runtime or blank runtime values.
<br><br>
Following these changes, starting with <b style="color:red;">v3.51.0</b> of the New Relic Terraform Provider, configuration of **new and existing** Synthetic monitors without runtime attributes (or) comprising runtime attributes signifying the legacy runtime <span style="color:red;">will be deemed invalid</span> (this applies to all Synthetic monitor resources, except `newrelic_synthetics_monitor` with type `SIMPLE`). If your monitors' configuration <span style="color:red;">is not updated with new runtime values</span>, you will see the consequences stated here. New Synthetic monitors created after August 26, 2024 already adhere to these restrictions, as part of the first phase of the EOL.
<br><br>
We kindly recommend that you upgrade your Synthetic Monitors to the new runtime as soon as possible <span style="color:red;">if they are still using the legacy runtime</span>, to avoid seeing the aforementioned consequences. Please check out [this guide](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/guides/synthetics_legacy_runtime_eol_migration_guide) in the documentation of the Terraform Provider (specifically, the table at the bottom of the guide, if you're looking for updates to be made to the configuration of Synthetic monitors) and [this announcement](https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm) for more details on the EOL, actions needed, relevant resources, and more.
<br><br>
You would not be affected by the EOL if your Synthetic monitors' Terraform configuration comprises new runtime values.

### Nested `location private` blocks

All nested `location_private` blocks support the following common arguments:

* `guid` - (Required) The unique identifier for the Synthetics private location in New Relic.
* `vse_password` - (Optional) The location's Verified Script Execution password, only necessary if Verified Script Execution is enabled for the location.

### Nested `steps` blocks

All nested `steps` blocks support the following common arguments:

* `ordinal` - (Required) The position of the step within the script ranging from 0-100.
* `type` - (Required) Name of the tag key. Valid values are `ASSERT_ELEMENT`, `ASSERT_MODAL`, `ASSERT_TEXT`, `ASSERT_TITLE`, `CLICK_ELEMENT`, `DISMISS_MODAL`, `DOUBLE_CLICK_ELEMENT`, `HOVER_ELEMENT`, `NAVIGATE`, `SECURE_TEXT_ENTRY`, `SELECT_ELEMENT`, `TEXT_ENTRY`.
* `values` - (Optional) The metadata values related to the step.

### Nested `tag` blocks

All nested `tag` blocks support the following common arguments:

* `key` - (Required) Name of the tag key.
* `values` - (Required) Values associated with the tag key.

## Additional Examples

### Create a monitor with a private location

The below example shows how you can define a private location and attach it to a monitor.

-> **NOTE:** It can take up to 10 minutes for a private location to become available.

```hcl
resource "newrelic_synthetics_private_location" "foo" {
  name                      = "Sample Private Location"
  description               = "Sample Private Location Description"
  verified_script_execution = true
}

resource "newrelic_synthetics_step_monitor" "foo" {
  name                 = "Sample Step Monitor"
  period               = "EVERY_6_HOURS"
  status               = "ENABLED"
  runtime_type         = "CHROME_BROWSER"
  runtime_type_version = "100"
  devices              = ["DESKTOP", "MOBILE_PORTRAIT", "TABLET_LANDSCAPE"]
  browsers             = ["CHROME"]
  location_private {
    guid         = newrelic_synthetics_private_location.foo.id
    vse_password = "secret"
  }
  steps {
    ordinal = 0
    type    = "NAVIGATE"
    values  = ["https://google.com"]
  }
  tag {
    key    = "some_key"
    values = ["some_value"]
  }
}
```
## Attributes Reference

The following attributes are exported:

* `id` - The ID (GUID) of the synthetics step monitor.
* `period_in_minutes` - The interval in minutes at which Synthetic monitor should run.

## Import

A step monitor can be imported using its GUID, using the following command.

```bash
$ terraform import newrelic_synthetics_step_monitor.monitor <guid>
```
