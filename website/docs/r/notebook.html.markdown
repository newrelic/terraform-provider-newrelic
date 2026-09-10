---
layout: "newrelic"
page_title: "New Relic: newrelic_notebook"
sidebar_current: "docs-newrelic-resource-notebook"
description: |-
  Create and manage New Relic Notebooks.
---

# Resource: newrelic_notebook

Use this resource to create and manage [New Relic Notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/introduction-notebooks/).

Notebooks let you combine live NRQL queries, visualizations, and markdown narrative into a single shareable document. This resource manages the full lifecycle of a notebook — its title and full block content — using the **declarative UI** format that the Notebooks platform uses natively.

## Content schema

Notebook content is expressed as a JSON document following the declarative UI envelope:

```json
{
  "type": "declarative",
  "version": 1,
  "content": [ <one or more container blocks> ]
}
```

Each **container** groups one or more visualization widgets with a shared layout:

```json
{
  "type": "container",
  "props": { "layout": "stack" },
  "content": [ <widget blocks> ]
}
```

Each **widget** wraps a single visualization:

```json
{
  "type": "widget",
  "props": {},
  "content": {
    "type": "visualization",
    "id": "<viz-id>",
    "props": { <visualization-specific props> }
  }
}
```

The `id` field identifies the visualization type. All supported visualization types and their props are documented in the [Nested widget blocks](#nested-widget-blocks) section below.

## Choosing a content mode

You must specify exactly one of `content` or `content_json`. They are mutually exclusive:

| Mode | When to use |
|---|---|
| `content` | Authoring notebooks directly in Terraform. Uses `jsonencode({})` for structured HCL that produces field-level diffs in `terraform plan`. |
| `content_json` | Working from JSON exported out of the New Relic UI or stored in a file. Produces line-level diffs of the normalized content. |

---

## Example Usage

<details>
  <summary>Minimal notebook - single markdown block</summary>

```hcl
resource "newrelic_notebook" "example" {
  title = "Incident Response Notes"

  content = jsonencode({
    type    = "declarative"
    version = 1
    content = [
      {
        type  = "container"
        props = { layout = "stack" }
        content = [
          {
            type  = "widget"
            props = {}
            content = {
              type = "visualization"
              id   = "viz.markdown"
              props = {
                text = "## Summary\n\nAdd investigation notes here."
              }
            }
          }
        ]
      }
    ]
  })
}
```

</details>

<details>
  <summary>Incident investigation runbook - markdown, billboard, line, area, table (content mode)</summary>

A realistic notebook for a production incident: a header, KPI billboard with thresholds, a timeline line chart, database latency area chart, drill-down table, and a root-cause markdown summary.

```hcl
resource "newrelic_notebook" "incident_runbook" {
  title = "Production API Incident - Investigation"

  content = jsonencode({
    type    = "declarative"
    version = 1
    content = [
      {
        type  = "container"
        props = { layout = "stack" }
        content = [
          {
            type  = "widget"
            props = {}
            content = {
              type = "visualization"
              id   = "viz.markdown"
              props = {
                text = "# Production API Incident\n\n**Status**: Resolved  \n**Duration**: 46 minutes  \n**Impact**: ~15% of API requests failed"
              }
            }
          },
          {
            type  = "widget"
            props = { title = "Peak error rate" }
            content = {
              type = "visualization"
              id   = "viz.billboard"
              props = {
                nrqlQueries = [
                  {
                    accountIds = [var.account_id]
                    query      = "SELECT percentage(count(*), WHERE httpResponseCode >= 400) AS 'Error Rate %' FROM Transaction WHERE appName = 'api-production' SINCE '2024-10-15 14:00:00' UNTIL '2024-10-15 16:00:00'"
                  }
                ]
                thresholdsWithSeriesOverrides = {
                  thresholds = [
                    { to = 1,           severity = "success"  }
                    { from = 1, to = 5, severity = "warning"  }
                    { from = 5,         severity = "critical" }
                  ]
                }
              }
            }
          },
          {
            type  = "widget"
            props = { title = "Error rate over time" }
            content = {
              type = "visualization"
              id   = "viz.line"
              props = {
                nrqlQueries = [
                  {
                    accountIds = [var.account_id]
                    query      = "SELECT percentage(count(*), WHERE httpResponseCode >= 400) FROM Transaction WHERE appName = 'api-production' TIMESERIES 1 minute SINCE '2024-10-15 14:00:00' UNTIL '2024-10-15 16:00:00'"
                  }
                ]
                legend = { enabled = true }
                yAxisLeft = { zero = true }
              }
            }
          },
          {
            type  = "widget"
            props = { title = "Database latency" }
            content = {
              type = "visualization"
              id   = "viz.area"
              props = {
                nrqlQueries = [
                  {
                    accountIds = [var.account_id]
                    query      = "SELECT average(duration) AS 'Avg (ms)', percentile(duration, 95) AS 'P95 (ms)' FROM DatabaseSample TIMESERIES 5 minutes SINCE '2024-10-15 14:00:00' UNTIL '2024-10-15 16:00:00'"
                  }
                ]
              }
            }
          },
          {
            type  = "widget"
            props = { title = "Top affected endpoints" }
            content = {
              type = "visualization"
              id   = "viz.table"
              props = {
                nrqlQueries = [
                  {
                    accountIds = [var.account_id]
                    query      = "SELECT count(*) AS 'Errors', average(duration)*1000 AS 'Avg ms' FROM Transaction WHERE httpResponseCode >= 400 FACET request.uri SINCE '2024-10-15 14:00:00' UNTIL '2024-10-15 16:00:00' ORDER BY count(*) DESC LIMIT 10"
                  }
                ]
                initialSorting = { name = "Errors", direction = "desc" }
              }
            }
          },
          {
            type  = "widget"
            props = {}
            content = {
              type = "visualization"
              id   = "viz.markdown"
              props = {
                text = "## Root Cause\n\n- Database connection pool exhausted under traffic spike\n\n## Action Items\n\n- [ ] Raise DB connection pool limit\n- [ ] Add rate limiting on `/api/users/register`\n- [ ] Alert on pool utilization > 80%"
              }
            }
          }
        ]
      }
    ]
  })
}
```

</details>

<details>
  <summary>Weekly service health review - billboard, line, area, stacked-bar, pie (content_json mode)</summary>

**`notebooks/weekly-health.json`**

```json
{
  "type": "declarative",
  "version": 1,
  "content": [
    {
      "type": "container",
      "props": { "layout": "stack" },
      "content": [
        {
          "type": "widget",
          "props": {},
          "content": {
            "type": "visualization",
            "id": "viz.markdown",
            "props": {
              "text": "# Weekly Service Health Review\n\n**Period**: Last 7 days vs. prior week"
            }
          }
        },
        {
          "type": "widget",
          "props": { "title": "Apdex (7 days)" },
          "content": {
            "type": "visualization",
            "id": "viz.billboard",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [1234567],
                  "query": "SELECT apdex(duration, 0.5) AS 'Apdex' FROM Transaction WHERE appName IN ('web-frontend', 'api-backend') SINCE 7 days ago"
                }
              ],
              "thresholdsWithSeriesOverrides": {
                "thresholds": [
                  { "from": 0.9, "severity": "success" },
                  { "from": 0.7, "to": 0.9, "severity": "warning" },
                  { "to": 0.7, "severity": "critical" }
                ]
              }
            }
          }
        },
        {
          "type": "widget",
          "props": { "title": "P95 latency (vs. prior week)" },
          "content": {
            "type": "visualization",
            "id": "viz.line",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [1234567],
                  "query": "SELECT percentile(duration, 95) FROM Transaction WHERE appName IN ('web-frontend', 'api-backend') FACET appName TIMESERIES 1 day SINCE 7 days ago COMPARE WITH 1 week ago"
                }
              ],
              "legend": { "enabled": true, "position": "bottom" }
            }
          }
        },
        {
          "type": "widget",
          "props": { "title": "Request volume" },
          "content": {
            "type": "visualization",
            "id": "viz.area",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [1234567],
                  "query": "SELECT count(*) FROM Transaction WHERE appName IN ('web-frontend', 'api-backend') FACET appName TIMESERIES 1 hour SINCE 7 days ago"
                }
              ]
            }
          }
        },
        {
          "type": "widget",
          "props": { "title": "Response codes over time" },
          "content": {
            "type": "visualization",
            "id": "viz.stacked-bar",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [1234567],
                  "query": "SELECT count(*) FROM Transaction WHERE appName IN ('web-frontend', 'api-backend') FACET httpResponseCode TIMESERIES 1 hour SINCE 7 days ago"
                }
              ],
              "legend": { "enabled": true }
            }
          }
        },
        {
          "type": "widget",
          "props": { "title": "Error distribution by status" },
          "content": {
            "type": "visualization",
            "id": "viz.pie",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [1234567],
                  "query": "SELECT count(*) FROM Transaction WHERE error IS true FACET httpResponseCode SINCE 7 days ago"
                }
              ],
              "facet": { "showOtherSeries": true }
            }
          }
        }
      ]
    }
  ]
}
```

**`main.tf`**

```hcl
resource "newrelic_notebook" "weekly_review" {
  title        = "Weekly Service Health Review"
  content_json = file("${path.module}/notebooks/weekly-health.json")
}
```

</details>

<details>
  <summary>Iterating to create multiple notebooks (for_each)</summary>

```hcl
variable "services" {
  type    = set(string)
  default = ["checkout", "payments", "inventory"]
}

resource "newrelic_notebook" "per_service" {
  for_each = var.services
  title    = "${each.value} runbook"

  content = jsonencode({
    type    = "declarative"
    version = 1
    content = [
      {
        type  = "container"
        props = { layout = "stack" }
        content = [
          {
            type  = "widget"
            props = {}
            content = {
              type  = "visualization"
              id    = "viz.markdown"
              props = { text = "# ${each.value}\n\nAdd runbook steps here." }
            }
          },
          {
            type  = "widget"
            props = { title = "Error rate - ${each.value}" }
            content = {
              type = "visualization"
              id   = "viz.billboard"
              props = {
                nrqlQueries = [{ accountIds = [var.account_id], query = "SELECT percentage(count(*), WHERE error IS true) AS 'Error %' FROM Transaction WHERE appName = '${each.value}' SINCE 1 hour ago" }]
              }
            }
          }
        ]
      }
    ]
  })
}
```

</details>

---

## Argument Reference

* `title` - (Required) The title of the notebook. Must be unique within the organization.
* `content` - (Optional) The notebook body as an HCL object using `jsonencode({...})`. Produces field-level diffs in `terraform plan`. Mutually exclusive with `content_json`.
* `content_json` - (Optional) The notebook body as a raw JSON string. Use when working from a UI export or a file. Produces line-level diffs of normalized content. Mutually exclusive with `content`.
* `organization_id` - (Computed) The New Relic organization ID. Resolved automatically from the provider credentials.

## Attributes Reference

* `guid` - The unique entity identifier (GUID) assigned to the notebook by New Relic.
* `blob_id` - The blob identifier of the current notebook content. Updated after every Terraform-managed write. Used internally to detect when content has changed.

---

## Nested `container` blocks

The `content` array at the top level contains one or more **container** blocks.

The following arguments are supported in a container block's `props`:

* `layout` - (Optional) Layout algorithm for the widgets. Accepted values: `"stack"` (stacks widgets vertically, default), `"grid"`. Defaults to `"stack"` when `props` is omitted.

---

## Nested `widget` blocks

Each widget block within a container's `content` array supports this structure:

```json
{
  "type": "widget",
  "props": { <widget-level props> },
  "content": {
    "type": "visualization",
    "id": "<viz-id>",
    "props": { <visualization-specific props> }
  }
}
```

### Widget-level `props`

Set directly on the `widget` object's `props`, shared across all widget types:

* `title` - (Optional) A label displayed above the rendered chart.
* `ignoreTimeRange` - (Optional, **widget-level prop**) This is equivalent to `platformOptions.ignoreTimeRange` on some viz types. Prefer setting it directly via the viz prop `platformOptions.ignoreTimeRange`.

---

## Supported visualization types

| `id` | Display name | Description | Required NRQL shape |
|---|---|---|---|
| `viz.line` | Line | Time series as lines | `TIMESERIES` recommended |
| `viz.area` | Area | Filled area under line(s) | `TIMESERIES` recommended |
| `viz.stacked-bar` | Bar | Stacked vertical bars (TIMESERIES or FACET) | Both modes |
| `viz.bar` | Simplified Bar | Simple horizontal bar (FACET only) | `FACET` recommended |
| `viz.stacked-horizontal-bar` | Horizontal Stacked Bar | Horizontal bars stacked by FACET | `FACET` recommended |
| `viz.pie` | Pie | Proportional segments | `FACET` recommended |
| `viz.table` | Table | Data rows and columns | `FACET` recommended |
| `viz.billboard` | Billboard | Single large metric with threshold coloring | Single aggregation |
| `viz.gauge` | Gauge | Dial or bar showing value vs scale | Single aggregation |
| `viz.histogram` | Histogram | Distribution buckets | `histogram()` function |
| `viz.heatmap` | Heatmap | 2D distribution grid | `histogram()` with `FACET` |
| `viz.scatter` | Scatter | X-Y scatter plot | Two-value `FACET` result |
| `viz.apdex` | Apdex | Apdex score display | Apdex-compatible |
| `viz.bullet` | Bullet | Progress bar vs target | Single aggregation |
| `viz.funnel` | Funnel | Conversion funnel | `funnel()` function |
| `viz.event-feed` | Event Feed | Scrolling event list | Any |
| `viz.json` | JSON | Raw JSON output | Any |
| `viz.sparkline` | Sparkline | Compact trend line with axes | `TIMESERIES` |
| `viz.sparkline-lite` | Sparkline Lite | Ultra-compact sparkline, no axes | `TIMESERIES` |
| `viz.markdown` | Markdown | Rich text block | N/A |

---

### `viz.line` — Line chart

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
* `alertQueries` - (Optional) Array of alert violation query objects. Fetches alert violation data and renders warning/critical bands as event bars beneath the line. Only available for `viz.line`. Each object requires: `accountIds` (number[]), `violationId` (number), `duration` (number, ms), `endTime` (number, epoch ms).
* `sqlQueries` - (Optional) Array of SQL query objects via Federated Data Source (FDS). Each object requires: `query` (string), and `accountId` (number) or `accountIds` (number[]).
* `legend.enabled` - (Optional) Show or hide the chart legend. Default: `true`.
* `legend.position` - (Optional) Position of the legend. Valid values: `"bottom"` (default), `"left"`, `"right"`, `null`.
* `facet.showOtherSeries` - (Optional) Show the aggregated "Other" group for `FACET` queries. Default: `false`.
* `yAxisLeft.min` - (Optional) Minimum value for the left Y axis.
* `yAxisLeft.max` - (Optional) Maximum value for the left Y axis.
* `yAxisLeft.zero` - (Optional) Force zero as the axis origin. Default: `true`.
* `yAxisLeft.scale` - (Optional) Scale type for the left Y axis. Valid values: `"linear"` (default), `"logarithmic"`.
* `yAxisRight.min` - (Optional) Minimum value for the right Y axis.
* `yAxisRight.max` - (Optional) Maximum value for the right Y axis.
* `yAxisRight.zero` - (Optional) Force zero on the right Y axis. Default: `true`.
* `yAxisRight.scale` - (Optional) Scale type. Valid values: `"linear"` (default), `"logarithmic"`.
* `yAxisRight.series[].name` - (Optional) Series names to plot against the right Y axis instead of the left. A series can only belong to one axis.
* `nullValues.nullValue` - (Optional) How to handle null data points globally. Valid values: `"default"`, `"zero"`, `"preserve"`, `"remove"`.
* `nullValues.seriesOverrides[].seriesName` - (Optional) Series name to apply the override to.
* `nullValues.seriesOverrides[].nullValue` - (Optional) Per-series null value strategy.
* `colors.colorPalette` - (Optional) Color palette. Valid values: `"consistent"` (default), `"dynamic"`, `null`.
* `colors.seriesOverrides[].seriesName` - (Optional) Series name to apply a custom color to.
* `colors.seriesOverrides[].color` - (Optional) Color for the series (RGB hex, e.g. `"#FF0000"`).
* `units.unit` - (Optional) Default data unit for all series. See [Valid `units.unit` values](#valid-unitsunit-values).
* `units.seriesOverrides[].seriesName` - (Optional) Series name to apply a unit to.
* `units.seriesOverrides[].unit` - (Optional) Unit for the series.
* `thresholds.isLabelVisible` - (Optional) Show threshold labels on the chart. Default: `false`.
* `thresholds.thresholds[].name` - (Optional) Label for this threshold.
* `thresholds.thresholds[].from` - (Optional) Lower bound of the threshold range.
* `thresholds.thresholds[].to` - (Optional) Upper bound of the threshold range.
* `thresholds.thresholds[].severity` - (Optional) Color for data in this range. Valid values: `"critical"`, `"severe"`, `"warning"`, `"success"`, `"unavailable"`.
* `chartStyles.lineInterpolation` - (Optional) How data points are connected. Valid values: `"linear"` (default), `"smooth"`, `"stepBefore"`, `"stepAfter"`.
* `chartStyles.gradient.enabled` - (Optional) Enable gradient fill. Default: `false`.
* `chartTypes.seriesOverrides[].seriesName` - (Optional) Series name to render as a different chart type.
* `chartTypes.seriesOverrides[].chartType` - (Optional) Override chart type for a specific series. Valid values: `"line"`, `"area"`.
* `tooltip.mode` - (Optional) Tooltip display behaviour. Valid values: `"single"` (default), `"all"`, `"hidden"`.
* `platformOptions.ignoreTimeRange` - (Optional) Override the notebook time picker with the query's own time range. Default: `false`.
* `refreshRate.frequency` - (Optional) Auto-refresh interval in milliseconds, or `"auto"`. Default: `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Transaction throughput" },
  "content": {
    "type": "visualization",
    "id": "viz.line",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT rate(count(*), 1 minute) FROM Transaction FACET appName TIMESERIES 5 minutes SINCE 3 hours ago" }
      ],
      "legend": { "enabled": true, "position": "bottom" },
      "yAxisLeft": { "zero": true, "min": 0 },
      "nullValues": { "nullValue": "zero" },
      "chartStyles": { "lineInterpolation": "smooth" },
      "thresholds": {
        "thresholds": [
          { "name": "Warning", "from": 1000, "to": 2000, "severity": "warning" },
          { "name": "Critical", "from": 2000, "severity": "critical" }
        ]
      }
    }
  }
}
```

---

### `viz.area` — Area chart

**Props (`content.props`):**

Same as `viz.line` except:
- `alertQueries` is not supported
- `yAxisRight` is not supported
- `nullValues` supports: `"default"`, `"zero"`, `"preserve"` (no `"remove"`)
- Additional prop: `chartStyles.stacked.enabled` - (Optional) Enable stacked mode. Default: `true`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Memory utilisation" },
  "content": {
    "type": "visualization",
    "id": "viz.area",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT average(memoryUsedPercent) FROM SystemSample FACET hostname TIMESERIES 5 minutes SINCE 3 hours ago" }
      ],
      "yAxisLeft": { "zero": true, "min": 0, "max": 100 },
      "chartStyles": { "stacked": { "enabled": false } }
    }
  }
}
```

---

### `viz.stacked-bar` — Bar chart (stacked/timeseries)

Vertical bar chart supporting both TIMESERIES and categorical FACET data. Stacked layout with optional thresholds.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `sqlQueries` - (Optional) Array of SQL query objects.
* `legend.enabled` - (Optional) Default: `true`.
* `legend.position` - (Optional) `"bottom"` | `"left"` | `"right"` | `null`. Default: `"bottom"`.
* `facet.showOtherSeries` - (Optional) Default: `false`.
* `colors.colorPalette` - (Optional) `"consistent"` | `"dynamic"` | `null`. Default: `"consistent"`.
* `colors.seriesOverrides[].seriesName` / `.color`
* `units.unit` / `units.seriesOverrides[]`
* `thresholds.isLabelVisible` / `thresholds.thresholds[].name/from/to/severity`
* `chartStyles.gradient.enabled` - (Optional) Default: `false`.
* `chartStyles.stacked.enabled` - (Optional) Default: `true`.
* `chartTypes.seriesOverrides[].seriesName` / `.chartType` — Valid: `"line"`, `"area"`.
* `nullValues.nullValue` / `nullValues.seriesOverrides[]` — Values: `"default"`, `"zero"`, `"preserve"`.
* `tooltip.mode` - (Optional) `"single"` | `"all"` | `"hidden"`. Default: `"single"`.
* `yAxisLeft.min` / `.max` / `.zero` / `.scale`
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Response codes over time" },
  "content": {
    "type": "visualization",
    "id": "viz.stacked-bar",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT count(*) FROM Transaction FACET httpResponseCode TIMESERIES 5 minutes SINCE 3 hours ago" }
      ],
      "legend": { "enabled": true }
    }
  }
}
```

