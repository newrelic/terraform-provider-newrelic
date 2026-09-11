---
layout: "newrelic"
page_title: "New Relic: newrelic_notebook"
sidebar_current: "docs-newrelic-datasource-notebook"
description: |-
  Look up a New Relic Notebook by GUID and optionally retrieve its full content.
---

# Data Source: newrelic_notebook

Use this data source to look up an existing New Relic Notebook by its entity GUID.

By default the data source fetches only **metadata** (title, organization ID, blob ID) via a single NerdGraph call. Set `fetch_content = true` to also retrieve the full notebook body from the Blob Storage API.

## Example Usage

### Metadata only (default)

```hcl
data "newrelic_notebook" "example" {
  guid = "MTIxOTEy..."
}

output "notebook_title" {
  value = data.newrelic_notebook.example.title
}
```

### With full content

```hcl
data "newrelic_notebook" "example" {
  guid          = "MTIxOTEy..."
  fetch_content = true
}

output "notebook_body" {
  value = data.newrelic_notebook.example.content
}
```

### Referencing a notebook created by the resource

```hcl
resource "newrelic_notebook" "runbook" {
  title = "Incident Runbook"
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
              props = { text = "## Runbook\n\nInvestigation steps." }
            }
          }
        ]
      }
    ]
  })
}

data "newrelic_notebook" "runbook_read" {
  guid          = newrelic_notebook.runbook.guid
  fetch_content = true
}
```

## Argument Reference

* `guid` - (Required) The unique entity identifier (GUID) of the notebook to look up.
* `fetch_content` - (Optional) When `true`, the full notebook body is retrieved from the Blob Storage API and exposed in the `content` attribute. Defaults to `false`.

## Attributes Reference

* `title` - The title of the notebook.
* `organization_id` - The New Relic organization ID the notebook belongs to.
* `blob_id` - The blob identifier of the current notebook content. This is the immutable ID of the specific content version stored in the Blob Storage API.
* `content` - The notebook body as a normalized JSON string (alphabetically sorted keys, 2-space indentation). Only populated when `fetch_content = true`; empty string otherwise. The content follows the declarative UI schema — see the [resource documentation](../r/notebook.html.markdown) for the full schema reference.
