---
layout: "newrelic"
page_title: "New Relic: newrelic_infra_alert_condition"
sidebar_current: "docs-newrelic-resource-infra-alert-condition"
description: |-
  Create and manage an Infrastructure alert condition for a policy in New Relic.
---

# Resource: newrelic\_infra_alert\_condition

Use this resource to create and manage Infrastructure alert conditions in New Relic.

## Example Usage

```hcl
resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_infra_alert_condition" "high_disk_usage" {
  policy_id = newrelic_alert_policy.foo.policy_id

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

  warning {
    duration      = 10
    value         = 80
    time_function = "all"
  }
}

resource "newrelic_infra_alert_condition" "high_db_conn_count" {
  policy_id = newrelic_alert_policy.foo.policy_id

  name       = "High database connection count"
  type       = "infra_metric"
  event      = "DatastoreSample"
  select     = "provider.databaseConnections.Average"
  comparison = "above"
  where      = "(`hostname` LIKE '%db%')"
  integration_provider = "RdsDbInstance"

  critical {
    duration      = 25
    value         = 90
    time_function = "all"
  }
}

resource "newrelic_infra_alert_condition" "process_not_running" {
  policy_id = newrelic_alert_policy.foo.policy_id

  name             = "Process not running (/usr/bin/ruby)"
  type             = "infra_process_running"
  comparison       = "equal"
  process_where    = "`commandName` = '/usr/bin/ruby'"

  critical {
    duration      = 5
    value         = 0
  }
}

resource "newrelic_infra_alert_condition" "host_not_reporting" {
  policy_id = newrelic_alert_policy.foo.policy_id

  name       = "Host not reporting"
  type       = "infra_host_not_reporting"
  where      = "(`hostname` LIKE '%frontend%')"

  critical {
    duration = 5
  }
}
```

## Argument Reference

The following arguments are supported:

  * `policy_id` - (Required) The ID of the alert policy where this condition should be used.
  * `name` - (Required) The Infrastructure alert condition's name.
  * `type` - (Required) The type of Infrastructure alert condition.  Valid values are  `infra_process_running`, `infra_metric`, and `infra_host_not_reporting`.
  * `event` - (Required) The metric event; for example, `SystemSample` or `StorageSample`.  Supported by the `infra_metric` condition type.
  * `select` - (Required) The attribute name to identify the metric being targeted; for example, `cpuPercent`, `diskFreePercent`, or `memoryResidentSizeBytes`.  The underlying API will automatically populate this value for Infrastructure integrations (for example `diskFreePercent`), so make sure to explicitly include this value to avoid diff issues.  Supported by the `infra_metric` condition type.
  * `comparison` - (Required) The operator used to evaluate the threshold value.  Valid values are `above`, `below`, and `equal`.  Supported by the `infra_metric` and `infra_process_running` condition types.
  * `critical` - (Required) Identifies the threshold parameters for opening a critical alert violation. See [Thresholds](#thresholds) below for details.
  * `warning` - (Optional) Identifies the threshold parameters for opening a warning alert violation. See [Thresholds](#thresholds) below for details.
  * `enabled` - (Optional) Whether the condition is turned on or off.  Valid values are `true` and `false`.  Defaults to `true`.
  * `where` - (Optional) If applicable, this identifies any Infrastructure host filters used; for example: `hostname LIKE '%cassandra%'`.
  * `process_where` - (Optional) Any filters applied to processes; for example: `commandName = 'java'`.  Supported by the `infra_process_running` condition type.
  * `integration_provider` - (Optional) For alerts on integrations, use this instead of `event`.  Supported by the `infra_metric` condition type.
  * `runbook_url` - (Optional) Runbook URL to display in notifications.
  * `violation_close_timer` - (Optional) Determines how much time will pass before a violation is automatically closed. Setting the time limit to 0 prevents a violation from being force-closed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The ID of the Infrastructure alert condition.
  * `created_at` - The timestamp the alert condition was created.
  * `updated_at` - The timestamp the alert condition was last updated.

## Thresholds

The `critical` and `warning` threshold mapping supports the following arguments:

  * `duration` - (Required) Identifies the number of minutes the threshold must be passed or met for the alert to trigger. Threshold durations must be between 1 and 60 minutes (inclusive).
  * `value` - (Optional) Threshold value, computed against the `comparison` operator. Supported by `infra_metric` and `infra_process_running` alert condition types.
  * `time_function` - (Optional) Indicates if the condition needs to be sustained or to just break the threshold once; `all` or `any`. Supported by the `infra_metric` alert condition type.


## Import

Infrastructure alert conditions can be imported using a composite ID of `<policy_id>:<condition_id>`, e.g.

```
$ terraform import newrelic_infra_alert_condition.main 12345:67890
```
