---
layout: "newrelic"
page_title: "New Relic: newrelic_pathpoint_flow"
sidebar_current: "docs-newrelic-resource-pathpoint-flow"
description: |-
  Create and manage a New Relic Pathpoint flow.
---

# Resource: newrelic\_pathpoint\_flow

Use this resource to create, update, and delete a New Relic Pathpoint flow.

Pathpoint is an enterprise platform tracker that models the health of the specific stages and steps that make up your business processes (e.g. checkout, payment, fulfilment). It provides a real-time view of the entire value chain and gives you immediate visibility into issues, allowing you to understand how they impact the overall customer journey.

A New Relic User API key is required to provision this resource. Set the `api_key` attribute in the `provider` block or the `NEW_RELIC_API_KEY` environment variable with your User API key.

-> **NOTE:** The `version` attribute is managed automatically by this resource and must not be set manually. It is updated in state after every create/update operation and is sent back to the API on every subsequent update for optimistic concurrency control. If the Pathpoint flow is modified outside of Terraform (e.g. from the New Relic UI), the `version` in state will become stale and subsequent `terraform apply` runs will fail with a version conflict. In that case, use `terraform import` to re-sync state.

## Example Usage

### Basic Flow

```hcl
resource "newrelic_pathpoint_flow" "checkout" {
  account_id       = 12345678
  name             = "Checkout Flow"
  description      = "End-to-end checkout pipeline"
  health_rollup    = "AUTOMATIC_ROLL_UP"
  refresh_interval = "FIVE_MINUTES"

  stages {
    name          = "Frontend"
    health_rollup = "AUTOMATIC_ROLL_UP"

    levels {
      steps {
        name = "Login Page"
        entity_search_query {
          query = "domain='BROWSER' AND name='Login'"
        }
      }

      steps {
        name = "Cart Page"
        entity_search_query {
          query = "domain='BROWSER' AND name='Cart'"
        }
      }
    }
  }

  stages {
    name          = "Backend"
    health_rollup = "AUTOMATIC_ROLL_UP"

    levels {
      steps {
        name = "Order Service"
        entity_search_query {
          query = "domain='APM' AND name='OrderService'"
        }
      }
    }
  }
}
```

### Flow with KPIs

```hcl
resource "newrelic_pathpoint_flow" "checkout" {
  account_id       = 12345678
  name             = "Checkout Flow with KPIs"
  health_rollup    = "AUTOMATIC_ROLL_UP"
  refresh_interval = "ONE_MINUTE"

  kpis {
    name        = "Order Success Rate"
    description = "Percentage of orders completed successfully"
    category    = "Revenue"

    query {
      from = "Transaction"

      select {
        aggregation_type = "COUNT"
        alias            = "orders"
      }

      time_window {
        relative_range {
          since = "SIXTY_MINUTES"
        }
      }

      where = "transactionType='Web' AND name='checkout'"
    }
  }

  stages {
    name          = "Payment"
    health_rollup = "AUTOMATIC_ROLL_UP"

    stage_kpis {
      name     = "Payment Errors"
      category = "Reliability"

      query {
        from = "TransactionError"

        select {
          aggregation_type = "COUNT"
          alias            = "errors"
        }

        time_window {
          custom_range = "SINCE 30 minutes ago"
        }
      }
    }

    levels {
      steps {
        name = "Payment Gateway"

        config {
          health_rollup   = "WORST_STATUS_WINS"
          threshold_type  = "FIXED"
          threshold_value = 1
        }

        entity_search_query {
          query = "domain='APM' AND name='PaymentGateway'"
        }

        signals {
          guid = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"
          name = "PaymentGateway"
          type = "ENTITY"
        }
      }
    }
  }
}
```

### Flow with Stage Relationships

