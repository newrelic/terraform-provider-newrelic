---
layout: "newrelic"
page_title: "New Relic: newrelic_synthetics_script_monitor"
sidebar_current: "docs-newrelic-resource-synthetics-script-monitor"
description: |-
Create and manage a Synthetics script monitor in New Relic.
---

# Resource: newrelic\_synthetics\_script\_monitor

Use this resource to create and manage New Relic synthetics script monitor.

## Example Usage

##### Type: `SCRIPT_API`
```hcl
    resource "newrelic_synthetics_script_monitor" "foo" {
     name = "SCRIPT_MONITOR"
     type = "SCRIPT_API"
     locations_public = ["AP_SOUTH_1","AP_EAST_1"]
     period = "EVERY_6_HOURS"
     status = "ENABLED"
     script = "console.log('terraform integration test updated')"
     script_language = "JAVASCRIPT"
     runtime_type = "NODE_API"
     runtime_type_version = "16.10"
     tag {
         key = "some_key"
         values = ["some_value"]
     }
}
```
##### Type: `SCRIPT_BROWSER`
```hcl
resource "newrelic_synthetics_script_monitor" "bar" {
			enable_screenshot_on_failure_and_script = false
			locations_public  = ["AP_SOUTH_1","AP_EAST_1"]
			name  = "SCRIPT_BROWSER"
			period  = "EVERY_HOUR"
			runtime_type_version  = "100"
			runtime_type  = "CHROME_BROWSER"
			script_language = "JAVASCRIPT"
			status  = "DISABLED"
			type  = "SCRIPT_BROWSER"
			script  = "$browser.get('https://one.newrelic.com')"
			tag {
				key = "Name"
				values  = ["scriptedMonitor","hello"]
			}
		}
```

## Argument Reference

The following are the common arguments supported for `SCRIPT_API` and `SCRIPT_BROWSER` monitors:

* `name` - (Required) The name for the monitor.
* `type` - (Required) The plaintext representing the monitor script.
* `locations_public` - (Required) The locations the monitor will run from.
* `period` - (Required) The interval at which the monitor runs in minutes.
* `runtime_type` - (Required) The runtime that the monitor will use to run jobs.
* `runtime_type_version` - (Required) The specific version of the runtime type selected.
* `script_language` - (Optional) The programing language that should execute the script.
* `status` - (Required) The run state of the monitor.
* `script` - (Required) The script that the monitor runs.
* `key` - (Required) Name of the tag key.
* `values` - (Required) Values associated with the tag key.

The `SCRIPTED_BROWSER` monitor type supports the following additional arguments:

* `enable_screenshot_on_failure_and_script` - (Optional) Capture a screenshot during job execution

```
Warning: This resource will use the account ID linked to your API key. At the moment it is not possible to dynamically set the account ID.
```

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Synthetics monitor that the script is attached to.

## Import

Synthetics monitor scripts can be imported using the `guid`, e.g.

```bash
$ terraform import newrelic_synthetics_monitor_script.bar <guid>
```

