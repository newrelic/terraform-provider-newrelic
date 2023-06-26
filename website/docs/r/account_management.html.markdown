---
layout: "newrelic"
page_title: "New Relic: newrelic_account_management"
sidebar_current: "docs-newrelic-resource-account-management"
description: |-
  Create and manage  sub accounts in New Relic.
---

# Resource: newrelic\_account\_management

Use this resource to create and manage New Relic sub accounts.

-> **WARNING:** The `newrelic_account_management` resource will only create/update but won't delete a sub account. Please visit our documentation on  [`Account Management`](https://docs.newrelic.com/docs/apis/nerdgraph/examples/manage-accounts-nerdgraph/#delete) for more information .

## Example Usage

##### Create Account
```hcl
resource "newrelic_account_management" "foo"{
	name=  "Test Account Name"
	region= "us01"
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the Account.
  * `region` - (Required) The region code of the account.  One of: `us01`, `eu01`.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the account created.

## Import

Accounts can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_account_management.foo <id>
```