---

### `viz.bar` — Simplified Bar chart

Simple horizontal bar chart for categorical FACET comparisons. Fewer props than `viz.stacked-bar`.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `sqlQueries` - (Optional) Array of SQL query objects.
* `facet.showOtherSeries` - (Optional) Default: `false`.
* `colors.colorPalette` - (Optional) Default: `"consistent"`.
* `colors.seriesOverrides[].seriesName` / `.color`
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Top error transactions" },
  "content": {
    "type": "visualization",
    "id": "viz.bar",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT count(*) FROM Transaction WHERE error IS true FACET name SINCE 1 hour ago ORDER BY count(*) DESC LIMIT 10" }
      ]
    }
  }
}
```

---

### `viz.stacked-horizontal-bar` — Horizontal Stacked Bar

Horizontal bars stacked by FACET. Minimal configuration.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `facet.showOtherSeries` - (Optional) Default: `false`.
* `legend.enabled` - (Optional) Default: `true`.
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Errors by service" },
  "content": {
    "type": "visualization",
    "id": "viz.stacked-horizontal-bar",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT count(*) FROM Transaction WHERE error IS true FACET appName SINCE 7 days ago" }
      ]
    }
  }
}
```

---

### `viz.pie` — Pie chart

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `sqlQueries` - (Optional) Array of SQL query objects.
* `facet.showOtherSeries` - (Optional) Default: **`true`** (pie default differs from other viz types).
* `legend.enabled` - (Optional) Default: `true`.
* `colors.colorPalette` - (Optional) Default: `"consistent"`.
* `colors.seriesOverrides[].seriesName` / `.color`
* `chartStyles.gradient.enabled` - (Optional) Default: `false`.
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Traffic by HTTP status" },
  "content": {
    "type": "visualization",
    "id": "viz.pie",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT count(*) FROM Transaction FACET httpResponseCode SINCE 1 hour ago" }
      ],
      "facet": { "showOtherSeries": true }
    }
  }
}
```

---

### `viz.table` — Table

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `sqlQueries` - (Optional) Array of SQL query objects.
* `facet.showOtherSeries` - (Optional) Default: `false`.
* `initialSorting` - (Optional) Default sort configuration. Object with: `name` (string, column name) and `direction` (`"asc"` | `"desc"`).
* `hiddenColumns[].columnName` - (Optional) Columns to hide from the rendered table. Matches either the raw column name from the query (e.g. `"timestamp"`) or the display label (e.g. `"Timestamp"`).
* `thresholds[].columnName` - (Optional) Column to apply the threshold to.
* `thresholds[].from` - (Optional) Lower bound.
* `thresholds[].to` - (Optional) Upper bound.
* `thresholds[].severity` - (Optional) `"critical"` | `"severe"` | `"warning"` | `"success"` | `"unavailable"`.
* `dataFormatters` - (Optional) Custom data format configuration (see the vizco documentation for format details).
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Slowest transactions" },
  "content": {
    "type": "visualization",
    "id": "viz.table",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT average(duration)*1000 AS 'Avg ms', max(duration)*1000 AS 'Max ms', count(*) AS 'Calls' FROM Transaction FACET name SINCE 1 hour ago" }
      ],
      "initialSorting": { "name": "Avg ms", "direction": "desc" },
      "hiddenColumns": [{ "columnName": "Calls" }],
      "thresholds": [
        { "columnName": "Avg ms", "from": 500, "severity": "warning" },
        { "columnName": "Avg ms", "from": 1000, "severity": "critical" }
      ]
    }
  }
}
```

