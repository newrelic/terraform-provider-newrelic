---
layout: "newrelic"
page_title: "New Relic: newrelic_account_management"
sidebar_current: "docs-newrelic-resource-account-management"
description: |-
  Create and manage  sub accounts in New Relic.
---

# Resource: newrelic\_account\_management

Use this resource to create and manage New Relic sub accounts.

## Example Usage

##### Create Account
```hcl
resource "newrelic_account_management" "foo"{
	name = "Test Account Name"
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the Account.

-> **NOTE** <span style="color:red;">Starting <b>v3.95.0</b> of the New Relic Terraform Provider, the `region` argument on `newrelic_account_management` is deprecated and will be removed in a future major release.</span><br><br>Every New Relic organization is now tied to a specific region, and any sub-account you create is automatically placed in the region of the organization that owns your API key. The `regionCode` field on the underlying `accountManagementCreateAccount` mutation has been deprecated upstream, and this provider no longer forwards it to the API.<br><br>Setting `region` in your configuration is still accepted for backward compatibility, but it has <b>no effect</b> on where the account is created. Please <span style="color:tomato;">stop setting it in new configurations</span>, and remove it from existing ones when it's convenient.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the account created.

## Import

Accounts can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_account_management.foo <id>
```

