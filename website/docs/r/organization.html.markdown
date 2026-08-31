---
layout: "newrelic"
page_title: "New Relic: newrelic_organization"
sidebar_current: "docs-newrelic-resource-organization"
description: |-
  Create and manage a New Relic organization.
---

# Resource: newrelic\_organization

Use this resource to create and manage New Relic organizations.

~> **NOTE:** Deleting an organization is not supported via the API. Removing this resource from Terraform configuration will only remove it from state — the organization itself will not be deleted in New Relic.

## Example Usage

```hcl
resource "newrelic_organization" "foo" {
  name = "My New Organization"
}
```

### With a New Managed Account

```hcl
resource "newrelic_organization" "foo" {
  name = "My New Organization"

  new_managed_account {
    name        = "My New Account"
    region_code = "US01"
  }
}
```

### With a Shared Account

```hcl
resource "newrelic_organization" "foo" {
  name = "My New Organization"

  shared_account {
    account_id       = 12345678
    limiting_role_id = 1234
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the organization.
* `customer_id` - (Optional) The customer ID to associate with the new organization. This value cannot be changed after creation.
* `new_managed_account` - (Optional) A nested block that describes a new managed account to create within the organization. Only one `new_managed_account` block is permitted. This value cannot be changed after creation. See [Nested new_managed_account blocks](#nested-new_managed_account-blocks) below for details.
* `shared_account` - (Optional) A nested block that describes an account to share with the new organization. Only one `shared_account` block is permitted. This value cannot be changed after creation. See [Nested shared_account blocks](#nested-shared_account-blocks) below for details.

### Nested `new_managed_account` blocks

* `name` - (Optional) The name of the new account to be created. This value cannot be changed after creation.
* `region_code` - (Optional) The region code for the account. One of: `EU01`, `US01`. This value cannot be changed after creation.

### Nested `shared_account` blocks

* `account_id` - (Required) The ID of the account to share with the new organization. This value cannot be changed after creation.
* `limiting_role_id` - (Optional) The limiting role ID the new organization will be granted for the shared account. This value cannot be changed after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the organization creation job.
* `job_id` - The job ID of the organization creation task.

## Import

New Relic organizations can be imported using the organization ID:

```
terraform import newrelic_organization.foo <organization_id>
```

~> **NOTE:** Sensitive data and certain read-only fields may not be returned from the underlying API and may not be set in state when importing.