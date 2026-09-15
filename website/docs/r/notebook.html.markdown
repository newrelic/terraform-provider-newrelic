---
layout: "newrelic"
page_title: "New Relic: newrelic_notebook"
sidebar_current: "docs-newrelic-resource-notebook"
description: |-
  Create and manage New Relic Notebooks.
---

# Resource: newrelic\_notebook

Use this resource to create and manage [New Relic Notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/introduction-notebooks/).

Notebooks are shareable documents that combine live NRQL queries, visualizations, and Markdown narrative in a single view. The resource manages the notebook's title and its full body content using the declarative UI format used natively by the platform.

-> **NOTE:** The notebook body must be a valid declarative UI document. See [Content Format and Schema](#content-format-and-schema) for the required envelope structure and supported visualization types.

See additional [examples](#additional-examples).

## Example Usage

```hcl
resource "newrelic_notebook" "incident_runbook" {
  title = "Production API Incident - Investigation"

  content = jsonencode({
    type    = "declarative"
    version = 1
    content = [{
      type  = "container"
      props = { layout = "stack" }
      content = [
        {
          type  = "widget"
          content = {
            type = "visualization"
            id   = "viz.markdown"
            props = {
              text = "# Production API Incident\n\n**Status**: Resolved  \n**Impact**: ~15% of API requests failed\n\nUse this notebook to walk through what happened and track action items."
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
              nrqlQueries = [{
                accountIds = [var.account_id]
                query      = "SELECT percentage(count(*), WHERE httpResponseCode >= 400) AS 'Error Rate %' FROM Transaction WHERE appName = 'api-production' SINCE '2024-10-15 14:00:00' UNTIL '2024-10-15 16:00:00'"
              }]
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
              nrqlQueries = [{
                accountIds = [var.account_id]
                query      = "SELECT percentage(count(*), WHERE httpResponseCode >= 400) FROM Transaction WHERE appName = 'api-production' TIMESERIES 1 minute SINCE '2024-10-15 14:00:00' UNTIL '2024-10-15 16:00:00'"
              }]
              legend    = { enabled = true }
              yAxisLeft = { zero = true }
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
              nrqlQueries = [{
                accountIds = [var.account_id]
                query      = "SELECT count(*) AS 'Errors', average(duration)*1000 AS 'Avg ms' FROM Transaction WHERE httpResponseCode >= 400 FACET request.uri SINCE '2024-10-15 14:00:00' UNTIL '2024-10-15 16:00:00' LIMIT 10"
              }]
              initialSorting = { name = "Errors", direction = "desc" }
            }
          }
        },
        {
          type  = "widget"
          content = {
            type = "visualization"
            id   = "viz.markdown"
            props = {
              text = "## Root Cause\n\n- Database connection pool exhausted\n\n## Action Items\n\n- [ ] Raise DB connection pool limit\n- [ ] Alert on pool utilization > 80%"
            }
          }
        }
      ]
    }]
  })
}
```

## Argument Reference

The following arguments are supported:

  * `title` - (Required) The title of the notebook. Must be unique within the organization.
  * `content` - (Required) The notebook body as a JSON string. Accepts the raw JSON exported from the New Relic UI or produced by `jsonencode({...})`. Produces line-level diffs of the normalized JSON in `terraform plan`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `guid` - The unique entity identifier (GUID) of the notebook in New Relic.
  * `blob_id` - The blob identifier of the current notebook content, updated after each Terraform-managed write.
  * `organization_id` - The New Relic organization ID the notebook belongs to. Resolved automatically from the provider credentials.

## Import

Notebooks can be imported using their entity GUID:

```
$ terraform import newrelic_notebook.example <guid>
```

After importing, run `terraform plan`. The plan will show no changes if the `content` in your configuration matches the normalized content fetched from the API.

---

## Content Format and Schema

Every notebook body must follow the declarative UI envelope:

```json
{
  "type": "declarative",
  "version": 1,
  "content": [ "<container blocks>" ]
}
```

### Nested `container` blocks

The top-level `content` array contains one or more container blocks. A container groups widgets with a shared layout.

The following argument is supported in a container's `props`:

  * `layout` - (Optional) The layout for the widgets. Valid values are `"stack"` (default) and `"grid"`.

### Nested `widget` blocks

Each widget inside a container's `content` array follows this structure:

```json
{
  "type": "widget",
  "props": { "<widget-level props>" },
  "content": {
    "type": "visualization",
    "id": "<viz-id>",
    "props": { "<visualization props>" }
  }
}
```

#### Widget-level `props`

The `props` key at the widget level has different requirements depending on the visualization type:

  * **`viz.markdown` widgets**: `props` must be **absent entirely** from the widget object. Do not include `"props": {}` on a markdown widget.
  * **All other viz widgets**: `props` must be present and must include `title`. The `title` key is **required** (can be an empty string `""`).

The following arguments are available on the widget's `props` object for non-markdown visualizations:

  * `title` - (Required for non-markdown widgets) A label shown above the rendered chart. May be an empty string `""`.
  * `ignoreTimeRange` - Equivalent to `platformOptions.ignoreTimeRange` in some visualizations. Prefer setting `platformOptions.ignoreTimeRange` directly on the visualization props.

### Supported visualization types

The `id` field identifies the chart type. The following values are supported:

| `id` | Display name | Typical NRQL shape |
|---|---|---|
| `viz.markdown` | Markdown | N/A |
| `viz.line` | Line | `TIMESERIES` |
| `viz.area` | Area | `TIMESERIES` |
| `viz.stacked-bar` | Bar (stacked) | `TIMESERIES` or `FACET` |
| `viz.bar` | Bar (simplified) | `FACET` |
| `viz.stacked-horizontal-bar` | Horizontal stacked bar | `FACET` |
| `viz.pie` | Pie | `FACET` |
| `viz.table` | Table | `FACET` |
| `viz.billboard` | Billboard | Single aggregation |
| `viz.gauge` | Gauge | Single aggregation |
| `viz.histogram` | Histogram | `histogram()` function |
| `viz.heatmap` | Heatmap | `histogram()` with `FACET` |
| `viz.scatter` | Scatter | Two-value `FACET` |
| `viz.apdex` | Apdex | Apdex-compatible |
| `viz.bullet` | Bullet | Single aggregation |
| `viz.funnel` | Funnel | `funnel()` function |
| `viz.event-feed` | Event feed | Any |
| `viz.json` | JSON | Any |
| `viz.sparkline` | Sparkline | `TIMESERIES` |
| `viz.sparkline-lite` | Sparkline Lite | `TIMESERIES` |

---

### `viz.markdown` — Markdown

Renders a text block supporting GitHub-flavoured Markdown.

-> **NOTE:** The widget-level `props` key must be **absent entirely** on `viz.markdown` widgets. Do not include `"props": {}` — omit the key altogether.

  * `text` - (Required) The Markdown source string. Supports headings, bold, italic, code, lists, task lists (`- [ ]`), tables, and links. Use `\n` for newlines within the string.

**Example:**
```json
{
  "type": "widget",
  "content": {
    "type": "visualization",
    "id": "viz.markdown",
    "props": {
      "text": "## Investigation Notes\n\n- [ ] Check alert history\n- [ ] Review recent deploys"
    }
  }
}
```

---

### `viz.line` — Line chart

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`). Example: `"props": { "title": "My chart" }`.

  * `nrqlQueries` - (Required) Array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `alertQueries` - (Optional) Array of alert violation objects to overlay warning/critical bands on the chart. Only available for `viz.line`. Each object requires `accountIds` (number[]), `violationId` (number), `duration` (number, ms), `endTime` (number, epoch ms).
  * `sqlQueries` - (Optional) Array of SQL query objects for Federated Data Source (FDS). Each requires `query` (string) and `accountId` (number) or `accountIds` (number[]).
  * `legend.enabled` - (Optional) Show or hide the legend. Default: `true`.
  * `legend.position` - (Optional) Legend position. Valid values: `"bottom"` (default), `"left"`, `"right"`, `null`.
  * `facet.showOtherSeries` - (Optional) Show the "Other" group for `FACET` queries. Default: `false`.
  * `yAxisLeft.min` / `yAxisLeft.max` - (Optional) Fixed minimum and maximum for the left Y axis.
  * `yAxisLeft.zero` - (Optional) Force zero as the axis origin. Default: `true`.
  * `yAxisLeft.scale` - (Optional) Axis scale. Valid values: `"linear"` (default), `"logarithmic"`.
  * `yAxisRight.min` / `yAxisRight.max` / `yAxisRight.zero` / `yAxisRight.scale` - (Optional) Right Y axis configuration. Same valid values as the left axis.
  * `yAxisRight.series[].name` - (Optional) Series names to plot against the right Y axis.
  * `nullValues.nullValue` - (Optional) How to handle null data points. Valid values: `"default"`, `"zero"`, `"preserve"`, `"remove"`.
  * `nullValues.seriesOverrides[].seriesName` / `.nullValue` - (Optional) Per-series null value override.
  * `colors.colorPalette` - (Optional) Color palette. Valid values: `"consistent"` (default), `"dynamic"`, `null`.
  * `colors.seriesOverrides[].seriesName` / `.color` - (Optional) Per-series color (RGB hex, e.g. `"#FF0000"`).
  * `units.unit` - (Optional) Default data unit for all series. See [Valid `units.unit` values](#valid-unitsunit-values).
  * `units.seriesOverrides[].seriesName` / `.unit` - (Optional) Per-series unit override.
  * `thresholds.isLabelVisible` - (Optional) Show threshold labels on the chart. Default: `false`.
  * `thresholds.thresholds[].name` / `.from` / `.to` / `.severity` - (Optional) Horizontal threshold bands. See [Nested `thresholds` blocks](#nested-thresholds-blocks-line-area-stacked-bar).
  * `chartStyles.lineInterpolation` - (Optional) Line interpolation. Valid values: `"linear"` (default), `"smooth"`, `"stepBefore"`, `"stepAfter"`.
  * `chartStyles.gradient.enabled` - (Optional) Gradient fill. Default: `false`.
  * `chartTypes.seriesOverrides[].seriesName` / `.chartType` - (Optional) Override chart type per series. Valid chartType values: `"line"`, `"area"`.
  * `tooltip.mode` - (Optional) Tooltip behaviour. Valid values: `"single"` (default), `"all"`, `"hidden"`.
  * `platformOptions.ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Default: `false`.
  * `refreshRate.frequency` - (Optional) Auto-refresh interval in milliseconds or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

---

### `viz.area` — Area chart

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

Supports the same props as `viz.line`, with the following differences:

  * `alertQueries` and `yAxisRight` are not available.
  * `nullValues` supports: `"default"`, `"zero"`, `"preserve"` (not `"remove"`).
  * Additional prop: `chartStyles.stacked.enabled` - (Optional) Enable stacked layout. Default: `true`.

---

### `viz.stacked-bar` — Bar chart (stacked / timeseries)

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

Supports all `viz.line` props except `alertQueries`, `yAxisRight`, `nullValues.remove`, `chartStyles.lineInterpolation`, `chartStyles.gradient`, and `chartTypes`. Additional prop: `chartStyles.stacked.enabled` - (Optional) Default: `true`.

---

### `viz.bar` — Bar chart (simplified)

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

A simpler bar chart for categorical FACET comparisons.

  * `nrqlQueries` - (Required)
  * `sqlQueries` - (Optional)
  * `facet.showOtherSeries` - (Optional) Default: `false`.
  * `colors.colorPalette` / `colors.seriesOverrides[]` - (Optional)
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.stacked-horizontal-bar` — Horizontal stacked bar

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

  * `nrqlQueries` - (Required)
  * `facet.showOtherSeries` - (Optional) Default: `false`.
  * `legend.enabled` - (Optional) Default: `true`.
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.pie` — Pie chart

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

  * `nrqlQueries` - (Required)
  * `sqlQueries` - (Optional)
  * `facet.showOtherSeries` - (Optional) Default: `true` (pie default differs from other types).
  * `legend.enabled` - (Optional) Default: `true`.
  * `colors.colorPalette` / `colors.seriesOverrides[]` - (Optional)
  * `chartStyles.gradient.enabled` - (Optional) Default: `false`.
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.table` — Table

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

  * `nrqlQueries` - (Required)
  * `sqlQueries` - (Optional)
  * `facet.showOtherSeries` - (Optional) Default: `false`.
  * `initialSorting` - (Optional) Default sort: object with `name` (string, column name) and `direction` (`"asc"` or `"desc"`).
  * `hiddenColumns[].columnName` - (Optional) Columns to hide.
  * `thresholds[].columnName` / `.from` / `.to` / `.severity` - (Optional) See [Nested `thresholds` blocks (table)](#nested-thresholds-blocks-table).
  * `dataFormatters` - (Optional) Custom data format configuration.
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.billboard` — Billboard

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

Displays a single large metric value with color-coded threshold ranges.

  * `nrqlQueries` - (Required)
  * `sqlQueries` - (Optional)
  * `thresholdsWithSeriesOverrides.thresholds[].from` / `.to` / `.severity` - (Optional) Threshold ranges. See [Nested `thresholdsWithSeriesOverrides` blocks](#nested-thresholdswithseriesoverrides-blocks-billboard).
  * `thresholdsWithSeriesOverrides.seriesOverrides[].seriesName` / `.from` / `.to` / `.severity` - (Optional) Per-series threshold overrides.
  * `billboardSettings.visual.alignment` - (Optional) Valid values: `"auto"`, `"stacked"`, `"inline"`.
  * `billboardSettings.visual.display` - (Optional) Valid values: `"auto"`, `"all"`, `"value"`, `"label"`, `"none"`.
  * `billboardSettings.gridOptions.columns` / `.label` / `.value` - (Optional) Grid layout options for multi-value billboards.
  * `billboardSettings.link.url` / `.title` / `.newTab` - (Optional) Clickable link on the billboard.
  * `units.unit` / `units.seriesOverrides[]` - (Optional)
  * `dataFormatters` - (Optional)
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.gauge` — Gauge

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

  * `nrqlQueries` - (Required)
  * `sqlQueries` - (Optional)
  * `gaugeSettings.displayMode` - (Optional) Valid values: `"arc"` (default), `"circular"`, `"bar"`.
  * `gaugeSettings.min` / `gaugeSettings.max` - (Optional) Scale bounds. Default max: `100`.
  * `gaugeSettings.display` - (Optional) Valid values: `"auto"` (default), `"all"`, `"value"`, `"name"`, `"none"`.
  * `gaugeSettings.showThresholdMarkers` - (Optional) Default: `true`.
  * `gaugeSettings.showThresholdLabels` - (Optional) Default: `false`.
  * `gaugeSettings.thresholdsColorMode` - (Optional) Valid values: `"solid"` (default), `"gradient"`.
  * `thresholds.thresholds[].from` / `.severity` - (Optional)
  * `units.unit` / `units.seriesOverrides[]` - (Optional)
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.histogram` — Histogram

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

  * `nrqlQueries` - (Required) Query must use `histogram(attr)` or `histogram(attr, width: N, buckets: N)`.
  * `legend.enabled` / `legend.position` - (Optional)
  * `colors.colorPalette` / `colors.seriesOverrides[]` - (Optional)
  * `yAxisLeft.min` / `.max` / `.zero` / `.scale` - (Optional)
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.heatmap` — Heatmap

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

  * `nrqlQueries` - (Required) Requires `histogram()` with a `FACET`.
  * `facet.showOtherSeries` - (Optional) Default: `false`.
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.scatter` — Scatter plot

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

  * `nrqlQueries` - (Required) Query should return two numeric values and a `FACET`.
  * `facet.showOtherSeries` - (Optional)
  * `legend.enabled` / `legend.position` - (Optional)
  * `nullValues.nullValue` / `nullValues.seriesOverrides[]` - (Optional)
  * `yAxisLeft.min` / `.max` / `.zero` / `.scale` - (Optional)
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.apdex` — Apdex

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

  * `nrqlQueries` - (Required)
  * `legend.enabled` / `legend.position` - (Optional)
  * `tooltip.mode` - (Optional)
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.bullet` — Bullet chart

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

  * `nrqlQueries` - (Required)
  * `limit` - (Required) The target value shown as the goal line.
  * `sqlQueries` - (Optional)
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.funnel` — Funnel chart

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

  * `nrqlQueries` - (Required) Requires the `funnel()` NRQL function.
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.event-feed` — Event feed

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

  * `nrqlQueries` - (Required)
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.json` — JSON

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

  * `nrqlQueries` - (Required)
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.sparkline` — Sparkline

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

  * `nrqlQueries` - (Required)
  * `facet.showOtherSeries` - (Optional)
  * `nullValues.nullValue` / `nullValues.seriesOverrides[]` - (Optional)
  * `chartStyles.lineInterpolation` - (Optional)
  * `yAxisLeft.min` / `.max` / `.zero` - (Optional)
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.sparkline-lite` — Sparkline Lite

Ultra-compact sparkline with no axis labels. Supports the same props as `viz.sparkline` except `yAxisLeft`.

---

## Nested `nrqlQueries` blocks

All query-based visualizations accept a `nrqlQueries` array. Each item supports:

  * `accountIds` - (Required) An array of one or more New Relic account IDs.
  * `query` - (Required) A valid NRQL query string.
  * `offset` - (Optional) Number of milliseconds to offset the query time window.

## Nested `thresholds` blocks (line, area, stacked-bar)

  * `name` - (Optional) A label for the threshold.
  * `from` - (Optional) Lower bound of the threshold range (inclusive).
  * `to` - (Optional) Upper bound of the threshold range (inclusive).
  * `severity` - (Required) Valid values: `"success"`, `"warning"`, `"severe"`, `"critical"`, `"unavailable"`.

## Nested `thresholds` blocks (table)

  * `columnName` - (Required) The column to apply the threshold to.
  * `from` - (Optional) Lower bound.
  * `to` - (Optional) Upper bound.
  * `severity` - (Required) Valid values: `"success"`, `"warning"`, `"severe"`, `"critical"`, `"unavailable"`.

## Nested `thresholdsWithSeriesOverrides` blocks (billboard)

  * `thresholds[].from` / `.to` / `.severity` - (Optional) Global threshold ranges.
  * `seriesOverrides[].seriesName` / `.from` / `.to` / `.severity` - (Optional) Per-series overrides.

## Valid `units.unit` values

`APDEX`, `BITS`, `BITS_PER_MS`, `BITS_PER_SECOND`, `BYTES`, `BYTES_PER_MS`, `BYTES_PER_SECOND`, `CELSIUS`, `COUNT`, `DOLLAR`, `HERTZ`, `MS`, `PAGES_PER_SECOND`, `PERCENTAGE`, `REQUESTS_PER_SECOND`, `REQUESTS_PER_MINUTE`, `SECONDS`, `TIMESTAMP`

## Valid `refreshRate.frequency` values

| Value | Interval |
|---|---|
| `"auto"` | Platform default |
| `0` | No refresh |
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

## Additional Examples

### Load notebook from a JSON file

```hcl
resource "newrelic_notebook" "weekly_review" {
  title        = "Weekly Service Health Review"
  content = file("${path.module}/notebooks/weekly-health.json")
}
```

### One notebook per service

```hcl
variable "services" {
  type    = set(string)
  default = ["checkout", "payments", "inventory"]
}

resource "newrelic_notebook" "runbooks" {
  for_each = var.services
  title    = "${each.value} runbook"

  content = jsonencode({
    type    = "declarative"
    version = 1
    content = [{
      type  = "container"
      props = { layout = "stack" }
      content = [
        {
          type  = "widget"
          content = {
            type = "visualization", id = "viz.markdown"
            props = { text = "# ${each.value}\n\nAdd runbook steps here." }
          }
        },
        {
          type  = "widget"
          props = { title = "Error rate" }
          content = {
            type = "visualization", id = "viz.billboard"
            props = {
              nrqlQueries = [{ accountIds = [var.account_id], query = "SELECT percentage(count(*), WHERE error IS true) AS 'Error %' FROM Transaction WHERE appName = '${each.value}' SINCE 1 hour ago" }]
            }
          }
        }
      ]
    }]
  })
}
```
