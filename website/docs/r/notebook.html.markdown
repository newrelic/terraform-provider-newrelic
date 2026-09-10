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

A realistic notebook for walking through a production incident: a header, a KPI billboard with thresholds, a timeline line chart, a database latency area chart, a drill-down table of affected endpoints, and a markdown root-cause summary.

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
                legendEnabled = true
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
  <summary>Weekly service health review - billboard, line, area, bar, pie (content_json mode)</summary>

A weekly review notebook stored as a JSON file, combining an Apdex billboard, P95 latency comparison, request volume, error breakdown, and a markdown analysis section.

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
          "props": { "title": "P95 latency by service (vs. prior week)" },
          "content": {
            "type": "visualization",
            "id": "viz.line",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [1234567],
                  "query": "SELECT percentile(duration, 95) AS 'P95 (ms)' FROM Transaction WHERE appName IN ('web-frontend', 'api-backend') FACET appName TIMESERIES 1 day SINCE 7 days ago COMPARE WITH 1 week ago"
                }
              ],
              "legendEnabled": true
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
          "props": { "title": "Errors by service" },
          "content": {
            "type": "visualization",
            "id": "viz.bar",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [1234567],
                  "query": "SELECT count(*) FROM Transaction WHERE error IS true AND appName IN ('web-frontend', 'api-backend') FACET appName SINCE 7 days ago"
                }
              ]
            }
          }
        },
        {
          "type": "widget",
          "props": { "title": "Error distribution by HTTP status" },
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
              "facetShowOtherSeries": true
            }
          }
        },
        {
          "type": "widget",
          "props": {},
          "content": {
            "type": "visualization",
            "id": "viz.markdown",
            "props": {
              "text": "## Analysis\n\n### Wins\n- API P95 down 35% after DB index added Monday\n\n### Concerns\n- Frontend P95 up 15% WoW — investigate asset bundle regression"
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
  <summary>All visualization types - one notebook per viz type</summary>

```hcl
# viz.line
resource "newrelic_notebook" "line_example" {
  title = "Line Chart Example"
  content = jsonencode({
    type = "declarative", version = 1
    content = [{
      type = "container", props = { layout = "stack" }
      content = [{
        type = "widget", props = { title = "Transaction throughput" }
        content = { type = "visualization", id = "viz.line"
          props = {
            nrqlQueries     = [{ accountIds = [var.account_id], query = "SELECT rate(count(*), 1 minute) FROM Transaction FACET appName TIMESERIES 5 minutes SINCE 3 hours ago" }]
            legendEnabled   = true
            yAxisLeft       = { zero = true, min = 0 }
            nullValues      = { nullValue = "zero" }
            facetShowOtherSeries = true
          }
        }
      }]
    }]
  })
}

# viz.area
resource "newrelic_notebook" "area_example" {
  title = "Area Chart Example"
  content = jsonencode({
    type = "declarative", version = 1
    content = [{
      type = "container", props = { layout = "stack" }
      content = [{
        type = "widget", props = { title = "Memory usage" }
        content = { type = "visualization", id = "viz.area"
          props = {
            nrqlQueries   = [{ accountIds = [var.account_id], query = "SELECT average(memoryUsedPercent) FROM SystemSample FACET hostname TIMESERIES 5 minutes SINCE 3 hours ago" }]
            legendEnabled = true
          }
        }
      }]
    }]
  })
}

# viz.bar
resource "newrelic_notebook" "bar_example" {
  title = "Bar Chart Example"
  content = jsonencode({
    type = "declarative", version = 1
    content = [{
      type = "container", props = { layout = "stack" }
      content = [{
        type = "widget", props = { title = "Top transaction errors" }
        content = { type = "visualization", id = "viz.bar"
          props = {
            nrqlQueries          = [{ accountIds = [var.account_id], query = "SELECT count(*) FROM Transaction WHERE error IS true FACET name SINCE 1 hour ago ORDER BY count(*) DESC LIMIT 10" }]
            facetShowOtherSeries = false
          }
        }
      }]
    }]
  })
}

# viz.pie
resource "newrelic_notebook" "pie_example" {
  title = "Pie Chart Example"
  content = jsonencode({
    type = "declarative", version = 1
    content = [{
      type = "container", props = { layout = "stack" }
      content = [{
        type = "widget", props = { title = "Requests by HTTP status" }
        content = { type = "visualization", id = "viz.pie"
          props = {
            nrqlQueries          = [{ accountIds = [var.account_id], query = "SELECT count(*) FROM Transaction FACET httpResponseCode SINCE 1 hour ago" }]
            facetShowOtherSeries = true
          }
        }
      }]
    }]
  })
}

# viz.table
resource "newrelic_notebook" "table_example" {
  title = "Table Example"
  content = jsonencode({
    type = "declarative", version = 1
    content = [{
      type = "container", props = { layout = "stack" }
      content = [{
        type = "widget", props = { title = "Slowest transactions" }
        content = { type = "visualization", id = "viz.table"
          props = {
            nrqlQueries = [{ accountIds = [var.account_id], query = "SELECT average(duration)*1000 AS 'Avg ms', max(duration)*1000 AS 'Max ms', count(*) AS 'Calls' FROM Transaction FACET name SINCE 1 hour ago ORDER BY average(duration) DESC LIMIT 20" }]
          }
        }
      }]
    }]
  })
}

# viz.billboard
resource "newrelic_notebook" "billboard_example" {
  title = "Billboard Example"
  content = jsonencode({
    type = "declarative", version = 1
    content = [{
      type = "container", props = { layout = "stack" }
      content = [{
        type = "widget", props = { title = "Error rate" }
        content = { type = "visualization", id = "viz.billboard"
          props = {
            nrqlQueries = [{ accountIds = [var.account_id], query = "SELECT percentage(count(*), WHERE error IS true) AS 'Error %' FROM Transaction SINCE 1 hour ago" }]
            thresholdsWithSeriesOverrides = {
              thresholds = [
                { to = 1,           severity = "success"  }
                { from = 1, to = 5, severity = "warning"  }
                { from = 5,         severity = "critical" }
              ]
            }
          }
        }
      }]
    }]
  })
}

# viz.stacked-bar
resource "newrelic_notebook" "stacked_bar_example" {
  title = "Stacked Bar Example"
  content = jsonencode({
    type = "declarative", version = 1
    content = [{
      type = "container", props = { layout = "stack" }
      content = [{
        type = "widget", props = { title = "Response codes over time" }
        content = { type = "visualization", id = "viz.stacked-bar"
          props = {
            nrqlQueries = [{ accountIds = [var.account_id], query = "SELECT count(*) FROM Transaction FACET httpResponseCode TIMESERIES 5 minutes SINCE 3 hours ago" }]
            legendEnabled = true
          }
        }
      }]
    }]
  })
}

# viz.histogram
resource "newrelic_notebook" "histogram_example" {
  title = "Histogram Example"
  content = jsonencode({
    type = "declarative", version = 1
    content = [{
      type = "container", props = { layout = "stack" }
      content = [{
        type = "widget", props = { title = "Duration distribution" }
        content = { type = "visualization", id = "viz.histogram"
          props = {
            nrqlQueries = [{ accountIds = [var.account_id], query = "SELECT histogram(duration, width: 0.1, buckets: 20) FROM Transaction SINCE 1 hour ago" }]
          }
        }
      }]
    }]
  })
}

# viz.heatmap
resource "newrelic_notebook" "heatmap_example" {
  title = "Heatmap Example"
  content = jsonencode({
    type = "declarative", version = 1
    content = [{
      type = "container", props = { layout = "stack" }
      content = [{
        type = "widget", props = { title = "Duration heatmap by app" }
        content = { type = "visualization", id = "viz.heatmap"
          props = {
            nrqlQueries = [{ accountIds = [var.account_id], query = "SELECT histogram(duration) FROM Transaction FACET appName SINCE 1 hour ago" }]
          }
        }
      }]
    }]
  })
}

# viz.funnel
resource "newrelic_notebook" "funnel_example" {
  title = "Funnel Example"
  content = jsonencode({
    type = "declarative", version = 1
    content = [{
      type = "container", props = { layout = "stack" }
      content = [{
        type = "widget", props = { title = "Conversion funnel" }
        content = { type = "visualization", id = "viz.funnel"
          props = {
            nrqlQueries = [{ accountIds = [var.account_id], query = "SELECT funnel(session, WHERE name = 'HomeView' AS 'Home', WHERE name = 'CheckoutView' AS 'Checkout', WHERE name = 'ConfirmationView' AS 'Purchase') FROM PageView SINCE 1 day ago" }]
          }
        }
      }]
    }]
  })
}

# viz.scatter
resource "newrelic_notebook" "scatter_example" {
  title = "Scatter Plot Example"
  content = jsonencode({
    type = "declarative", version = 1
    content = [{
      type = "container", props = { layout = "stack" }
      content = [{
        type = "widget", props = { title = "Duration vs Apdex by app" }
        content = { type = "visualization", id = "viz.scatter"
          props = {
            nrqlQueries = [{ accountIds = [var.account_id], query = "SELECT apdex(duration, t: 0.5), average(duration) FROM Transaction FACET appName LIMIT 10 SINCE 1 hour ago" }]
          }
        }
      }]
    }]
  })
}

# viz.bullet
resource "newrelic_notebook" "bullet_example" {
  title = "Bullet Chart Example"
  content = jsonencode({
    type = "declarative", version = 1
    content = [{
      type = "container", props = { layout = "stack" }
      content = [{
        type = "widget", props = { title = "Requests vs target" }
        content = { type = "visualization", id = "viz.bullet"
          props = {
            nrqlQueries = [{ accountIds = [var.account_id], query = "SELECT rate(count(*), 1 minute) AS 'Requests/min' FROM Transaction SINCE 5 minutes ago" }]
            limit = 1000
          }
        }
      }]
    }]
  })
}

# viz.sparkline
resource "newrelic_notebook" "sparkline_example" {
  title = "Sparkline Example"
  content = jsonencode({
    type = "declarative", version = 1
    content = [{
      type = "container", props = { layout = "stack" }
      content = [{
        type = "widget", props = { title = "Throughput sparkline" }
        content = { type = "visualization", id = "viz.sparkline"
          props = {
            nrqlQueries = [{ accountIds = [var.account_id], query = "SELECT count(*) FROM Transaction TIMESERIES AUTO SINCE 1 hour ago" }]
          }
        }
      }]
    }]
  })
}

# viz.markdown (text block)
resource "newrelic_notebook" "markdown_example" {
  title = "Markdown Example"
  content = jsonencode({
    type = "declarative", version = 1
    content = [{
      type = "container", props = { layout = "stack" }
      content = [{
        type  = "widget"
        props = {}
        content = {
          type = "visualization"
          id   = "viz.markdown"
          props = {
            text = "# Runbook\n\nUse **bold**, _italic_, `code`, and [links](https://one.newrelic.com).\n\n## Checklist\n\n- [ ] Verify dashboards\n- [ ] Check alert policies"
          }
        }
      }]
    }]
  })
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
              props = { text = "# ${each.value}\n\nAdd runbook steps for this service here." }
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
* `content` - (Optional) The notebook body as an HCL object using `jsonencode({...})`. Produces field-level diffs in `terraform plan`. Mutually exclusive with `content_json`. Must follow the declarative UI schema described in [Content schema](#content-schema).
* `content_json` - (Optional) The notebook body as a raw JSON string. Use when working from a UI export or a file. Produces line-level diffs of normalized content. Mutually exclusive with `content`. Must follow the declarative UI schema.
* `organization_id` - (Computed) The New Relic organization ID. Resolved automatically from the provider credentials.

## Attributes Reference

* `guid` - The unique entity identifier (GUID) assigned to the notebook by New Relic.
* `blob_id` - The blob identifier of the current notebook content. Updated after every Terraform-managed write. Used internally to detect when content has changed.

---

## Nested `container` blocks

The `content` array at the top level of the declarative UI document contains one or more **container** blocks. A container groups a set of widgets with a shared layout.

The following arguments are supported in a container block's `props`:

* `layout` - (Optional) The layout algorithm for the widgets in this container. Accepted values: `"stack"` (stacks widgets vertically, default), `"grid"` (arranges widgets in a grid). Defaults to `"stack"` when `props` is omitted entirely.

---

## Nested `widget` blocks

Each widget block within a container's `content` array supports the following structure:

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

These are set directly on the `widget` object's `props`, not inside the `content.props`. They are optional and shared across all widget types:

* `title` - (Optional) A label displayed above the rendered chart. Corresponds to the **Name** field in the UI's chart customization panel.
* `ignoreTimeRange` - (Optional) When `true`, the widget uses the time range specified in its own NRQL query instead of inheriting the notebook's global time picker. Defaults to `false`.

### Supported visualization types

The following visualization IDs are supported. Each maps to a chart type in the New Relic query builder and the Notebooks UI.

| `id` | Description | Required NRQL shape |
|---|---|---|
| `viz.markdown` | Text block with Markdown support | N/A |
| `viz.line` | Line chart | `TIMESERIES` recommended |
| `viz.area` | Area chart | `TIMESERIES` recommended |
| `viz.bar` | Bar chart | `FACET` recommended |
| `viz.stacked-bar` | Stacked bar chart | `FACET` or `TIMESERIES` |
| `viz.pie` | Pie chart | `FACET` recommended |
| `viz.table` | Data table | `FACET` recommended |
| `viz.billboard` | Single-value display | Single aggregation |
| `viz.histogram` | Distribution histogram | `histogram()` function |
| `viz.heatmap` | Heatmap | `histogram()` with `FACET` |
| `viz.funnel` | Conversion funnel | `funnel()` function |
| `viz.scatter` | Scatter plot | Two-value result |
| `viz.bullet` | Bullet chart | Single aggregation |
| `viz.sparkline` | Sparkline (mini line) | `TIMESERIES` |

---

### `viz.markdown` — Markdown text block

Renders a text block supporting [GitHub-flavored Markdown](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/add-queries-to-notebooks/).

**Props (`content.props`):**

* `text` - (Required) The Markdown source string. Supports headings, bold, italic, inline code, fenced code blocks, ordered and unordered lists, task lists (`- [ ]`), tables, and links. Newlines within the string must be expressed as `\n`.

**Example:**
```json
{
  "type": "widget",
  "props": {},
  "content": {
    "type": "visualization",
    "id": "viz.markdown",
    "props": {
      "text": "## Investigation Notes\n\nAdd your analysis here.\n\n- [ ] Check alert history\n- [ ] Review recent deploys"
    }
  }
}
```

---

### `viz.line` — Line chart

Plots one or more NRQL time series as line(s).

**Props (`content.props`):**

* `nrqlQueries` - (Required) An array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks) below.
* `legendEnabled` - (Optional) Show or hide the chart legend. Defaults to `true`.
* `facetShowOtherSeries` - (Optional) Show the aggregated "Other" group when a `FACET` returns more than the maximum number of series. Defaults to `false`.
* `yAxisLeft` - (Optional) Configuration for the left Y axis. See [Nested `yAxisLeft` / `yAxisRight` blocks](#nested-yaxisleft--yaxisright-blocks) below.
* `yAxisRight` - (Optional) Configuration for a second (right) Y axis. See [Nested `yAxisLeft` / `yAxisRight` blocks](#nested-yaxisleft--yaxisright-blocks) below.
* `nullValues` - (Optional) Specifies how null data points are rendered. See [Nested `nullValues` blocks](#nested-nullvalues-blocks) below.
* `colors` - (Optional) Per-series color overrides. See [Nested `colors` blocks](#nested-colors-blocks) below.
* `units` - (Optional) Per-series unit overrides. See [Nested `units` blocks](#nested-units-blocks) below.
* `thresholds` - (Optional) Horizontal threshold lines drawn on the chart. See [Nested `thresholds` blocks (line/area/stacked-bar)](#nested-thresholds-blocks-linearea-and-stacked-bar) below.
* `refreshRate` - (Optional) Data refresh interval in milliseconds. See [Valid `refreshRate` values](#valid-refreshrate-values) below.
* `tooltip` - (Optional) Tooltip display mode. See [Nested `tooltip` blocks](#nested-tooltip-blocks) below.
* `lineInterpolation` - (Optional) How data points are connected. Valid values: `"linear"` (default), `"smooth"`, `"stepBefore"`, `"stepAfter"`.
* `markers` - (Optional) Vertical reference lines drawn at specific timestamps or values on the chart. See [Nested `markers` blocks](#nested-markers-blocks) below.
* `chartTypeOverrides` - (Optional) Override the chart type for individual series within the same query result. See [Nested `chartTypeOverrides` blocks](#nested-charttypeoverrides-blocks) below.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Throughput over time" },
  "content": {
    "type": "visualization",
    "id": "viz.line",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [1234567],
          "query": "SELECT rate(count(*), 1 minute) FROM Transaction FACET appName TIMESERIES 5 minutes SINCE 3 hours ago"
        }
      ],
      "legendEnabled": true,
      "yAxisLeft": { "zero": true, "min": 0 },
      "nullValues": { "nullValue": "zero" },
      "lineInterpolation": "smooth"
    }
  }
}
```

---

### `viz.area` — Area chart

Plots time series with filled areas below the lines. Useful for showing volume.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `legendEnabled` - (Optional) Show or hide the legend. Defaults to `true`.
* `facetShowOtherSeries` - (Optional) Show aggregated "Other" group. Defaults to `false`.
* `yAxisLeft` - (Optional) Left Y axis configuration.
* `nullValues` - (Optional) Null value rendering.
* `colors` - (Optional) Per-series color overrides.
* `units` - (Optional) Per-series unit overrides.
* `thresholds` - (Optional) Horizontal threshold lines.
* `refreshRate` - (Optional) Refresh interval in ms.
* `tooltip` - (Optional) Tooltip mode.
* `lineInterpolation` - (Optional) Line interpolation. Valid values: `"linear"` (default), `"smooth"`, `"stepBefore"`, `"stepAfter"`.
* `markers` - (Optional) Vertical reference lines at specific timestamps or values. See [Nested `markers` blocks](#nested-markers-blocks) below.

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
        {
          "accountIds": [1234567],
          "query": "SELECT average(memoryUsedPercent) FROM SystemSample FACET hostname TIMESERIES 5 minutes SINCE 3 hours ago"
        }
      ],
      "legendEnabled": true,
      "yAxisLeft": { "zero": true, "min": 0, "max": 100 }
    }
  }
}
```

---

### `viz.bar` — Bar chart

Compares values across discrete categories.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `legendEnabled` - (Optional) Show or hide the legend. Defaults to `true`.
* `facetShowOtherSeries` - (Optional) Show aggregated "Other" group. Defaults to `false`.
* `colors` - (Optional) Per-series color overrides.
* `refreshRate` - (Optional) Refresh interval in ms.
* `tooltip` - (Optional) Tooltip mode.

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
        {
          "accountIds": [1234567],
          "query": "SELECT count(*) FROM Transaction WHERE error IS true FACET name SINCE 1 hour ago ORDER BY count(*) DESC LIMIT 10"
        }
      ],
      "facetShowOtherSeries": false
    }
  }
}
```

