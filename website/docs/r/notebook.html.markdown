---
layout: "newrelic"
page_title: "New Relic: newrelic_notebook"
sidebar_current: "docs-newrelic-resource-notebook"
description: |-
  Create and manage New Relic Notebooks.
---

# Resource: newrelic\_notebook

Use this resource to create and manage [New Relic Notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/introduction-notebooks/).

Notebooks are shareable documents that combine live NRQL queries, visualizations, and Markdown narrative in a single view - the kind of thing you'd reach for when investigating an incident, reviewing weekly service health, or building an onboarding guide for a new team member. This resource manages both the notebook title and its full body content using the same declarative UI format the New Relic platform uses natively.

-> **NOTE:** The notebook body must be a valid declarative UI document. See [Content Format and Schema](#content-format-and-schema) for the required envelope structure and supported visualization types.

See [Examples](#examples) at the bottom of this page for ready-to-use notebook templates.

## Example Usage

The recommended approach is to supply the notebook body as a raw JSON string. You can use a `file()` reference for larger notebooks, or an inline heredoc for smaller ones.

```hcl
resource "newrelic_notebook" "incident_runbook" {
  title = "Production API Incident - Investigation"

  content = <<-JSON
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
                "props": {
                  "text": "# Production API Incident\n\n**Status**: Resolved\n**Impact**: ~15% of API requests failed\n\nUse this notebook to walk through what happened and track action items."
                }
              }
            },
            {
              "type": "widget",
              "props": { "title": "Peak error rate" },
              "content": {
                "type": "visualization",
                "id": "viz.billboard",
                "props": {
                  "nrqlQueries": [
                    {
                      "accountIds": [1234567],
                      "query": "SELECT percentage(count(*), WHERE httpResponseCode >= 400) AS 'Error Rate %' FROM Transaction WHERE appName = 'api-production' SINCE 1 hour ago"
                    }
                  ],
                  "thresholdsWithSeriesOverrides": {
                    "thresholds": [
                      { "to": 1,   "severity": "success"  },
                      { "from": 1, "to": 5, "severity": "warning"  },
                      { "from": 5, "severity": "critical" }
                    ]
                  }
                }
              }
            },
            {
              "type": "widget",
              "props": { "title": "Error rate over time" },
              "content": {
                "type": "visualization",
                "id": "viz.line",
                "props": {
                  "nrqlQueries": [
                    {
                      "accountIds": [1234567],
                      "query": "SELECT percentage(count(*), WHERE httpResponseCode >= 400) FROM Transaction WHERE appName = 'api-production' TIMESERIES 1 minute SINCE 1 hour ago"
                    }
                  ],
                  "legend":    { "enabled": true },
                  "yAxisLeft": { "zero": true },
                  "nullValues": { "nullValue": "zero" }
                }
              }
            }
          ]
        }
      ]
    }
  JSON
}
```

## Argument Reference

The following arguments are supported:

  * `title` - (Required) The title of the notebook. Must be unique within the organization.
  * `content` - (Required) The notebook body as a JSON string. Accepts the raw JSON exported from the New Relic UI, a `file()` reference, or a `jsonencode({...})` expression. Produces line-level diffs of the normalized JSON in `terraform plan`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `guid` - The unique entity identifier (GUID) of the notebook in New Relic.
  * `blob_id` - The blob identifier of the current notebook content, updated after each Terraform-managed write.
  * `organization_id` - The New Relic organization ID the notebook belongs to. Resolved automatically from the provider credentials.

---

## Content Format and Schema

A notebook body is a JSON document with a fixed three-level structure. The first two levels are always written the same way. The third level is where all your widgets go.

-> **WARNING:** The `//` annotations in the snippet below are for illustration only and are **not valid JSON**. Do not include them in your actual notebook body.

```json
{
  // Level 1 - always fixed
  // This structure is identical across every notebook. Do not modify these fields.

  "type": "declarative",
  "version": 1,
  "content": [
    {
      // Level 2 - always fixed
      // This structure is identical across every notebook. Do not modify these fields.

      "type": "container",
      "props": { "layout": "stack" },
      "content": [

        // Level 3 - your widgets go here. Add as many as you need.

        {
          "type": "widget",
          "content": {
            "type": "visualization",
            "id": "viz.markdown",
            "props": { "text": "## Section header\n\nAdd narrative, context, or action items here." }
          }
        },
        {
          "type": "widget",
          "props": { "title": "Error rate" },
          "content": {
            "type": "visualization",
            "id": "viz.billboard",
            "props": { ... }
          }
        },
        {
          "type": "widget",
          "props": { "title": "Throughput over time" },
          "content": {
            "type": "visualization",
            "id": "viz.line",
            "props": { ... }
          }
        }
      ]
    }
  ]
}
```

### Understanding the structure

**Levels 1 and 2 are always identical.** You always write them exactly as shown above - they cannot be customized beyond what is described here.

  * **Level 1 - Document envelope (fixed)**
    * `type` must always be `"declarative"`.
    * `version` must always be the integer `1`.
    * `content` is an array containing exactly **one** container object (Level 2). The New Relic Notebooks UI only renders the first container - place all widgets inside it.

  * **Level 2 - Container (fixed)**
    * `type` must always be `"container"`.
    * `props` must always be `{ "layout": "stack" }`.
    * `content` is the array where your widgets live. This is the only attribute in Level 2 that varies - add as many widget objects as you need.

**Level 3 is where you do your work.** Every object in the Level 2 `content` array is a widget. Each widget has three attributes:

  * **Level 3 - Widgets (customizable)**
    * `type` must always be `"widget"`.
    * `props` *(widget-level, controls the title label above the chart)* - The rules differ by chart type. Note: this `props` is **not** the same as `content.props` below - it only controls the widget's display title, not the chart configuration.
      * **`viz.markdown` only:** `props` must be **absent entirely** from the widget object. Including even an empty `"props": {}` will fail validation at `terraform plan` time.
      * **All other chart types:** `props` must be present and must include a `title` key. An empty string `""` is valid as a title. Example: `"props": { "title": "Error rate" }`.
    * `content` - The visualization to render. Three sub-fields are always required:
      * `type` must be `"visualization"`.
      * `id` identifies the chart type - for example `"viz.markdown"`, `"viz.line"`, `"viz.area"`. See [Supported visualization types](#supported-visualization-types) for the full list.
      * `props` *(content-level, chart configuration)* - Holds the chart-specific configuration such as NRQL queries, axis settings, and thresholds. This is **separate** from the widget-level `props` described above - these props define how the chart renders its data, not the display title.

**Putting it together - a quick reference:**

The following summarizes the steps to create a notebook that follows the required structure and passes validation without errors.

1. Start with Level 1: `"type": "declarative"`, `"version": 1`, and a `"content"` array with one object.
2. Inside that object, add Level 2: `"type": "container"`, `"props": { "layout": "stack" }`, and a `"content"` array.
3. Inside the Level 2 `"content"` array, add as many Level 3 widget objects as you need. Each widget must have:
   - `"type": "widget"` - always.
   - `"props": { "title": "..." }` *(widget-level - display title only, not chart config)* - Required for all chart types. **Omit entirely** for `viz.markdown` widgets.
   - `"content"` containing `"type": "visualization"`, `"id": "<chart-type>"`, and `"props": { ... }` *(content-level - chart configuration such as queries, thresholds, axis options)*.

---

## Supported Visualization Types

The `id` field inside a widget's `content` object identifies which chart type to render. The following visualization types are supported:

### Visualization type reference

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

### Widget props reference

Each visualization type below lists its supported `props`. These go inside the widget's `content.props` object.

---

### `viz.markdown` - Markdown

Renders a text block supporting GitHub-flavoured Markdown. Use it to add context, section headers, runbook steps, or action-item checklists between your charts.

-> **NOTE:** The `props` key must be **absent entirely** from the widget object for `viz.markdown`. Including `"props": {}` will fail validation. See [Understanding the structure](#understanding-the-structure) for both a code example and an explanation of why this restriction exists.

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

### `viz.line` - Line chart

The workhorse for time-series data. Use it when you want to show how a metric changes over time and need full control over axes, thresholds, and per-series styling.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `alertQueries` - (Optional) Alert violation objects to overlay warning/critical bands on the chart. Only available for `viz.line`. Each object supports:
    * `accountIds` - (Required) Account IDs the alert query runs against.
    * `violationId` - (Required) The violation ID to overlay.
    * `duration` - (Required) Duration in milliseconds.
    * `endTime` - (Required) End time as an epoch millisecond timestamp.
  * `legend` - (Optional) Controls the chart legend.
    * `enabled` - (Optional) Show or hide the legend. Defaults to `true`.
    * `position` - (Optional) Where to place the legend. Valid values are `"bottom"` (default), `"left"`, or `"right"`.
  * `facet` - (Optional) Controls FACET grouping behaviour.
    * `showOtherSeries` - (Optional) Show the "Other" group for `FACET` queries. Defaults to `false`.
  * `yAxisLeft` - (Optional) Configuration for the left Y axis.
    * `min` - (Optional) Fixed minimum value for the axis.
    * `max` - (Optional) Fixed maximum value for the axis.
    * `zero` - (Optional) Force zero as the axis origin. Defaults to `true`.
    * `scale` - (Optional) Axis scale. Valid values are `"linear"` (default) and `"logarithmic"`.
  * `yAxisRight` - (Optional) Configuration for the right Y axis. Supports the same sub-keys as `yAxisLeft`, plus:
    * `series` - (Optional) List of series names to plot against the right axis.
      * `name` - (Required) The series name to bind to the right axis.
  * `nullValues` - (Optional) Controls how null data points are rendered.
    * `nullValue` - (Optional) Default null handling. Valid values: `"default"`, `"zero"`, `"preserve"`, `"remove"`.
    * `seriesOverrides` - (Optional) Per-series null value overrides.
      * `seriesName` - (Required) The series to override.
      * `nullValue` - (Required) The null handling for that series.
  * `colors` - (Optional) Controls chart coloring.
    * `colorPalette` - (Optional) Color palette. Valid values are `"consistent"` (default) and `"dynamic"`.
    * `seriesOverrides` - (Optional) Per-series color overrides.
      * `seriesName` - (Required) The series to override.
      * `color` - (Required) RGB hex color, e.g. `"#FF0000"`.
  * `units` - (Optional) Data units for display formatting.
    * `unit` - (Optional) Default unit for all series. See [Valid `units.unit` values](#valid-unitsunit-values).
    * `seriesOverrides` - (Optional) Per-series unit overrides.
      * `seriesName` - (Required) The series to override.
      * `unit` - (Required) The unit for that series.
  * `thresholds` - (Optional) Horizontal threshold bands drawn across the chart. See [Nested `thresholds` blocks (line, area, stacked-bar)](#nested-thresholds-blocks-line-area-stacked-bar).
    * `isLabelVisible` - (Optional) Show threshold labels on the chart. Defaults to `false`.
    * `thresholds` - (Optional) Array of threshold band definitions.
      * `name` - (Optional) A label for the threshold band.
      * `from` - (Optional) Lower bound of the threshold range (inclusive).
      * `to` - (Optional) Upper bound of the threshold range (inclusive).
      * `severity` - (Required) Color coding. Valid values: `"success"`, `"warning"`, `"severe"`, `"critical"`, `"unavailable"`.
  * `chartStyles` - (Optional) Visual style overrides for the chart.
    * `lineInterpolation` - (Optional) Line interpolation mode. Valid values: `"linear"` (default), `"smooth"`, `"stepBefore"`, `"stepAfter"`.
    * `gradient` - (Optional) Gradient fill under the line.
      * `enabled` - (Optional) Enable gradient fill. Defaults to `false`.
  * `chartTypes` - (Optional) Override the chart type per series.
    * `seriesOverrides` - (Optional) Per-series chart type overrides.
      * `seriesName` - (Required) The series to override.
      * `chartType` - (Required) Valid values: `"line"`, `"area"`.
  * `tooltip` - (Optional) Tooltip display behaviour.
    * `mode` - (Optional) Valid values: `"single"` (default), `"all"`, `"hidden"`.
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Throughput by Application (requests / min)"
  },
  "content": {
    "type": "visualization",
    "id": "viz.line",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT rate(count(*), 1 minute) AS 'rpm' FROM Transaction FACET appName TIMESERIES 5 minutes SINCE 3 hours ago"
        }
      ],
      "legend": {
        "enabled": true,
        "position": "bottom"
      },
      "yAxisLeft": {
        "min": 0,
        "zero": true
      },
      "nullValues": {
        "nullValue": "zero"
      },
      "chartStyles": {
        "lineInterpolation": "smooth"
      },
      "tooltip": {
        "mode": "all"
      },
      "colors": {
        "colorPalette": "consistent"
      },
      "platformOptions": {
        "ignoreTimeRange": false
      },
      "refreshRate": {
        "frequency": 60000
      }
    }
  }
}
```

---

### `viz.area` - Area chart

Great for showing cumulative totals or filled time-series where the area under the line carries meaning.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

Supports the same props as `viz.line`, with the following differences:

  * `alertQueries` and `yAxisRight` are not available.
  * `nullValues` supports `"default"`, `"zero"`, `"preserve"` (not `"remove"`).
  * Additional prop: `chartStyles`
    * `stacked` - (Optional) Stack multiple series on top of each other.
      * `enabled` - (Optional) Enable stacked layout. Defaults to `true`.

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Response Time Percentiles (ms)"
  },
  "content": {
    "type": "visualization",
    "id": "viz.area",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT percentile(duration * 1000, 50) AS 'P50 ms', percentile(duration * 1000, 95) AS 'P95 ms', percentile(duration * 1000, 99) AS 'P99 ms' FROM Transaction TIMESERIES 5 minutes SINCE 3 hours ago"
        }
      ],
      "legend": {
        "enabled": true,
        "position": "bottom"
      },
      "yAxisLeft": {
        "zero": true
      },
      "nullValues": {
        "nullValue": "zero"
      },
      "chartStyles": {
        "stacked": {
          "enabled": false
        }
      },
      "colors": {
        "seriesOverrides": [
          {
            "seriesName": "P50 ms",
            "color": "#11A600"
          },
          {
            "seriesName": "P95 ms",
            "color": "#FFB951"
          },
          {
            "seriesName": "P99 ms",
            "color": "#BF0016"
          }
        ]
      },
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.stacked-bar` - Bar chart (stacked / timeseries)

Use this when you want to compare part-to-whole relationships over time - for example, traffic broken down by endpoint or host.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

Supports all `viz.line` props except `alertQueries`, `yAxisRight`, and the `"remove"` null value. Additional prop:

  * `chartStyles` - (Optional) Visual style configuration.
    * `stacked` - (Optional) Stack configuration.
      * `enabled` - (Optional) Enable stacked layout. Defaults to `true`.

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Request Volume by App (stacked)"
  },
  "content": {
    "type": "visualization",
    "id": "viz.stacked-bar",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT count(*) FROM Transaction FACET appName TIMESERIES 10 minutes SINCE 3 hours ago"
        }
      ],
      "legend": {
        "enabled": true
      },
      "yAxisLeft": {
        "zero": true
      },
      "nullValues": {
        "nullValue": "zero"
      },
      "chartStyles": {
        "stacked": {
          "enabled": true
        }
      },
      "thresholds": {
        "isLabelVisible": true,
        "thresholds": [
          {
            "name": "High traffic",
            "from": 5000,
            "severity": "warning"
          }
        ]
      },
      "platformOptions": {
        "ignoreTimeRange": false
      },
      "refreshRate": {
        "frequency": 60000
      }
    }
  }
}
```

---

### `viz.bar` - Bar chart (simplified)

A simpler bar chart for categorical FACET comparisons - when you want counts or averages across a dimension and a timeseries axis isn't needed.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `facet` - (Optional) Controls FACET grouping behaviour.
    * `showOtherSeries` - (Optional) Show the "Other" group for `FACET` queries. Defaults to `false`.
  * `colors` - (Optional) Controls chart coloring.
    * `colorPalette` - (Optional) Color palette. Valid values are `"consistent"` (default) and `"dynamic"`.
    * `seriesOverrides` - (Optional) Per-series color overrides.
      * `seriesName` - (Required) The series to override.
      * `color` - (Required) RGB hex color, e.g. `"#FF0000"`.
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Top Transaction Types by Count"
  },
  "content": {
    "type": "visualization",
    "id": "viz.bar",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT count(*) FROM Transaction FACET name SINCE 1 hour ago LIMIT 10"
        }
      ],
      "facet": {
        "showOtherSeries": false
      },
      "colors": {
        "colorPalette": "consistent"
      },
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.stacked-horizontal-bar` - Horizontal stacked bar

