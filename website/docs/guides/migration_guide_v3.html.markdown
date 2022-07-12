---
layout: "newrelic"
page_title: "New Relic Terraform Provider v3.x Migration Guide"
sidebar_current: "docs-newrelic-provider-v3-migration-guide"
description: |-
  Use this guide to update the New Relic Terraform Provider from v2.x to v3.x
---

## Upgrade to v3.x of the New Relic Terraform Provider

Version 3.x of the provider uses a new underlying API for Synthetics. This results in some changes that will need to be made to existing Synthetics resource to keep them compatible with the new API.

### Migrating script Synthetics monitor resources

The previous `newrelic_synthetics_monitor_resourece` has been split into two new resources: `newrelic_synthetics_monitor_resource` and `newrelic_synthetics_script_monitor_resource`. Previously a monitor script had to be defined separately from a Synthetics Monitor using the `newrelic_synthetics_monitor_script_resource`. In this new version a script is attached directly to a scripted monitor. See the below example for how to migrate to version 3:

Previous  
```hcl
resource "newrelic_synthetics_monitor" "foo" {
  name = "monitor-name"
  type = "SCRIPT_BROWSER"
  frequency = 1
  status = "DISABLED"
  locations = ["AWS_US_EAST_1"]
  uri = "https://google.com"
}
resource "newrelic_synthetics_monitor_script" "foo_script" {
  monitor_id = newrelic_synthetics_monitor.foo.id
  text = "console.log("hello, world")"
	location {
		name = "AWS_US_EAST_1"
	}
}
```

Current
```hcl
resource "newrelic_synthetics_script_monitor" "foo" {
  name	          =	"monitor-name"
  type	          =	"SCRIPT_API"
 	location_public	=	["US_EAST_1"]
 	period	        =	"EVERY_HOUR"
 	status	        =	"ENABLED"
 	script	        =	"console.log("hello, world")"
 	tag {
 		key	          =	"some_key"
 		values	      =	["some_value"]
 	}
}
```

#### Steps to migrate to new resource

1. Move the value in `text` from `newrelic_synthetics_monitor_script` to `script` in `newrelic_synthetics_script_monitor`
2. Remove `AWS_` from the location name, e.g. `AWS_US_EAST_1` becomes `US_EAST_1`