---

### `viz.stacked-bar` — Stacked bar chart

Like a bar chart but with facets stacked within each bar, showing composition.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `legendEnabled` - (Optional) Show or hide the legend. Defaults to `true`.
* `facetShowOtherSeries` - (Optional) Show aggregated "Other" group. Defaults to `false`.
* `colors` - (Optional) Per-series color overrides.
* `thresholds` - (Optional) Horizontal threshold lines.
* `refreshRate` - (Optional) Refresh interval in ms.
* `tooltip` - (Optional) Tooltip mode.

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
        {
          "accountIds": [1234567],
          "query": "SELECT count(*) FROM Transaction FACET httpResponseCode TIMESERIES 5 minutes SINCE 3 hours ago"
        }
      ],
      "legendEnabled": true
    }
  }
}
```

---

### `viz.pie` — Pie chart

Displays proportional data as a pie. Best for 5-7 categories.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `facetShowOtherSeries` - (Optional) Show aggregated "Other" group. Defaults to `false`.
* `colors` - (Optional) Per-series color overrides.
* `legendEnabled` - (Optional) Show or hide the legend. Defaults to `true`.
* `refreshRate` - (Optional) Refresh interval in ms.

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
        {
          "accountIds": [1234567],
          "query": "SELECT count(*) FROM Transaction FACET httpResponseCode SINCE 1 hour ago"
        }
      ],
      "facetShowOtherSeries": true
    }
  }
}
```