Useful for comparing proportions across categories when your label text is long and reads better horizontally.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `facet` - (Optional) Controls FACET grouping behaviour.
    * `showOtherSeries` - (Optional) Show the "Other" group for `FACET` queries. Defaults to `false`.
  * `legend` - (Optional) Controls the chart legend.
    * `enabled` - (Optional) Show or hide the legend. Defaults to `true`.
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Log Volume by Level"
  },
  "content": {
    "type": "visualization",
    "id": "viz.stacked-horizontal-bar",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT count(*) AS 'Log Lines' FROM Log FACET level SINCE 3 hours ago"
        }
      ],
      "legend": {
        "enabled": true
      },
      "facet": {
        "showOtherSeries": false
      },
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.pie` - Pie chart

Good for showing how a whole is divided among a small number of categories. If you have more than 7-8 slices, a bar chart is usually clearer.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `facet` - (Optional) Controls FACET grouping behaviour.
    * `showOtherSeries` - (Optional) Show the "Other" group for `FACET` queries. Defaults to `true` (differs from other chart types).
  * `legend` - (Optional) Controls the chart legend.
    * `enabled` - (Optional) Show or hide the legend. Defaults to `true`.
  * `colors` - (Optional) Controls chart coloring.
    * `colorPalette` - (Optional) Color palette. Valid values are `"consistent"` (default) and `"dynamic"`.
    * `seriesOverrides` - (Optional) Per-series color overrides.
      * `seriesName` - (Required) The series to override.
      * `color` - (Required) RGB hex color, e.g. `"#FF0000"`.
  * `chartStyles` - (Optional) Visual style configuration.
    * `gradient` - (Optional) Gradient fill for pie slices.
      * `enabled` - (Optional) Enable gradient. Defaults to `false`.
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Disk Utilisation by Host"
  },
  "content": {
    "type": "visualization",
    "id": "viz.pie",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT average(diskUsedPercent) AS 'Disk Used %' FROM StorageSample FACET hostname SINCE 30 minutes ago"
        }
      ],
      "legend": {
        "enabled": true
      },
      "facet": {
        "showOtherSeries": false
      },
      "colors": {
        "colorPalette": "consistent"
      },
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.table` - Table

