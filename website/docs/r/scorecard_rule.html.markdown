---
layout: "newrelic"
page_title: "New Relic: newrelic_scorecard_rule"
sidebar_current: "docs-newrelic-resource-scorecard-rule"
description: |-
  Create and manage New Relic Scorecard Rules.
---

# Resource: newrelic\_scorecard\_rule

Use this resource to create, update, and delete [New Relic Scorecard Rules](https://docs.newrelic.com/docs/service-architecture-intelligence/scorecards/getting-started/).

A Scorecard Rule is a NRQL-based check that evaluates a binary score (0 = fail, 1 = pass) for each entity in your account on a recurring schedule. Rules are **independent entities** — they exist outside any scorecard and are attached to a [`newrelic_scorecard`](scorecard.html) via that scorecard's `rule_ids` attribute.

-> **CONSTRAINT: One rule per scorecard.** Each rule can only belong to **one scorecard at a time**. The New Relic API enforces this at the collection level. If you attempt to attach a rule that is already assigned to another scorecard, the apply will fail. To reassign a rule: remove its GUID from the current scorecard's `rule_ids`, apply that change, then add it to the new scorecard.

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

  tags {
    key    = "team"
    values = ["platform"]
  }
}
```

See additional [examples](#additional-examples) below.

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the rule.
  * `description` - (Optional) A description of what this rule measures.
  * `enabled` - (Required) Whether the rule is active. Set to `false` to pause evaluation without deleting the rule.
  * `run_interval` - (Optional) How frequently the NRQL query runs, in minutes. Accepted values: `60` (hourly), `360` (6 hours), `720` (12 hours), `1440` (daily). Omit to use the API default.
  * `nrql_engine` - (Required) A block defining the NRQL query that produces the score. See [Nested `nrql_engine` block](#nested-nrql_engine-block) below.
  * `impact_weight` - (Optional) A non-negative integer weighting this rule's contribution to the overall scorecard score relative to other rules. Omit for equal weighting.
  * `progress_level` - (Optional) The `id` of a progress tier defined in the parent [`newrelic_scorecard`](scorecard.html). Assigns the rule to a visual grouping in the UI. The value must exactly match a `progress_levels.id` on the scorecard (e.g. `"red"`, `"amber"`, `"green"`). Omit to leave the rule ungrouped. See [Progress Level Relationship](#progress-level-relationship) below.
  * `tags` - (Optional) One or more `tags` blocks assigning key-value metadata to this rule. Tags prefixed with `nr.` are managed by New Relic and are preserved automatically. Each block supports:
    * `key` - (Required) The tag key.
    * `values` - (Required) One or more tag values.

### Nested `nrql_engine` block

  * `accounts` - (Required) A list of New Relic account IDs the query runs against.
  * `query` - (Required) A NRQL query that returns a numeric `score` column with value `0` (fail) or `1` (pass) per entity, using `FACET id` or `FACET entity.guid` to attribute results. Example: `SELECT if(count(*) > 0, 1, 0) AS 'score' FROM Transaction FACET entity.guid LIMIT MAX SINCE 1 hour ago`.
  * `join_accounts` - (Optional) Additional account IDs whose data is joined into the query for cross-account checks. Must not overlap with `accounts`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The entity GUID of the scorecard rule.
  * `organization_id` - The organisation UUID, resolved automatically from the provider account.

## Progress Level Relationship

The `progress_level` field links a rule to a tier in the parent scorecard's `progress_levels` block. The value must exactly match one of those tier IDs.

```hcl
resource "newrelic_scorecard" "engineering" {
  name = "Engineering Standards"

  progress_levels {
    id             = "foundation"
    name           = "Foundation"
    hex_color_code = "#CC4400"
  }
  progress_levels {
    id             = "excellence"
    name           = "Excellence"
    hex_color_code = "#00CC44"
  }

  rule_ids = [
    newrelic_scorecard_rule.alert_configured.id,
    newrelic_scorecard_rule.low_latency.id,
  ]
}

resource "newrelic_scorecard_rule" "alert_configured" {
  name           = "Alerts configured"
  enabled        = true
  run_interval   = 1440
  progress_level = "foundation"  # groups under the Foundation tier

  nrql_engine {
    accounts = [var.account_id]
    query    = "SELECT if(latest(alertSeverity) != 'NOT_CONFIGURED', 1, 0) AS 'score' FROM Entity WHERE type = 'APM-APPLICATION' FACET id LIMIT MAX SINCE 1 day ago"
  }
}

resource "newrelic_scorecard_rule" "low_latency" {
  name           = "P95 latency below 500ms"
  enabled        = true
  run_interval   = 60
  progress_level = "excellence"  # groups under the Excellence tier

  nrql_engine {
    accounts = [var.account_id]
    query    = "SELECT if(percentile(duration, 95) < 0.5, 1, 0) AS 'score' FROM Transaction FACET entity.guid LIMIT MAX SINCE 1 hour ago"
  }
}
```

Rules are visually grouped under their tier label in the scorecard detail view. Tiers are not hard gates — an entity can pass an `"excellence"` rule without first passing all `"foundation"` rules. Because each rule belongs to only one scorecard, the tier assignment is always unambiguous.

## Additional Examples

### Disabled rule (paused without deletion)

```hcl
resource "newrelic_scorecard_rule" "latency_threshold" {
  name         = "P95 latency below 500ms"
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
resource "newrelic_scorecard_rule" "cross_account_errors" {
  name         = "Cross-account error rate below 1%"
  enabled      = true
  run_interval = 60

  nrql_engine {
    accounts      = [var.primary_account_id]
    join_accounts = [var.secondary_account_id]
    query         = "SELECT if(percentage(count(*), WHERE error IS true) < 1.0, 1, 0) AS 'score' FROM Transaction FACET entity.guid LIMIT MAX SINCE 1 hour ago"
  }
}
```

### Rule with impact weighting

```hcl
resource "newrelic_scorecard_rule" "critical_alert_check" {
  name          = "Critical alerts configured"
  enabled       = true
  run_interval  = 1440
  impact_weight = 3  # counts 3× relative to rules without a weight

  nrql_engine {
    accounts = [var.account_id]
    query    = "SELECT if(latest(alertSeverity) = 'CRITICAL', 1, 0) AS 'score' FROM Entity WHERE type = 'APM-APPLICATION' FACET id LIMIT MAX SINCE 1 day ago"
  }
}
```

## Import

Scorecard rules can be imported using the entity GUID:

```bash
$ terraform import newrelic_scorecard_rule.example <guid>
```
