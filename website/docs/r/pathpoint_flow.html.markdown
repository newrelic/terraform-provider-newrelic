---
layout: "newrelic"
page_title: "New Relic: newrelic_pathpoint_flow"
sidebar_current: "docs-newrelic-resource-pathpoint-flow"
description: |-
  Create and manage a New Relic Pathpoint flow.
---

# Resource: newrelic\_pathpoint\_flow

-> **LIMITED PREVIEW:** This resource is in limited preview and is only available for accounts that have been granted access. Features and behavior may change before general availability.

Pathpoint maps the health of your technical systems onto the business journeys they support. Each Flow in Pathpoint represents one journey — checkout, authentication, or onboarding — broken into stages, so when something goes wrong you can see which part of the customer journey the problem affects.

Use this resource to create, read, update, and delete a New Relic Pathpoint flow.



A New Relic User API key is required to provision this resource. Set the `api_key` attribute in the `provider` block or the `NEW_RELIC_API_KEY` environment variable with your User API key.

-> **NOTE:** Any manual changes made outside Terraform (e.g., via the New Relic UI) will be automatically overridden on the next terraform apply. Review your plan output carefully, if you want to keep external changes, [export the existing flow](#exporting-an-existing-flow)

## Example Usage

A flow is made up of `stages`, each stage made up of `levels`, each level made up of `steps`. Each step contains `signals` — entities, alerts, or entities discovered dynamically via `entity_search_query` — and its health is derived from those signals.

### Stages

Stages are ordered by their position in the configuration array and represent major phases of a business process — for example `Frontend`, `Payment`, and `Fulfilment` in an e-commerce checkout flow

The example below shows two stages with their most common options — `health_rollup` (defaults to `AUTOMATIC_ROLL_UP`), `is_excluded`, `related` (sequential chain hints), and `link` (runbook URL). Stages without a `related` block are sequential by default.

```hcl
stages {
  name          = "Frontend"
  health_rollup = "AUTOMATIC_ROLL_UP"  # default: rolls up from levels → steps → signals
  is_excluded   = false                # default: participates in flow health
  link          = "https://runbooks.example.com/checkout/frontend"

  related {
    source = false  # first stage: no incoming connection
    target = true   # has outgoing connection to the next stage
  }

  levels {
    steps {
      name = "Login Page"
      signals {
        guid = "GUID1"
        name = "Login Service"
        type = "ENTITY"
      }
    }
  }
}

stages {
  name = "Payment"
  link = "https://runbooks.example.com/checkout/payment"

  related {
    source = true  # connected to the previous stage
    target = false # last stage: no outgoing connection
  }

  levels {
    steps {
      name = "Payment API"
      entity_search_query {
        query = "accountId=1234 AND domain='APM' AND name='PaymentService'"
      }
    }
  }
}
```

### Steps

A step's signals can come from three sources:

- **Dynamic query** (`entity_search_query`): the API runs the filter and auto-discovers matching entities at each refresh. The search is scoped to the accounts in `scoped_accounts`; if that is not set, it falls back to the account the flow belongs to. To target a specific account, include `accountId` directly in the filter query.
- **Entity signal** (`signals` with `type = "ENTITY"`): a specific New Relic entity pinned by its GUID. Use when you always want the same exact entity regardless of naming changes.
- **Alert signal** (`signals` with `type = "ALERT"`): an alert condition pinned by its entity GUID. Use when step health should be driven by an alert policy rather than entity telemetry.

-> **NOTE:** Alert signals are not yet supported in this Limited Preview — the `type = "ALERT"` signal shown below is included to illustrate the shape of the config.

```hcl
steps {
  name = "Login Page"

  # Dynamic query: API auto-discovers entities matching the filter at each refresh.
  entity_search_query {
    query = "accountId=1234 AND domain='BROWSER' AND name='Login'"
  }

  # Entity signal: a specific New Relic entity pinned by GUID.
  signals {
    guid = "GUID1"
    name = "Cart Service"
    type = "ENTITY"
  }

  # Alert signal: an alert condition pinned by its entity GUID.
  signals {
    guid = "GUID2"
    name = "Checkout Error Rate"
    type = "ALERT"
  }
}
```

### KPIs

KPIs are numeric metrics displayed as scorecards on the flow. You define a NRQL query — `from` points at a New Relic event type (e.g. `Transaction`, `JavaScriptError`), and Pathpoint synthesizes it into a single aggregated metric value.

- **Flow-level `kpis`** — visible across the entire flow, shown above all stages.
- **Stage-level `stage_kpis`** — scoped to a single stage, shown on that stage's  only.

Both use the same schema.

```hcl
# Flow-level KPI: visible across the entire flow.
# `from` is a New Relic event type — Pathpoint synthesizes it into a single metric value.
kpis {
  name        = "Order Success Rate"
  description = "Percentage of orders completed successfully"
  category    = "Revenue"

  query {
    from  = "Transaction"        # event type — the raw data source
    where = "name='checkout'"

    select {
      aggregation_type = "COUNT"
      alias            = "orders"
    }
  }
}

# Stage-level KPI: scoped to one stage only, shown on that stage's card.
# Useful when a stage has a distinct business metric separate from the flow's global KPIs.
stage_kpis {
  name     = "Payment Errors"
  category = "Reliability"

  query {
    from = "TransactionError"    # event type — synthesized into an error-count metric

    select {
      aggregation_type = "COUNT"
      alias            = "errors"
    }
  }
}
```

-> **NOTE:** Cross-account KPIs — setting `account_id` on a `kpis`/`stage_kpis` block to an account other than the flow's own — are not yet supported in this Limited Preview.

### Health

Health can be evaluated at the flow, stage, and step level.

#### Flow-level health

`health_rollup` on the flow controls how the flow's overall health is derived:

- `AUTOMATIC_ROLL_UP` (the default) — health rolls up automatically from the flow's stages.

  - `is_excluded` controls whether the stage contributes to the flow's overall health calculation:
    - `true` — the stage is excluded from the flow's health rollup. Its levels, steps, and signals are still evaluated and shown in the Pathpoint UI, but the stage does not affect the flow's overall health status. Useful for stages under construction or temporarily removed from the scope.
    - `false` (default) — the stage participates in the flow's health rollup normally.

- `ALERT_CONDITIONS` — health is tied directly to the flow's KPI alert conditions instead of stage rollup. Typically paired with flow-level `kpis`.

```hcl
resource "newrelic_pathpoint_flow" "checkout" {
  name          = "Checkout Flow"
  health_rollup = "AUTOMATIC_ROLL_UP"
  # ...
}
```

#### Stage-level health

`health_rollup` controls how a stage's health is derived:

- `AUTOMATIC_ROLL_UP` (the default) — health rolls up automatically from that stage's levels, through their steps and signals.
- `ALERT_CONDITIONS` — health is tied directly to the stage's KPI alert conditions instead of step signals. Use this when a stage's health status must reflect the underlying business outcomes. For example “e-mail open rate kpi” associated with Marketing Campaigns stage
 
```hcl
stages {
  name          = "Revenue"
  health_rollup = "AUTOMATIC_ROLL_UP"
  is_excluded   = false
  # ...
}
```

#### Step-level health

A step's `config` block controls how its health is derived from its signals:

- `health_rollup` — `WORST_STATUS_WINS` (the default) marks the step unhealthy if any signal is unhealthy; `BEST_STATUS_WINS` marks it healthy if any signal is healthy.
  - `is_excluded` — can be set on a step, an individual `signals` entry, or `entity_search_query` to remove that item from health calculation without deleting it. Defaults to `false`.
  - `threshold_type` / `threshold_value` — instead of an all-or-nothing rollup, require a `FIXED` count or `PERCENTAGE` of signals to be healthy before the step is considered healthy.

```hcl
steps {
  name        = "Checkout API"
  is_excluded = false

  config {
    health_rollup   = "WORST_STATUS_WINS"
    threshold_type  = "PERCENTAGE"
    threshold_value = 80
  }

  entity_search_query {
    query       = "accountId=1234 AND domain='APM' AND name LIKE 'checkout-%'"
    is_excluded = false
  }

  signals {
    guid        = "GUID3"
    name        = "Deprecated Checkout Alert"
    type        = "ALERT"
    is_excluded = true
  }
}
```

### Example

The example below is a **Checkout** flow with 2 stages (`Revenue` and `Frontend`), a flow-level KPI (`Order Success Rate`), a stage-level KPI (`Payment Errors`), and a step that uses all three signal types — an entity signal (`GUID1`), an alert signal (`GUID2`), and a dynamic `entity_search_query`.

```hcl
resource "newrelic_pathpoint_flow" "checkout" {
  account_id       = 1234
  name             = "Checkout Flow"
  description      = "End-to-end checkout pipeline"
  refresh_interval = "FIVE_MINUTES"  # defaults to FIVE_MINUTES if not set

  kpis {
    name        = "Order Success Rate"
    description = "Percentage of orders completed successfully"
    category    = "Revenue"

    query {
      from  = "Transaction"
      where = "name='checkout'"

      select {
        aggregation_type = "COUNT"
        alias            = "orders"
      }
    }
  }

  stages {
    name          = "Revenue"
    health_rollup = "ALERT_CONDITIONS"
    link          = "https://runbooks.example.com/checkout/revenue"

    related {
      source = false
      target = true
    }

    stage_kpis {
      name     = "Payment Errors"
      category = "Reliability"

      query {
        from = "TransactionError"

        select {
          aggregation_type = "COUNT"
          alias            = "errors"
        }
      }
    }

    levels {
      steps {
        name = "Order Service"
        entity_search_query {
          query = "accountId=123 AND domain='APM' AND name='OrderService'"
        }
      }
    }
  }

  stages {
    name = "Frontend"
    link = "https://runbooks.example.com/checkout/frontend"
    # health_rollup defaults to AUTOMATIC_ROLL_UP: health rolls up from the step signals below.

    related {
      source = true
      target = false
    }

    levels {
      steps {
        name = "Login Page"

        # Entity signal: a specific New Relic entity pinned by GUID.
        signals {
          guid = "GUID1"
          name = "Cart Service"
          type = "ENTITY"
        }

        # Alert signal: an alert condition pinned by its entity GUID.
        signals {
          guid = "GUID2"
          name = "Checkout Error Rate"
          type = "ALERT"
        }

        # Dynamic query: API auto-discovers entities matching the filter at each refresh.
        entity_search_query {
          query = "accountId=1234 AND domain='BROWSER' AND name='Login'"
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
* `refresh_interval` - (Optional) How often the flow, stage, level, and step health statuses are refreshed. Defaults to `FIVE_MINUTES` if not set. Valid values: `ONE_MINUTE`, `FIVE_MINUTES`, `TEN_MINUTES`, `FIFTEEN_MINUTES`, `THIRTY_MINUTES`.
* `kpis` - (Optional) A list of Key Performance Indicators tracked at the flow level. See [Nested `kpis` blocks](#nested-kpis-blocks) below for details.
* `stages` - (Optional) An ordered list of stages that make up this flow. Maximum 50 stages. A flow can be created without stages and stages can be added later. See [Nested `stages` blocks](#nested-stages-blocks) below for details.

### Nested `kpis` blocks

KPIs are numeric metrics derived from NRQL queries, displayed as scorecards above the flow or stage. Flow-level KPIs are visible across all stages; stage-level KPIs (`stage_kpis`) are scoped to a single stage.

* `name` - (Required) The display name of the KPI.
* `description` - (Optional) A short description explaining what the KPI measures.
* `category` - (Optional) A label used to group related KPIs (e.g. `Revenue`, `Reliability`).
* `account_id` - (Optional) The account whose data this KPI queries. Defaults to the flow's account ID.
* `query` - (Required) The NRQL query definition for this KPI. See [Nested `query` blocks](#nested-query-blocks) below for details.
* `metric_query` - (Computed) The resolved NRQL metric query string synthesized by the API from the `query` block. Read-only.

-> **NOTE:** Cross-account KPIs — setting `account_id` to an account other than the flow's own account — are not yet supported in this Limited Preview.

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

* `query` - (Required) A New Relic entity filter expression, e.g. `accountId=1234 AND domain='APM' AND name='OrderService'`.
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

-> **NOTE:** `type = "ALERT"` is not yet supported in this Limited Preview.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `guid` - The entity GUID assigned to this Pathpoint flow in New Relic.
* `version` - The last-updated epoch-millisecond timestamp used for optimistic concurrency control. This is managed automatically and must not be modified. It is persisted to state after every create/update and sent to the API on each subsequent update.
* `stages.#.id` - The internal workload ID of the stage. Populated after creation and used to identify stages on updates.
* `stages.#.levels.#.id` - The internal workload ID of the level. Populated after creation and used to identify levels on updates.
* `stages.#.levels.#.steps.#.id` - The internal workload ID of the step. Populated after creation and used to identify steps on updates.
* `kpis.#.id` - The internal ID of the flow-level KPI. Populated after creation.
* `kpis.#.metric_query` - The resolved NRQL metric query string synthesized from the KPI's `query` block. This is a read-only computed value set by the API.
* `stages.#.stage_kpis.#.id` - The internal ID of the stage-level KPI. Populated after creation.

-> **NOTE:** On update, the provider matches each item in the new configuration to its existing counterpart by name first, falling back to position in the list. To avoid unexpected ID associations, avoid renaming and reordering items in the same `terraform apply`.

## Import

New Relic Pathpoint flows can be imported using the flow's entity GUID, e.g.

```bash
$ terraform import newrelic_pathpoint_flow.checkout GUID1
```
-> **NOTE:** After importing, run `terraform plan` to verify the state matches the existing configuration. The provider will read the current flow configuration from the API and populate all attributes in state.

## Exporting an Existing Flow

If a flow already exists in your account (created through the UI, or by another user), you can pull its Terraform configuration directly from the Pathpoint UI instead of hand-writing it:

1. Open the flow in the Pathpoint UI.
2. Click the overflow menu (**•••**) in the flow view
3. Select **View as code** to get the full `newrelic_pathpoint_flow` resource block for that pathpoint flow.

-> **TIP:** Treat the exported configuration as a starting point. Pair it with [`terraform import`](#import) using the flow's GUID so Terraform state matches the live resource before your next `apply` — otherwise Terraform will try to recreate what already exists.