The right choice when you need to show raw rows, multi-column comparisons, or sortable ranked lists. Pairs well with `initialSorting` and column-level thresholds.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `facet` - (Optional) Controls FACET grouping behaviour.
    * `showOtherSeries` - (Optional) Show the "Other" group for `FACET` queries. Defaults to `false`.
  * `initialSorting` - (Optional) Default sort applied when the table first renders.
    * `name` - (Required) Column name to sort by.
    * `direction` - (Required) Sort direction. Valid values: `"asc"`, `"desc"`.
  * `hiddenColumns` - (Optional) Columns to suppress from the rendered table.
    * `columnName` - (Required) The column name to hide, as it appears in the query results.
  * `thresholds` - (Optional) Column-level color thresholds. See [Nested `thresholds` blocks (table)](#nested-thresholds-blocks-table). Each entry supports:
    * `columnName` - (Required) The column to apply the threshold to.
    * `from` - (Optional) Lower bound of the range.
    * `to` - (Optional) Upper bound of the range.
    * `severity` - (Required) Valid values: `"success"`, `"warning"`, `"severe"`, `"critical"`, `"unavailable"`.
  * `dataFormatters` - (Optional) Custom per-column data format configuration.
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Host Health - CPU, Memory and Disk"
  },
  "content": {
    "type": "visualization",
    "id": "viz.table",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT average(cpuPercent) AS 'Avg CPU %', max(cpuPercent) AS 'Max CPU %', average(memoryUsedPercent) AS 'Avg Mem %', average(diskUsedPercent) AS 'Disk Used %' FROM SystemSample FACET hostname SINCE 30 minutes ago LIMIT 20"
        }
      ],
      "initialSorting": {
        "name": "Avg CPU %",
        "direction": "desc"
      },
      "thresholds": [
        {
          "columnName": "Avg CPU %",
          "from": 80,
          "severity": "critical"
        },
        {
          "columnName": "Avg CPU %",
          "from": 60,
          "to": 80,
          "severity": "warning"
        },
        {
          "columnName": "Avg Mem %",
          "from": 90,
          "severity": "critical"
        },
        {
          "columnName": "Disk Used %",
          "from": 85,
          "severity": "critical"
        }
      ],
      "facet": {
        "showOtherSeries": false
      },
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.billboard` - Billboard

Displays a single large metric value with color-coded threshold ranges. Perfect for at-a-glance status tiles at the top of a notebook.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `facet` - (Optional) Controls FACET grouping behaviour.
    * `showOtherSeries` - (Optional) Show the "Other" group for `FACET` queries. Defaults to `false`.
  * `thresholdsWithSeriesOverrides` - (Optional) Threshold configuration for the billboard. See [Nested `thresholdsWithSeriesOverrides` blocks (billboard)](#nested-thresholdswithseriesoverrides-blocks-billboard).
    * `thresholds` - (Optional) Global threshold ranges applied to all series.
      * `from` - (Optional) Lower bound of the range.
      * `to` - (Optional) Upper bound of the range.
      * `severity` - (Required) Valid values: `"success"`, `"warning"`, `"severe"`, `"critical"`, `"unavailable"`.
    * `seriesOverrides` - (Optional) Per-series threshold overrides for multi-value billboards.
      * `seriesName` - (Required) The series to override.
      * `from` - (Optional) Lower bound.
      * `to` - (Optional) Upper bound.
      * `severity` - (Required) Severity level for this series.
  * `billboardSettings` - (Optional) Fine-grained layout and linking options.
    * `visual` - (Optional) Controls the visual presentation of the billboard.
      * `alignment` - (Optional) Alignment of values on the billboard. Valid values: `"auto"`, `"stacked"`, `"inline"`.
      * `display` - (Optional) What to show on the billboard. Valid values: `"auto"`, `"all"`, `"value"`, `"label"`, `"none"`.
    * `gridOptions` - (Optional) Grid layout for multi-value billboards.
      * `columns` - (Optional) Number of columns in the grid.
      * `label` - (Optional) Label display size configuration.
      * `value` - (Optional) Value display size configuration.
    * `link` - (Optional) A clickable link attached to the billboard.
      * `url` - (Required) The destination URL.
      * `title` - (Optional) Link text shown on the billboard.
      * `newTab` - (Optional) Open the link in a new tab. Defaults to `false`.
  * `units` - (Optional) Data units for display formatting.
    * `unit` - (Optional) Default unit for all series. See [Valid `units.unit` values](#valid-unitsunit-values).
    * `seriesOverrides` - (Optional) Per-series unit overrides.
      * `seriesName` - (Required) The series to override.
      * `unit` - (Required) The unit for that series.
  * `dataFormatters` - (Optional) Custom per-column data format configuration.
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Request Volume and Error Rate"
  },
  "content": {
    "type": "visualization",
    "id": "viz.billboard",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT rate(count(*), 1 minute) AS 'Requests / min', percentage(count(*), WHERE error IS true) AS 'Error Rate %', count(*) AS 'Total Requests' FROM Transaction SINCE 1 hour ago"
        }
      ],
      "thresholdsWithSeriesOverrides": {
        "thresholds": [
          {
            "to": 1,
            "severity": "success"
          },
          {
            "from": 1,
            "to": 5,
            "severity": "warning"
          },
          {
            "from": 5,
            "severity": "critical"
          }
        ],
        "seriesOverrides": [
          {
            "seriesName": "Error Rate %",
            "from": 1,
            "severity": "critical"
          }
        ]
      },
      "billboardSettings": {
        "visual": {
          "alignment": "inline",
          "display": "all"
        }
      },
      "facet": {
        "showOtherSeries": false
      },
      "platformOptions": {
        "ignoreTimeRange": false
      },
      "refreshRate": {
        "frequency": 60000
      }
    }
  }
}
```

