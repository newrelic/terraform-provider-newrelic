---
layout: "newrelic"
page_title: "New Relic: newrelic_scorecard_rule"
sidebar_current: "docs-newrelic-resource-scorecard-rule"
description: |-
  Create and manage New Relic Scorecard Rules.
---

# Resource: newrelic\_scorecard\_rule

Use this resource to create, update, and delete [New Relic Scorecard Rules](https://docs.newrelic.com/docs/service-architecture-intelligence/scorecards/getting-started/).

A Scorecard Rule is a NRQL-based check that evaluates a binary score (0 or 1) for each entity in your account. Rules are standalone resources that can be attached to one or more [`newrelic_scorecard`](scorecard.html) resources via `rule_ids`.

-> **NOTE:** The deprecated `schedule` field is not supported. Use `run_interval` (minutes) together with `enabled` to control when and whether a rule runs.

## Example Usage

```hcl
resource "newrelic_scorecard_rule" "alert_coverage" {
  name         = "APM alert coverage"
  description  = "All APM services must have alerts configured"
  enabled      = true
  run_interval = 1440

  nrql_engine {
    accounts = [var.account_id]
    query    = "SELECT if(latest(alertSeverity) != 'NOT_CONFIGURED', 1, 0) AS 'score' FROM Entity WHERE type = 'APM-APPLICATION' FACET id LIMIT MAX SINCE 1 day ago"
  }

  tags = ["team:platform", "purpose:observability-standards"]
}
```

See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the rule.
  * `description` - (Optional) A description of what this rule measures. Can be cleared by setting to an empty string `""`.
  * `enabled` - (Required) Whether the rule is active and collecting scores. Set to `false` to pause the rule without deleting it.
  * `run_interval` - (Required) How frequently (in minutes) the rule's NRQL query is executed. Common values: `60` (hourly), `720` (12h), `1440` (daily).
  * `nrql_engine` - (Required) A nested block defining the NRQL query that produces the score. See [Nested `nrql_engine` block](#nested-nrql_engine-block) below.
  * `tags` - (Optional) A list of tags in `"key:value"` format. Tags managed by New Relic (prefixed with `nr.`) are preserved automatically.
  * `organization_id` - (Optional, Computed) The NGEP organization UUID. Resolved automatically from the provider credentials if omitted.

### Nested `nrql_engine` block

The `nrql_engine` block supports the following arguments:

  * `accounts` - (Required) A list of New Relic account IDs the query runs against.
  * `query` - (Required) A NRQL query that must return a numeric `score` column with value `0` (fail) or `1` (pass) for each entity, with `FACET id` or `FACET entity.guid` to attribute results to individual entities. Example: `SELECT if(count(*) > 0, 1, 0) AS 'score' FROM Transaction FACET entity.guid LIMIT MAX SINCE 1 hour ago`.
  * `join_accounts` - (Optional) Additional account IDs whose data is joined into the query for cross-account checks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The entity GUID of the scorecard rule.

## Additional Examples

### Disabled rule (paused, not deleted)

```hcl
resource "newrelic_scorecard_rule" "latency_threshold" {
  name         = "P95 latency < 500ms"
  enabled      = false
  run_interval = 720

  nrql_engine {
    accounts = [var.account_id]
    query    = "SELECT if(percentile(duration, 95) < 0.5, 1, 0) AS 'score' FROM Transaction FACET entity.guid LIMIT MAX SINCE 1 hour ago"
  }
}
```

### Cross-account rule using `join_accounts`

```hcl
resource "newrelic_scorecard_rule" "cross_account_check" {
  name         = "Cross-account error rate"
  enabled      = true
  run_interval = 60

  nrql_engine {
    accounts      = [var.primary_account_id]
    join_accounts = [var.secondary_account_id]
    query         = "SELECT if(percentage(count(*), WHERE error IS true) < 1.0, 1, 0) AS 'score' FROM Transaction FACET entity.guid LIMIT MAX SINCE 1 hour ago"
  }
}
```

## Import

Scorecard rules can be imported using the entity GUID:

```bash
$ terraform import newrelic_scorecard_rule.example <guid>
```
