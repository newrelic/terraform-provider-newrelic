---
layout: "newrelic"
page_title: "New Relic Terraform Provider Alert Conditions Migration Guide"
sidebar_current: "docs-newrelic-provider-alert-conditions-migration-guide"
description: |-
  Use this guide to migrate from the deprecated 'newrelic_synthetics_alert_condition' and 'newrelic_infra_alert_condition' onto the 'newrelic_nrql_alert_condition' resource.
---

## Migrating to NRQL Alert Conditions

Certain subtypes of Alert Conditions (Synthetics Alert Condition and Infra Alert Condition) have been removed in favor of NRQL Alert Conditions.

Users wanting to migrate alert conditions will need to make a few adjustments to their configuration, by following the examples outlined below.

### Migrating from Synthetics Alert Conditions to NRQL Alert Conditions

The following example illustrates changing over from a synthetics alert condition, i.e. [newrelic_synthetics_alert_condition](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/nrql_alert_condition) to an NRQL-based alert condition using the [newrelic_nrql_alert_condition](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/nrql_alert_condition) resource.

It may be obtained from the example that the ID of the monitor specified against the argument `monitor_id` in the synthetics alert condition may be used in NRQL with `SyntheticCheck` to convert this into an NRQL-based alert condition.

Example newrelic_synthetics_alert_condition:
```
resource "newrelic_synthetics_alert_condition" "foo" {
  policy_id = newrelic_alert_policy.foo.id

  name        = "foo synthetics condition"
  monitor_id  = newrelic_synthetics_monitor.foo.id
  runbook_url = "https://www.example.com"
}
```

Migration to newrelic_nrql_alert_condition:
```
resource "newrelic_nrql_alert_condition" "foo" {
  policy_id                      = newrelic_alert_policy.foo.id
  type                           = "static"
  name                           = "foo synthetics condition"
  description                    = "Alert when transactions are taking too long"
  runbook_url                    = "https://www.example.com"
  enabled                        = true
  violation_time_limit_seconds   = 3600
  fill_option                    = "static"
  fill_value                     = 1.0
  aggregation_window             = 60
  aggregation_method             = "event_flow"
  aggregation_delay              = 120
  expiration_duration            = 120
  open_violation_on_expiration   = true
  close_violations_on_expiration = true
  slide_by                       = 30

  nrql {
    query = "SELECT filter(count(*), WHERE result = 'FAILED') FROM SyntheticCheck WHERE NOT isMuted AND entityGuid IN ('${newrelic_synthetics_monitor.foo.id}') FACET entityGuid, location"
  }

  critical {
    operator              = "above"
    threshold             = 0
    threshold_duration    = 60
    threshold_occurrences = "ALL"
  }
}
```

### Migrating from Infra Alert Conditions to NRQL Alert Conditions

The following examples illustrate changing over from infra alert conditions, i.e. [newrelic_infra_alert_condition](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/infra_alert_condition) to NRQL-based alert conditions using the resource [newrelic_nrql_alert_condition](https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/nrql_alert_condition).

#### newrelic_infra_alert_condition: High Disk Usage

The following example illustrates changing an infra alert condition for `High Disk Usage` to an NRQL-based alert condition.

Example newrelic_infra_alert_condition:
```
resource "newrelic_infra_alert_condition" "high_disk_usage" {
  policy_id = newrelic_alert_policy.foo.id

  name        = "High disk usage"
  description = "Warning if disk usage goes above 80% and critical alert if goes above 90%"
  type        = "infra_metric"
  event       = "StorageSample"
  select      = "diskUsedPercent"
  comparison  = "above"
  where       = "(hostname LIKE '%frontend%')"

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
```

Migration to newrelic_nrql_alert_condition:
```
resource "newrelic_nrql_alert_condition" "foo" {
  policy_id                      = newrelic_alert_policy.foo.id
  type                           = "static"
  name                           = "High disk usage"
  description                    = "Warning if disk usage goes above 80% and critical alert if goes above 90%"
  enabled                        = true
  violation_time_limit_seconds   = 3600
  fill_option                    = "static"
  fill_value                     = 1.0
  aggregation_window             = 60
  aggregation_method             = "event_flow"
  aggregation_delay              = 120
  expiration_duration            = 120
  open_violation_on_expiration   = true
  close_violations_on_expiration = true
  slide_by                       = 30

  nrql {
    query = "SELECT average(diskUsedPercent) FROM StorageSample WHERE (hostname LIKE '%frontend%') FACET entityAndMountPoint, mountPoint"
  }

  critical {
    operator              = "above"
    threshold             = 90
    threshold_duration    = 1500
    threshold_occurrences = "ALL"
  }

  warning {
    operator              = "above"
    threshold             = 80
    threshold_duration    = 600
    threshold_occurrences = "ALL"
  }
}
```

#### newrelic_infra_alert_condition: High DB Connection Count

The following example illustrates changing an infra alert condition for `High DB Connection Count` to an NRQL-based alert condition.

