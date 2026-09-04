---
layout: "newrelic"
page_title: "New Relic: newrelic_notebook"
sidebar_current: "docs-newrelic-resource-notebook"
description: |-
  Create and manage New Relic Notebooks.
---

# Resource: newrelic_notebook

Use this resource to create and manage [New Relic Notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/introduction-notebooks/).

Notebooks let you combine live NRQL queries, visualizations, and markdown narrative into a single shareable document. This resource manages the full lifecycle of a notebook, including its title and block content.

## Choosing a content mode

You must specify exactly one of `content` or `content_json`. They are mutually exclusive:

| Mode | When to use |
|---|---|
| `content` | Authoring notebooks directly in Terraform. Uses `jsonencode({})` for structured HCL that produces field-level diffs in `terraform plan`. |
| `content_json` | Working from JSON exported out of the New Relic UI or stored in a file. Produces line-level diffs of the normalized content. |

---

## Example Usage

### Using `content` with a markdown block

The simplest notebook - a single text block authored in HCL.

```hcl
resource "newrelic_notebook" "incident_notes" {
  title           = "Incident Response Notes"

  content = jsonencode({
    version = "1"
    blocks = [
      {
        type = "widget"
        content = {
          type = "visualization"
          id   = "viz.markdown"
          props = {
            text = "## Summary\n\nAdd investigation notes here."
          }
        }
      }
    ]
  })
}
```

### Using `content` with multiple widget types

A notebook mixing a markdown header, a billboard metric, and a time-series chart.

```hcl
resource "newrelic_notebook" "service_overview" {
  title           = "Service Health Overview"

  content = jsonencode({
    version = "1"
    blocks = [
      {
        type = "widget"
        content = {
          type = "visualization"
          id   = "viz.markdown"
          props = {
            text = "# Service Health\n\nLive metrics for the checkout service."
          }
        }
      },
      {
        type = "widget"
        content = {
          type = "visualization"
          id   = "viz.billboard"
          props = {
            title = "Error rate (last hour)"
            nrqlQueries = [
              {
                accountIds = [var.account_id]
                query      = "SELECT percentage(count(*), WHERE error IS true) FROM Transaction SINCE 1 hour ago"
              }
            ]
            thresholdsWithSeriesOverrides = {
              thresholds = [
                { to = 1,  severity = "success" },
                { from = 1, to = 5, severity = "warning" },
                { from = 5, severity = "critical" }
              ]
            }
          }
        }
      },
      {
        type = "widget"
        content = {
          type = "visualization"
          id   = "viz.line"
          props = {
            title = "Throughput over time"
            nrqlQueries = [
              {
                accountIds = [var.account_id]
                query      = "SELECT count(*) FROM Transaction TIMESERIES AUTO"
              }
            ]
          }
        }
      }
    ]
  })
}
```

### Using `content` with a container block

Group related widgets visually using a `container` with `layout = "stack"`.

```hcl
resource "newrelic_notebook" "investigation" {
  title           = "DB Investigation"

  content = jsonencode({
    version = "1"
    blocks = [
      {
        type = "widget"
        content = {
          type  = "visualization"
          id    = "viz.markdown"
          props = { text = "## Database Performance\n\nQuery latency and throughput." }
        }
      },
      {
        type = "container"
        props = { layout = "stack" }
        content = [
          {
            type = "widget"
            content = {
              type = "visualization"
              id   = "viz.line"
              props = {
                nrqlQueries = [{ accountIds = [var.account_id], query = "SELECT average(duration) FROM DatabaseSample TIMESERIES" }]
              }
            }
          },
          {
            type = "widget"
            content = {
              type = "visualization"
              id   = "viz.bar"
              props = {
                nrqlQueries = [{ accountIds = [var.account_id], query = "SELECT count(*) FROM DatabaseSample FACET queryType" }]
              }
            }
          }
        ]
      }
    ]
  })
}
```

### Using `content_json` from a file

Paste JSON exported from the New Relic Notebooks UI directly into a file and reference it. The JSON below is the exact equivalent of the multi-widget `content` example above, so you can see the 1:1 parity between the two modes.

**`notebooks/service-health.json`**

