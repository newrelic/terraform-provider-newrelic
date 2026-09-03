---
layout: "newrelic"
page_title: "New Relic: newrelic_notebook"
sidebar_current: "docs-newrelic-resource-notebook"
description: |-
  Create and manage New Relic Notebooks.
---

# Resource: newrelic_notebook

Use this resource to create and manage [New Relic Notebooks](https://docs.newrelic.com/docs/query-your-data/explore-query-data/notebooks/introduction-notebooks/).

A notebook is a document that combines live New Relic queries, markdown text, and visualizations into a single shareable view. The `content` field holds the full notebook body as a JSON string and is sent to the Blob Storage API without any schema expansion, so the resource never needs updates when the notebook format evolves.

## Example Usage

### Using `jsonencode` (recommended)

Using `jsonencode({...})` lets you write the notebook content as an HCL object. Terraform evaluates the expression at plan time, which means `terraform plan` output shows individual attribute changes rather than a single opaque string diff.

```hcl
resource "newrelic_notebook" "example" {
  title = "My Notebook"

  content = jsonencode({
    version = "1"
    blocks = [
      {
        type = "widget"
        content = {
          type  = "visualization"
          id    = "viz.markdown"
          props = {
            text = "## Hello from Terraform"
          }
        }
      },
      {
        type = "widget"
        content = {
          type  = "visualization"
          id    = "viz.line"
          props = {
            accountIds = [var.account_id]
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

### Using a JSON file

```hcl
resource "newrelic_notebook" "from_file" {
  title   = "Notebook from file"
  content = file("${path.module}/notebook.json")
}
```

### Specifying an organization ID explicitly

```hcl
resource "newrelic_notebook" "example" {
  title           = "Scoped Notebook"
  organization_id = var.organization_id

  content = jsonencode({
    version = "1"
    blocks  = []
  })
}
```

## Argument Reference

The following arguments are supported:

* `title` - (Required) The title of the notebook.
* `content` - (Required) The notebook content as a JSON string. The provider normalises the JSON (sorted keys, consistent indentation) before storing it in state, so cosmetic formatting differences — such as reordering keys or changing indentation — do not produce a diff in `terraform plan`. Use `jsonencode({...})` for the most granular plan output.
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

## Diff Behaviour and the `jsonencode` Pattern

The `content` field is stored in state as canonical JSON (alphabetically sorted keys, 2-space indentation). Two JSON strings that differ only in whitespace, indentation, or key ordering are treated as equal and produce no diff.

When you use `jsonencode({...})`, Terraform's plan output shows changes at the individual attribute level:

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

When you supply a raw JSON string (e.g. via `file()`), Terraform shows a line-level unified diff of the normalised content. Only the lines that actually changed are highlighted.