---

### `viz.table` — Table

Displays query results as a tabular list. Supports multi-column data.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects.
* `facetShowOtherSeries` - (Optional) Show aggregated "Other" group. Defaults to `false`.
* `refreshRate` - (Optional) Refresh interval in ms.

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
        {
          "accountIds": [1234567],
          "query": "SELECT average(duration)*1000 AS 'Avg ms', max(duration)*1000 AS 'Max ms', count(*) AS 'Calls' FROM Transaction FACET name SINCE 1 hour ago ORDER BY average(duration) DESC LIMIT 20"
        }
      ]
    }
  }
}
```

---

### `viz.billboard` — Billboard

Displays a single large metric value with optional color-coded thresholds. Use for KPIs and summary statistics.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects. The query should return a single aggregated value or a small number of faceted values.
* `thresholdsWithSeriesOverrides` - (Optional) Threshold ranges that control the color of the displayed value. See [Nested `thresholdsWithSeriesOverrides` blocks (billboard)](#nested-thresholdswithseriesoverrides-blocks-billboard) below.
* `refreshRate` - (Optional) Refresh interval in ms.

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
        {
          "accountIds": [1234567],
          "query": "SELECT percentage(count(*), WHERE error IS true) AS 'Error %' FROM Transaction SINCE 1 hour ago"
        }
      ],
      "thresholdsWithSeriesOverrides": {
        "thresholds": [
          { "to": 1,           "severity": "success"  },
          { "from": 1, "to": 5, "severity": "warning"  },
          { "from": 5,          "severity": "critical" }
        ]
      }
    }
  }
}
```

