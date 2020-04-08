---
layout: "newrelic"
page_title: "New Relic: newrelic_plugin"
sidebar_current: "docs-newrelic-datasource-plugin"
description: |-
  Looks up the information about a plugin in New Relic.
---

# Data Source: newrelic\_plugin

Use this data source to get information about a specific installed plugin in New Relic. More information on Terraform's data sources can be found [here](https://www.terraform.io/docs/configuration/data-sources.html).

Each plugin published to New Relic's Plugin Central is assigned a [GUID](https://docs.newrelic.com/docs/plugins/plugin-developer-resources/planning-your-plugin/parts-plugin#guid). Once you have installed a plugin into your account it is assigned an ID. This account-specific ID is required when creating Plugins alert conditions.

## Example Usage

```hcl
data "newrelic_plugin" "foo" {
  guid = "com.example.my-plugin"
}

resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_plugins_alert_condition" "foo" {
  policy_id          = newrelic_alert_policy.foo.id
  name               = "foo"
  metric             = "Component/Summary/Consumers[consumers]"
  plugin_id          = data.newrelic_plugin.foo.id
  plugin_guid        = data.newrelic_plugin.foo.guid
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

* `guid` - (Required) The GUID of the plugin in New Relic.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the installed plugin instance.