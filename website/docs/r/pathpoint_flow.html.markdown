---
layout: "newrelic"
page_title: "New Relic: newrelic_pathpoint_flow"
sidebar_current: "docs-newrelic-resource-pathpoint-flow"
description: |-
  Create and manage a New Relic Pathpoint flow.
---

# Resource: newrelic\_pathpoint\_flow

-> **LIMITED PREVIEW:** This resource is in limited preview and is only available for accounts that have been granted access. Features and behavior may change before general availability.

Use this resource to create, read, update, and delete a New Relic Pathpoint flow.

Pathpoint is New Relic's business journey observability product. It helps map your technical system health to the business flows that matter most — so when something goes wrong, you can immediately see where in the customer journey the problem is occurring.

Each Flow in Pathpoint represents a business journey — like checkout, authentication, or onboarding — broken into a set of stages. Pathpoint shows you the health of each stage, helping you move faster from alert to impact to resolution.

A New Relic User API key is required to provision this resource. Set the `api_key` attribute in the `provider` block or the `NEW_RELIC_API_KEY` environment variable with your User API key.

-> **NOTE:** The `version` attribute is managed automatically by this resource and must not be set manually. It is refreshed from the API on every `terraform plan` and `terraform apply`. Any changes made outside of Terraform (e.g. from the New Relic UI) will be detected but **overridden** on the next apply. Review the plan output carefully before applying to avoid unintentionally overwriting manual changes.

## Example Usage

### Basic Flow

A minimal flow demonstrating the three ways to populate a step with signals.

`health_rollup` is omitted here — both the flow and each stage default to `AUTOMATIC_ROLL_UP`, which means health is derived by rolling up the status of child objects automatically. `refresh_interval` defaults to `THIRTY_MINUTES`.

Steps can be populated in three ways:
- **Dynamic query** (`entity_search_query`): the API runs the filter and auto-discovers matching entities at each refresh. The search is scoped to the accounts in `scoped_accounts`; if that is not set, it falls back to the account the flow belongs to. To target a specific account, include `accountId` directly in the filter query.
- **Entity signal** (`type = "ENTITY"`): a specific New Relic entity pinned by its GUID. Use when you always want the same exact entity regardless of naming changes.
- **Alert signal** (`type = "ALERT"`): an alert condition pinned by its entity GUID. Use when step health should be driven by an alert policy rather than entity telemetry.

