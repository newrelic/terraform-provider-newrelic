---
layout: "newrelic"
page_title: "New Relic: newrelic_pipeline_cloud_rule"
sidebar_current: "docs-newrelic-pipeline-cloud-rule"
description: |-
  Use this resource to create and manage a New Relic Pipeline Cloud Rule.
---

# newrelic\_pipeline\_cloud\_rule

Use this resource to create and manage a New Relic Pipeline Cloud Rule.

## Example Usage

```hcl
resource "newrelic_pipeline_cloud_rule" "foo" {
  account_id  = "12345"
  name        = "My Rule"
  description = "My rule description"
  nrql        = "SELECT * FROM Log WHERE message like '%my-app%'"
}
```

## Argument Reference

The following arguments are supported:

*   `account_id` - (Optional) The account ID where the Pipeline Cloud rule will be created.
*   `name` - (Required) The name of the rule. This must be unique within an account.
*   `nrql` - (Required) The NRQL query that defines which data will be processed by this pipeline cloud rule.
*   `description` - (Optional) Provides additional information about the rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

*   `id` - The ID of the pipeline cloud rule.

## Import

Pipeline cloud rules can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_pipeline_cloud_rule.foo <id>
```
