---
layout: 'newrelic'
page_title: 'New Relic: newrelic_nrql_alert_condition'
sidebar_current: 'docs-newrelic-resource-nrql-alert-condition'
description: |-
  Create and manage a NRQL alert condition for a policy in New Relic.
---

# newrelic_nrql_alert_condition

Use this resource to create and manage NRQL alert conditions in New Relic.

## Example Usage

##### Type: `static` (default)
```hcl
resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_nrql_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name        = "foo"
  type        = "static"
  runbook_url = "https://www.example.com"
  enabled     = true

  term {
    duration      = 5
    operator      = "below"
    priority      = "critical"
    threshold     = "1"
    time_function = "all"
  }

  nrql {
    query       = "SELECT count(*) FROM SyntheticCheck WHERE monitorId = '<monitorId>'"
    since_value = "3"
  }

  value_function = "single_value"
}
```
See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

- `policy_id` - (Required) The ID of the policy where this condition should be used.
- `name` - (Required) The title of the condition
- `type` - (Optional) The type of the condition. Valid values are `static` or `outlier`. Defaults to `static`.
- `runbook_url` - (Optional) Runbook URL to display in notifications.
- `enabled` - (Optional) Whether to enable the alert condition. Valid values are `true` and `false`. Defaults to `true`.
- `term` - (Required) A list of terms for this condition. See [Terms](#terms) below for details.
- `nrql` - (Required) A NRQL query. See [NRQL](#nrql) below for details.
- `value_function` - (Optional) Possible values are `single_value`, `sum`.
- `expected_groups` - (Optional) Number of expected groups when using `outlier` detection.
- `ignore_overlap` - (Optional) Whether to look for a convergence of groups when using `outlier` detection.
- `violation_time_limit_seconds` - (Optional) Sets a time limit, in seconds, that will automatically force-close a long-lasting violation after the time limit you select.  Possible values are `3600`, `7200`, `14400`, `28800`, `43200`, and `86400`.

## Terms

The `term` mapping supports the following arguments:

- `duration` - (Required) In minutes, must be in the range of `1` to `120`, inclusive.
- `operator` - (Optional) `above`, `below`, or `equal`. Defaults to `equal`.
- `priority` - (Optional) `critical` or `warning`. Defaults to `critical`.
- `threshold` - (Required) Must be 0 or greater.
- `time_function` - (Required) `all` or `any`.

## NRQL

The `nrql` attribute supports the following arguments:

- `query` - (Required) The NRQL query to execute for the condition.
- `since_value` - (Required) The value to be used in the `SINCE <X> MINUTES AGO` clause for the NRQL query. Must be between `1` and `20`.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the NRQL alert condition.

## Additional Examples

##### Type: `outlier`
```hcl
resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_nrql_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name        = "outlier-example"
  runbook_url = "https://bar.example.com"
  enabled     = true

  term {
    duration      = 10
    operator      = "above"
    priority      = "critical"
    threshold     = "0.65"
    time_function = "all"
  }
  nrql {
    query       = "SELECT percentile(duration, 99) FROM Transaction FACET remote_ip"
    since_value = "3"
  }
  type            = "outlier"
  expected_groups = 2
  ignore_overlap  = true
}
```

## Import

Alert conditions can be imported using a composite ID of `<policy_id>:<condition_id>`, e.g.

```
$ terraform import newrelic_nrql_alert_condition.main 12345:67890
```

The actual values for `policy_id` and `condition_id` can be retrieved from the following URL when looking at the alert condition:

https://alerts.newrelic.com/accounts/<account_id>/policies/<policy_id>/conditions/<condition_id>/edit?selectedField=thresholds