```hcl
resource "newrelic_pathpoint_flow" "checkout" {
  account_id       = 12345678
  name             = "Checkout Flow"
  description      = "End-to-end checkout pipeline"
  refresh_interval = "FIVE_MINUTES"
  # health_rollup defaults to AUTOMATIC_ROLL_UP

  stages {
    name = "Frontend"
    # health_rollup defaults to AUTOMATIC_ROLL_UP

    levels {
      # Dynamic query: API auto-discovers entities matching the filter at each refresh.
      # Account scope: uses scoped_accounts if set, otherwise defaults to the flow's account.
      # Include accountId in the filter to target a specific account explicitly.
      steps {
        name = "Login Page"
        entity_search_query {
          query = "accountId=12345678 AND domain='BROWSER' AND name='Login'"
        }
      }

      # Entity signal: a specific New Relic entity pinned by GUID.
      steps {
        name = "Cart Service"
        signals {
          guid = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1"
          name = "Cart Service"
          type = "ENTITY"
        }
      }

      # Alert signal: an alert condition pinned by its entity GUID.
      steps {
        name = "Checkout Alert"
        signals {
          guid = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk2"
          name = "Checkout Error Rate"
          type = "ALERT"
        }
      }
    }
  }

  stages {
    name = "Backend"

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

Demonstrates flow-level and stage-level KPIs with all three `time_window` options: `relative_range` with a simple lookback, `relative_range` with `compare_against` for period-over-period comparison (equivalent to NRQL `COMPARE WITH`), and `custom_range` for a free-form NRQL time expression. Also shows a step `config` block for threshold-based health evaluation.

```hcl
resource "newrelic_pathpoint_flow" "checkout" {
  account_id       = 12345678
  name             = "Checkout Flow with KPIs"
  health_rollup    = "AUTOMATIC_ROLL_UP"
  refresh_interval = "ONE_MINUTE"

  # relative_range: simple lookback window.
  kpis {
    name        = "Order Success Rate"
    description = "Percentage of orders completed successfully"
    category    = "Revenue"

    query {
      from  = "Transaction"
      where = "transactionType='Web' AND name='checkout'"

      select {
        aggregation_type = "COUNT"
        alias            = "orders"
      }

      time_window {
        relative_range {
          since = "SIXTY_MINUTES"
        }
      }
    }
  }

  # relative_range with compare_against: period-over-period comparison (NRQL COMPARE WITH).
  kpis {
    name        = "Average Response Time"
    description = "Average response time vs the same window 7 days ago"
    category    = "Latency"

    query {
      from  = "Transaction"
      where = "transactionType='Web'"

      select {
        aggregation_type = "AVERAGE"
        attribute        = "duration"
        alias            = "avg_duration"
      }

      time_window {
        relative_range {
          since           = "TWENTY_FOUR_HOURS"
          compare_against = "SEVEN_DAYS"
        }
      }
    }
  }

  stages {
    name          = "Payment"
    health_rollup = "AUTOMATIC_ROLL_UP"

    # custom_range: free-form NRQL time expression for cases not covered by relative_range.
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

### Flow with Excluded Stages

Shows `is_excluded` at every level of the hierarchy — stage, step, entity search query, and individual signal — to remove items from health calculation without deleting them. Also demonstrates `health_rollup = "ALERT_CONDITIONS"` on a stage — when set, the stage's health is tied directly to its KPI alert conditions rather than entity telemetry, making it suitable for business-health tracking where you want the stage status to reflect whether business KPIs are breaching alert thresholds.

```hcl
resource "newrelic_pathpoint_flow" "pipeline" {
  account_id    = 12345678
  name          = "Order Pipeline"
  health_rollup = "AUTOMATIC_ROLL_UP"

  stages {
    name          = "Business Health"
    # ALERT_CONDITIONS ties this stage's health to its KPI alert conditions.
    # Use this when the stage represents a business outcome rather than system telemetry.
    health_rollup = "ALERT_CONDITIONS"

    stage_kpis {
      name     = "Revenue Rate"
      category = "Business"

      query {
        from  = "Transaction"
        where = "transactionType='Web' AND name='purchase'"

        select {
          aggregation_type = "COUNT"
          alias            = "purchases"
        }

        time_window {
          relative_range {
            since = "SIXTY_MINUTES"
          }
        }
      }
    }

    levels {
      steps {
        name = "Purchase Service"
        entity_search_query {
          query = "domain='APM' AND name='PurchaseService'"
        }
      }

      # Step excluded from health calculation — kept in config but not counted.
      steps {
        name        = "Deprecated Checkout"
        is_excluded = true

        entity_search_query {
          query = "domain='APM' AND name='DeprecatedCheckout'"
        }

        # Individual signal excluded — present for visibility but ignored in health rollup.
        signals {
          guid        = "MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk3"
          name        = "Deprecated Checkout Alert"
          type        = "ALERT"
          is_excluded = true
        }
      }
    }
  }

  stages {
    name        = "Legacy Processor"
    is_excluded = true

    levels {
      steps {
        name = "Legacy API"

        # entity_search_query excluded — query is retained but its results are ignored.
        entity_search_query {
          query       = "domain='APM' AND name='LegacyAPI'"
          is_excluded = true
        }
      }
    }
  }
}
```

### Flow with Stage Relationships

The `related` block is purely a UI hint — it tells the Pathpoint UI whether to render a stage as a sequential (arrow-shaped) node or a non-sequential (rectangular) node. Stages without a `related` block are treated as sequential by default.

The values follow the stage's position in the array:
- **First stage** — `source = false, target = true`: no incoming connection, has an outgoing connection to the next stage.
- **Middle stages** — `source = true, target = true`: connected on both sides.
- **Last stage** — `source = true, target = false`: has an incoming connection, no outgoing connection.

Stages that are part of an unbroken `source/target` chain render as sequential arrow shapes in the UI. Setting `related` does not affect health rollup — it only changes the visual shape.

```hcl
resource "newrelic_pathpoint_flow" "supply_chain" {
  account_id    = 12345678
  name          = "Supply Chain Flow"
  health_rollup = "AUTOMATIC_ROLL_UP"

  # First stage: no incoming connection, feeds into the next stage.
  stages {
    name = "Warehouse"

    related {
      source = false
      target = true
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

  # Middle stage: connected on both sides.
  stages {
    name = "Packaging"

    related {
      source = true
      target = true
    }

    levels {
      steps {
        name = "Packaging Service"
        entity_search_query {
          query = "domain='APM' AND name='PackagingService'"
        }
      }
    }
  }

  # Last stage: receives from the previous stage, no outgoing connection.
  stages {
    name = "Shipping"

    related {
      source = true
      target = false
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

KPIs are numeric metrics derived from NRQL queries, displayed as scorecards above the flow or stage. Flow-level KPIs are visible across all stages; stage-level KPIs (`stage_kpis`) are scoped to a single stage.

* `name` - (Required) The display name of the KPI.
* `description` - (Optional) A short description explaining what the KPI measures.
* `category` - (Optional) A label used to group related KPIs (e.g. `Revenue`, `Reliability`).
* `account_id` - (Optional) The account whose data this KPI queries. Defaults to the flow's account ID.
* `query` - (Required) The NRQL query definition for this KPI. See [Nested `query` blocks](#nested-query-blocks) below for details.

### Nested `query` blocks

Defines the NRQL query that backs a KPI. The result is a single aggregated value rendered as the KPI score.

* `from` - (Required) The NRQL event type to query (e.g. `Transaction`, `JavaScriptError`).
* `where` - (Optional) A NRQL `WHERE` clause to filter which events are included.
* `select` - (Required) The aggregation to compute. See [Nested `select` blocks](#nested-select-blocks) below for details.
* `time_window` - (Optional) The time window over which the KPI is evaluated. Provide either `custom_range` or `relative_range`, not both. See [Nested `time_window` blocks](#nested-time_window-blocks) below for details.

### Nested `select` blocks

Controls the NRQL aggregation function and which attribute to aggregate.

* `aggregation_type` - (Required) The aggregation function. Valid values: `AVERAGE`, `COUNT`, `HISTOGRAM`, `MAX`, `MIN`, `PERCENTILE`, `SUM`, `UNIQUE_COUNT`.
* `attribute` - (Optional) The event attribute to aggregate. Required for all functions except `COUNT`.
* `alias` - (Optional) A display name for the aggregated result.
* `threshold` - (Optional) A threshold value used by the selected aggregation function (e.g. the percentile for `PERCENTILE`).

### Nested `time_window` blocks

Scopes the KPI query to a specific time range. Use `custom_range` for free-form NRQL time expressions, or `relative_range` for predefined durations.

* `custom_range` - (Optional) A raw NRQL time expression, e.g. `SINCE 3 days ago COMPARE WITH 1 day ago`. Mutually exclusive with `relative_range`.
* `relative_range` - (Optional) A relative time window built from predefined durations. Mutually exclusive with `custom_range`. See [Nested `relative_range` blocks](#nested-relative_range-blocks) below for details.

### Nested `relative_range` blocks

Defines a named time window using fixed duration values instead of a raw NRQL string.

* `since` - (Required) How far back the KPI is evaluated (maps to NRQL `SINCE`). Valid values: `THIRTY_MINUTES`, `SIXTY_MINUTES`, `THREE_HOURS`, `SIX_HOURS`, `TWENTY_FOUR_HOURS`, `SEVEN_DAYS`, `THIRTY_DAYS`.
* `compare_against` - (Optional) An earlier window to compare the current result against (maps to NRQL `COMPARE WITH`). Valid values are the same as `since`.

### Nested `stages` blocks

Stages are the top-level groupings within a flow, representing major phases of a business process (e.g. `Frontend`, `Payment`, `Fulfilment`). Stages are ordered by their position in the configuration array. Each stage can carry its own KPIs and contains one or more levels.

* `name` - (Required) The display name of the stage.
* `health_rollup` - (Optional) How this stage's health is derived from its levels. Defaults to `AUTOMATIC_ROLL_UP`. Valid values: `ALERT_CONDITIONS`, `AUTOMATIC_ROLL_UP`.
* `is_excluded` - (Optional) When `true`, this stage is excluded from the flow's health calculation without being deleted. Defaults to `false`.
* `link` - (Optional) A URL to an external resource (e.g. runbook, wiki page) associated with this stage.
* `related` - (Optional) Controls the stage's visual shape and position in the sequential chain. See [Nested `related` blocks](#nested-related-blocks) below for details.
* `stage_kpis` - (Optional) KPIs scoped to this stage only. Uses the same schema as the top-level [`kpis` blocks](#nested-kpis-blocks).
* `levels` - (Optional) An ordered list of levels within this stage. Maximum 50 levels. See [Nested `levels` blocks](#nested-levels-blocks) below for details.

### Nested `related` blocks

Controls the visual shape of the stage in the Pathpoint UI and whether it participates in a sequential chain. This is a UI-only hint and does not affect health rollup. Stages without a `related` block are treated as sequential by default.

Set values based on the stage's position in the array: first stage uses `source = false, target = true`; last stage uses `source = true, target = false`; middle stages use `source = true, target = true`. Stages in an unbroken chain render as sequential arrow shapes.

* `source` - (Optional) When `true`, this stage has an incoming connection from the preceding stage. Defaults to `false`.
* `target` - (Optional) When `true`, this stage has an outgoing connection to the following stage. Defaults to `false`.

### Nested `levels` blocks

Levels group steps within a stage, allowing multiple parallel tracks of signals to be evaluated independently. A stage can have up to 50 levels.

* `steps` - (Optional) An ordered list of steps within this level. Maximum 50 steps. See [Nested `steps` blocks](#nested-steps-blocks) below for details.

### Nested `steps` blocks

A step represents a single unit of work or signal group within a level. Each step's health is derived from its signals — either auto-discovered via `entity_search_query`, explicitly listed via `signals`, or both.

* `name` - (Required) The display name of the step.
* `is_excluded` - (Optional) When `true`, this step is excluded from the level's health calculation without being deleted. Defaults to `false`.
* `link` - (Optional) A URL to an external resource associated with this step.
* `scoped_accounts` - (Optional) A list of account IDs to restrict this step's signal search to. Useful in multi-account setups to prevent signals from unrelated accounts being included.
* `entity_search_query` - (Optional) A dynamic filter that the API evaluates at each refresh to auto-discover matching signals. See [Nested `entity_search_query` blocks](#nested-entity_search_query-blocks) below for details.
* `config` - (Optional) Health evaluation thresholds for this step. See [Nested `config` blocks](#nested-config-blocks) below for details.
* `signals` - (Optional) Explicitly pinned signals for this step. See [Nested `signals` blocks](#nested-signals-blocks) below for details.

### Nested `entity_search_query` blocks

A dynamic filter query that the API evaluates at each refresh interval to automatically find and attach matching entities to the step. The search is scoped to the accounts listed in `scoped_accounts`; if `scoped_accounts` is not set, it defaults to the account the flow belongs to. To target a specific account regardless of `scoped_accounts`, include `accountId` directly in the filter expression.

* `query` - (Required) A New Relic entity filter expression, e.g. `accountId=12345678 AND domain='APM' AND name='OrderService'`.
* `is_excluded` - (Optional) When `true`, results from this query are excluded from the step's health calculation. Defaults to `false`.

### Nested `config` blocks

Configures how the step's health status is derived from its signals and what threshold triggers a status change.

* `health_rollup` - (Optional) How the step's health is derived from its signals. `WORST_STATUS_WINS` marks the step unhealthy if any signal is unhealthy; `BEST_STATUS_WINS` marks it healthy if any signal is healthy. Valid values: `BEST_STATUS_WINS`, `WORST_STATUS_WINS`.
* `threshold_type` - (Optional) Whether the threshold is a `FIXED` count or a `PERCENTAGE` of signals that must be healthy.
* `threshold_value` - (Optional) The numeric threshold value applied against `threshold_type` to determine step health.

### Nested `signals` blocks

Explicitly pins a specific New Relic entity or alert condition to the step by GUID. Use this when you always want the same signal regardless of naming or tagging changes.

* `guid` - (Required) The entity GUID of the signal to attach.
* `name` - (Optional) A display name for the signal as it appears in the step.
* `type` - (Optional) Whether the GUID refers to a monitored entity (`ENTITY`) or an alert condition (`ALERT`). Valid values: `ENTITY`, `ALERT`.
* `is_excluded` - (Optional) When `true`, this signal is excluded from the step's health calculation. Defaults to `false`.

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

-> **NOTE:** After importing, run `terraform plan` to verify the state matches the existing configuration. The provider will read the current flow configuration from the API and populate all attributes in state.
