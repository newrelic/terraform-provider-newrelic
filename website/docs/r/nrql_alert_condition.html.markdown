---
layout: "newrelic"
page_title: "New Relic: newrelic_nrql_alert_condition"
sidebar_current: "docs-newrelic-resource-nrql-alert-condition"
description: |-
  Create and manage a NRQL alert condition for a policy in New Relic.
---

# newrelic\_nrql_alert\_condition

## Example Usage

```hcl
resource "newrelic_alert_policy" "foo" {
  name = "foo"
}

resource "newrelic_nrql_alert_condition" "foo" {
  policy_id = "${newrelic_alert_policy.foo.id}"

  name        = "foo"
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

## Argument Reference

The following arguments are supported:

  * `policy_id` - (Required) The ID of the policy where this condition should be used.
  * `name` - (Required) The title of the condition
  * `runbook_url` - (Optional) Runbook URL to display in notifications.
  * `enabled` - (Optional) Set whether to enable the alert condition. Defaults to `true`.
  * `term` - (Required) A list of terms for this condition. See [Terms](#terms) below for details.
  * `nrql` - (Required) A NRQL query. See [NRQL](#nrql) below for details.
  * `value_function` - (Optional) Possible values are `single_value`, `sum`.

## Terms

The `term` mapping supports the following arguments:

  * `duration` - (Required) Query duration in minutes.
  * `operator` - (Optional) `above`, `below`, or `equal`.  Defaults to `equal`.
  * `priority` - (Optional) `critical` or `warning`.  Defaults to `critical`.
  * `threshold` - (Required) Must be 0 or greater.
  * `time_function` - (Required) `all` or `any`.

## NRQL

The `nrql` attribute supports the following arguments:

  * `query` - (Required) The NRQL query to execute for the condition.
  * `since_value` - (Required) The value to be used in the `SINCE <X> MINUTES AGO` clause for the NRQL query. Must be between `1` and `20`.

## Attributes Reference

The following attributes are exported:

  * `id` - The ID of the NRQL alert condition.

## Import

Alert conditions can be imported using the `id`, e.g.

```
$ terraform import newrelic_nrql_alert_condition.main 12345
```
