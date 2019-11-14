---
layout: "newrelic"
page_title: "New Relic: newrelic_infra_alert_condition"
sidebar_current: "docs-newrelic-resource-infra-alert-condition"
description: |-
  Create and manage an Infrastructure alert condition for a policy in New Relic.
---

# newrelic\_infra_alert\_condition

## Example Usage

```hcl
resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_infra_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name       = "High disk usage"
  type       = "infra_metric"
  event      = "StorageSample"
  select     = "diskUsedPercent"
  comparison = "above"
  where      = "(`hostname` LIKE '%frontend%')"

  critical {
    duration      = 25
    value         = 90
    time_function = "all"
  }
}
```

## Argument Reference

The following arguments are supported:

  * `policy_id` - (Required) The ID of the alert policy where this condition should be used.
  * `name` - (Required) The Infrastructure alert condition's name.
  * `enabled` - (Optional) Set whether to enable the alert condition. Defaults to `true`.
  * `type` - (Required) The type of Infrastructure alert condition: "infra_process_running", "infra_metric", or "infra_host_not_reporting".
  * `event` - (Required) The metric event; for example, system metrics, process metrics, storage metrics, or network metrics.
  * `select` - (Required) The attribute name to identify the type of metric condition; for example, "network", "process", "system", or "storage".
  * `comparison` - (Required) The operator used to evaluate the threshold value; "above", "below", "equal".
  * `critical` - (Required) Identifies the critical threshold parameters for triggering an alert notification. See [Thresholds](#thresholds) below for details.
  * `warning` - (Optional) Identifies the warning threshold parameters. See [Thresholds](#thresholds) below for details.
  * `where` - (Optional) Infrastructure host filter for the alert condition.
  * `process_where` - (Optional) Any filters applied to processes; for example: `"commandName = 'java'"`.
  * `integration_provider` - (Optional) For alerts on integrations, use this instead of `event`.
  * `runbook_url` - (Optional) Runbook URL to display in notifications.

## Thresholds

The `critical` and `warning` threshold mapping supports the following arguments:

  * `duration` - (Required) Identifies the number of minutes the threshold must be passed or met for the alert to trigger. Threshold durations must be between 1 and 60 minutes (inclusive).
  * `value` - (Optional) Threshold value, computed against the `comparison` operator. Supported by "infra_metric" and "infra_process_running" alert condition types.
  * `time_function` - (Optional) Indicates if the condition needs to be sustained or to just break the threshold once; `all` or `any`. Supported by the "infra_metric" alert condition type.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the Infrastructure alert condition.
