---
layout: "newrelic"
page_title: "New Relic: newrelic_notebook"
sidebar_current: "docs-newrelic-resource-notebook"
description: |-
  Create and manage New Relic Notebooks.
---

# Resource: newrelic\_notebook

Use this resource to create and manage [New Relic Notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/introduction-notebooks/).

Notebooks are shareable documents that combine live NRQL queries, visualizations, and Markdown narrative in a single view. The resource manages the notebook's title and its full block content using the declarative UI format used natively by the platform.

-> **NOTE:** The notebook body must be a valid declarative UI document. See [Content Format and Schema](#content-format-and-schema) for the required envelope structure and supported visualization types.

See additional [examples](#additional-examples).

## Example Usage

The primary recommended approach is to supply the notebook body as a raw JSON string using `content_json`. This keeps the JSON structure visible, makes it easy to copy content from the New Relic UI, and produces clean line-level diffs in `terraform plan`.

```hcl
resource "newrelic_notebook" "incident_runbook" {
  title = "Production API Incident - Investigation"

  content_json = jsonencode({
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
  * `content_json` - (Optional) The notebook body as a JSON string. Recommended for notebooks copied from the New Relic UI or loaded from a file using `file()`. Produces line-level diffs of the normalized JSON. Mutually exclusive with `content`.
  * `content` - (Optional) The notebook body expressed as an HCL object using `jsonencode({...})`. Produces field-level diffs in `terraform plan` at the expense of more verbose configuration. Mutually exclusive with `content_json`.

Exactly one of `content` or `content_json` must be specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `guid` - The unique entity identifier (GUID) of the notebook in New Relic.
  * `blob_id` - The blob identifier of the current notebook content, updated after each Terraform-managed write.
  * `organization_id` - The New Relic organization ID the notebook belongs to. Resolved automatically from the provider credentials.

## Import

Notebooks can be imported using their entity GUID. Append `:content` or `:content_json` to control which attribute is populated in state — this should match your Terraform configuration.

```
# Populates content_json (default)
$ terraform import newrelic_notebook.example <guid>
$ terraform import newrelic_notebook.example <guid>:content_json

# Populates content (for configs that use content = jsonencode({...}))
$ terraform import newrelic_notebook.example <guid>:content
```

After importing, run `terraform plan`. If the mode matches your configuration, the plan will show no changes.

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

  * **`viz.markdown` widgets**: `props` must be **absent entirely** from the widget object. Do not include `"props": {}` or `props = {}` on a markdown widget.
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

-> **NOTE:** The widget-level `props` key must be **absent entirely** on `viz.markdown` widgets. Do not include `"props": {}` or `props = {}` on a markdown widget — omit the key altogether.

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

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`). Example: `props = { title = "My chart" }` or `"props": {"title": ""}`.

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
  * `hiddenColumns[].columnName` - (Optional) Columns to hide. Accepts the raw column name from the query (e.g. `"timestamp"`) or its display label.
  * `thresholds[].columnName` - (Optional) Column to apply the threshold to.
  * `thresholds[].from` / `.to` / `.severity` - (Optional) See [Nested `thresholds` blocks (table)](#nested-thresholds-blocks-table).
  * `dataFormatters` - (Optional) Custom data format configuration.
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.billboard` — Billboard

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

Displays a single large metric value with color-coded threshold ranges.

  * `nrqlQueries` - (Required) Query should return a single aggregated value or a small number of faceted values.
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

Displays a value as a dial, ring, or horizontal bar gauge.

  * `nrqlQueries` - (Required)
  * `sqlQueries` - (Optional)
  * `gaugeSettings.displayMode` - (Optional) Gauge style. Valid values: `"arc"` (default), `"circular"`, `"bar"`.
  * `gaugeSettings.min` / `gaugeSettings.max` - (Optional) Scale bounds. Default max: `100`.
  * `gaugeSettings.display` - (Optional) What to show inside the gauge. Valid values: `"auto"` (default), `"all"`, `"value"`, `"name"`, `"none"`.
  * `gaugeSettings.showThresholdMarkers` - (Optional) Show the threshold color ring/strip. Default: `true`.
  * `gaugeSettings.showThresholdLabels` - (Optional) Show value labels at threshold boundaries. Default: `false`.
  * `gaugeSettings.thresholdsColorMode` - (Optional) Valid values: `"solid"` (default), `"gradient"`.
  * `thresholds.thresholds[].from` / `.severity` - (Optional) Threshold zones.
  * `thresholds.seriesOverrides[].seriesName` - (Optional) Per-series threshold overrides.
  * `units.unit` / `units.seriesOverrides[]` - (Optional)
  * `dataFormatters` - (Optional)
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.histogram` — Histogram

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

Requires the `histogram()` NRQL function.

  * `nrqlQueries` - (Required) Query must use `histogram(attr)` or `histogram(attr, width: N, buckets: N)`.
  * `legend.enabled` / `legend.position` - (Optional)
  * `colors.colorPalette` / `colors.seriesOverrides[]` - (Optional)
  * `chartStyles.gradient.enabled` - (Optional)
  * `thresholds.isLabelVisible` / `thresholds.thresholds[]` - (Optional)
  * `yAxisLeft.min` / `.max` / `.zero` / `.scale` - (Optional)
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.heatmap` — Heatmap

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

Requires `histogram()` with a `FACET`.

  * `nrqlQueries` - (Required)
  * `facet.showOtherSeries` - (Optional) Default: `false`.
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.scatter` — Scatter plot

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

  * `nrqlQueries` - (Required) Query should return two numeric values and a `FACET`.
  * `facet.showOtherSeries` - (Optional)
  * `legend.enabled` / `legend.position` - (Optional)
  * `colors.colorPalette` / `colors.seriesOverrides[]` - (Optional)
  * `units.unit` / `units.seriesOverrides[]` - (Optional)
  * `nullValues.nullValue` / `nullValues.seriesOverrides[]` - (Optional) Values: `"default"`, `"zero"`, `"preserve"`, `"remove"`.
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

Requires the `funnel()` NRQL function.

  * `nrqlQueries` - (Required)
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

Renders the raw JSON output of the NRQL query. Useful for debugging.

  * `nrqlQueries` - (Required)
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.sparkline` — Sparkline

-> **NOTE:** The widget-level `props` must be present and must include `title` (can be `""`).

Compact time-series line chart with Y-axis labels.

  * `nrqlQueries` - (Required)
  * `facet.showOtherSeries` - (Optional)
  * `colors.colorPalette` / `colors.seriesOverrides[]` - (Optional)
  * `units.unit` / `units.seriesOverrides[]` - (Optional)
  * `nullValues.nullValue` / `nullValues.seriesOverrides[]` - (Optional) Values: `"default"`, `"zero"`, `"preserve"`, `"remove"`.
  * `chartStyles.lineInterpolation` - (Optional)
  * `tooltip.mode` - (Optional)
  * `yAxisLeft.min` / `.max` / `.zero` - (Optional)
  * `platformOptions.ignoreTimeRange` - (Optional)
  * `refreshRate.frequency` - (Optional)

---

### `viz.sparkline-lite` — Sparkline Lite

Ultra-compact sparkline with no axis labels. Supports the same props as `viz.sparkline` except `yAxisLeft`.

---

## Nested `nrqlQueries` blocks

All query-based visualizations accept a `nrqlQueries` array. Each item supports:

  * `accountIds` - (Required) An array of one or more New Relic account IDs, e.g. `[1234567]`. For cross-account queries: `[1234567, 7654321]`.
  * `query` - (Required) A valid NRQL query string.
  * `offset` - (Optional) Number of milliseconds to offset the query time window.

-> **NOTE:** If you query an account to which you do not have access, the widget renders with a "data inaccessible" message. No Terraform error is raised.

Most visualizations also accept `sqlQueries` as an alternative data source for Federated Data Source (FDS) connections. Each SQL query object requires `query` (string) and `accountId` (number) or `accountIds` (number[]).

## Nested `thresholds` blocks (line, area, stacked-bar)

The `thresholds.thresholds` array defines horizontal bands drawn on the chart:

  * `name` - (Optional) A label for the threshold.
  * `from` - (Optional) Lower bound of the threshold range (inclusive).
  * `to` - (Optional) Upper bound of the threshold range (inclusive).
  * `severity` - (Required) Color applied to data within this range. Valid values: `"success"`, `"warning"`, `"severe"`, `"critical"`, `"unavailable"`.
  * `isLabelVisible` - (Optional) Always show the threshold label. Specified at the parent `thresholds` level, not per-threshold.

## Nested `thresholds` blocks (table)

The `thresholds` array for `viz.table` applies coloring per column:

  * `columnName` - (Required) The column to apply the threshold to.
  * `from` - (Optional) Lower bound.
  * `to` - (Optional) Upper bound.
  * `severity` - (Required) Valid values: `"success"`, `"warning"`, `"severe"`, `"critical"`, `"unavailable"`.

## Nested `thresholdsWithSeriesOverrides` blocks (billboard)

  * `thresholds[].from` / `.to` / `.severity` - (Optional) Global threshold ranges applied to all series.
  * `seriesOverrides[].seriesName` / `.from` / `.to` / `.severity` - (Optional) Per-series threshold overrides for multi-value billboards.

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

### Weekly service health review (content\_json from file)

```hcl
resource "newrelic_notebook" "weekly_review" {
  title        = "Weekly Service Health Review"
  content_json = file("${path.module}/notebooks/weekly-health.json")
}
```

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
          "content": {
            "type": "visualization",
            "id": "viz.markdown",
            "props": { "text": "# Weekly Service Health\n\n**Period**: Last 7 days vs. prior week" }
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
                { "accountIds": [1234567], "query": "SELECT apdex(duration, 0.5) AS 'Apdex' FROM Transaction WHERE appName IN ('web-frontend', 'api-backend') SINCE 7 days ago" }
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
          "props": { "title": "P95 latency by service" },
          "content": {
            "type": "visualization",
            "id": "viz.line",
            "props": {
              "nrqlQueries": [
                { "accountIds": [1234567], "query": "SELECT percentile(duration, 95) FROM Transaction WHERE appName IN ('web-frontend', 'api-backend') FACET appName TIMESERIES 1 day SINCE 7 days ago COMPARE WITH 1 week ago" }
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
                { "accountIds": [1234567], "query": "SELECT count(*) FROM Transaction WHERE appName IN ('web-frontend', 'api-backend') FACET appName TIMESERIES 1 hour SINCE 7 days ago" }
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
                { "accountIds": [1234567], "query": "SELECT count(*) FROM Transaction WHERE appName IN ('web-frontend', 'api-backend') FACET httpResponseCode TIMESERIES 1 hour SINCE 7 days ago" }
              ],
              "legend": { "enabled": true }
            }
          }
        },
        {
          "type": "widget",
          "props": { "title": "Error distribution" },
          "content": {
            "type": "visualization",
            "id": "viz.pie",
            "props": {
              "nrqlQueries": [
                { "accountIds": [1234567], "query": "SELECT count(*) FROM Transaction WHERE error IS true FACET httpResponseCode SINCE 7 days ago" }
              ],
              "facet": { "showOtherSeries": true }
            }
          }
        },
        {
          "type": "widget",
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

### Using `content = jsonencode({...})` (HCL authoring)

The `content` attribute with `jsonencode` is an alternative to `content_json`. It produces more granular field-level plan output but results in more verbose Terraform configuration. Use it when authoring notebooks entirely in Terraform without an external JSON file.

```hcl
resource "newrelic_notebook" "service_health" {
  title = "Service Health - HCL Example"

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
            props = { text = "# Service Health\n\nKey metrics for the checkout service." }
          }
        },
        {
          type  = "widget"
          props = { title = "Throughput" }
          content = {
            type = "visualization"
            id   = "viz.line"
            props = {
              nrqlQueries = [{
                accountIds = [var.account_id]
                query      = "SELECT rate(count(*), 1 minute) FROM Transaction WHERE appName = 'checkout' FACET appName TIMESERIES 5 minutes SINCE 3 hours ago"
              }]
              legend = { enabled = true }
            }
          }
        }
      ]
    }]
  })
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

  content_json = jsonencode({
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

### Infrastructure metrics notebook

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
          "content": {
            "type": "visualization",
            "id": "viz.markdown",
            "props": { "text": "# Infrastructure Metrics\n\nHost-level golden signals." }
          }
        },
        {
          "type": "widget",
          "props": { "title": "Hosts reporting" },
          "content": {
            "type": "visualization",
            "id": "viz.billboard",
            "props": {
              "nrqlQueries": [{ "accountIds": [1234567], "query": "SELECT uniqueCount(hostname) AS 'Hosts' FROM SystemSample SINCE 5 minutes ago" }]
            }
          }
        },
        {
          "type": "widget",
          "props": { "title": "CPU % by host" },
          "content": {
            "type": "visualization",
            "id": "viz.line",
            "props": {
              "nrqlQueries": [{ "accountIds": [1234567], "query": "SELECT average(cpuPercent) FROM SystemSample FACET hostname TIMESERIES 5 minutes SINCE 3 hours ago LIMIT 10" }],
              "yAxisLeft": { "zero": true, "min": 0, "max": 100 }
            }
          }
        },
        {
          "type": "widget",
          "props": { "title": "Memory %" },
          "content": {
            "type": "visualization",
            "id": "viz.area",
            "props": {
              "nrqlQueries": [{ "accountIds": [1234567], "query": "SELECT average(memoryUsedPercent) FROM SystemSample FACET hostname TIMESERIES 5 minutes SINCE 3 hours ago LIMIT 10" }]
            }
          }
        },
        {
          "type": "widget",
          "props": { "title": "Disk space by mount" },
          "content": {
            "type": "visualization",
            "id": "viz.pie",
            "props": {
              "nrqlQueries": [{ "accountIds": [1234567], "query": "SELECT latest(diskUsedPercent) FROM StorageSample FACET mountPoint SINCE 10 minutes ago" }]
            }
          }
        }
      ]
    }
  ]
}
```

```hcl
resource "newrelic_notebook" "infra_metrics" {
  title        = "Infrastructure Metrics"
  content_json = file("${path.module}/notebooks/infra-metrics.json")
}
```
