---
layout: 'newrelic'
page_title: 'New Relic: newrelic_pathpoint_flow'
sidebar_current: 'docs-newrelic-resource-pathpoint-flow'
description: |-
  Create and manage pathpoint flow in New Relic.
---

# Resource: newrelic_pathpoint_flow

Use this resource to create and manage New Relic pathpoint flow.

## Example Usage

```hcl
resource "newrelic_pathpoint_flow" "foo" {
  account_id = <your_account_id>
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The New Relic account ID to operate on. Defaults to the account set in the provider.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

  * `id` - The ID of the pathpoint flow.
  * `account_id` - The New Relic account ID to operate on. Defaults to the account set in the provider.

## Import

Pathpoint Flow can be imported using the `<id>`, e.g.

```
$ terraform import newrelic_pathpoint_flow.foo MjUyMDUyOHxBUE18QVBQTElDQVRJT058MjE1MDM3Nzk1
```