Example newrelic_infra_alert_condition:
```
resource "newrelic_infra_alert_condition" "high_db_conn_count" {
  policy_id = newrelic_alert_policy.foo.id

  name        = "High database connection count"
  description = "Critical alert when the number of database connections goes above 90"
  type        = "infra_metric"
  event       = "DatastoreSample"
  select      = "provider.databaseConnections.Average"
  comparison  = "above"
  where       = "(hostname LIKE '%db%')"
  integration_provider = "RdsDbInstance"

  critical {
    duration      = 25
    value         = 90
    time_function = "all"
  }
}
```

Migration to newrelic_nrql_alert_condition:
```
resource "newrelic_nrql_alert_condition" "foo" {
  policy_id                      = newrelic_alert_policy.foo.id
  type                           = "static"
  name                           = "High database connection count"
  description                    = "Critical alert when the number of database connections goes above 90"
  enabled                        = true
  violation_time_limit_seconds   = 3600
  fill_option                    = "static"
  fill_value                     = 1.0
  aggregation_window             = 60
  aggregation_method             = "event_flow"
  aggregation_delay              = 120
  expiration_duration            = 120
  open_violation_on_expiration   = true
  close_violations_on_expiration = true
  slide_by                       = 30

  nrql {
    query = "SELECT `average`(`aws.rds.DatabaseConnections`) FROM Metric WHERE metricName IN ('aws.rds.DatabaseConnections') AND (`hostname` LIKE '%db%') FACET entity.guid"
  }

  critical {
    operator              = "above"
    threshold             = 90
    threshold_duration    = 1500
    threshold_occurrences = "ALL"
  }
}
```

#### newrelic_infra_alert_condition: Process Not Running

The following example illustrates changing an infra alert condition for `Process Not Running` to an NRQL-based alert condition.

Example newrelic_infra_alert_condition:
```
resource "newrelic_infra_alert_condition" "process_not_running" {
  policy_id = newrelic_alert_policy.foo.id

  name             = "Process not running (/usr/bin/ruby)"
  description      = "Critical alert when ruby isn't running"
  type             = "infra_process_running"
  comparison       = "equal"
  where            = "hostname = 'web01'"
  process_where    = "commandName = '/usr/bin/ruby'"

  critical {
    duration = 5
    value    = 0
  }
}
```

Migration to newrelic_nrql_alert_condition:
```
resource "newrelic_nrql_alert_condition" "foo" {
  policy_id                      = newrelic_alert_policy.foo.id
  type                           = "static"
  name                           = "Process not running (/usr/bin/ruby)"
  description                    = "Critical alert when ruby isn't running"
  enabled                        = true
  violation_time_limit_seconds   = 3600
  fill_option                    = "static"
  fill_value                     = 1.0
  aggregation_window             = 60
  aggregation_method             = "event_flow"
  aggregation_delay              = 120
  expiration_duration            = 120
  open_violation_on_expiration   = true
  close_violations_on_expiration = true
  slide_by                       = 30

  nrql {
    query = "SELECT filter(uniqueCount(processId), WHERE commandName = '/usr/bin/ruby') FROM ProcessSample WHERE hostname IS NOT NULL AND hostname = 'web01' FACET entityGuid"
  }

  critical {
    operator              = "equals"
    threshold             = 0
    threshold_duration    = 300
    threshold_occurrences = "ALL"
  }
}
```

#### newrelic_infra_alert_condition: Host Not Reporting

The following example illustrates changing an infra alert condition for `Host Not Reporting` to an NRQL-based alert condition.

Example newrelic_infra_alert_condition:
```
resource "newrelic_infra_alert_condition" "host_not_reporting" {
  policy_id = newrelic_alert_policy.foo.id

  name        = "Host not reporting"
  description = "Critical alert when the host is not reporting"
  type        = "infra_host_not_reporting"
  where       = "(hostname LIKE '%frontend%')"

  critical {
    duration = 5
  }
}
```

Migration to newrelic_nrql_alert_condition:
```
resource "newrelic_nrql_alert_condition" "foo" {
  policy_id                      = newrelic_alert_policy.foo.id
  type                           = "static"
  name                           = "Host not reporting"
  description                    = "Alert when the host is not reporting"
  enabled                        = true
  violation_time_limit_seconds   = 3600
  fill_option                    = "static"
  fill_value                     = 1.0
  aggregation_window             = 60
  aggregation_method             = "event_flow"
  aggregation_delay              = 120
  expiration_duration            = 120
  open_violation_on_expiration   = true
  close_violations_on_expiration = true
  slide_by                       = 30

  nrql {
    query = "SELECT `count`(`host.cpuPercent`) FROM `Metric` WHERE (((`metricName` = 'host.cpuPercent') AND NOT (`host`.`hostname` IS NULL)) AND (`host.hostname` LIKE '%frontend%')) FACET `entity`.`guid`"
  }

  critical {
    operator              = "equals"
    threshold             = 0
    threshold_duration    = 300
    threshold_occurrences = "ALL"
  }
}
```
