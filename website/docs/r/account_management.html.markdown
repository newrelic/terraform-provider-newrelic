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

-> **`region` is deprecated and silently ignored.** New Relic organizations are now single-region — every account created is placed in the region of the organization that owns the caller's API key. The `regionCode` argument on the underlying `accountManagementCreateAccount` mutation has been deprecated upstream by the IAM team, and the provider no longer forwards this value to the API. If you set `region` in your configuration, the value is accepted for backward compatibility but has **no effect** on where the account is created; Terraform will print a deprecation warning at plan time. The field will be removed in a future major release.

## Example Usage

##### Create Account (recommended — omit `region`)
```hcl
resource "newrelic_account_management" "foo" {
  name = "Test Account Name"
}
```

##### Legacy multi-region organization (region argument is now a no-op, kept for back-compat)
```hcl
resource "newrelic_account_management" "foo" {
  name   = "Test Account Name"
  region = "us01"   # DEPRECATED and silently ignored — leave unset for new configurations
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) The name of the Account.
  * `region` - (Optional, **Deprecated, silently ignored**) The region code of the account. Accepted values: `us01`, `eu01`, `jp01`, or omitted. **This field no longer has any effect** — accounts are always created in the region of the organization associated with the caller's API key. It is retained only for backward compatibility with existing configurations and will be removed in a future major release.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the account created.

## Import

Accounts can be imported using the `id`, e.g.

```bash
$ terraform import newrelic_account_management.foo <id>
```

