---
layout: "newrelic"
page_title: "New Relic: newrelic_organization"
sidebar_current: "docs-newrelic-resource-organization"
description: |-
  Create and manage a New Relic organization.
---

# Resource: newrelic\_organization

Use this resource to create and manage New Relic organizations.

~> **NOTE:** Deleting an organization is not supported via the New Relic API. Removing this resource from your Terraform configuration will only remove it from state; the organization itself will not be deleted.

## Example Usage

```hcl
resource "newrelic_organization" "foo" {
  name = "My Organization"

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
* `customer_id` - (Optional) The customer ID for the organization. This field is immutable; changing it forces a new resource to be created.
* `new_managed_account` - (Optional) A nested block describing a new managed account to create alongside the organization. This field is immutable; changing it forces a new resource to be created. Only one `new_managed_account` block is permitted. See [Nested new_managed_account blocks](#nested-new_managed_account-blocks) below for details.
* `shared_account` - (Optional) A nested block describing an account to share with the new organization. This field is immutable; changing it forces a new resource to be created. Only one `shared_account` block is permitted. See [Nested shared_account blocks](#nested-shared_account-blocks) below for details.

### Nested `new_managed_account` blocks

* `name` - (Optional) The name of the new managed account to be created. This field is immutable; changing it forces a new resource to be created.
* `region_code` - (Optional) The region code for the account to be created. This field is immutable; changing it forces a new resource to be created. One of: `EU01`, `US01`.

### Nested `shared_account` blocks

* `account_id` - (Required) The ID of the account to share with the new organization. This field is immutable; changing it forces a new resource to be created.
* `limiting_role_id` - (Optional) The limiting role ID the new organization will be granted on the shared account. This field is immutable; changing it forces a new resource to be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The job ID of the organization creation task.
* `job_id` - The job ID of the organization creation task.

## Import

New Relic organizations can be imported using the organization ID.

```
terraform import newrelic_organization.foo <organization_id>
```

~> **NOTE:** After importing, run `terraform plan` to verify the state matches your configuration. Some fields (such as `customer_id`, `new_managed_account`, and `shared_account`) are only used at creation time and may not be reflected in the imported state.