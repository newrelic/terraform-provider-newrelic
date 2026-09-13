---
layout: "newrelic"
page_title: "New Relic: newrelic_notebook"
sidebar_current: "docs-newrelic-datasource-notebook"
description: |-
  Look up a New Relic Notebook by GUID.
---

# Data Source: newrelic\_notebook

Use this data source to retrieve information about an existing New Relic Notebook by its entity GUID.

By default, only NerdGraph **metadata** is fetched (title, organization ID, blob ID). Set `fetch_content = true` to also retrieve the full notebook body from the Blob Storage API, which is useful when you need to pass the content to another resource or inspect it in outputs.

## Example Usage

```hcl
data "newrelic_notebook" "example" {
  guid = "MTIxOTEy..."
}

output "notebook_title" {
  value = data.newrelic_notebook.example.title
}
```

### With content

```hcl
data "newrelic_notebook" "example" {
  guid          = "MTIxOTEy..."
  fetch_content = true
}

output "notebook_body" {
  value = data.newrelic_notebook.example.content
}
```

### Referencing a managed notebook

```hcl
resource "newrelic_notebook" "runbook" {
  title = "Checkout Runbook"
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
          props = { text = "## Runbook" }
        }
      }]
    }]
  })
}

data "newrelic_notebook" "runbook_read" {
  guid          = newrelic_notebook.runbook.guid
  fetch_content = true
}
```

## Argument Reference

The following arguments are supported:

- `guid` - (Required) The unique entity identifier (GUID) of the notebook to look up.
- `fetch_content` - (Optional) When `true`, the full notebook body is retrieved from the Blob Storage API and exposed in the `content` attribute. Defaults to `false`. Set this to `true` only when you need the content, as it incurs an additional API call.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `title` - The title of the notebook.
- `organization_id` - The New Relic organization ID the notebook belongs to.
- `blob_id` - The blob identifier of the current notebook content version.
- `content` - The notebook body as a normalized JSON string. Only populated when `fetch_content = true`; empty string otherwise. See the [resource documentation](../r/notebook.html.markdown) for the content schema.
