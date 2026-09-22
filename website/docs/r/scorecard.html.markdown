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

-> **NOTE:** Rules are managed as separate [`newrelic_scorecard_rule`](scorecard_rule.html) resources. Use `rule_ids` to attach them. Each rule can only belong to **one scorecard at a time** — the API rejects attaching a rule that is already in another scorecard's collection.

-> **NOTE:** `progress_levels` are set at create time only. Adding, removing, or changing the *content* of a level requires resource recreation. **Reordering** the `progress_levels` blocks in your configuration does not require recreation — only content changes do. If omitted, the organization's default progress levels are applied.

-> **NOTE:** To assign a rule to a specific tier, set `progress_level` on the [`newrelic_scorecard_rule`](scorecard_rule.html) to match one of the `id` values in the `progress_levels` block below (e.g. `progress_level = "red"`). See [Progress Level Relationship](scorecard_rule.html#progress-level-relationship) for a full example.

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
  tags        = ["team:platform", "env:production"]

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

See additional [examples](#additional-examples).

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the scorecard. Must not be empty.
  * `description` - (Optional) A description of the scorecard's purpose. Can be cleared by setting to an empty string `""`.
  * `tags` - (Optional) A set of tags in `"key:value"` format. Order does not matter. Tags managed by New Relic (prefixed with `nr.`) are preserved automatically and must not be included here.
  * `progress_levels` - (Optional) One or more nested blocks defining the scorecard's scoring tiers. Adding, removing, or changing level values requires resource recreation; reordering existing blocks does not. See [Nested `progress_levels` blocks](#nested-progress_levels-blocks) below.
  * `rule_ids` - (Optional) A set of `newrelic_scorecard_rule` entity GUIDs to attach to this scorecard. Each rule can only belong to one scorecard — removing a GUID detaches the rule (without deleting it) so it can be re-attached elsewhere.
  * `organization_id` - (Optional, Computed) The NGEP organization UUID. Resolved automatically from the provider credentials if omitted.

### Nested `progress_levels` blocks

Each `progress_levels` block defines one scoring tier. The `id` values are also used by `newrelic_scorecard_rule.progress_level` to assign rules to tiers.

  * `id` - (Required) A machine identifier for this tier, referenced by rule's `progress_level` (e.g. `"red"`, `"amber"`, `"green"`). Changing this value requires resource recreation.
  * `name` - (Required) The display label shown in the New Relic UI (e.g. `"Needs Work"`). Changing this value requires resource recreation.
  * `description` - (Optional) A short description of what this tier means. Changing this value requires resource recreation.
  * `hex_color_code` - (Optional) Hex color code for the tier indicator badge, e.g. `"#FF0000"`. Must be 4–9 characters. Changing this value requires resource recreation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The entity GUID of the scorecard.
  * `rules_collection_id` - The GUID of the auto-created rules collection. This collection is managed by the provider via `rule_ids` and should not be modified directly.

## Additional Examples

### Scorecard with no custom progress levels

When `progress_levels` is omitted the organization's default levels are applied automatically.

```hcl
resource "newrelic_scorecard" "minimal" {
  name = "Service Health Check"

  rule_ids = [newrelic_scorecard_rule.alert_coverage.id]
}
```

### Attaching multiple rules

```hcl
resource "newrelic_scorecard" "platform" {
  name = "Platform Standards"
  tags = ["team:platform"]

  rule_ids = [
    newrelic_scorecard_rule.alert_coverage.id,
    newrelic_scorecard_rule.latency_threshold.id,
    newrelic_scorecard_rule.error_rate.id,
  ]
}
```

## Import

Scorecards can be imported using the entity GUID:

```bash
$ terraform import newrelic_scorecard.example <guid>
```