---

### `viz.billboard` — Billboard

Single large metric value with color-coded threshold ranges.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `sqlQueries` - (Optional) Array of SQL query objects.
* `thresholdsWithSeriesOverrides.thresholds[].from` - (Optional) Lower bound of threshold range.
* `thresholdsWithSeriesOverrides.thresholds[].to` - (Optional) Upper bound of threshold range.
* `thresholdsWithSeriesOverrides.thresholds[].severity` - (Optional) `"critical"` | `"severe"` | `"warning"` | `"success"` | `"unavailable"`.
* `thresholdsWithSeriesOverrides.seriesOverrides[].seriesName` - (Optional) Series name to apply custom thresholds.
* `thresholdsWithSeriesOverrides.seriesOverrides[].from` / `.to` / `.severity` - (Optional) Per-series threshold.
* `billboardSettings.visual.alignment` - (Optional) Layout of value and label. Valid values: `"auto"`, `"stacked"`, `"inline"`.
* `billboardSettings.visual.display` - (Optional) What to show. Valid values: `"auto"`, `"all"`, `"value"`, `"label"`, `"none"`.
* `billboardSettings.gridOptions.columns` - (Optional) Number of columns in the grid (for multi-value billboards).
* `billboardSettings.gridOptions.label` - (Optional) Font size of the label in pixels.
* `billboardSettings.gridOptions.value` - (Optional) Font size of the value in pixels.
* `billboardSettings.link.url` - (Optional) URL to navigate to when the billboard is clicked.
* `billboardSettings.link.title` - (Optional) Link text.
* `billboardSettings.link.newTab` - (Optional) Open link in a new tab. Default: `false`.
* `chartStyles.lineInterpolation` - (Optional) Line interpolation if the billboard shows a trend sparkline.
* `facet.showOtherSeries` - (Optional) Default: `false`.
* `units.unit` / `units.seriesOverrides[]`
* `dataFormatters` - (Optional) Custom data format configuration.
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Error rate" },
  "content": {
    "type": "visualization",
    "id": "viz.billboard",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT percentage(count(*), WHERE error IS true) AS 'Error %' FROM Transaction SINCE 1 hour ago" }
      ],
      "thresholdsWithSeriesOverrides": {
        "thresholds": [
          { "to": 1,            "severity": "success"  },
          { "from": 1, "to": 5, "severity": "warning"  },
          { "from": 5,          "severity": "critical" }
        ]
      },
      "billboardSettings": {
        "visual": { "alignment": "inline", "display": "auto" }
      }
    }
  }
}
```

---

### `viz.gauge` — Gauge

Dial, ring, or bar gauge with threshold coloring.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `sqlQueries` - (Optional) Array of SQL query objects.
* `gaugeSettings.displayMode` - (Optional) Gauge style. Valid values: `"arc"` (270° dial, default), `"circular"` (360° ring), `"bar"` (horizontal bar).
* `gaugeSettings.min` - (Optional) Minimum value of the gauge scale.
* `gaugeSettings.max` - (Optional) Maximum value of the gauge scale. Default: `100`.
* `gaugeSettings.display` - (Optional) What to show inside the gauge. Valid values: `"auto"` (default), `"all"`, `"value"`, `"name"`, `"none"`.
* `gaugeSettings.showThresholdMarkers` - (Optional) Show the threshold color ring/strip. Default: `true`.
* `gaugeSettings.showThresholdLabels` - (Optional) Show value labels at threshold boundaries. Default: `false`.
* `gaugeSettings.thresholdsColorMode` - (Optional) How threshold zones are colored. Valid values: `"solid"` (default), `"gradient"`.
* `thresholds.thresholds[].from` - (Optional) Lower bound.
* `thresholds.thresholds[].severity` - (Optional) `"critical"` | `"severe"` | `"warning"` | `"success"` | `"unavailable"`.
* `thresholds.seriesOverrides[].seriesName` - (Optional) Per-series threshold overrides.
* `colors.colorPalette` - (Optional) Default: `"consistent"`.
* `units.unit` / `units.seriesOverrides[]`
* `dataFormatters` - (Optional) Custom data format configuration.
* `facet.showOtherSeries` - (Optional) Default: `false`.
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "CPU utilisation" },
  "content": {
    "type": "visualization",
    "id": "viz.gauge",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT average(cpuPercent) FROM SystemSample SINCE 5 minutes ago" }
      ],
      "gaugeSettings": {
        "displayMode": "arc",
        "min": 0,
        "max": 100,
        "showThresholdMarkers": true,
        "thresholdsColorMode": "solid"
      },
      "thresholds": {
        "thresholds": [
          { "from": 80, "severity": "warning" },
          { "from": 90, "severity": "critical" }
        ]
      }
    }
  }
}
```

