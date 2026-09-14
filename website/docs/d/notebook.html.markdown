---
layout: "newrelic"
page_title: "New Relic: newrelic_notebook"
sidebar_current: "docs-newrelic-datasource-notebook"
description: |-
  Look up a New Relic Notebook by GUID.
---

# Data Source: newrelic\_notebook

Use this data source to look up an existing [New Relic Notebook](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/introduction-notebooks/) by its entity GUID.

## Example Usage

```hcl
data "newrelic_notebook" "example" {
  guid          = "MTIxOTEy..."
  fetch_content = true
}

output "notebook_content" {
  value = data.newrelic_notebook.example.content_json
}
```

## Argument Reference

  * `guid` - (Required) The entity GUID of the notebook to look up.
  * `fetch_content` - (Optional) When `true`, fetches the full notebook body and populates `content_json`. Default: `false`.

## Attributes Reference

  * `title` - The title of the notebook.
  * `organization_id` - The New Relic organization ID the notebook belongs to.
  * `blob_id` - The blob identifier of the current notebook content.
  * `content_json` - The notebook body as a normalized JSON string. Only populated when `fetch_content` is `true`; empty string otherwise.