---

### `viz.gauge` - Gauge

Shows a single value against a min/max scale, making it easy to communicate how close you are to a limit or target.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `colors` - (Optional) Controls chart coloring.
    * `colorPalette` - (Optional) Color palette. Valid values are `"consistent"` (default) and `"dynamic"`.
    * `seriesOverrides` - (Optional) Per-series color overrides.
      * `seriesName` - (Required) The series to override.
      * `color` - (Required) RGB hex color, e.g. `"#FF0000"`.
  * `facet` - (Optional) Controls FACET grouping behaviour.
    * `showOtherSeries` - (Optional) Show the "Other" group for `FACET` queries. Defaults to `false`.
  * `gaugeSettings` - (Optional) Appearance and scale options for the gauge.
    * `displayMode` - (Optional) Gauge shape. Valid values: `"arc"` (default), `"circular"`, `"bar"`.
    * `min` - (Optional) Lower bound of the scale.
    * `max` - (Optional) Upper bound of the scale. Defaults to `100`.
    * `display` - (Optional) What to render inside the gauge. Valid values: `"auto"` (default), `"all"`, `"value"`, `"name"`, `"none"`.
    * `showThresholdMarkers` - (Optional) Show tick marks at threshold boundaries. Defaults to `true`.
    * `showThresholdLabels` - (Optional) Show labels at threshold boundaries. Defaults to `false`.
    * `thresholdsColorMode` - (Optional) How threshold colors are applied. Valid values: `"solid"` (default), `"gradient"`.
  * `thresholds` - (Optional) Threshold zones displayed on the gauge scale.
    * `thresholds` - (Optional) Array of threshold zone entries.
      * `from` - (Optional) Lower bound of the zone.
      * `severity` - (Required) Severity level for this zone.
  * `units` - (Optional) Data units for display formatting.
    * `unit` - (Optional) Default unit for all series. See [Valid `units.unit` values](#valid-unitsunit-values).
    * `seriesOverrides` - (Optional) Per-series unit overrides.
      * `seriesName` - (Required) The series to override.
      * `unit` - (Required) The unit for that series.
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Memory Utilisation (%)"
  },
  "content": {
    "type": "visualization",
    "id": "viz.gauge",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT average(memoryUsedPercent) AS 'Memory %' FROM SystemSample SINCE 30 minutes ago"
        }
      ],
      "gaugeSettings": {
        "displayMode": "arc",
        "min": 0,
        "max": 100,
        "display": "auto",
        "showThresholdMarkers": true,
        "showThresholdLabels": false,
        "thresholdsColorMode": "solid"
      },
      "thresholds": {
        "thresholds": [
          {
            "from": 0,
            "severity": "success"
          },
          {
            "from": 60,
            "severity": "warning"
          },
          {
            "from": 80,
            "severity": "critical"
          }
        ]
      },
      "colors": {
        "colorPalette": "consistent"
      },
      "units": {
        "unit": "PERCENTAGE"
      },
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.histogram` - Histogram

Visualizes the distribution of a numeric attribute across buckets. Your query must use the `histogram()` NRQL function.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. Query must use `histogram(attr)` or `histogram(attr, width: N, buckets: N)`. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `legend` - (Optional) Controls the chart legend.
    * `enabled` - (Optional) Show or hide the legend. Defaults to `true`.
    * `position` - (Optional) Where to place the legend. Valid values are `"bottom"` (default), `"left"`, or `"right"`.
  * `colors` - (Optional) Controls chart coloring.
    * `colorPalette` - (Optional) Color palette. Valid values are `"consistent"` (default) and `"dynamic"`.
    * `seriesOverrides` - (Optional) Per-series color overrides.
      * `seriesName` - (Required) The series to override.
      * `color` - (Required) RGB hex color, e.g. `"#FF0000"`.
  * `yAxisLeft` - (Optional) Configuration for the left Y axis.
    * `min` - (Optional) Fixed minimum value for the axis.
    * `max` - (Optional) Fixed maximum value for the axis.
    * `zero` - (Optional) Force zero as the axis origin. Defaults to `true`.
    * `scale` - (Optional) Axis scale. Valid values are `"linear"` (default) and `"logarithmic"`.
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Transaction Duration Distribution"
  },
  "content": {
    "type": "visualization",
    "id": "viz.histogram",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT histogram(duration * 1000, 0.5, 20) FROM Transaction SINCE 1 hour ago"
        }
      ],
      "legend": {
        "enabled": true
      },
      "yAxisLeft": {
        "zero": true
      },
      "colors": {
        "colorPalette": "consistent"
      },
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.heatmap` - Heatmap

