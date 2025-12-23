---
layout: "newrelic"
page_title: "New Relic: newrelic_alert_compound_condition (BETA PREVIEW)"
sidebar_current: "docs-newrelic-resource-compound-alert-condition"
description: |-
  Create and manage compound alert conditions in New Relic. This feature is in Beta Preview.
---

# Resource: newrelic_compound_alert_condition (BETA PREVIEW)

Use this resource to create and manage compound alert conditions in New Relic. Compound conditions allow you to combine multiple alert conditions using logical expressions (AND, OR) to create more sophisticated alerting logic.

## Example Usage

### Basic Compound Condition (AND)

```hcl
resource "newrelic_alert_policy" "example" {
  name = "my-policy"
}

# Create component NRQL conditions
resource "newrelic_nrql_alert_condition" "high_response_time" {
  policy_id = newrelic_alert_policy.example.id
  name      = "High Response Time"
  enabled   = true

  nrql {
    query = "SELECT average(duration) FROM Transaction WHERE appName = 'MyApp'"
  }

  critical {
    operator              = "above"
    threshold             = 5.0
    threshold_duration    = 300
    threshold_occurrences = "all"
  }

  violation_time_limit_seconds = 3600
}

resource "newrelic_nrql_alert_condition" "high_error_rate" {
  policy_id = newrelic_alert_policy.example.id
  name      = "High Error Rate"
  enabled   = true

  nrql {
    query = "SELECT percentage(count(*), WHERE error IS true) FROM Transaction WHERE appName = 'MyApp'"
  }

  critical {
    operator              = "above"
    threshold             = 5.0
    threshold_duration    = 300
    threshold_occurrences = "all"
  }

  violation_time_limit_seconds = 3600
}

# Create alert compound condition combining both
resource "newrelic_alert_compound_condition" "critical_service_health" {
  policy_id          = newrelic_alert_policy.example.id
  name               = "Critical Service Health"
  enabled            = true
  trigger_expression = "A AND B"
  runbook_url        = "https://example.com/runbooks/critical-health"
  threshold_duration = 120

  component_conditions {
    id    = newrelic_nrql_alert_condition.high_response_time.id
    alias = "A"
  }

  component_conditions {
    id    = newrelic_nrql_alert_condition.high_error_rate.id
    alias = "B"
  }

  facet_matching_behavior = "FACETS_IGNORED"
}
```

### Complex Condition with Three Components

```hcl
resource "newrelic_nrql_alert_condition" "high_cpu" {
  policy_id = newrelic_alert_policy.example.id
  name      = "High CPU"
  enabled   = true

  nrql {
    query = "SELECT average(cpuPercent) FROM SystemSample WHERE hostname = 'myhost'"
  }

  critical {
    operator              = "above"
    threshold             = 80.0
    threshold_duration    = 300
    threshold_occurrences = "all"
  }

  violation_time_limit_seconds = 3600
}

resource "newrelic_nrql_alert_condition" "high_memory" {
  policy_id = newrelic_alert_policy.example.id
  name      = "High Memory"
  enabled   = true

  nrql {
    query = "SELECT average(memoryUsedPercent) FROM SystemSample WHERE hostname = 'myhost'"
  }

  critical {
    operator              = "above"
    threshold             = 85.0
    threshold_duration    = 300
    threshold_occurrences = "all"
  }

  violation_time_limit_seconds = 3600
}

resource "newrelic_nrql_alert_condition" "disk_full" {
  policy_id = newrelic_alert_policy.example.id
  name      = "Disk Full"
  enabled   = true

  nrql {
    query = "SELECT average(diskUsedPercent) FROM SystemSample WHERE hostname = 'myhost'"
  }

  critical {
    operator              = "above"
    threshold             = 90.0
    threshold_duration    = 300
    threshold_occurrences = "all"
  }

  violation_time_limit_seconds = 3600
}

resource "newrelic_alert_compound_condition" "complex" {
  policy_id          = newrelic_alert_policy.example.id
  name               = "Complex Infrastructure Alert"
  enabled            = true
  trigger_expression = "(A AND B) OR C"

  component_conditions {
    id    = newrelic_nrql_alert_condition.high_cpu.id
    alias = "A"
  }

  component_conditions {
    id    = newrelic_nrql_alert_condition.high_memory.id
    alias = "B"
  }

  component_conditions {
    id    = newrelic_nrql_alert_condition.disk_full.id
    alias = "C"
  }
}
```

### With Facet Matching

