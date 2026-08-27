---
layout: "newrelic"
page_title: "New Relic: newrelic_pathpoint"
sidebar_current: "docs-newrelic-resource-pathpoint"
description: |-
  Create and manage a Pathpoint flow in New Relic.
---

# Resource: newrelic\_pathpoint

Use this resource to create and manage New Relic Pathpoint flows. Pathpoint is a New Relic app that lets you model your business flow and understand how each stage impacts the overall health of your system. Details regarding Pathpoint can be found [here](https://docs.newrelic.com/docs/new-relic-solutions/business-observability/introduction-pathpoint/).

## Example Usage

```hcl
resource "newrelic_pathpoint" "example" {
  account_id       = 12345678
  name             = "My Pathpoint Flow"
  category         = "Checkout"
  description      = "Monitors the end-to-end checkout flow"
  health_rollup    = "AUTOMATIC_ROLL_UP"
  refresh_interval = "ONE_MINUTE"

  kpis {
    name        = "Flow KPI"
    account_id  = 12345678

    query {
      from  = "Transaction"
      where = "appName = 'MyApp'"

      select {
        aggregation_type = "COUNT"
      }

      time_window {
        relative_range {
          since           = "SIXTY_MINUTES"
          compare_against = "TWENTY_FOUR_HOURS"
        }
      }
    }
  }

  stages {
    name         = "Browsing"
    health_rollup = "AUTOMATIC_ROLL_UP"
    is_excluded  = false
    link         = "https://example.com/browsing"

    related {
      source = false
      target = true
    }

    levels {
      steps {
        name        = "Product Page"
        is_excluded = false
        link        = "https://example.com/product"

        scoped_accounts = [12345678]

        config {
          health_rollup   = "WORST_STATUS_WINS"
          threshold_type  = "PERCENTAGE"
          threshold_value = 90
        }

        entity_search_query {
          query       = "domain = 'APM' AND type = 'APPLICATION'"
          is_excluded = false
        }

        signals {
          guid        = "MXxBUE18QVBQTElDQVRJT058MTIz"
          name        = "My App"
          type        = "ENTITY"
          is_excluded = false
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The New Relic account ID where the Pathpoint flow will be created. Defaults to the account associated with the API key used. Forces a new resource if changed.
* `name` - (Required) The display name of the Pathpoint flow.
* `category` - (Optional) A category used to group flows (e.g. Marketing, Checkout).
* `description` - (Optional) A description of the Pathpoint flow.
* `health_rollup` - (Optional) Health rollup strategy for the flow, derived from its stages. One of: `ALERT_CONDITIONS`, `AUTOMATIC_ROLL_UP`.
* `refresh_interval` - (Optional) How often the flow, stage, level, and step health statuses are fetched. One of: `FIFTEEN_MINUTES`, `FIVE_MINUTES`, `ONE_MINUTE`, `TEN_MINUTES`, `THIRTY_MINUTES`.
* `kpis` - (Optional) A list of KPIs to track at the flow level. See [Nested `kpis` blocks](#nested-kpis-blocks) below for details.
* `stages` - (Optional) The ordered list of stages that make up this flow. See [Nested `stages` blocks](#nested-stages-blocks) below for details.

### Nested `kpis` blocks

* `name` - (Required) The display name of the KPI.
* `account_id` - (Optional) The account ID this KPI belongs to.
* `category` - (Optional) Optional category to group KPIs.
* `description` - (Optional) Optional description of this KPI.
* `query` - (Optional) The NRQL query definition for this KPI. See [Nested `query` blocks](#nested-query-blocks) below for details.

### Nested `query` blocks

* `from` - (Required) The data source to query from (e.g., `Transaction`, `Metric`, `Log`).
* `where` - (Optional) Optional WHERE clause conditions to filter the data.
* `select` - (Optional) The SELECT clause defining what to aggregate and return. See [Nested `select` blocks](#nested-select-blocks) below for details.
* `time_window` - (Optional) The time window over which the KPI is evaluated. See [Nested `time_window` blocks](#nested-time_window-blocks) below for details.

### Nested `select` blocks

* `aggregation_type` - (Required) The aggregation function to apply to the attribute. One of: `AVERAGE`, `COUNT`, `HISTOGRAM`, `MAX`, `MIN`, `PERCENTILE`, `SUM`, `UNIQUE_COUNT`.
* `alias` - (Optional) Optional alias for the aggregated value in the query result.
* `attribute` - (Optional) The attribute name to aggregate. Required for all functions except `COUNT`.
* `threshold` - (Optional) The threshold used in the selected function.

### Nested `time_window` blocks

* `custom_range` - (Optional) A raw NRQL time fragment, e.g. `SINCE 3 days ago COMPARE WITH 1 day ago`. Mutually exclusive with `relative_range`.
* `relative_range` - (Optional) A relative window built from predefined duration values. Mutually exclusive with `custom_range`. See [Nested `relative_range` blocks](#nested-relative_range-blocks) below for details.

### Nested `relative_range` blocks

* `since` - (Required) How far back the KPI is evaluated (maps to NRQL `SINCE`). One of: `SEVEN_DAYS`, `SIXTY_MINUTES`, `SIX_HOURS`, `THIRTY_DAYS`, `THIRTY_MINUTES`, `THREE_HOURS`, `TWENTY_FOUR_HOURS`.
* `compare_against` - (Optional) The earlier window the current window is compared against (maps to NRQL `COMPARE WITH`). One of: `SEVEN_DAYS`, `SIXTY_MINUTES`, `SIX_HOURS`, `THIRTY_DAYS`, `THIRTY_MINUTES`, `THREE_HOURS`, `TWENTY_FOUR_HOURS`.

### Nested `stages` blocks

* `name` - (Required) The display name of the stage.
* `health_rollup` - (Optional) Health rollup strategy for the stage, derived from its levels. One of: `ALERT_CONDITIONS`, `AUTOMATIC_ROLL_UP`.
* `is_excluded` - (Optional) When true, this stage is excluded from the flow health calculation. Defaults to `false`.
* `link` - (Optional) Optional URL to an external resource related to this stage.
* `related` - (Optional) Defines the relationship role of a stage within a flow. See [Nested `related` blocks](#nested-related-blocks) below for details.
* `stage_kpis` - (Optional) KPIs tracked at the stage level. See [Nested `kpis` blocks](#nested-kpis-blocks) above for details.
* `levels` - (Optional) The ordered list of levels within this stage. See [Nested `levels` blocks](#nested-levels-blocks) below for details.

### Nested `related` blocks

* `source` - (Optional) When true, this stage acts as a source to other stages. Defaults to `false`.
* `target` - (Optional) When true, this stage acts as a target to other stages. Defaults to `false`.

### Nested `levels` blocks

* `steps` - (Optional) The ordered list of steps within this level. See [Nested `steps` blocks](#nested-steps-blocks) below for details.

### Nested `steps` blocks

* `name` - (Required) The display name of the step.
* `is_excluded` - (Optional) When true, this step is excluded from the level health calculation. Defaults to `false`.
* `link` - (Optional) Optional URL to an external resource related to this step.
* `scoped_accounts` - (Optional) Account IDs whose entities are included in scope for this step.
* `config` - (Optional) Health evaluation configuration for this step. See [Nested `config` blocks](#nested-config-blocks) below for details.
* `entity_search_query` - (Optional) The filter query used to fetch signals for this step. See [Nested `entity_search_query` blocks](#nested-entity_search_query-blocks) below for details.
* `signals` - (Optional) Entity signals associated with this step. See [Nested `signals` blocks](#nested-signals-blocks) below for details.

### Nested `config` blocks

* `health_rollup` - (Optional) How the step health is rolled up from its signals. One of: `BEST_STATUS_WINS`, `WORST_STATUS_WINS`.
* `threshold_type` - (Optional) Whether the threshold value is a fixed number or a percentage. One of: `FIXED`, `PERCENTAGE`.
* `threshold_value` - (Optional) The numeric threshold value used to evaluate step health.

### Nested `entity_search_query` blocks

* `query` - (Required) Filter query for signals in this step, e.g. `domain='NR1' AND type='APPLICATION'`.
* `is_excluded` - (Optional) When true, this query is excluded from the step health evaluation. Defaults to `false`.

### Nested `signals` blocks

* `guid` - (Required) The entity GUID of the signal.
* `name` - (Optional) The display name of the signal.
* `type` - (Optional) Whether this GUID belongs to an entity or an alert condition. One of: `ALERT`, `ENTITY`.
* `is_excluded` - (Optional) When true, this signal is excluded from the step health evaluation. Defaults to `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The entity GUID of the Pathpoint flow (same as `guid`).
* `guid` - The unique entity GUID assigned to this Pathpoint flow.
* `flow_id` - The workload ID of the flow entity.
* `version` - The last updated timestamp of the Pathpoint, used for version control.
* `health_status` - The computed health status of the flow. One of: `DEGRADED`, `DISRUPTED`, `OPERATIONAL`, `UNKNOWN`.
* `message` - Error details when one or more stages, levels, or steps fail during a create or update operation.
* `excluded_kpis` - Names of KPIs that were excluded during create or update.

The following computed attributes are also exported on nested blocks:

* `stages[*].stage_id` - The workload ID of the stage entity.
* `stages[*].health_status` - The computed health status of the stage.
* `stages[*].levels[*].level_id` - The workload ID of the level entity.
* `stages[*].levels[*].health_status` - The computed health status of the level.
* `stages[*].levels[*].steps[*].step_id` - The workload ID of the step entity.
* `stages[*].levels[*].steps[*].health_status` - The computed health status of the step.
* `kpis[*].kpi_id` - The unique identifier of the KPI.
* `stages[*].stage_kpis[*].kpi_id` - The unique identifier of the stage KPI.

## Import

Pathpoint flows can be imported using the entity GUID:

```
terraform import newrelic_pathpoint.example <guid>
```

~> **NOTE:** The `account_id` is encoded in the entity GUID. After import, run `terraform state show newrelic_pathpoint.example` to view all attributes.

## Additional Information

More information about Pathpoint can be found in the New Relic [documentation](https://docs.newrelic.com/docs/new-relic-solutions/business-observability/introduction-pathpoint/).