---
layout: "newrelic"
page_title: "New Relic: newrelic_alert_compound_condition"
sidebar_current: "docs-newrelic-datasource-alert-compound-condition"
description: |-
  Looks up a compound alert condition in New Relic.
---

# Data Source: newrelic_alert_compound_condition

Use this data source to retrieve information about a specific compound alert condition in New Relic.

## Example Usage

```hcl
data "newrelic_alert_compound_condition" "example" {
  id = "123456"
}

# Use the retrieved condition data
output "condition_name" {
  value = data.newrelic_alert_compound_condition.example.name
}

output "trigger_expression" {
  value = data.newrelic_alert_compound_condition.example.trigger_expression
}
```

## Argument Reference

The following arguments are supported:

- `id` - (Required) The ID of the compound alert condition to retrieve.
- `account_id` - (Optional) The New Relic account ID to operate on. Defaults to the account ID set in your environment variable `NEW_RELIC_ACCOUNT_ID`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `policy_id` - The ID of the policy associated with this compound alert condition.
- `name` - The name of the compound alert condition.
- `enabled` - Whether or not the compound alert condition is enabled.
- `trigger_expression` - Expression that defines how component condition evaluations are combined. Valid operators are 'AND', 'OR', 'NOT'.
- `component_conditions` - The list of NRQL conditions that are combined in this compound alert condition. Each component condition has:
  - `id` - The ID of the component NRQL condition.
  - `alias` - The identifier used in the trigger_expression.
- `facet_matching_behavior` - How the compound condition takes into account the component conditions' facets during evaluation.
- `runbook_url` - Runbook URL for the condition.
- `threshold_duration` - The duration, in seconds, that the trigger expression must be true before the compound alert condition will activate.
- `entity_guid` - The unique entity identifier of the alert compound condition.