```hcl
resource "newrelic_nrql_alert_condition" "high_throughput_per_host" {
  policy_id = newrelic_alert_policy.example.id
  name      = "High Throughput Per Host"
  enabled   = true

  nrql {
    query = "SELECT rate(count(*), 1 minute) FROM Transaction FACET host"
  }

  critical {
    operator              = "above"
    threshold             = 1000.0
    threshold_duration    = 300
    threshold_occurrences = "all"
  }

  violation_time_limit_seconds = 3600
}

resource "newrelic_nrql_alert_condition" "high_error_rate_per_host" {
  policy_id = newrelic_alert_policy.example.id
  name      = "High Error Rate Per Host"
  enabled   = true

  nrql {
    query = "SELECT percentage(count(*), WHERE error IS true) FROM Transaction FACET host"
  }

  critical {
    operator              = "above"
    threshold             = 5.0
    threshold_duration    = 300
    threshold_occurrences = "all"
  }

  violation_time_limit_seconds = 3600
}

resource "newrelic_alert_compound_condition" "with_facets" {
  policy_id               = newrelic_alert_policy.example.id
  name                    = "Host-Specific Alert"
  enabled                 = true
  trigger_expression      = "A AND B"
  facet_matching_behavior = "FACETS_MATCH"

  component_conditions {
    id    = newrelic_nrql_alert_condition.high_throughput_per_host.id
    alias = "A"
  }

  component_conditions {
    id    = newrelic_nrql_alert_condition.high_error_rate_per_host.id
    alias = "B"
  }
}
```

## Argument Reference

The following arguments are supported:

- `policy_id` - (Required) The ID of the policy where this alert compound condition should be used.
- `name` - (Required) The title of the compound alert condition.
- `trigger_expression` - (Required) Expression that defines how component condition evaluations are combined. Valid operators are 'AND', 'OR', 'NOT'. For more complex expressions, use parentheses. Use the aliases from `component_conditions` to build expressions like `"A AND B"`, `"A OR B"`, `"(A AND B) OR C"`, or `"A AND (B OR C) AND NOT (D AND E)"`.
- `component_conditions` - (Required) The list of conditions to be combined. Each component condition must be enabled. Must include at least 2. See [Component Conditions](#component-conditions) below for details.
- `enabled` - (Optional) Whether or not the compound alert condition is enabled. Defaults to `true`.
- `account_id` - (Optional) The New Relic account ID for managing your compound alert conditions. Defaults to the account ID set in your environment variable `NEW_RELIC_ACCOUNT_ID`.
- `facet_matching_behavior` - (Optional) How the compound condition will take into account the component conditions' facets during evaluation. Valid values are:
  - `FACETS_IGNORED` - (Default) Facets are not taken into consideration when determining when the compound alert condition activates
  - `FACETS_MATCH` - The compound alert condition will activate only when shared facets have matching values
- `runbook_url` - (Optional) Runbook URL to display in notifications.
- `threshold_duration` - (Optional) The duration, in seconds, that the trigger expression must be true before the compound alert condition will activate. Between 30-86400 seconds.

### Component Conditions

The `component_conditions` block supports the following arguments:

- `id` - (Required) The ID of the existing alert condition to use as a component.
- `alias` - (Required) The identifier that will be used in the compound alert condition's `trigger_expression` (e.g., 'b', 'b', 'c', 'd', 'e').

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the compound alert condition.

## Import

Compound alert conditions can be imported using the condition ID, e.g.

```bash
$ terraform import newrelic_alert_compound_condition.main 789012
```

## Additional Information

### Understanding Trigger Expressions

Trigger expressions define the logical conditions under which your alert compound condition will activate. Valid operators are:

- **AND** - Both conditions must be true
- **OR** - Either condition must be true
- **NOT** - Negates a condition
- **Parentheses** - Group conditions for complex logic

Examples:

- `"A AND B"` - Activate when both A and B are in violation
- `"A OR B"` - Activate when either A or B is in violation
- `"A AND NOT B"` - Activate when A is in violation but B is not
- `"(A AND B) OR C"` - Activate when both A and B are in violation, OR when C is in violation
- `"A AND (B OR C) AND NOT D"` - Activate when A is in violation AND either B or C is in violation AND D is not in violation

### Facet Matching Behavior

When your component NRQL conditions use FACET clauses:

- **FACETS_IGNORED** (Default) - Facets are not taken into consideration when determining when the compound alert condition activates. If component conditions have violations (on any facet), the compound alert condition will activate based on the trigger expression.
- **FACETS_MATCH** - The compound alert condition will activate only when shared facets have matching values. For example, if condition A fires for `host="server-1"` and condition B fires for `host="server-2"`, the compound alert condition will NOT activate because the facet values don't match.

### Threshold Duration

The `threshold_duration` parameter controls how long the trigger expression must remain true before the compound alert condition will activate.
