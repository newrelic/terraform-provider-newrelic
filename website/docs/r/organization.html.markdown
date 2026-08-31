---
layout: "newrelic"
page_title: "New Relic: newrelic_organization"
sidebar_current: "docs-newrelic-resource-organization"
description: |-
  Create and manage a New Relic organization.
---

# Resource: newrelic\_organization

Use this resource to create and manage New Relic organizations.

~> **NOTE:** Deleting an organization is not supported via the New Relic API. Destroying this resource will only remove it from Terraform state.

## Example Usage

```hcl
resource "newrelic_organization" "foo" {
  name = "My Organization"
}
```

### With a new managed account

```hcl
resource "newrelic_organization" "foo" {
  name        = "My Organization"
  customer_id = "my-customer-id"

  new_managed_account {
    name        = "My New Account"
    region_code = "US01"
  }
}
```

### With a shared account

```hcl
resource "newrelic_organization" "foo" {
  name = "My Organization"

  shared_account {
    account_id       = 12345678
    limiting_role_id = 1234
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the organization.
* `customer_id` - (Optional) The customer ID associated with the organization. Forces a new resource if changed.
* `new_managed_account` - (Optional) A nested block that describes a new managed account to create alongside the organization. Only one `new_managed_account` block is permitted. Forces a new resource if changed. See [Nested new_managed_account blocks](#nested-new_managed_account-blocks) below for details.
* `shared_account` - (Optional) A nested block that describes an account to share with the new organization. Only one `shared_account` block is permitted. Forces a new resource if changed. See [Nested shared_account blocks](#nested-shared_account-blocks) below for details.

### Nested `new_managed_account` blocks

* `name` - (Optional) The name of the new account to be created. Forces a new resource if changed.
* `region_code` - (Optional) The region code for the account. One of: `EU01`, `US01`. Forces a new resource if changed.

### Nested `shared_account` blocks

* `account_id` - (Required) The ID of the account to share with the new organization. Forces a new resource if changed.
* `limiting_role_id` - (Optional) The limiting role ID the new organization will be granted for the shared account. Forces a new resource if changed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `job_id` - The job ID of the organization creation task.
* `org_id` - The ID of the organization.

## Import

New Relic organizations can be imported using the organization ID:

```
terraform import newrelic_organization.foo <organization_id>
```

~> **NOTE:** After importing, run `terraform state show newrelic_organization.foo` to view the current state and copy the relevant fields into your configuration.