---

### `viz.histogram` — Histogram

Shows the distribution of values across automatically computed buckets. Requires the `histogram()` NRQL function.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects. The query must use `histogram(attribute)` or `histogram(attribute, width: N, buckets: N)`.
* `colors` - (Optional) Per-series color overrides.
* `refreshRate` - (Optional) Refresh interval in ms.

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
        {
          "accountIds": [1234567],
          "query": "SELECT histogram(duration, width: 0.1, buckets: 20) FROM Transaction SINCE 1 hour ago"
        }
      ]
    }
  }
}
```

---

### `viz.heatmap` — Heatmap

Displays a histogram-style distribution broken down by a `FACET`. Useful for comparing distributions across entities.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects. The query must use `histogram()` with a `FACET`.
* `colors` - (Optional) Per-series color overrides.
* `refreshRate` - (Optional) Refresh interval in ms.

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
        {
          "accountIds": [1234567],
          "query": "SELECT histogram(duration) FROM Transaction FACET appName SINCE 1 hour ago"
        }
      ]
    }
  }
}
```

---

### `viz.funnel` — Funnel chart

Visualizes a conversion funnel across sequential steps. Requires the `funnel()` NRQL function.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects. The query must use the `funnel()` function.
* `colors` - (Optional) Per-series color overrides.
* `refreshRate` - (Optional) Refresh interval in ms.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "User conversion funnel" },
  "content": {
    "type": "visualization",
    "id": "viz.funnel",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [1234567],
          "query": "SELECT funnel(session, WHERE name = 'HomeView' AS 'Home', WHERE name = 'ProductView' AS 'Product', WHERE name = 'CartView' AS 'Cart', WHERE name = 'ConfirmationView' AS 'Purchase') FROM PageView SINCE 1 day ago"
        }
      ]
    }
  }
}
```

---

### `viz.scatter` — Scatter plot

Plots data points on an X-Y plane. The query should return two numeric values per facet.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects. The query should return two numeric attributes and a `FACET`.
* `colors` - (Optional) Per-series color overrides.
* `refreshRate` - (Optional) Refresh interval in ms.
* `tooltip` - (Optional) Tooltip mode.

**Example:**
```json
{
  "type": "widget",
  "props": { "title": "Duration vs Apdex by app" },
  "content": {
    "type": "visualization",
    "id": "viz.scatter",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [1234567],
          "query": "SELECT apdex(duration, t: 0.5), average(duration) FROM Transaction FACET appName LIMIT 20 SINCE 1 hour ago"
        }
      ]
    }
  }
}
```

---

### `viz.bullet` — Bullet chart

Displays a value as a horizontal bar against a target limit. The target is set via the `limit` prop.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects. The query should return a single aggregated value.
* `limit` - (Required) The target value shown as the bullet chart's goal line.
* `refreshRate` - (Optional) Refresh interval in ms.

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
        {
          "accountIds": [1234567],
          "query": "SELECT rate(count(*), 1 minute) AS 'Requests/min' FROM Transaction SINCE 5 minutes ago"
        }
      ],
      "limit": 5000
    }
  }
}
```