```hcl
resource "newrelic_pathpoint_flow" "supply_chain" {
  account_id    = 12345678
  name          = "Supply Chain Flow"
  health_rollup = "AUTOMATIC_ROLL_UP"

  stages {
    name = "Warehouse"

    related {
      source = true
      target = false
    }

    levels {
      steps {
        name = "Inventory Check"
        entity_search_query {
          query = "domain='APM' AND name='InventoryService'"
        }
      }
    }
  }

  stages {
    name = "Shipping"

    related {
      source = false
      target = true
    }

    levels {
      steps {
        name = "Dispatch Service"
        entity_search_query {
          query = "domain='APM' AND name='DispatchService'"
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The New Relic account ID that owns this Pathpoint flow. Defaults to the provider account ID.
* `name` - (Required) The display name of the Pathpoint flow.
* `description` - (Optional) A brief description of the flow.
* `category` - (Optional) A category used to group flows (e.g. `Marketing`, `Checkout`).
* `health_rollup` - (Optional) Health rollup strategy for the flow, derived from its stages. Valid values: `ALERT_CONDITIONS`, `AUTOMATIC_ROLL_UP`.
* `refresh_interval` - (Optional) How often the flow, stage, level, and step health statuses are refreshed. Valid values: `ONE_MINUTE`, `FIVE_MINUTES`, `TEN_MINUTES`, `FIFTEEN_MINUTES`, `THIRTY_MINUTES`.
* `kpis` - (Optional) A list of Key Performance Indicators tracked at the flow level. See [Nested `kpis` blocks](#nested-kpis-blocks) below for details.
* `stages` - (Required) An ordered list of stages that make up this flow. Maximum 50 stages. See [Nested `stages` blocks](#nested-stages-blocks) below for details.

### Nested `kpis` blocks

* `name` - (Required) The display name of the KPI.
* `description` - (Optional) A description of the KPI.
* `category` - (Optional) A category to group KPIs.
* `account_id` - (Optional) The account ID this KPI belongs to. Defaults to the flow's account ID.
* `query` - (Required) The NRQL query definition for this KPI. See [Nested `query` blocks](#nested-query-blocks) below for details.

### Nested `query` blocks

* `from` - (Required) The NRQL data source to query (e.g. `Transaction`).
* `where` - (Optional) A NRQL `WHERE` clause to filter data.
* `select` - (Required) The `SELECT` clause defining what to aggregate. See [Nested `select` blocks](#nested-select-blocks) below for details.
* `time_window` - (Optional) The time window over which the KPI is evaluated. Provide either `custom_range` or `relative_range`, not both. See [Nested `time_window` blocks](#nested-time_window-blocks) below for details.

### Nested `select` blocks

* `aggregation_type` - (Required) The aggregation function. Valid values: `AVERAGE`, `COUNT`, `HISTOGRAM`, `MAX`, `MIN`, `PERCENTILE`, `SUM`, `UNIQUE_COUNT`.
* `alias` - (Optional) An alias for the aggregated value in query results.
* `attribute` - (Optional) The attribute name to aggregate. Required for all functions except `COUNT`.
* `threshold` - (Optional) The threshold value used in the selected function.

### Nested `time_window` blocks

* `custom_range` - (Optional) A raw NRQL time fragment, e.g. `SINCE 3 days ago COMPARE WITH 1 day ago`. Mutually exclusive with `relative_range`.
* `relative_range` - (Optional) A relative time window built from predefined durations. Mutually exclusive with `custom_range`. See [Nested `relative_range` blocks](#nested-relative_range-blocks) below for details.

### Nested `relative_range` blocks

* `since` - (Required) How far back the KPI is evaluated (maps to NRQL `SINCE`). Valid values: `SIXTY_MINUTES`, `THREE_HOURS`, `SIX_HOURS`, `TWENTY_FOUR_HOURS`, `SEVEN_DAYS`, `THIRTY_DAYS`, `THIRTY_MINUTES`.
* `compare_against` - (Optional) The earlier window to compare against (maps to NRQL `COMPARE WITH`). Valid values are the same as `since`.

### Nested `stages` blocks

* `name` - (Required) The display name of the stage.
* `health_rollup` - (Optional) Health rollup strategy for the stage. Valid values: `ALERT_CONDITIONS`, `AUTOMATIC_ROLL_UP`.
* `is_excluded` - (Optional) When `true`, this stage is excluded from flow health calculation. Defaults to `false`.
* `link` - (Optional) An optional URL to an external resource related to this stage.
* `related` - (Optional) Defines the relationship role of this stage within the flow. See [Nested `related` blocks](#nested-related-blocks) below for details.
* `stage_kpis` - (Optional) A list of KPIs tracked at the stage level. Uses the same schema as the top-level [`kpis` blocks](#nested-kpis-blocks).
* `levels` - (Optional) An ordered list of levels within this stage. Maximum 50 levels. See [Nested `levels` blocks](#nested-levels-blocks) below for details.

### Nested `related` blocks

* `source` - (Optional) When `true`, this stage acts as a source to other stages. Defaults to `false`.
* `target` - (Optional) When `true`, this stage acts as a target from other stages. Defaults to `false`.

### Nested `levels` blocks

* `steps` - (Optional) An ordered list of steps within this level. Maximum 50 steps. See [Nested `steps` blocks](#nested-steps-blocks) below for details.

### Nested `steps` blocks

* `name` - (Required) The display name of the step.
* `is_excluded` - (Optional) When `true`, this step is excluded from the level health calculation. Defaults to `false`.
* `link` - (Optional) An optional URL to an external resource related to this step.
* `scoped_accounts` - (Optional) A list of account IDs whose data is scoped to this step.
* `entity_search_query` - (Optional) A filter query used to automatically fetch signals for this step. See [Nested `entity_search_query` blocks](#nested-entity_search_query-blocks) below for details.
* `config` - (Optional) Health evaluation configuration for this step. See [Nested `config` blocks](#nested-config-blocks) below for details.
* `signals` - (Optional) A list of entity signals explicitly associated with this step. See [Nested `signals` blocks](#nested-signals-blocks) below for details.

### Nested `entity_search_query` blocks

* `query` - (Required) A filter query for signals, e.g. `domain='NR1' AND type='APPLICATION'`.
* `is_excluded` - (Optional) When `true`, this query is excluded from the step health calculation. Defaults to `false`.

### Nested `config` blocks

* `health_rollup` - (Optional) How the step health is derived from its signals. Valid values: `BEST_STATUS_WINS`, `WORST_STATUS_WINS`.
* `threshold_type` - (Optional) Whether the threshold is a `FIXED` number or a `PERCENTAGE`.
* `threshold_value` - (Optional) The numeric threshold value used to evaluate step health.

### Nested `signals` blocks

* `guid` - (Required) The entity GUID of the signal.
* `name` - (Optional) The display name of the signal.
* `type` - (Optional) Whether this GUID belongs to an entity or an alert condition. Valid values: `ENTITY`, `ALERT`.
* `is_excluded` - (Optional) When `true`, this signal is excluded from the step health calculation. Defaults to `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `guid` - The entity GUID assigned to this Pathpoint flow in New Relic.
* `version` - The last-updated epoch-millisecond timestamp used for optimistic concurrency control. This is managed automatically and must not be modified. It is persisted to state after every create/update and sent to the API on each subsequent update.
* `stages.#.id` - The internal workload ID of the stage. Populated after creation and used to identify stages on updates.
* `stages.#.levels.#.id` - The internal workload ID of the level. Populated after creation and used to identify levels on updates.
* `stages.#.levels.#.steps.#.id` - The internal workload ID of the step. Populated after creation and used to identify steps on updates.

## Import

New Relic Pathpoint flows can be imported using the flow's entity GUID, e.g.

```bash
$ terraform import newrelic_pathpoint_flow.checkout MjUyMDUyOHxOUjF8UFRIUFRTfDEyMzQ1
```

-> **NOTE:** After importing, run `terraform plan` to verify the state matches the existing configuration. Since the Pathpoint API has no read/get query, the imported state will only contain the GUID and version. Apply the resource once to fully populate all attributes in state.