```json
{
  "version": "1",
  "blocks": [
    {
      "type": "widget",
      "content": {
        "type": "visualization",
        "id": "viz.markdown",
        "props": {
          "text": "# Service Health\n\nLive metrics for the checkout service."
        }
      }
    },
    {
      "type": "widget",
      "content": {
        "type": "visualization",
        "id": "viz.billboard",
        "props": {
          "title": "Error rate (last hour)",
          "nrqlQueries": [
            {
              "accountIds": [1234567],
              "query": "SELECT percentage(count(*), WHERE error IS true) FROM Transaction SINCE 1 hour ago"
            }
          ],
          "thresholdsWithSeriesOverrides": {
            "thresholds": [
              { "to": 1, "severity": "success" },
              { "from": 1, "to": 5, "severity": "warning" },
              { "from": 5, "severity": "critical" }
            ]
          }
        }
      }
    },
    {
      "type": "widget",
      "content": {
        "type": "visualization",
        "id": "viz.line",
        "props": {
          "title": "Throughput over time",
          "nrqlQueries": [
            {
              "accountIds": [1234567],
              "query": "SELECT count(*) FROM Transaction TIMESERIES AUTO"
            }
          ]
        }
      }
    }
  ]
}
```

**`main.tf`**

```hcl
resource "newrelic_notebook" "service_overview" {
  title           = "Service Health Overview"
  content_json    = file("${path.module}/notebooks/service-health.json")
}
```

### Using `content_json` with an inline JSON string

For notebooks that are generated programmatically or pulled from another data source.

```hcl
locals {
  notebook_body = jsonencode({
    version = "1"
    blocks = [
      {
        type = "widget"
        content = {
          type  = "visualization"
          id    = "viz.markdown"
          props = { text = "# Generated notebook\n\nCreated by Terraform on ${timestamp()}." }
        }
      }
    ]
  })
}

resource "newrelic_notebook" "generated" {
  title           = "Auto-generated Notebook"
  content_json    = local.notebook_body
}
```

### Iterating to create multiple notebooks from a list

Use `for_each` to create a notebook per service.

```hcl
variable "services" {
  type = set(string)
  default = ["checkout", "payments", "inventory"]
}

resource "newrelic_notebook" "per_service" {
  for_each        = var.services
  title           = "${each.value} runbook"

  content = jsonencode({
    version = "1"
    blocks = [
      {
        type = "widget"
        content = {
          type  = "visualization"
          id    = "viz.markdown"
          props = { text = "# ${each.value}\n\nAdd runbook steps for this service here." }
        }
      }
    ]
  })
}
```

---

## Argument Reference

* `title` - (Required) The title of the notebook. Must be unique within the organization.
* `content` - (Optional) The notebook body, expressed as an HCL object using `jsonencode({...})`. Terraform evaluates the expression at plan time, producing field-level diffs. Mutually exclusive with `content_json`.
* `content_json` - (Optional) The notebook body as a raw JSON string. Use when working from a UI export or a file. Produces line-level diffs of normalized content. Mutually exclusive with `content`.
* `organization_id` - (Computed) The New Relic organization ID. Resolved automatically from the provider credentials.

## Attributes Reference

* `guid` - The unique entity identifier (GUID) assigned to the notebook by New Relic.
* `version` - The current revision counter for the notebook content. Starts at 1 and increments on every content write. Used internally to skip unnecessary Blob API reads when the content has not changed.

## Import

Notebooks can be imported by GUID. Optionally append `:content` or `:content_json` to control which field is populated in state, matching your Terraform configuration.

```
# Default - imports into content_json (for configs using content_json = file(...) or inline JSON)
$ terraform import newrelic_notebook.example <guid>
$ terraform import newrelic_notebook.example <guid>:content_json

# Import into content field (for configs using content = jsonencode({...}))
$ terraform import newrelic_notebook.example <guid>:content
```

After importing, run `terraform plan`. If the imported state and your config use the same mode, the plan will show no changes. If they differ, the plan surfaces the difference so you can reconcile your configuration.

## Plan Diff Behavior

Both content fields store a normalized form of the JSON in state (alphabetically sorted keys, 2-space indentation). This means:

- Reformatting your HCL or JSON file without changing any values produces **no diff** in `terraform plan`.
- Changing a single widget property shows **only that property** as changed.
- Externally modifying the notebook in the UI causes the changed fields to surface **precisely** in the next `terraform plan`.
