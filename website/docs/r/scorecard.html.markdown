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

-> **NOTE:** Rules attached to a scorecard are managed as separate [`newrelic_scorecard_rule`](scorecard_rule.html) resources. Use the `rule_ids` attribute to attach existing rules to this scorecard.

-> **NOTE:** `progress_levels` are set at create time only. Changes to progress levels require destroying and re-creating the scorecard. If omitted, the organization's default progress levels are applied.

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

  * `name` - (Required) The name of the scorecard.
  * `description` - (Optional) A description of the scorecard's purpose. Can be cleared by setting to an empty string `""`.
  * `tags` - (Optional) A list of tags in `"key:value"` format to assign to the scorecard. Tags managed by New Relic (prefixed with `nr.`) are preserved automatically and must not be included here.
  * `progress_levels` - (Optional) One or more nested blocks defining the scorecard's scoring tiers. Changes require resource recreation. See [Nested `progress_levels` blocks](#nested-progress_levels-blocks) below for details.
  * `rule_ids` - (Optional) A set of `newrelic_scorecard_rule` entity GUIDs to attach to this scorecard. Removing a GUID detaches the rule without deleting it — rules are standalone resources that can be shared across scorecards.
  * `organization_id` - (Optional, Computed) The NGEP organization UUID. Resolved automatically from the provider credentials if omitted.

### Nested `progress_levels` blocks

Each `progress_levels` block supports the following arguments. All fields force resource recreation when changed.

  * `id` - (Required) A unique identifier for this level within the scorecard (e.g. `"red"`, `"amber"`, `"green"`).
  * `name` - (Required) The display name shown in the New Relic UI (e.g. `"Needs Work"`).
  * `description` - (Optional) A short description of what this level means.
  * `hex_color_code` - (Optional) Hex color code for the level indicator, e.g. `"#FF0000"`. Must be 4–9 characters.

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
