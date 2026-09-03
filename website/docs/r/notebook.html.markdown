---
layout: "newrelic"
page_title: "New Relic: newrelic_notebook"
sidebar_current: "docs-newrelic-resource-notebook"
description: |-
  Create and manage New Relic Notebooks.
---

# Resource: newrelic_notebook

Use this resource to create and manage [New Relic Notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/introduction-notebooks/).

A notebook is a document that combines live New Relic queries, markdown text, and visualizations into a shareable view. Notebook content is stored in the New Relic Blob Storage API; this resource sends the content body verbatim without schema expansion, so it stays compatible with future changes to the notebook format.

## Content modes

Exactly one of `content` or `content_json` must be specified. They are mutually exclusive.

## Example Usage

### `content` - HCL object syntax (recommended)

Use `jsonencode({...})` to write notebook content as an HCL object literal. Terraform evaluates the expression at plan time, producing individual attribute-level diffs in `terraform plan` output.

```hcl
resource "newrelic_notebook" "example" {
  title           = "My Notebook"
  organization_id = var.organization_id

  content = jsonencode({
    version = "1"
    blocks = [
      {
        type = "widget"
        content = {
          type  = "visualization"
          id    = "viz.markdown"
          props = { text = "## Hello from Terraform" }
        }
      },
      {
        type = "widget"
        content = {
          type = "visualization"
          id   = "viz.line"
          props = {
            nrqlQueries = [{
              accountIds = [var.account_id]
              query      = "SELECT count(*) FROM Transaction TIMESERIES AUTO"
            }]
          }
        }
      }
    ]
  })
}
```

### `content_json` - raw JSON string

Use `content_json` when pasting JSON exported from the New Relic UI (via **Copy JSON**) or loading from a file. Terraform still shows line-level diffs of the normalised content.

```hcl
resource "newrelic_notebook" "from_file" {
  title           = "Notebook from file"
  organization_id = var.organization_id
  content_json    = file("${path.module}/notebook.json")
}
```

## Argument Reference

The following arguments are supported:

* `title` - (Required) The title of the notebook. Must be unique within the organization.
* `content` - (Optional) The notebook content as an HCL object, expressed using `jsonencode({...})`. Produces the most granular `terraform plan` diffs. Mutually exclusive with `content_json`.
* `content_json` - (Optional) The notebook content as a raw JSON string. Use when pasting content exported from the New Relic UI or loading from a file. Mutually exclusive with `content`.
* `organization_id` - (Optional, Computed) The New Relic organization ID. When omitted, the provider resolves it automatically from the authenticated account.
* `account_id` - (Optional, Computed) The New Relic account ID. Defaults to the account configured in the provider.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `guid` - The unique entity identifier (GUID) of the notebook in New Relic.

## Import

Notebooks can be imported using the entity GUID:

```
$ terraform import newrelic_notebook.example <guid>
```

## Diff Behaviour

Both `content` and `content_json` are stored in state as canonical JSON (alphabetically sorted keys, 2-space indentation). Two JSON strings that differ only in whitespace, indentation, or key ordering are treated as equal and produce no diff.

**With `content = jsonencode({...})`**, Terraform shows individual attribute-level diffs:

```
~ content = jsonencode(
    ~ {
        ~ blocks = [
            ~ {
                ~ content = {
                    ~ props = {
                        ~ text = "## Hello" -> "## Hello, world"
                      }
                  }
              },
          ]
      }
  )
```

**With `content_json`**, Terraform shows a line-level diff of the normalised JSON. Only lines that actually changed are highlighted - unchanged blocks appear as context without `+`/`-` markers.