Renders a two-dimensional density map. Requires `histogram()` combined with a `FACET` to produce the row-and-bucket structure.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. Query must use `histogram()` with a `FACET` clause. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `facet` - (Optional) Controls FACET grouping behaviour.
    * `showOtherSeries` - (Optional) Show the "Other" group for `FACET` queries. Defaults to `false`.
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Duration Distribution by App"
  },
  "content": {
    "type": "visualization",
    "id": "viz.heatmap",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT histogram(duration * 1000, 0.5, 20) FROM Transaction FACET appName SINCE 3 hours ago"
        }
      ],
      "facet": {
        "showOtherSeries": false
      },
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.scatter` - Scatter plot

Plots two numeric measures against each other to reveal correlations. Your query should return two numeric values and a `FACET` for the data points.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. Query should return two numeric values and a `FACET`. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `facet` - (Optional) Controls FACET grouping behaviour.
    * `showOtherSeries` - (Optional) Show the "Other" group for `FACET` queries. Defaults to `false`.
  * `legend` - (Optional) Controls the chart legend.
    * `enabled` - (Optional) Show or hide the legend. Defaults to `true`.
    * `position` - (Optional) Where to place the legend. Valid values are `"bottom"` (default), `"left"`, or `"right"`.
  * `nullValues` - (Optional) Controls how null data points are rendered.
    * `nullValue` - (Optional) Default null handling. Valid values: `"default"`, `"zero"`, `"preserve"`, `"remove"`.
    * `seriesOverrides` - (Optional) Per-series null value overrides.
      * `seriesName` - (Required) The series to override.
      * `nullValue` - (Required) The null handling for that series.
  * `yAxisLeft` - (Optional) Configuration for the left Y axis.
    * `min` - (Optional) Fixed minimum value for the axis.
    * `max` - (Optional) Fixed maximum value for the axis.
    * `zero` - (Optional) Force zero as the axis origin. Defaults to `true`.
    * `scale` - (Optional) Axis scale. Valid values are `"linear"` (default) and `"logarithmic"`.
  * `units` - (Optional) Data units for display formatting.
    * `unit` - (Optional) Default unit for all series. See [Valid `units.unit` values](#valid-unitsunit-values).
    * `seriesOverrides` - (Optional) Per-series unit overrides.
      * `seriesName` - (Required) The series to override.
      * `unit` - (Required) The unit for that series.
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "CPU % vs Memory % by Host"
  },
  "content": {
    "type": "visualization",
    "id": "viz.scatter",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT average(cpuPercent) AS 'Avg CPU %', average(memoryUsedPercent) AS 'Avg Mem %' FROM SystemSample FACET hostname SINCE 1 hour ago TIMESERIES 5 minutes"
        }
      ],
      "legend": {
        "enabled": true
      },
      "yAxisLeft": {
        "zero": true
      },
      "facet": {
        "showOtherSeries": false
      },
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.apdex` - Apdex