---

### `viz.sparkline` — Sparkline

A compact mini line chart without axes or labels. Useful for quick trend indicators inline with other content.

**Props (`content.props`):**

* `nrqlQueries` - (Required) Array of NRQL query objects. The query should use `TIMESERIES`.
* `refreshRate` - (Optional) Refresh interval in ms.

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
        {
          "accountIds": [1234567],
          "query": "SELECT count(*) FROM Transaction TIMESERIES AUTO SINCE 1 hour ago"
        }
      ]
    }
  }
}
```

---

## Nested `nrqlQueries` blocks

All query-based visualizations accept a `nrqlQueries` array. Each item in the array supports the following:

* `accountIds` - (Required) An array of one or more New Relic account IDs to query against, e.g. `[1234567]`. For cross-account queries, include multiple IDs: `[1234567, 7654321]`.
* `query` - (Required) A valid NRQL query string. See [Introduction to NRQL](https://docs.newrelic.com/docs/nrql/get-started/introduction-nrql-new-relics-query-language/) for help.

-> **NOTE:** If you query an account to which you do not have access, the widget will render with a "data inaccessible" message. No Terraform error is raised.

Multiple `nrqlQueries` objects can be provided on query-based widgets to overlay results from different accounts or queries on the same chart.

---

## Nested `yAxisLeft` / `yAxisRight` blocks

Applies to: `viz.line`, `viz.area`.

Both blocks accept the same fields:

* `zero` - (Optional) When `true`, the axis always starts at zero even if the minimum data value is higher. Defaults to `true`.
* `min` - (Optional) The minimum value to display on the axis. Overrides `zero` if set.
* `max` - (Optional) The maximum value to display on the axis. Set to `0` to auto-scale.
* `scale` - (Optional) The scale type of the axis. Valid values: `"linear"` (default), `"log"` (logarithmic).

For `yAxisRight` only:

* `series` - (Optional) An array of series names (strings) to plot against the right Y axis instead of the left.

**Example:**
```json
{
  "yAxisLeft": { "zero": true, "min": 0, "max": 100, "scale": "linear" },
  "yAxisRight": { "zero": true, "min": 0, "series": ["P99 latency"] }
}
```

---

## Nested `thresholds` blocks (line, area, and stacked-bar)

Applies to: `viz.line`, `viz.area`, `viz.stacked-bar`.

The `thresholds` key is an array of threshold range objects that draw horizontal bands on the chart:

* `name` - (Optional) A label for the threshold.
* `from` - (Optional) The lower bound of the threshold range (inclusive). Omit to start from negative infinity.
* `to` - (Optional) The upper bound of the threshold range (inclusive). Omit to extend to positive infinity.
* `severity` - (Required) The color applied to data within the threshold range. Valid values: `"success"` (green), `"warning"` (yellow), `"severe"` (orange), `"critical"` (red), `"unavailable"` (grey).
* `isLabelVisible` - (Optional) When `true`, always shows the threshold label on the chart regardless of whether the data crosses it. Defaults to `false`.

**Example:**
```json
{
  "thresholds": [
    { "name": "OK",       "to": 1,            "severity": "success"  },
    { "name": "Warning",  "from": 1,  "to": 5, "severity": "warning"  },
    { "name": "Critical", "from": 5,            "severity": "critical" }
  ]
}
```

---

## Nested `thresholdsWithSeriesOverrides` blocks (billboard)

Applies to: `viz.billboard`.

The billboard threshold format is different from line/area. The `thresholdsWithSeriesOverrides` key contains a `thresholds` array and an optional `seriesOverrides` array:

### `thresholds`

Each threshold defines a value range and the color to apply:

* `from` - (Optional) Lower bound (inclusive). Omit to start from negative infinity.
* `to` - (Optional) Upper bound (inclusive). Omit to extend to positive infinity.
* `severity` - (Required) The color applied when the value falls in this range. Valid values: `"success"` (green), `"warning"` (yellow), `"severe"` (orange), `"critical"` (red), `"unavailable"` (grey).

### `seriesOverrides`

When a billboard query returns multiple faceted values, `seriesOverrides` lets you apply different thresholds to specific series by name:

* `seriesName` - (Required) The series name (facet value) to override.
* `thresholds` - (Required) An array of threshold objects in the same format as above.

**Example:**
```json
{
  "thresholdsWithSeriesOverrides": {
    "thresholds": [
      { "to": 1,            "severity": "success"  },
      { "from": 1, "to": 5, "severity": "warning"  },
      { "from": 5,           "severity": "critical" }
    ],
    "seriesOverrides": [
      {
        "seriesName": "EU error rate",
        "thresholds": [
          { "to": 2,            "severity": "success"  },
          { "from": 2, "to": 8, "severity": "warning"  },
          { "from": 8,           "severity": "critical" }
        ]
      }
    ]
  }
}
```

---

## Nested `nullValues` blocks

Applies to: `viz.line`, `viz.area`.

Controls how null (missing) data points are handled.

* `nullValue` - (Optional) The strategy for the entire chart. Valid values:
  - `"default"` — Platform default (typically shows a gap).
  - `"remove"` — Removes the data point; line skips over the gap.
  - `"preserve"` — Keeps the null as a visible gap in the line.
  - `"zero"` — Treats null as zero.
* `seriesOverrides` - (Optional) An array of per-series overrides. Each object supports:
  - `seriesName` - (Required) The series name (string) to override.
  - `nullValue` - (Required) One of the values above, applied only to this series.

**Example:**
```json
{
  "nullValues": {
    "nullValue": "zero",
    "seriesOverrides": [
      { "seriesName": "P99 latency", "nullValue": "preserve" }
    ]
  }
}
```

---

## Nested `colors` blocks

Applies to: `viz.line`, `viz.area`, `viz.bar`, `viz.stacked-bar`, `viz.pie`, `viz.histogram`, `viz.heatmap`, `viz.funnel`, `viz.scatter`.

* `color` - (Optional) A default color applied to all series. Accepted formats: RGB hex (`"#ff0000"`), named colors.
* `seriesOverrides` - (Optional) An array of per-series overrides. Each object supports:
  - `seriesName` - (Required) The series name (string) to override.
  - `color` - (Required) The color for this series. Accepted formats: RGB hex or named colors.

**Example:**
```json
{
  "colors": {
    "color": "#0070F0",
    "seriesOverrides": [
      { "seriesName": "errors", "color": "#FF0000" },
      { "seriesName": "warnings", "color": "#FFA500" }
    ]
  }
}
```

---

## Nested `units` blocks

Applies to: `viz.line`, `viz.area`.

Adds unit labels to Y axis values and chart tooltips.

* `unit` - (Optional) The default unit string applied to all series (e.g. `"ms"`, `"%"`, `"MB"`).
* `seriesOverrides` - (Optional) An array of per-series unit overrides. Each object supports:
  - `seriesName` - (Required) The series name (string) to override.
  - `unit` - (Required) The unit string for this series.

**Example:**
```json
{
  "units": {
    "unit": "ms",
    "seriesOverrides": [
      { "seriesName": "Throughput", "unit": "rpm" }
    ]
  }
}
```

---

## Nested `tooltip` blocks

Applies to: `viz.line`, `viz.area`, `viz.bar`, `viz.stacked-bar`, `viz.scatter`.

* `mode` - (Required) How tooltips are displayed when hovering over the chart. Valid values:
  - `"all"` — Show a tooltip for all series at the hovered time.
  - `"single"` — Show a tooltip for the closest single series only.
  - `"hidden"` — Disable tooltips entirely.

**Example:**
```json
{
  "tooltip": { "mode": "all" }
}
```

---

## Nested `markers` blocks

Applies to: `viz.line`, `viz.area`.

Vertical reference lines drawn at specific values or timestamps. Use markers to annotate events such as deploys or incidents directly on a time series chart.

* `data` - (Required) An array of marker objects. Each marker supports:
  - `label` - (Required) The text label shown on the marker line.
  - `value` - (Optional) A numeric Y-axis value at which to draw the marker (for threshold-style markers).
  - `timestamp` - (Optional) A Unix timestamp in milliseconds at which to draw the marker (for time-based annotations).
  - `color` - (Optional) The color of the marker line. Accepted formats: RGB hex (e.g. `"#FF0000"`).

**Example:**
```json
{
  "markers": {
    "data": [
      { "label": "Deploy v2.1", "timestamp": 1728000000000, "color": "#FF6600" },
      { "label": "Rollback",    "timestamp": 1728003600000, "color": "#FF0000" }
    ]
  }
}
```

---

## Nested `chartTypeOverrides` blocks

Applies to: `viz.line`, `viz.area`, `viz.stacked-bar`.

Overrides the chart type for individual series within the same query result. This allows, for example, rendering a specific FACET series as a line while the rest appear as areas.

* An array of override objects. Each object supports:
  - `seriesName` - (Required) The series name (facet value or alias) to override.
  - `chartType` - (Required) The visualization type for this series. Valid values: `"LINE"`, `"AREA"`, `"BAR"`, `"STACKED_BAR"`, `"SCATTER"`.

**Example:**
```json
{
  "chartTypeOverrides": [
    { "seriesName": "P99 latency", "chartType": "LINE" },
    { "seriesName": "P50 latency", "chartType": "AREA" }
  ]
}
```

---

## Valid `refreshRate` values

Applies to all query-based visualization types.

The `refreshRate` prop sets how frequently the widget re-executes its NRQL query and re-renders the chart. The value is specified in milliseconds:

| Value | Refresh interval |
|---|---|
| `0` | No automatic refresh |
| `5000` | Every 5 seconds |
| `30000` | Every 30 seconds |
| `60000` | Every 60 seconds (1 minute) |
| `300000` | Every 5 minutes |
| `1800000` | Every 30 minutes |
| `3600000` | Every 60 minutes (1 hour) |
| `10800000` | Every 3 hours |
| `43200000` | Every 12 hours |
| `86400000` | Every 24 hours |

When omitted, the widget uses the platform default refresh rate.

---

## Further reading

The Blob API stores notebook content verbatim. Terraform tracks the entire blob. New Relic does not assign special meaning to any field outside of the declarative UI envelope, so you can include arbitrary top-level metadata; Terraform will diff those fields exactly like any other part of the content if they change.

For additional reference on visualization types, props, and NRQL:

- [Visualizations in notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/visualizations-in-notebooks/) — chart types available in notebooks
- [Blob Storage API for notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/blob-storage-api-for-notebooks/) — the underlying REST API used by this resource
- [Notebook examples](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/notebooks-examples/) — worked multi-block notebook JSON from the docs team
- [NR1 SDK: Chart components](https://docs.newrelic.com/docs/new-relic-solutions/build-nr-ui/sdk-component/charts/Charts/) — the SDK chart components notebooks viz IDs map to
- [NerdGraph: Create dashboard widgets](https://docs.newrelic.com/docs/apis/nerdgraph/examples/create-widgets-dashboards-api/) — `props` schema reference for each widget type (shared with notebooks)
- [Introduction to NRQL](https://docs.newrelic.com/docs/nrql/get-started/introduction-nrql-new-relics-query-language/) — writing NRQL queries

---

## Import

Notebooks can be imported by GUID. Optionally append `:content` or `:content_json` to control which field is populated in state, matching your Terraform configuration.

```
# Default - imports into content_json (for configs using content_json = file(...) or inline JSON)
$ terraform import newrelic_notebook.example <guid>
$ terraform import newrelic_notebook.example <guid>:content_json

# Import into content field (for configs using content = jsonencode({...}))
$ terraform import newrelic_notebook.example <guid>:content
```

After importing, run `terraform plan`. If the imported state and your config use the same mode, the plan will show no changes.

## Plan Diff Behavior

Both content fields store a normalized form of the JSON in state (alphabetically sorted keys, 2-space indentation). This means:

- Reformatting your HCL or JSON file without changing any values produces **no diff** in `terraform plan`.
- Changing a single widget property shows **only that property** as changed.
- Externally modifying the notebook in the UI causes the changed fields to surface **precisely** in the next `terraform plan`.
