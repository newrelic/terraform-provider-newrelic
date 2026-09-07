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

<details>
  <summary>Minimal notebook — single markdown block</summary>

```hcl
resource "newrelic_notebook" "incident_notes" {
  title = "Incident Response Notes"

  content = jsonencode({
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

</details>

<details>
  <summary>Incident investigation runbook — markdown, billboards, line chart, area chart, table (content mode)</summary>

A realistic notebook that walks through a production incident: a header with context, KPI billboards, a timeline chart, a drill-down table, a database-latency area chart, and a markdown root-cause summary with action items.

```hcl
resource "newrelic_notebook" "incident_runbook" {
  title = "Production API Incident — Investigation"

  content = jsonencode({
    blocks = [
      {
        type = "widget"
        content = {
          type = "visualization"
          id   = "viz.markdown"
          props = {
            text = "# Production API Incident\n\n**Status**: Resolved  \n**Duration**: 46 minutes (14:32–15:18 UTC)  \n**Impact**: ~15% of API requests failed\n\nThis notebook walks through what happened, identifies the root cause, and tracks action items."
          }
        }
      },
      {
        type = "widget"
        content = {
          type = "visualization"
          id   = "viz.billboard"
          props = {
            title = "Peak error rate during incident window"
            nrqlQueries = [
              {
                accountIds = [var.account_id]
                query      = "SELECT percentage(count(*), WHERE httpResponseCode >= 400) AS 'Error Rate %' FROM Transaction WHERE appName = 'api-production' SINCE '2024-10-15 14:00:00' UNTIL '2024-10-15 16:00:00'"
              }
            ]
            thresholdsWithSeriesOverrides = {
              thresholds = [
                { to = 1,  severity = "success" },
                { from = 1, to = 5,  severity = "warning" },
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
            title = "Error rate over time (1-minute granularity)"
            nrqlQueries = [
              {
                accountIds = [var.account_id]
                query      = "SELECT percentage(count(*), WHERE httpResponseCode >= 400) AS 'Error Rate %' FROM Transaction WHERE appName = 'api-production' TIMESERIES 1 minute SINCE '2024-10-15 14:00:00' UNTIL '2024-10-15 16:00:00'"
              }
            ]
          }
        }
      },
      {
        type = "widget"
        content = {
          type = "visualization"
          id   = "viz.table"
          props = {
            title = "Top affected endpoints"
            nrqlQueries = [
              {
                accountIds = [var.account_id]
                query      = "SELECT count(*) AS 'Errors', average(duration)*1000 AS 'Avg Duration (ms)' FROM Transaction WHERE appName = 'api-production' AND httpResponseCode >= 400 FACET request.uri SINCE '2024-10-15 14:00:00' UNTIL '2024-10-15 16:00:00' ORDER BY count(*) DESC LIMIT 10"
              }
            ]
          }
        }
      },
      {
        type = "widget"
        content = {
          type = "visualization"
          id   = "viz.area"
          props = {
            title = "Database query latency during incident"
            nrqlQueries = [
              {
                accountIds = [var.account_id]
                query      = "SELECT average(duration) AS 'Avg (ms)', percentile(duration, 95) AS 'P95 (ms)' FROM DatabaseSample WHERE host = 'prod-db-01' TIMESERIES 5 minutes SINCE '2024-10-15 14:00:00' UNTIL '2024-10-15 16:00:00'"
              }
            ]
          }
        }
      },
      {
        type = "widget"
        content = {
          type = "visualization"
          id   = "viz.markdown"
          props = {
            text = "## Root Cause\n\n- **14:32** — Error rates rose on `/api/users`; database connection pool exhausted under traffic spike\n- **14:45** — Database scaling initiated\n- **15:18** — Service fully recovered after pool limit raised\n\n## Action Items\n\n- [ ] Increase DB connection pool limit (`max_connections` → 500)\n- [ ] Add rate limiting on `/api/users/register`\n- [ ] Set up proactive alerting for DB connection pool utilization > 80%"
          }
        }
      }
    ]
  })
}
```

</details>

<details>
  <summary>Weekly service health review — markdown, billboard, line, area, bar, pie charts (content_json mode from file)</summary>

A realistic weekly review notebook stored as a JSON file. Combines a narrative header, an Apdex billboard, P95 latency comparison with prior week, request volume by service, error breakdown, and a markdown analysis section.

**`notebooks/weekly-health.json`**

```json
{
  "blocks": [
    {
      "type": "widget",
      "content": {
        "type": "visualization",
        "id": "viz.markdown",
        "props": {
          "text": "# Weekly Service Health Review\n\n**Period**: Last 7 days vs. prior week  \n**Services**: web-frontend, api-backend, auth-service\n\nThis notebook tracks P50/P95 response times, error rates, and throughput across all production services."
        }
      }
    },
    {
      "type": "widget",
      "content": {
        "type": "visualization",
        "id": "viz.billboard",
        "props": {
          "title": "Overall Apdex (7 days)",
          "nrqlQueries": [
            {
              "accountIds": [1234567],
              "query": "SELECT apdex(duration, 0.5) AS 'Apdex' FROM Transaction WHERE appName IN ('web-frontend', 'api-backend', 'auth-service') SINCE 7 days ago"
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
      "content": {
        "type": "visualization",
        "id": "viz.line",
        "props": {
          "title": "P95 response time by service (vs. prior week)",
          "nrqlQueries": [
            {
              "accountIds": [1234567],
              "query": "SELECT percentile(duration, 95) AS 'P95 (ms)' FROM Transaction WHERE appName IN ('web-frontend', 'api-backend', 'auth-service') FACET appName TIMESERIES 1 day SINCE 7 days ago COMPARE WITH 1 week ago"
            }
          ]
        }
      }
    },
    {
      "type": "widget",
      "content": {
        "type": "visualization",
        "id": "viz.area",
        "props": {
          "title": "Request volume over time",
          "nrqlQueries": [
            {
              "accountIds": [1234567],
              "query": "SELECT count(*) AS 'Requests' FROM Transaction WHERE appName IN ('web-frontend', 'api-backend', 'auth-service') FACET appName TIMESERIES 1 hour SINCE 7 days ago"
            }
          ]
        }
      }
    },
    {
      "type": "widget",
      "content": {
        "type": "visualization",
        "id": "viz.bar",
        "props": {
          "title": "Error count by service",
          "nrqlQueries": [
            {
              "accountIds": [1234567],
              "query": "SELECT count(*) AS 'Errors' FROM Transaction WHERE appName IN ('web-frontend', 'api-backend', 'auth-service') AND error IS true FACET appName SINCE 7 days ago"
            }
          ]
        }
      }
    },
    {
      "type": "widget",
      "content": {
        "type": "visualization",
        "id": "viz.pie",
        "props": {
          "title": "Error distribution by HTTP status code",
          "nrqlQueries": [
            {
              "accountIds": [1234567],
              "query": "SELECT count(*) FROM Transaction WHERE error IS true AND appName IN ('web-frontend', 'api-backend', 'auth-service') FACET httpResponseCode SINCE 7 days ago"
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
          "text": "## Analysis\n\n### Wins ✅\n- API backend P95 down from 650 ms to 420 ms after the DB index added on Monday\n- Auth service error rate down 50% week-over-week\n\n### Areas of Concern ⚠️\n- Web frontend P95 up 15% WoW — investigate asset bundle size regression\n- DB query timeouts up 25% — review slow query log before next release"
        }
      }
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
  <summary>Iterating to create multiple notebooks from a list (for_each)</summary>

```hcl
variable "services" {
  type    = set(string)
  default = ["checkout", "payments", "inventory"]
}

resource "newrelic_notebook" "per_service" {
  for_each = var.services
  title    = "${each.value} runbook"

  content = jsonencode({
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

</details>

---

## Widget types and block content schema

The `content` and `content_json` fields accept any valid notebook JSON. Each entry in `blocks` is a widget — a single visualization or markdown text block.

A widget block has this shape:

```json
{
  "type": "widget",
  "content": {
    "type": "visualization",
    "id": "<viz-id>",
    "props": { }
  }
}
```

`<viz-id>` is a visualization identifier such as `viz.markdown`, `viz.line`, `viz.area`, `viz.bar`, `viz.pie`, `viz.table`, or `viz.billboard`. The `props` schema differs per visualization type — query-based widgets accept `nrqlQueries`, while `viz.markdown` accepts `text`.

The Blob API stores your JSON verbatim — no fields are added, removed, or transformed server-side. Terraform tracks the entire blob. New Relic itself does not assign special meaning to any field outside of `blocks`, so you can include arbitrary top-level metadata (such as a `"version"` string for your own change-tracking); Terraform will diff those fields exactly like any other part of the content if they change.

For the full list of supported chart types, their `props` schemas, and worked examples:

- [Visualizations in notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/visualizations-in-notebooks/) — chart types and best practices
- [Create widgets with NerdGraph](https://docs.newrelic.com/docs/apis/nerdgraph/examples/create-widgets-dashboards-api/) — `props` reference for each widget type
- [Notebook examples](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/notebooks-examples/) — worked multi-block notebook JSON
- [Blob Storage API for notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/blob-storage-api-for-notebooks/) — the underlying REST API used by this resource

---

## Argument Reference

* `title` - (Required) The title of the notebook. Must be unique within the organization.
* `content` - (Optional) The notebook body, expressed as an HCL object using `jsonencode({...})`. Terraform evaluates the expression at plan time, producing field-level diffs. Mutually exclusive with `content_json`. The only required key inside the object is `blocks`; any additional top-level fields are stored verbatim by the Blob API and tracked in Terraform state.
* `content_json` - (Optional) The notebook body as a raw JSON string. Use when working from a UI export or a file. Produces line-level diffs of normalized content. Mutually exclusive with `content`. The only required key is `"blocks"`; any additional top-level fields are stored verbatim by the Blob API and tracked in Terraform state.
* `organization_id` - (Computed) The New Relic organization ID. Resolved automatically from the provider credentials.

## Attributes Reference

* `guid` - The unique entity identifier (GUID) assigned to the notebook by New Relic.
* `blob_id` - The blob identifier of the current notebook content. Updated after every Terraform-managed write. Used internally to detect when content has changed, avoiding unnecessary Blob Storage reads on plans where nothing has changed.

## Import

Notebooks can be imported by GUID. Optionally append `:content` or `:content_json` to control which field is populated in state, matching your Terraform configuration.

```
# Default — imports into content_json (for configs using content_json = file(...) or inline JSON)
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