Renders an Apdex score widget. Use it when you're tracking application performance against a satisfaction threshold.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `legend` - (Optional) Controls the chart legend.
    * `enabled` - (Optional) Show or hide the legend. Defaults to `true`.
    * `position` - (Optional) Where to place the legend. Valid values are `"bottom"` (default), `"left"`, or `"right"`.
  * `tooltip` - (Optional) Tooltip display behaviour.
    * `mode` - (Optional) Valid values: `"single"` (default), `"all"`, `"hidden"`.
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Apdex Score Over Time"
  },
  "content": {
    "type": "visualization",
    "id": "viz.apdex",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT apdex(duration, 0.5) AS 'Apdex' FROM Transaction SINCE 3 hours ago TIMESERIES 5 minutes"
        }
      ],
      "legend": {
        "enabled": true,
        "position": "bottom"
      },
      "tooltip": {
        "mode": "single"
      },
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.bullet` - Bullet chart

Compares an actual value against a target goal. The `limit` sets the goal line - without it the chart won't render correctly.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `limit` - (Required) The target value shown as the goal line on the chart.
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Requests vs Target"
  },
  "content": {
    "type": "visualization",
    "id": "viz.bullet",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT count(*) AS 'Requests' FROM Transaction SINCE 1 hour ago"
        }
      ],
      "limit": 50000,
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.funnel` - Funnel chart

Tracks how users or events progress through sequential steps. Requires the `funnel()` NRQL function.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. Query must use the `funnel()` NRQL function. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Transaction Funnel"
  },
  "content": {
    "type": "visualization",
    "id": "viz.funnel",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT funnel(session, WHERE name = 'step1' AS 'Step 1', WHERE name = 'step2' AS 'Step 2') FROM PageView SINCE 1 hour ago"
        }
      ],
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.event-feed` - Event feed

Streams individual events as a scrollable list. Useful when you want to show raw log entries or transaction events rather than aggregations.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Recent Transactions"
  },
  "content": {
    "type": "visualization",
    "id": "viz.event-feed",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT * FROM Transaction SINCE 10 minutes ago LIMIT 20"
        }
      ],
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.json` - JSON

Renders the raw JSON payload returned by your query. Handy for debugging complex nested event structures.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Transaction Metrics (JSON)"
  },
  "content": {
    "type": "visualization",
    "id": "viz.json",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT count(*) AS 'count', average(duration)*1000 AS 'avg_ms', percentile(duration*1000, 95) AS 'p95_ms' FROM Transaction FACET appName SINCE 1 hour ago"
        }
      ],
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.sparkline` - Sparkline

A compact inline time-series with no axis labels. Good for dense notebooks where you want trend at a glance without the full chart chrome.

This chart type requires a widget-level `props` object with a `title` key. See [Understanding the structure](#understanding-the-structure).

  * `nrqlQueries` - (Required) Array of NRQL query objects. See [Nested `nrqlQueries` blocks](#nested-nrqlqueries-blocks).
  * `facet` - (Optional) Controls FACET grouping behaviour.
    * `showOtherSeries` - (Optional) Show the "Other" group for `FACET` queries. Defaults to `false`.
  * `nullValues` - (Optional) Controls how null data points are rendered.
    * `nullValue` - (Optional) Default null handling. Valid values: `"default"`, `"zero"`, `"preserve"`, `"remove"`.
    * `seriesOverrides` - (Optional) Per-series null value overrides.
      * `seriesName` - (Required) The series to override.
      * `nullValue` - (Required) The null handling for that series.
  * `chartStyles` - (Optional) Visual style overrides.
    * `lineInterpolation` - (Optional) Line interpolation mode. Valid values: `"linear"` (default), `"smooth"`, `"stepBefore"`, `"stepAfter"`.
  * `yAxisLeft` - (Optional) Configuration for the left Y axis.
    * `min` - (Optional) Fixed minimum value for the axis.
    * `max` - (Optional) Fixed maximum value for the axis.
    * `zero` - (Optional) Force zero as the axis origin. Defaults to `true`.
  * `platformOptions` - (Optional) Platform-level rendering options.
    * `ignoreTimeRange` - (Optional) Use the query's own time range instead of the notebook time picker. Defaults to `false`.
  * `refreshRate` - (Optional) Auto-refresh configuration.
    * `frequency` - (Optional) Refresh interval in milliseconds, or `"auto"`. See [Valid `refreshRate.frequency` values](#valid-refreshratefrequency-values).

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Throughput Trend"
  },
  "content": {
    "type": "visualization",
    "id": "viz.sparkline",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT rate(count(*), 1 minute) AS 'rpm' FROM Transaction TIMESERIES 5 minutes SINCE 3 hours ago"
        }
      ],
      "chartStyles": {
        "lineInterpolation": "smooth"
      },
      "nullValues": {
        "nullValue": "zero"
      },
      "yAxisLeft": {
        "zero": true
      },
      "colors": {
        "colorPalette": "consistent"
      },
      "tooltip": {
        "mode": "single"
      },
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

### `viz.sparkline-lite` - Sparkline Lite

Ultra-compact sparkline with no axis labels. Supports the same props as `viz.sparkline` except `yAxisLeft`.

**Example:**
```json
{
  "type": "widget",
  "props": {
    "title": "Log Volume Trend"
  },
  "content": {
    "type": "visualization",
    "id": "viz.sparkline-lite",
    "props": {
      "nrqlQueries": [
        {
          "accountIds": [
            1234567
          ],
          "query": "SELECT count(*) AS 'Log lines' FROM Log TIMESERIES 5 minutes SINCE 3 hours ago"
        }
      ],
      "chartStyles": {
        "lineInterpolation": "stepAfter"
      },
      "nullValues": {
        "nullValue": "zero"
      },
      "colors": {
        "colorPalette": "consistent"
      },
      "platformOptions": {
        "ignoreTimeRange": false
      }
    }
  }
}
```

---

## Nested Blocks Reference

The following sections describe the structure of blocks that appear repeatedly across multiple visualization types.

### Nested `nrqlQueries` blocks

All query-based visualizations accept a `nrqlQueries` array. Each entry in the array supports:

  * `accountIds` - (Required) An array of one or more New Relic account IDs to run the query against.
  * `query` - (Required) A valid NRQL query string.
  * `offset` - (Optional) Number of milliseconds to offset the query time window.

### Nested `thresholds` blocks (line, area, stacked-bar)

Used with the `thresholds.thresholds` array inside `viz.line`, `viz.area`, and `viz.stacked-bar` to draw horizontal threshold bands across the chart.

  * `name` - (Optional) A label shown on the threshold band.
  * `from` - (Optional) Lower bound of the threshold range (inclusive).
  * `to` - (Optional) Upper bound of the threshold range (inclusive).
  * `severity` - (Required) Color coding applied to data in this range. Valid values: `"success"`, `"warning"`, `"severe"`, `"critical"`, `"unavailable"`.

### Nested `thresholds` blocks (table)

Used with the `thresholds` array inside `viz.table` to apply color coding per column.

  * `columnName` - (Required) The column to apply the threshold to, as it appears in the query results.
  * `from` - (Optional) Lower bound of the range.
  * `to` - (Optional) Upper bound of the range.
  * `severity` - (Required) Color coding for values in this range. Valid values: `"success"`, `"warning"`, `"severe"`, `"critical"`, `"unavailable"`.

### Nested `thresholdsWithSeriesOverrides` blocks (billboard)

Used with `viz.billboard` to control color coding for the displayed metric values.

  * `thresholds` - (Optional) Global threshold ranges applied to all series.
    * `from` - (Optional) Lower bound of the range.
    * `to` - (Optional) Upper bound of the range.
    * `severity` - (Required) Valid values: `"success"`, `"warning"`, `"severe"`, `"critical"`, `"unavailable"`.
  * `seriesOverrides` - (Optional) Per-series threshold overrides for multi-value billboards.
    * `seriesName` - (Required) The series name to override.
    * `from` - (Optional) Lower bound.
    * `to` - (Optional) Upper bound.
    * `severity` - (Required) Severity level for this series.

## Valid Values Reference

### Valid `units.unit` values

`APDEX`, `BITS`, `BITS_PER_MS`, `BITS_PER_SECOND`, `BYTES`, `BYTES_PER_MS`, `BYTES_PER_SECOND`, `CELSIUS`, `COUNT`, `DOLLAR`, `HERTZ`, `MS`, `PAGES_PER_SECOND`, `PERCENTAGE`, `REQUESTS_PER_SECOND`, `REQUESTS_PER_MINUTE`, `SECONDS`, `TIMESTAMP`

### Valid `refreshRate.frequency` values

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

## Examples

Ready-to-use notebook templates. Expand any example to see the full configuration.

<details>
<summary>Load notebook from a JSON file (recommended for larger notebooks)</summary>

### Loading from a file

The cleanest approach for anything beyond a simple notebook. Keep your JSON in a dedicated file and reference it from Terraform. This makes it easy to copy content directly from the New Relic UI's "Copy JSON" option.

```hcl
resource "newrelic_notebook" "weekly_review" {
  title   = "Weekly Service Health Review"
  content = file("${path.module}/notebooks/weekly-health.json")
}
```

</details>

<details>
<summary>One notebook per service using for_each (HCL authoring)</summary>

### Per-service runbook notebooks

Use `for_each` to generate a runbook notebook for each service in your estate. The `jsonencode({...})` expression lets you embed Terraform expressions like `each.value` and `var.account_id` directly inside the notebook body.

```hcl
variable "services" {
  type    = set(string)
  default = ["checkout", "payments", "inventory"]
}

