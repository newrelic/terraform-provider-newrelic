---
layout: "newrelic"
page_title: "New Relic: newrelic_plugins_alert_condition"
sidebar_current: "docs-newrelic-resource-plugins-alert-condition"
description: |-
  Create and manage a Plugins alert condition for a policy in New Relic.
---

# Resource: newrelic\_plugins\_alert\_condition

Use this resource to create and manage plugins alert conditions in New Relic.

## Example Usage

```hcl
data "newrelic_plugin" "foo" {
  guid = "com.example.my-plugin"
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

  * `policy_id` - (Required) The ID of the policy where this condition should be used.
  * `name` - (Required) The title of the condition. Must be between 1 and 64 characters, inclusive.
  * `metric` - (Required) The metric field accepts parameters based on the `type` set.
  * `plugin_id` - (Required) The ID of the installed plugin instance which produces the metric.
  * `plugin_guid` - (Required) The GUID of the plugin which produces the metric.
  * `runbook_url` - (Optional) Runbook URL to display in notifications.
  * `term` - (Required) A list of terms for this condition. See [Terms](#terms) below for details.

## Terms

The `term` mapping supports the following arguments:

  * `duration` - (Required) In minutes, must be in the range of `5` to `120`, inclusive.
  * `operator` - (Optional) `above`, `below`, or `equal`.  Defaults to `equal`.
  * `priority` - (Optional) `critical` or `warning`.  Defaults to `critical`.
  * `threshold` - (Required) Must be 0 or greater.
  * `time_function` - (Required) `all` or `any`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The ID of the alert condition.

## Import

Alert conditions can be imported using the `id`, e.g.

```
$ terraform import newrelic_plugins_alert_condition.main 12345
```
