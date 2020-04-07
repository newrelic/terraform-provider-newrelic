---
layout: 'newrelic'
page_title: 'New Relic: newrelic_nrql_alert_condition'
sidebar_current: 'docs-newrelic-resource-nrql-alert-condition'
description: |-
  Create and manage a NRQL alert condition for a policy in New Relic.
---

# Resource: newrelic_nrql_alert_condition

Use this resource to create and manage NRQL alert conditions in New Relic.

-> **IMPORTANT!** Version 2.0.0 of the New Relic Terraform Provider introduces some [additional requirements](/docs/providers/newrelic/index.html) for configuring the provider.
<br><br>
Before upgrading to version 2.0.0 or later, it is recommended to upgrade to the most recent 1.x version of the provider (version 1.16.0) and ensure that your environment successfully runs `terraform plan` without unexpected changes.

## Example Usage

##### Type: `static` (default)
```hcl
resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_nrql_alert_condition" "foo" {
  policy_id            = newrelic_alert_policy.foo.id
  type                 = "static"
  name                 = "foo"
  description          = "Alert when transactions are taking too long"
  runbook_url          = "https://www.example.com"
  enabled              = true
  value_function       = "single_value"
  violation_time_limit = "one_hour"

  nrql {
    query             = "SELECT average(duration) FROM Transaction where appName = 'Your App'"
    evaluation_offset = 3
  }

  term {
    operator              = "above"
    priority              = "critical"
    threshold             = 5.5
    threshold_duration    = 300
    threshold_occurrences = "ALL"
  }

  term {
    operator              = "above"
    priority              = "warning"
    threshold             = 3.5
    threshold_duration    = 600
    threshold_occurrences = "ALL"
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

- `account_id` - (Optional) The New Relic account ID of the account you wish to create the condition. Defaults to the account ID set in your environment variable `NEWRELIC_ACCOUNT_ID`.
- `baseline_direction` - (Optional) The baseline direction of a _baseline_ NRQL alert condition. Valid values are: `lower_only`, `upper_and_lower`, `upper_only` (case insensitive).
- `description` - (Optional) The description of the NRQL alert condition.
- `policy_id` - (Required) The ID of the policy where this condition should be used.
- `name` - (Required) The title of the condition.
- `type` - (Optional) The type of the condition. Valid values are `static`, `baseline`, or `outlier`. Defaults to `static`.
- `runbook_url` - (Optional) Runbook URL to display in notifications.
- `enabled` - (Optional) Whether to enable the alert condition. Valid values are `true` and `false`. Defaults to `true`.
- `nrql` - (Required) A NRQL query. See [NRQL](#nrql) below for details.
- `term` - (Required) A list of terms for this condition. See [Terms](#terms) below for details.
- `value_function` - (Optional) Possible values are `single_value`, `sum` (case insensitive). Defaults to `single_value`.
- `expected_groups` - (Optional) Number of expected groups when using `outlier` detection.
- `ignore_overlap` - (Optional) Whether to look for a convergence of groups when using `outlier` detection.
- `violation_time_limit` - (Optional) Sets a time limit, in hours, that will automatically force-close a long-lasting violation after the time limit you select. Possible values are `ONE_HOUR`, `TWO_HOURS`, `FOUR_HOURS`, `EIGHT_HOURS`, `TWELVE_HOURS`, `TWENTY_FOUR_HOURS` (case insensitive).
- `violation_time_limit_seconds` - (Optional) **DEPRECATED:** Use `violation_time_limit` instead. Sets a time limit, in seconds, that will automatically force-close a long-lasting violation after the time limit you select. Possible values are `3600`, `7200`, `14400`, `28800`, `43200`, and `86400`.

## NRQL

The `nrql` block supports the following arguments:

- `query` - (Required) The NRQL query to execute for the condition.
- `evaluation_offset` - (Optional) Represented in minutes and must be within 1-20 minutes (inclusive). NRQL queries are evaluated in one-minute time windows. The start time depends on this value. It's recommended to set this to 3 minutes. An offset of less than 3 minutes will trigger violations sooner, but you may see more false positives and negatives due to data latency. With `evaluation_offset` set to 3 minutes, the NRQL time window applied to your query will be: `SINCE 3 minutes ago UNTIL 2 minutes ago`.
- `since_value` - (Optional)  **DEPRECATED:** Use `evaluation_offset` instead. The value to be used in the `SINCE <X> minutes ago` clause for the NRQL query. Must be between 1-20 (inclusive).

## Terms

NRQL alert conditions support up to two terms. At least one `term` must have `priority` set to `critical` and the second optional `term` must have `priority` set to `warning`.

The `term` block the following arguments:

- `duration` - (Required) In minutes, must be in the range of `1` to `120`, inclusive.
- `operator` - (Optional) `above`, `below`, or `equal`. Defaults to `equal`. Note that when using a `type` of `outlier`, the only valid option here is `above`.
- `priority` - (Optional) `critical` or `warning`. Defaults to `critical`.
- `threshold` - (Required) The value which will trigger a violation. Must be `0` or greater.
- `threshold_duration` - (Optional) The duration of time, in seconds, that the threshold must violate for in order to create a violation. Value must be a multiple of 60.
<br>For _baseline_ NRQL alert conditions, the value must be within 120-3600 seconds (inclusive).
<br>For _static_ NRQL alert conditions, the value must be within 120-7200 seconds (inclusive).

- `threshold_occurrences` - (Optional) The criteria for how many data points must be in violation for the specified threshold duration. Valid values are: `all` or `at_least_once` (case insensitive).
- `duration` - (Optional) **DEPRECATED:** Use `threshold_duration` instead. The duration of time, in _minutes_, that the threshold must violate for in order to create a violation. Must be within 1-120 (inclusive).
- `time_function` - (Optional) **DEPRECATED:** Use `threshold_occurrences` instead. The criteria for how many data points must be in violation for the specified threshold duration. Valid values are: `all` or `any`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the NRQL alert condition. This is a composite ID with the format `<policy_id>:<condition_id>` - e.g. `538291:6789035`.

## Additional Examples


##### Type: `baseline`

Example baseline NRQL alert condition for alerting when transaction durations are above a specified threshold and dynamically adjusts based on data trends.

```hcl
resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_nrql_alert_condition" "foo" {
  type                 = "baseline"
  name                 = "foo"
  policy_id            = newrelic_alert_policy.foo.id
  description          = "Alert when transactions are taking too long"
  enabled              = true
  runbook_url          = "https://www.example.com"
  violation_time_limit = "one_hour"

  # baseline type only
  baseline_direction = "upper_only"

  nrql {
    query             = "SELECT percentile(duration, 95) FROM Transaction WHERE appName = 'ExampleAppName'"
    evaluation_offset = 3
  }

  term {
    operator              = "above"
    priority              = "critical"
    threshold             = 5.5
    threshold_duration    = 300
    threshold_occurrences = "all"
  }

  term {
    operator              = "above"
    priority              = "warning"
    threshold             = 3.5
    threshold_duration    = 600
    threshold_occurrences = "all"
  }
}
```

##### Type: `outlier`

Example outlier detection: Get notified if members of a group deviate by a key metric.

-> **NOTE:** The `outlier` NQRL alert condition type currently does not support new schema attributes introduced in v2.0.0 of the New Relic Terraform provider.

```hcl
resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_nrql_alert_condition" "foo" {
  type                         = "outlier"
  name                         = "foo"
  policy_id                    = newrelic_alert_policy.foo.id
  enabled                      = true
  runbook_url                  = "https://www.example.com"
  violation_time_limit_seconds = 7200

  # outlier type only
  expected_groups = 2

  # outlier type only
	ignore_overlap = true

  nrql {
    query       = "SELECT percentile(duration, 95) FROM Transaction WHERE appName = 'ExampleAppName' FACET host"
    since_value = "3"
  }

  term {
    operator      = "above"
    priority      = "critical"
    threshold     = 0.065
    duration      = 5
    time_function = "all"
  }

  term {
    operator      = "above"
    priority      = "warning"
    threshold     = 0.035
    duration      = 10
    time_function = "all"
  }
}
```

## Import


Alert conditions can be imported using a composite ID of `<policy_id>:<condition_id>:<conditionType>`, e.g.

```
// For `baseline` conditions
$ terraform import newrelic_nrql_alert_condition.foo 538291:6789035:baseline

// For `static` conditions
$ terraform import newrelic_nrql_alert_condition.foo 538291:6789035:static

// For `outlier` conditions
$ terraform import newrelic_nrql_alert_condition.foo 538291:6789035:outlier
```

~> **NOTE:** The value of `conditionType` in the import composite ID must be a valid condition type - `static`, `baseline`, or `outlier.`

The actual values for `policy_id` and `condition_id` can be retrieved from the following New Relic URL when viewing the NRQL alert condition you want to import:

<small>alerts.newrelic.com/accounts/**\<account_id\>**/policies/**\<policy_id\>**/conditions/**\<condition_id\>**/edit</small>
