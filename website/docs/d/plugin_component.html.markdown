---
layout: "newrelic"
page_title: "New Relic: newrelic_plugin_component"
sidebar_current: "docs-newrelic-datasource-plugin-component"
description: |-
  Looks up the information about a plugin component in New Relic.
---

# newrelic\_plugin\_component

Use this data source to get information about a single plugin component in New Relic.

Each plugin component reporting into to New Relic is assigned a unique ID. Once you have a plugin component reporting data into your account, its component ID can be used to create Plugins Alert Conditions.

## Example Usage

```hcl
data "newrelic_plugin" "foo" {
  guid = "com.example.my-plugin"
}

data "newrelic_plugin_component" "foo" {
  plugin_id = "${data.newrelic_plugin.foo.id}"
  name = "My Plugin Component"
}

resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_plugins_alert_condition" "foo" {
  policy_id          = "${newrelic_alert_policy.foo.id}"
  name               = "foo"
  metric             = "Component/Summary/Consumers[consumers]"
  plugin_id          = "${data.newrelic_plugin.foo.id}"
  plugin_guid        = "${data.newrelic_plugin.foo.guid}"
  entities           = ["${data.newrelic_plugin_component.foo.id}"]
  value_function     = "average"
  metric_description = "Queue consumers"

  term {
    duration      = 5
    operator      = "below"
    priority      = "critical"
    threshold     = "0.75"
    time_function = "all"
  }
}
```

## Argument Reference

The following arguments are supported:

* `plugin_id` - (Required) The ID of the plugin instance this component belongs to.
* `name` - (Required) The name of the plugin component.

## Attributes Reference
* `id` - The ID of the plugin component.
* `health_status` - The health status of the plugin component.