---

### `viz.histogram` — Histogram

Distribution of values across buckets. Requires the `histogram()` NRQL function.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects. Query must use `histogram(attr)` or `histogram(attr, width: N, buckets: N)`.
* `legend.enabled` - (Optional) Default: `true`.
* `legend.position` - (Optional) Default: `"bottom"`.
* `colors.colorPalette` - (Optional) Default: `"consistent"`.
* `colors.seriesOverrides[].seriesName` / `.color`
* `chartStyles.gradient.enabled` - (Optional) Default: `false`.
* `thresholds.isLabelVisible` / `thresholds.thresholds[].name/from/to/severity`
* `yAxisLeft.min` / `.max` / `.zero` / `.scale`
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Response duration distribution" },
  "content": {
    "type": "visualization",
    "id": "viz.histogram",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT histogram(duration, width: 0.1, buckets: 20) FROM Transaction SINCE 1 hour ago" }
      ]
    }
  }
}
```

---

### `viz.heatmap` — Heatmap

2D distribution grid, showing a `histogram()` broken down by `FACET`.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects. Query must use `histogram()` with `FACET`.
* `facet.showOtherSeries` - (Optional) Default: `false`.
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Duration distribution by app" },
  "content": {
    "type": "visualization",
    "id": "viz.heatmap",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT histogram(duration) FROM Transaction FACET appName SINCE 1 hour ago" }
      ]
    }
  }
}
```

