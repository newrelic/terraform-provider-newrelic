---
layout: "newrelic"
page_title: "New Relic: newrelic_organization"
sidebar_current: "docs-newrelic-resource-organization"
description: |-
  Create and manage a New Relic organization.
---

# Resource: newrelic\_organization

Use this resource to create and manage New Relic organizations.

~> **NOTE:** Deleting an organization is not supported via API. Destroying this resource in Terraform will only remove it from state.

## Example Usage

```hcl
resource "newrelic_organization" "foo" {
  name        = "My New Organization"
  customer_id = "customer-123"

  new_managed_account {
    name        = "My New Account"
    region_code = "US01"
  }

  shared_account {
    account_id       = 12345678
    limiting_role_id = 1234
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the organization.
* `customer_id` - (Optional) The customer ID for the organization.
* `new_managed_account` - (Optional) A nested block that describes a new managed account to create within the organization. Only one `new_managed_account` block is permitted per resource definition. See [Nested new_managed_account blocks](#nested-new_managed_account-blocks) below for details.
* `shared_account` - (Optional) A nested block that describes an account to share with the new organization. Only one `shared_account` block is permitted per resource definition. See [Nested shared_account blocks](#nested-shared_account-blocks) below for details.

### Nested `new_managed_account` blocks

* `name` - (Optional) The name of the new account to be created.
* `region_code` - (Optional) The region code for the account. One of: `EU01`, `US01`.

### Nested `shared_account` blocks

* `account_id` - (Required) The ID of the account to share with the new organization.
* `limiting_role_id` - (Optional) The limiting role ID the new organization will be granted for the shared account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The job ID of the organization creation task.
* `job_id` - The job ID of the organization creation task.

## Import

New Relic organizations can be imported using the organization ID:

```
terraform import newrelic_organization.foo <organization_id>
```

~> **NOTE:** Since organization deletion is not supported via the API, running `terraform destroy` on this resource will only remove it from Terraform state without making any API calls.