---
layout: "newrelic"
page_title: "New Relic: newrelic_scorecard"
sidebar_current: "docs-newrelic-resource-scorecard"
description: |-
  Create and manage New Relic Scorecards.
---

# Resource: newrelic\_scorecard

Use this resource to create, update, and delete [New Relic Scorecards](https://docs.newrelic.com/docs/service-architecture-intelligence/scorecards/getting-started/).

Scorecards let you define and track engineering quality standards across your organization by grouping rules that evaluate NRQL-based checks against your entities.

-> **NOTE:** Rules are managed as separate [`newrelic_scorecard_rule`](scorecard_rule.html) resources and attached via the `rule_ids` attribute. **Each rule can only belong to one scorecard at a time** — attaching a rule that is already assigned to another scorecard produces an error. Remove it from its current scorecard first.

-> **NOTE:** To assign a rule to a specific progress tier, set `progress_level` on the [`newrelic_scorecard_rule`](scorecard_rule.html) to match one of the `id` values in this scorecard's `progress_levels` block (e.g. `progress_level = "red"`). See [Progress Level Relationship](scorecard_rule.html#progress-level-relationship) for a full example.

## Example Usage

```hcl
resource "newrelic_scorecard_rule" "alert_coverage" {
  name         = "APM alert coverage"
  enabled      = true
  run_interval = 1440

  nrql_engine {
    accounts = [var.account_id]
    query    = "SELECT if(latest(alertSeverity) != 'NOT_CONFIGURED', 1, 0) AS 'score' FROM Entity WHERE type = 'APM-APPLICATION' FACET id LIMIT MAX SINCE 1 day ago"
  }
}

resource "newrelic_scorecard" "engineering" {
  name        = "Engineering Quality"
  description = "Tracks key observability standards across all APM services"

  tags {
    key    = "team"
    values = ["platform"]
  }
  tags {
    key    = "env"
    values = ["production"]
  }

  progress_levels {
    id             = "red"
    name           = "Needs Work"
    description    = "Score below 60%"
    hex_color_code = "#FF4444"
  }
  progress_levels {
    id             = "amber"
    name           = "Improving"
    description    = "Score between 60–80%"
    hex_color_code = "#FFAA00"
  }
  progress_levels {
    id             = "green"
    name           = "Healthy"
    description    = "Score above 80%"
    hex_color_code = "#00CC44"
  }

  rule_ids = [
    newrelic_scorecard_rule.alert_coverage.id,
  ]
}
```

See additional [examples](#additional-examples) below.

## Argument Reference

* `name` - (Required) The name of the scorecard. Must not be empty.
* `description` - (Optional) A description of the scorecard's purpose.
* `tags` - (Optional) One or more `tags` blocks assigning key-value metadata to this scorecard. Tags prefixed with `nr.` are managed by New Relic and are preserved automatically during updates. Each block supports:
  * `key` - (Required) The tag key.
  * `values` - (Required) One or more tag values.
* `progress_levels` - (Optional, Computed) One or more blocks defining the scorecard's scoring tiers (e.g. Red / Amber / Green). Progress levels can be added, updated, or removed in-place without recreating the scorecard. If omitted, the organization's default levels are applied and stored in state. See [Nested `progress_levels` blocks](#nested-progress_levels-blocks) below.
* `rule_ids` - (Optional) A set of `newrelic_scorecard_rule` entity GUIDs to attach to this scorecard. Removing a GUID detaches the rule without deleting it, making it available to attach to another scorecard.

### Nested `progress_levels` blocks

Each `progress_levels` block defines one scoring tier. The `id` value is referenced by `newrelic_scorecard_rule.progress_level` to visually group rules under that tier in the UI.

* `id` - (Required) A short identifier for this tier, used by rule's `progress_level` field (e.g. `"red"`, `"amber"`, `"green"`).
* `name` - (Required) The display label shown in the New Relic UI (e.g. `"Needs Work"`).
* `description` - (Optional) A short description of what this tier represents.
* `hex_color_code` - (Optional) Hex color code for the tier indicator badge, e.g. `"#FF4444"`. Must be 4–9 characters.

## Attributes Reference

In addition to all arguments above, the following computed attributes are exported:

* `id` - The entity GUID of the scorecard.
* `organization_id` - The organization UUID, resolved automatically from the provider account.
* `rules_collection_id` - The GUID of the auto-created rules collection. Managed by the provider via `rule_ids`; do not modify directly.

## Additional Examples

### Scorecard with no custom progress levels

When `progress_levels` is omitted, the organization's default levels are applied and stored in state automatically.

```hcl
resource "newrelic_scorecard" "minimal" {
  name     = "Service Health"
  rule_ids = [newrelic_scorecard_rule.alert_coverage.id]
}
```

### Attaching multiple rules

```hcl
resource "newrelic_scorecard" "platform_standards" {
  name = "Platform Standards"

  tags {
    key    = "team"
    values = ["platform"]
  }

  rule_ids = [
    newrelic_scorecard_rule.alert_coverage.id,
    newrelic_scorecard_rule.latency_threshold.id,
    newrelic_scorecard_rule.error_rate.id,
  ]
}
```

### Updating progress levels in place

Progress levels can be changed without recreating the scorecard.

```hcl
resource "newrelic_scorecard" "engineering" {
  name = "Engineering Standards"

  progress_levels {
    id             = "bronze"
    name           = "Bronze"
    description    = "Baseline requirements met"
    hex_color_code = "#CD7F32"
  }
  progress_levels {
    id             = "silver"
    name           = "Silver"
    description    = "Strong operational posture"
    hex_color_code = "#C0C0C0"
  }
  progress_levels {
    id             = "gold"
    name           = "Gold"
    description    = "Exemplary reliability and observability"
    hex_color_code = "#FFD700"
  }
}
```

## Import

Scorecards can be imported using the entity GUID:

```bash
$ terraform import newrelic_scorecard.example <guid>
```