---

### `viz.scatter` — Scatter plot

X-Y scatter plot. Query should return two numeric values and a `FACET`.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `facet.showOtherSeries` - (Optional) Default: `false`.
* `legend.enabled` - (Optional) Default: `true`.
* `legend.position` - (Optional) Default: `"bottom"`.
* `colors.colorPalette` - (Optional) Default: `"consistent"`.
* `colors.seriesOverrides[].seriesName` / `.color`
* `units.unit` / `units.seriesOverrides[]`
* `nullValues.nullValue` / `nullValues.seriesOverrides[]` — Values: `"default"`, `"zero"`, `"preserve"`, `"remove"`.
* `yAxisLeft.min` / `.max` / `.zero` / `.scale`
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Apdex vs avg duration" },
  "content": {
    "type": "visualization",
    "id": "viz.scatter",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT apdex(duration, t: 0.5), average(duration) FROM Transaction FACET appName LIMIT 20 SINCE 1 hour ago" }
      ]
    }
  }
}
```

---

### `viz.apdex` — Apdex

Dedicated Apdex score display with legend.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `legend.enabled` - (Optional) Default: `true`.
* `legend.position` - (Optional) Default: `"bottom"`.
* `tooltip.mode` - (Optional) `"single"` | `"all"` | `"hidden"`. Default: `"single"`.
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Overall Apdex" },
  "content": {
    "type": "visualization",
    "id": "viz.apdex",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT apdex(duration, 0.5) FROM Transaction WHERE appName IN ('web-frontend', 'api-backend') TIMESERIES 1 day SINCE 7 days ago" }
      ]
    }
  }
}
```