resource "newrelic_notebook" "runbooks" {
  for_each = var.services
  title    = "${each.value} - Runbook"

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
            props = { text = "# ${each.value}\n\nAdd investigation steps and runbook content here." }
          }
        },
        {
          type  = "widget"
          props = { title = "Error rate" }
          content = {
            type = "visualization"
            id   = "viz.billboard"
            props = {
              nrqlQueries = [{
                accountIds = [var.account_id]
                query      = "SELECT percentage(count(*), WHERE error IS true) AS 'Error %' FROM Transaction WHERE appName = '${each.value}' SINCE 1 hour ago"
              }]
            }
          }
        }
      ]
    }]
  })
}
```

</details>

<details>
<summary>APM health notebook - throughput, error rate, and latency</summary>

### APM Health

A notebook for monitoring a single application's golden signals: request rate, error percentage, and response time percentiles. Adapt the `appName` filter and account ID to your environment.

```json
{
  "type": "declarative",
  "version": 1,
  "content": [
    {
      "type": "container",
      "props": {
        "layout": "stack"
      },
      "content": [
        {
          "type": "widget",
          "content": {
            "type": "visualization",
            "id": "viz.markdown",
            "props": {
              "text": "# APM Health\n\nThroughput, error rate, and latency for **Dummy App Two Max**."
            }
          }
        },
        {
          "type": "widget",
          "props": {
            "title": "Throughput and Error Rate"
          },
          "content": {
            "type": "visualization",
            "id": "viz.billboard",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [
                    1234567
                  ],
                  "query": "SELECT rate(count(*), 1 minute) AS 'rpm', percentage(count(*), WHERE error IS true) AS 'Error %' FROM Transaction WHERE appName = 'Dummy App Two Max' SINCE 1 hour ago"
                }
              ],
              "thresholdsWithSeriesOverrides": {
                "thresholds": [
                  {
                    "from": 1,
                    "severity": "critical"
                  }
                ],
                "seriesOverrides": [
                  {
                    "seriesName": "Error %",
                    "from": 1,
                    "severity": "critical"
                  }
                ]
              }
            }
          }
        },
        {
          "type": "widget",
          "props": {
            "title": "Requests per minute"
          },
          "content": {
            "type": "visualization",
            "id": "viz.line",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [
                    1234567
                  ],
                  "query": "SELECT rate(count(*), 1 minute) AS 'rpm' FROM Transaction WHERE appName = 'Dummy App Two Max' TIMESERIES 5 minutes SINCE 3 hours ago"
                }
              ],
              "legend": {
                "enabled": true
              },
              "yAxisLeft": {
                "zero": true
              },
              "nullValues": {
                "nullValue": "zero"
              }
            }
          }
        },
        {
          "type": "widget",
          "props": {
            "title": "Response time - P50 / P95 / P99"
          },
          "content": {
            "type": "visualization",
            "id": "viz.area",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [
                    1234567
                  ],
                  "query": "SELECT percentile(duration * 1000, 50) AS 'P50 ms', percentile(duration * 1000, 95) AS 'P95 ms', percentile(duration * 1000, 99) AS 'P99 ms' FROM Transaction WHERE appName = 'Dummy App Two Max' TIMESERIES 5 minutes SINCE 3 hours ago"
                }
              ],
              "legend": {
                "enabled": true,
                "position": "bottom"
              },
              "yAxisLeft": {
                "zero": true
              },
              "nullValues": {
                "nullValue": "zero"
              },
              "colors": {
                "seriesOverrides": [
                  {
                    "seriesName": "P50 ms",
                    "color": "#11A600"
                  },
                  {
                    "seriesName": "P95 ms",
                    "color": "#FFB951"
                  },
                  {
                    "seriesName": "P99 ms",
                    "color": "#BF0016"
                  }
                ]
              }
            }
          }
        }
      ]
    }
  ]
}
```

Use this JSON with `content = file("apm-health.json")` or inline using a heredoc.

</details>

<details>
<summary>Infrastructure health notebook - CPU, memory, and disk across hosts</summary>

### Infrastructure Health

A notebook for host-level visibility across your environment: an at-a-glance table, CPU trend over time, and disk utilisation by host.

```json
{
  "type": "declarative",
  "version": 1,
  "content": [
    {
      "type": "container",
      "props": {
        "layout": "stack"
      },
      "content": [
        {
          "type": "widget",
          "content": {
            "type": "visualization",
            "id": "viz.markdown",
            "props": {
              "text": "# Infrastructure Health\n\nHost-level CPU, memory, and disk across the environment."
            }
          }
        },
        {
          "type": "widget",
          "props": {
            "title": "Host Health - CPU, Memory and Disk"
          },
          "content": {
            "type": "visualization",
            "id": "viz.table",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [
                    1234567
                  ],
                  "query": "SELECT average(cpuPercent) AS 'Avg CPU %', max(cpuPercent) AS 'Max CPU %', average(memoryUsedPercent) AS 'Avg Mem %', average(diskUsedPercent) AS 'Disk %' FROM SystemSample FACET hostname SINCE 30 minutes ago LIMIT 20"
                }
              ],
              "initialSorting": {
                "name": "Avg CPU %",
                "direction": "desc"
              },
              "thresholds": [
                {
                  "columnName": "Avg CPU %",
                  "from": 80,
                  "severity": "critical"
                },
                {
                  "columnName": "Avg CPU %",
                  "from": 60,
                  "to": 80,
                  "severity": "warning"
                },
                {
                  "columnName": "Avg Mem %",
                  "from": 90,
                  "severity": "critical"
                }
              ]
            }
          }
        },
        {
          "type": "widget",
          "props": {
            "title": "CPU by Host"
          },
          "content": {
            "type": "visualization",
            "id": "viz.stacked-bar",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [
                    1234567
                  ],
                  "query": "SELECT average(cpuPercent) FROM SystemSample FACET hostname TIMESERIES 5 minutes SINCE 1 hour ago"
                }
              ],
              "legend": {
                "enabled": true
              },
              "yAxisLeft": {
                "zero": true,
                "min": 0,
                "max": 100
              }
            }
          }
        },
        {
          "type": "widget",
          "props": {
            "title": "Disk Utilisation"
          },
          "content": {
            "type": "visualization",
            "id": "viz.pie",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [
                    1234567
                  ],
                  "query": "SELECT average(diskUsedPercent) FROM StorageSample FACET hostname SINCE 30 minutes ago"
                }
              ],
              "facet": {
                "showOtherSeries": false
              },
              "colors": {
                "colorPalette": "consistent"
              }
            }
          }
        }
      ]
    }
  ]
}
```

</details>

<details>
<summary>Incident investigation runbook - structured template with findings and action items</summary>

### Incident Investigation Runbook

A structured investigation template. Fill in the severity and on-call fields at the top, run the diagnostic queries, record findings in the text block, and track action items through to resolution.

```json
{
  "type": "declarative",
  "version": 1,
  "content": [
    {
      "type": "container",
      "props": {
        "layout": "stack"
      },
      "content": [
        {
          "type": "widget",
          "content": {
            "type": "visualization",
            "id": "viz.markdown",
            "props": {
              "text": "# Incident Investigation\n\n**Severity**: P1\n**Start time**: _fill in_\n**On-call**: _fill in_\n\nRun each query below, record findings in the text blocks, and track action items at the bottom."
            }
          }
        },
        {
          "type": "widget",
          "props": {
            "title": "Current error rate"
          },
          "content": {
            "type": "visualization",
            "id": "viz.billboard",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [
                    1234567
                  ],
                  "query": "SELECT percentage(count(*), WHERE error IS true) AS 'Error Rate %', count(*) AS 'Total Requests' FROM Transaction SINCE 15 minutes ago"
                }
              ],
              "thresholdsWithSeriesOverrides": {
                "thresholds": [
                  {
                    "to": 1,
                    "severity": "success"
                  },
                  {
                    "from": 1,
                    "to": 5,
                    "severity": "warning"
                  },
                  {
                    "from": 5,
                    "severity": "critical"
                  }
                ],
                "seriesOverrides": [
                  {
                    "seriesName": "Error Rate %",
                    "from": 1,
                    "severity": "critical"
                  }
                ]
              }
            }
          }
        },
        {
          "type": "widget",
          "props": {
            "title": "Error rate over the last hour"
          },
          "content": {
            "type": "visualization",
            "id": "viz.line",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [
                    1234567
                  ],
                  "query": "SELECT percentage(count(*), WHERE error IS true) FROM Transaction TIMESERIES 1 minute SINCE 1 hour ago"
                }
              ],
              "legend": {
                "enabled": false
              },
              "yAxisLeft": {
                "zero": true
              },
              "thresholds": {
                "isLabelVisible": true,
                "thresholds": [
                  {
                    "name": "Alert threshold",
                    "from": 5,
                    "severity": "critical"
                  }
                ]
              }
            }
          }
        },
        {
          "type": "widget",
          "props": {
            "title": "Slowest transactions"
          },
          "content": {
            "type": "visualization",
            "id": "viz.table",
            "props": {
              "nrqlQueries": [
                {
                  "accountIds": [
                    1234567
                  ],
                  "query": "SELECT average(duration)*1000 AS 'Avg ms', max(duration)*1000 AS 'Max ms', count(*) AS 'Calls' FROM Transaction FACET name SINCE 30 minutes ago ORDER BY average(duration) DESC LIMIT 20"
                }
              ],
              "initialSorting": {
                "name": "Avg ms",
                "direction": "desc"
              },
              "thresholds": [
                {
                  "columnName": "Avg ms",
                  "from": 500,
                  "severity": "critical"
                },
                {
                  "columnName": "Avg ms",
                  "from": 200,
                  "to": 500,
                  "severity": "warning"
                }
              ]
            }
          }
        },
        {
          "type": "widget",
          "content": {
            "type": "visualization",
            "id": "viz.markdown",
            "props": {
              "text": "## Findings\n\n_Record what you've found - what's causing the spike, affected services, timeline._\n\n## Action Items\n\n- [ ] Identify root cause\n- [ ] Notify stakeholders\n- [ ] Apply mitigation\n- [ ] Schedule post-mortem"
            }
          }
        }
      ]
    }
  ]
}
```

</details>

---

## Import

Notebooks can be imported using their entity GUID:

```
$ terraform import newrelic_notebook.example <guid>
```

After importing, run `terraform plan`. The plan will show no changes if the `content` in your configuration matches the normalized content fetched from the API.
