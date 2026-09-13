---
layout: "newrelic"
page_title: "New Relic: newrelic_notebook"
sidebar_current: "docs-newrelic-resource-notebook"
description: |-
  Create and manage New Relic Notebooks.
---

# Resource: newrelic\_notebook

Use this resource to create and manage [New Relic Notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/introduction-notebooks/).

Notebooks are shareable documents that combine live NRQL queries, visualizations, and markdown narrative. The resource manages the notebook's title and its full block content. Content is stored as a JSON document that follows the [declarative UI format](#content-format) used natively by the Notebooks platform.

## Example Usage

```hcl
resource "newrelic_notebook" "example" {
  title = "Service Health Overview"

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
                text = "# Service Health\n\nKey metrics for the checkout service."
              }
            }
          },
          {
            type  = "widget"
            props = { title = "Error rate" }
            content = {
              type = "visualization"
              id   = "viz.billboard"
              props = {
                nrqlQueries = [
                  {
                    accountIds = [var.account_id]
                    query      = "SELECT percentage(count(*), WHERE error IS true) AS 'Error %' FROM Transaction WHERE appName = 'checkout' SINCE 1 hour ago"
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
            props = { title = "Throughput" }
            content = {
              type = "visualization"
              id   = "viz.line"
              props = {
                nrqlQueries = [
                  {
                    accountIds = [var.account_id]
                    query      = "SELECT rate(count(*), 1 minute) FROM Transaction WHERE appName = 'checkout' FACET appName TIMESERIES 5 minutes SINCE 3 hours ago"
                  }
                ]
                legend = { enabled = true }
              }
            }
          }
        ]
      }
    ]
  })
}
```

## Argument Reference

The following arguments are supported:

- `title` - (Required) The title of the notebook. Must be unique within the organization.
- `content` - (Optional) The notebook body as an HCL object using `jsonencode({...})`. Produces field-level plan diffs so `terraform plan` shows exactly which widget property changed. Mutually exclusive with `content_json`. See [Content Format](#content-format).
- `content_json` - (Optional) The notebook body as a raw JSON string. Use when importing notebooks from the New Relic UI or loading from a file using `file()`. Produces a line-level diff of the normalized JSON. Mutually exclusive with `content`.

Exactly one of `content` or `content_json` must be specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `guid` - The unique entity identifier (GUID) of the notebook in New Relic.
- `blob_id` - The blob identifier of the current notebook content, updated after each Terraform-managed write.
- `organization_id` - The New Relic organization ID the notebook belongs to. Resolved automatically from the provider credentials.

## Import

Notebooks can be imported using their entity GUID. Append `:content` or `:content_json` to select which attribute is populated in state — this should match your Terraform configuration.

```
# Populates content_json (default; use when your config uses content_json)
$ terraform import newrelic_notebook.example <guid>

# Populates content (use when your config uses content = jsonencode({...}))
$ terraform import newrelic_notebook.example <guid>:content
```

After importing, run `terraform plan`. If the mode in state matches your configuration, the plan will show no changes.

## Additional Examples

### From a JSON file

Use `content_json` with `file()` to manage a notebook whose body lives in a separate JSON file. This is useful when exporting notebooks from the New Relic UI.

```hcl
resource "newrelic_notebook" "from_file" {
  title        = "Weekly Review"
  content_json = file("${path.module}/notebooks/weekly-review.json")
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
      content = [{
        type  = "widget"
        props = {}
        content = {
          type  = "visualization"
          id    = "viz.markdown"
          props = { text = "# ${each.value}\n\nAdd runbook steps here." }
        }
      }]
    }]
  })
}
```

## Content Format

Notebook content follows the **declarative UI** format. Every document must include `type`, `version`, and `content`:

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
            "id": "<viz-id>",
            "props": { }
          }
        }
      ]
    }
  ]
}
```

The `id` field identifies the visualization type. Supported values:

| `id` | Description |
|---|---|
| `viz.markdown` | Markdown text block |
| `viz.line` | Line chart |
| `viz.area` | Area chart |
| `viz.stacked-bar` | Stacked bar chart |
| `viz.bar` | Simplified bar chart |
| `viz.stacked-horizontal-bar` | Horizontal stacked bar |
| `viz.pie` | Pie chart |
| `viz.table` | Data table |
| `viz.billboard` | Single-value metric display |
| `viz.gauge` | Gauge (arc, circular, or bar) |
| `viz.histogram` | Distribution histogram |
| `viz.heatmap` | 2D distribution heatmap |
| `viz.scatter` | Scatter plot |
| `viz.apdex` | Apdex score display |
| `viz.bullet` | Bullet chart |
| `viz.funnel` | Conversion funnel |
| `viz.event-feed` | Scrolling event list |
| `viz.json` | Raw JSON output |
| `viz.sparkline` | Compact trend line |
| `viz.sparkline-lite` | Ultra-compact sparkline |

Query-based visualizations accept `nrqlQueries` (an array of `{ accountIds, query }` objects). `viz.markdown` accepts `text`. `viz.bullet` additionally requires `limit`. For the complete props reference for each visualization type, see the [New Relic SDK chart components documentation](https://docs.newrelic.com/docs/new-relic-solutions/build-nr-ui/sdk-component/charts/Charts/).

### Plan diff behaviour

Both `content` and `content_json` are stored in state as normalized JSON (alphabetically sorted keys, 2-space indentation). This means:

- Reformatting your HCL or JSON file without changing values produces **no diff**.
- Changing a single widget property shows **only that property** as changed.
- Modifications made directly in the New Relic UI are detected precisely on the next `terraform plan`.