---

### `viz.bullet` — Bullet chart

Horizontal progress bar showing a value against a target.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects. Should return a single aggregated value.
* `sqlQueries` - (Optional) Array of SQL query objects.
* `limit` - (Required) The target value shown as the bullet goal line.
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Requests/min vs target" },
  "content": {
    "type": "visualization",
    "id": "viz.bullet",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT rate(count(*), 1 minute) AS 'Requests/min' FROM Transaction SINCE 5 minutes ago" }
      ],
      "limit": 5000
    }
  }
}
```

---

### `viz.funnel` — Funnel chart

Conversion funnel across sequential steps. Requires the `funnel()` NRQL function.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects. Query must use the `funnel()` function.
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "User conversion" },
  "content": {
    "type": "visualization",
    "id": "viz.funnel",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT funnel(session, WHERE name = 'HomeView' AS 'Home', WHERE name = 'ProductView' AS 'Product', WHERE name = 'ConfirmationView' AS 'Purchase') FROM PageView SINCE 1 day ago" }
      ]
    }
  }
}
```

---

### `viz.event-feed` — Event Feed

Scrolling list of events returned by the NRQL query.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Recent errors" },
  "content": {
    "type": "visualization",
    "id": "viz.event-feed",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT message, error.class FROM Transaction WHERE error IS true SINCE 1 hour ago LIMIT 50" }
      ]
    }
  }
}
```

---

### `viz.json` — JSON

Renders the raw JSON output of the NRQL query. Useful for debugging.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Raw query output" },
  "content": {
    "type": "visualization",
    "id": "viz.json",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT * FROM Transaction LIMIT 5 SINCE 5 minutes ago" }
      ]
    }
  }
}
```

