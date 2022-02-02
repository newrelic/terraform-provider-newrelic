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
Before upgrading to version 2.0.0 or later, it is recommended to upgrade to the most recent 1.x version of the provider and ensure that your environment successfully runs `terraform plan` without unexpected changes.

## Example Usage

##### Type: `static` (default)
```hcl
resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_nrql_alert_condition" "foo" {
  account_id                     = <Your Account ID>
  policy_id                      = newrelic_alert_policy.foo.id
  type                           = "static"
  name                           = "foo"
  description                    = "Alert when transactions are taking too long"
  runbook_url                    = "https://www.example.com"
  enabled                        = true
  violation_time_limit_seconds   = 3600
  value_function                 = "single_value"
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
    query = "SELECT average(duration) FROM Transaction where appName = 'Your App'"
  }

  critical {
    operator              = "above"
    threshold             = 5.5
    threshold_duration    = 300
    threshold_occurrences = "ALL"
  }

  warning {
    operator              = "above"
    threshold             = 3.5
    threshold_duration    = 600
    threshold_occurrences = "ALL"
  }
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

- `account_id` - (Optional) The New Relic account ID of the account you wish to create the condition. Defaults to the account ID set in your environment variable `NEW_RELIC_ACCOUNT_ID`.
- `baseline_direction` - (Optional) The baseline direction of a _baseline_ NRQL alert condition. Valid values are: `lower_only`, `upper_and_lower`, `upper_only` (case insensitive).
- `description` - (Optional) The description of the NRQL alert condition.
- `policy_id` - (Required) The ID of the policy where this condition should be used.
- `name` - (Required) The title of the condition.
- `type` - (Optional) The type of the condition. Valid values are `static`, `baseline`, or `outlier` (deprecated). Defaults to `static`.
- `runbook_url` - (Optional) Runbook URL to display in notifications.
- `enabled` - (Optional) Whether to enable the alert condition. Valid values are `true` and `false`. Defaults to `true`.
- `nrql` - (Required) A NRQL query. See [NRQL](#nrql) below for details.
- `term` - (Optional) **DEPRECATED** Use `critical`, and `warning` instead.  A list of terms for this condition. See [Terms](#terms) below for details.
- `critical` - (Required) A list containing the `critical` threshold values. See [Terms](#terms) below for details.
- `warning` - (Optional) A list containing the `warning` threshold values. See [Terms](#terms) below for details.
- `value_function` - (Required if `type` is `static`, omit when `type` is `baseline` or `outlier` ) Possible values are `single_value`, `sum` (case insensitive).
- `expected_groups` - (Optional) Number of expected groups when using `outlier` detection.
- `open_violation_on_group_overlap` - (Optional) Whether or not to trigger a violation when groups overlap. Set to `true` if you want to trigger a violation when groups overlap. This argument is only applicable in `outlier` conditions.
- `ignore_overlap` - (Optional) **DEPRECATED:** Use `open_violation_on_group_overlap` instead, but use the inverse value of your boolean - e.g. if `ignore_overlap = false`, use `open_violation_on_group_overlap = true`. This argument sets whether to trigger a violation when groups overlap. If set to `true` overlapping groups will not trigger a violation. This argument is only applicable in `outlier` conditions.
- `violation_time_limit` - (Optional) **DEPRECATED:** Use `violation_time_limit_seconds` instead. Sets a time limit, in hours, that will automatically force-close a long-lasting violation after the time limit you select. Possible values are `ONE_HOUR`, `TWO_HOURS`, `FOUR_HOURS`, `EIGHT_HOURS`, `TWELVE_HOURS`, `TWENTY_FOUR_HOURS`, `THIRTY_DAYS` (case insensitive).<br>
<small>\***Note**: One of `violation_time_limit` _or_ `violation_time_limit_seconds` must be set, but not both.</small>

- `violation_time_limit_seconds` - (Optional) Sets a time limit, in seconds, that will automatically force-close a long-lasting violation after the time limit you select. The value must be between 300 seconds (5 minutes) to 2592000 seconds (30 days) (inclusive). <br>
<small>\***Note**: One of `violation_time_limit` _or_ `violation_time_limit_seconds` must be set, but not both.</small>

- `fill_option` - (Optional) Which strategy to use when filling gaps in the signal. Possible values are `none`, `last_value` or `static`. If `static`, the `fill_value` field will be used for filling gaps in the signal.
- `fill_value` - (Optional, required when `fill_option` is `static`) This value will be used for filling gaps in the signal.
- `aggregation_window` - (Optional) The duration of the time window used to evaluate the NRQL query, in seconds. The value must be at least 30 seconds, and no more than 15 minutes (900 seconds). Default is 60 seconds.
- `expiration_duration` - (Optional) The amount of time (in seconds) to wait before considering the signal expired.
- `open_violation_on_expiration` - (Optional) Whether to create a new violation to capture that the signal expired.
- `close_violations_on_expiration` - (Optional) Whether to close all open violations when the signal expires.
- `aggregation_method` - (Optional) Determines when we consider an aggregation window to be complete so that we can evaluate the signal for violations. Possible values are `cadence`, `event_flow` or `event_timer`. Default is `event_flow`. `aggregation_method` cannot be set with `nrql.evaluation_offset`.
- `aggregation_delay` - (Optional) How long we wait for data that belongs in each aggregation window. Depending on your data, a longer delay may increase accuracy but delay notifications. Use `aggregation_delay` with the `event_flow` and `cadence` methods. The maximum delay is 1200 seconds (20 minutes) when using `event_flow` and 3600 seconds (60 minutes) when using `cadence`. In both cases, the minimum delay is 0 seconds and the default is 120 seconds. `aggregation_delay` cannot be set with `nrql.evaluation_offset`.
- `aggregation_timer` - (Optional) How long we wait after each data point arrives to make sure we've processed the whole batch. Use `aggregation_timer` with the `event_timer` method. The timer value can range from 0 seconds to 1200 seconds (20 minutes); the default is 60 seconds. `aggregation_timer` cannot be set with `nrql.evaluation_offset`.
- `slide_by` - (Optional) Gathers data in overlapping time windows to smooth the chart line, making it easier to spot trends. The `slide_by` value is specified in seconds and must be smaller than and a factor of the `aggregation_window`. `slide_by` cannot be used with `outlier` NRQL conditions or `static` NRQL conditions using the `sum` `value_function`.

## NRQL

The `nrql` block supports the following arguments:

- `query` - (Required) The NRQL query to execute for the condition.
- `evaluation_offset` - (Optional) **DEPRECATED:** Use `aggregation_method` instead. Represented in minutes and must be within 1-20 minutes (inclusive). NRQL queries are evaluated based on their `aggregation_window` size. The start time depends on this value. It's recommended to set this to 3 windows. An offset of less than 3 windows will trigger violations sooner, but you may see more false positives and negatives due to data latency. With `evaluation_offset` set to 3 windows and an `aggregation_window` of 60 seconds, the NRQL time window applied to your query will be: `SINCE 3 minutes ago UNTIL 2 minutes ago`. `evaluation_offset` cannot be set with `aggregation_method`, `aggregation_delay`, or `aggregation_timer`.<br>
- `since_value` - (Optional)  **DEPRECATED:** Use `aggregation_method` instead. The value to be used in the `SINCE <X> minutes ago` clause for the NRQL query. Must be between 1-20 (inclusive). <br>

## Terms

~> **NOTE:** The direct use of the `term` has been deprecated, and users should use `critical` and `warning` instead.  What follows now applies to the named priority attributes for `critical` and `warning`, but for those attributes the priority is not allowed.

NRQL alert conditions support up to two terms. At least one `term` must have `priority` set to `critical` and the second optional `term` must have `priority` set to `warning`.

The `term` block supports the following arguments:

- `operator` - (Optional) Valid values are `above`, `below`, or `equals` (case insensitive). Defaults to `equals`. Note that when using a `type` of `outlier` or `baseline`, the only valid option here is `above`.
- `priority` - (Optional) `critical` or `warning`. Defaults to `critical`.
- `threshold` - (Required) The value which will trigger a violation. Must be `0` or greater.
<br>For _baseline_ NRQL alert conditions, the value must be in the range [1, 1000]. The value is the number of standard deviations from the baseline that the metric must exceed in order to create a violation.
- `threshold_duration` - (Optional) The duration, in seconds, that the threshold must violate in order to create a violation. Value must be a multiple of the `aggregation_window` (which has a default of 60 seconds).
<br>For _baseline_ and _outlier_ NRQL alert conditions, the value must be within 120-3600 seconds (inclusive).
<br>For _static_ NRQL alert conditions with the `sum` value function, the value must be within 120-7200 seconds (inclusive).
<br>For _static_ NRQL alert conditions with the `single_value` value function, the value must be within 60-7200 seconds (inclusive).

- `threshold_occurrences` - (Optional) The criteria for how many data points must be in violation for the specified threshold duration. Valid values are: `all` or `at_least_once` (case insensitive).
- `duration` - (Optional) **DEPRECATED:** Use `threshold_duration` instead. The duration of time, in _minutes_, that the threshold must violate for in order to create a violation. Must be within 1-120 (inclusive).
- `time_function` - (Optional) **DEPRECATED:** Use `threshold_occurrences` instead. The criteria for how many data points must be in violation for the specified threshold duration. Valid values are: `all` or `any`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the NRQL alert condition. This is a composite ID with the format `<policy_id>:<condition_id>` - e.g. `538291:6789035`.

## Additional Examples


##### Type: `baseline`

[Baseline NRQL alert conditions](https://docs.newrelic.com/docs/alerts/new-relic-alerts/defining-conditions/create-baseline-alert-conditions) are dynamic in nature and adjust to the behavior of your data. The example below demonstrates a baseline NRQL alert condition for alerting when transaction durations are above a specified threshold and dynamically adjusts based on data trends.

```hcl
resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_nrql_alert_condition" "foo" {
  type                         = "baseline"
  account_id                   = <Your Account ID>
  name                         = "foo"
  policy_id                    = newrelic_alert_policy.foo.id
  description                  = "Alert when transactions are taking too long"
  enabled                      = true
  runbook_url                  = "https://www.example.com"
  violation_time_limit_seconds = 3600
  aggregation_method           = "event_flow"
  aggregation_delay            = 120
  slide_by                     = 30

  # baseline type only
  baseline_direction = "upper_only"

  nrql {
    query = "SELECT percentile(duration, 95) FROM Transaction WHERE appName = 'ExampleAppName'"
  }

  critical {
    operator              = "above"
    threshold             = 5.5
    threshold_duration    = 300
    threshold_occurrences = "all"
  }

  warning {
    operator              = "above"
    threshold             = 3.5
    threshold_duration    = 600
    threshold_occurrences = "all"
  }
}
```

<br>

##### Type: `outlier`

In software development and operations, it is common to have a group consisting of members you expect to behave approximately the same. [Outlier detection](https://docs.newrelic.com/docs/alerts/new-relic-alerts/defining-conditions/outlier-detection-nrql-alert) facilitates alerting when the behavior of one or more common members falls outside a specified range expectation.

```hcl
resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_nrql_alert_condition" "foo" {
  type                         = "outlier"
  account_id                   = <Your Account ID>
  name                         = "foo"
  policy_id                    = newrelic_alert_policy.foo.id
  description                  = "Alert when outlier conditions occur"
  enabled                      = true
  runbook_url                  = "https://www.example.com"
  violation_time_limit_seconds = 3600
  aggregation_method           = "event_flow"
  aggregation_delay            = 120

  # Outlier only
  expected_groups = 2

  # Outlier only
	open_violation_on_group_overlap = true

  nrql {
    query = "SELECT percentile(duration, 95) FROM Transaction WHERE appName = 'ExampleAppName' FACET host"
  }

  critical {
    operator              = "above"
    threshold             = 0.002
    threshold_duration    = 600
    threshold_occurrences = "all"
  }

  warning {
    operator              = "above"
    threshold             = 0.0015
    threshold_duration    = 600
    threshold_occurrences = "all"
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

~> **NOTE:** The value of `conditionType` in the import composite ID must be a valid condition type - `static`, `baseline`, or `outlier.` Also note that deprecated arguments will *not* be set when importing.

The actual values for `policy_id` and `condition_id` can be retrieved from the following New Relic URL when viewing the NRQL alert condition you want to import:

<small>alerts.newrelic.com/accounts/**\<account_id\>**/policies/**\<policy_id\>**/conditions/**\<condition_id\>**/edit</small>

## Upgrade from 1.x to 2.x

There have been several deprecations in the `newrelic_nrql_alert_condition`
resource.  Users will need to make some updates in order to have a smooth
upgrade.

An example resource from 1.x might look like the following.

```hcl
resource "newrelic_nrql_alert_condition" "z" {
  policy_id = newrelic_alert_policy.z.id

  name                 = "zleslie-test"
  type                 = "static"
  runbook_url          = "https://localhost"
  enabled              = true
  value_function       = "sum"
  violation_time_limit = "TWENTY_FOUR_HOURS"

  critical {
    operator              = "above"
    threshold_duration    = 120
    threshold             = 3
    threshold_occurrences = "AT_LEAST_ONCE"
  }

  nrql {
    query = "SELECT count(*) FROM TransactionError WHERE appName like '%Dummy App%' FACET appName"
  }
}
```

After making the appropriate adjustments mentioned in the deprecation warnings,
the resource now looks like the following.

```hcl
resource "newrelic_nrql_alert_condition" "z" {
  policy_id = newrelic_alert_policy.z.id

  name                         = "zleslie-test"
  type                         = "static"
  runbook_url                  = "https://localhost"
  enabled                      = true
  value_function               = "sum"
  violation_time_limit_seconds = 86400

  term {
    priority      = "critical"
    operator      = "above"
    threshold     = 3
    duration      = 5
    time_function = "any"
  }

  nrql {
    query = "SELECT count(*) FROM TransactionError WHERE appName like '%Dummy App%' FACET appName"
  }
}
```