---

### `viz.sparkline` — Sparkline

Compact time-series line chart with Y-axis labels and full prop support.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects with `TIMESERIES`.
* `facet.showOtherSeries` - (Optional) Default: `false`.
* `legend.enabled` / `legend.position` (implicitly via sparkline display)
* `colors.colorPalette` / `colors.seriesOverrides[]`
* `units.unit` / `units.seriesOverrides[]`
* `nullValues.nullValue` / `nullValues.seriesOverrides[]` — Values: `"default"`, `"zero"`, `"preserve"`, `"remove"`.
* `chartStyles.lineInterpolation` - (Optional) `"linear"` | `"smooth"` | `"stepBefore"` | `"stepAfter"`. Default: `"linear"`.
* `tooltip.mode` - (Optional) `"single"` | `"all"` | `"hidden"`. Default: `"single"`.
* `yAxisLeft.min` / `.max` / `.zero`
* `platformOptions.ignoreTimeRange` - (Optional) Default: `false`.
* `refreshRate.frequency` - (Optional) Default: `"auto"`.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Throughput trend" },
  "content": {
    "type": "visualization",
    "id": "viz.sparkline",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT count(*) FROM Transaction TIMESERIES AUTO SINCE 1 hour ago" }
      ]
    }
  }
}
```

---

### `viz.sparkline-lite` — Sparkline Lite

Ultra-compact sparkline with no axis labels. Ideal as an inline trend indicator.

**Props (`content.props`):**

Same as `viz.sparkline` except `yAxisLeft` is not supported.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Throughput (mini)" },
  "content": {
    "type": "visualization",
    "id": "viz.sparkline-lite",
    "props": {
      "nrqlQueries": [
        { "accountIds": [1234567], "query": "SELECT count(*) FROM Transaction TIMESERIES AUTO SINCE 1 hour ago" }
      ]
    }
  }
}
```

---

### `viz.markdown` — Markdown text block

Renders Markdown-formatted text.

**Props (`content.props`):**

* `text` - (Required) Markdown source string, or an array of strings (each rendered on a new line). Supports headings, bold, italic, inline code, fenced code blocks, ordered and unordered lists, task lists (`- [ ]`), tables, and links.

**Example:**
```json
{
  "type": "widget",
  "props": {},
  "content": {
    "type": "visualization",
    "id": "viz.markdown",
    "props": {
      "text": "## Investigation Notes\n\n- [ ] Check alert history\n- [ ] Review recent deploys\n\n| Metric | Target |\n|---|---|\n| Error rate | < 1% |\n| P95 latency | < 500 ms |"
    }
  }
}
```

---

## Nested `nrqlQueries` blocks

All query-based visualizations accept a `nrqlQueries` array. Each item supports:

* `accountIds` - (Required) An array of one or more New Relic account IDs, e.g. `[1234567]`. For cross-account queries: `[1234567, 7654321]`.
* `query` - (Required) A valid NRQL query string. See [Introduction to NRQL](https://docs.newrelic.com/docs/nrql/get-started/introduction-nrql-new-relics-query-language/).
* `offset` - (Optional) Number of milliseconds to offset the query time window.

Most viz types also accept `sqlQueries` as an alternative data source for Federated Data Source (FDS) connections. Each SQL query object requires `query` (string) and `accountId` (number) or `accountIds` (number[]).

-> **NOTE:** If you query an account to which you do not have access, the widget will render with a "data inaccessible" message. No Terraform error is raised.

---

## Valid `units.unit` values

The following unit identifiers are accepted for `units.unit` and `units.seriesOverrides[].unit`:

`APDEX`, `BITS`, `BITS_PER_MS`, `BITS_PER_SECOND`, `BYTES`, `BYTES_PER_MS`, `BYTES_PER_SECOND`, `CELSIUS`, `COUNT`, `DOLLAR`, `HERTZ`, `MS`, `PAGES_PER_SECOND`, `PERCENTAGE`, `REQUESTS_PER_SECOND`, `REQUESTS_PER_MINUTE`, `SECONDS`, `TIMESTAMP`

---

## Valid `refreshRate.frequency` values

| Value | Refresh interval |
|---|---|
| `"auto"` | Platform default (based on query) |
| `0` | No automatic refresh |
| `5000` | 5 seconds |
| `30000` | 30 seconds |
| `60000` | 1 minute |
| `300000` | 5 minutes |
| `1800000` | 30 minutes |
| `3600000` | 1 hour |
| `10800000` | 3 hours |
| `43200000` | 12 hours |
| `86400000` | 24 hours |

---

## Further reading

- [Visualizations in notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/visualizations-in-notebooks/) — chart types available
- [Blob Storage API for notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/blob-storage-api-for-notebooks/) — the underlying REST API
- [Notebook examples](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/notebooks-examples/) — worked multi-block examples
- [NR1 SDK: Chart components](https://docs.newrelic.com/docs/new-relic-solutions/build-nr-ui/sdk-component/charts/Charts/) — SDK chart components notebooks viz IDs map to
- [Introduction to NRQL](https://docs.newrelic.com/docs/nrql/get-started/introduction-nrql-new-relics-query-language/) — writing NRQL queries

---

## Import

```
# Default - imports into content_json
$ terraform import newrelic_notebook.example <guid>
$ terraform import newrelic_notebook.example <guid>:content_json

# Import into content field
$ terraform import newrelic_notebook.example <guid>:content
```

After importing, run `terraform plan`. If the imported state and your config use the same mode, the plan will show no changes.

## Plan Diff Behavior

Both content fields store a normalized form of the JSON in state (alphabetically sorted keys, 2-space indentation). This means:

- Reformatting your HCL or JSON file without changing any values produces **no diff** in `terraform plan`.
- Changing a single widget property shows **only that property** as changed.
- Externally modifying the notebook in the UI causes the changed fields to surface **precisely** in the next `terraform plan`.
