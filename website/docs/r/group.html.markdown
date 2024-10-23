---
layout: 'newrelic'
page_title: 'New Relic: newrelic_group'
sidebar_current: 'docs-newrelic-resource-group'
description: |-
  Create and manage groups in New Relic.
---

# Resource: newrelic\_group

The `newrelic_group` resource facilitates creating, updating, and deleting groups in New Relic, while also enabling the addition and removal of users from these groups.

## Example Usage

```hcl
data "newrelic_authentication_domain" "foo" {
  name = "Test Authentication Domain"
}

resource "newrelic_group" "foo" {
  name                     = "Test Group"
  authentication_domain_id = data.newrelic_authentication_domain.foo.id
  user_ids                 = ["0001112222", "2221110000"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the group to be created.
* `authentication_domain_id` - (Required) The ID of the authentication domain to which the group to be created would belong.
* `user_ids` - (Optional) A list of IDs of users to be included in the group to be created.

-> **NOTE** The ID of an authentication domain can be retrieved using its name, via the data source `newrelic_authentication_domain`, as shown in the example above. Head over to the documentation of this data source for more details and examples.

~> **WARNING:** Changing the `authentication_domain_id` of a `newrelic_group` resource that has already been applied would result in a **replacement** of the resource – destruction of the existing resource, followed by the addition of a new resource with the specified configuration. This is due to the fact that updating the `authentication_domain_id` of an existing group is not supported.

## Attributes Reference

In addition to the attributes listed above, the following attribute is also exported by this resource:

* `id` - The ID of the created group.

## Additional Examples

### Updating User Group Membership Management in Terraform

### Overview
There is a potential race condition within Terraform when managing user accounts and their respective group memberships. A user might be deleted before Terraform disassociates them from a user group. This can lead to an error during terraform apply because the user ID no longer exists when the group resource is being updated.

### Recommended Solution
To address this and ensure proper sequential execution of resource updates, it is recommended to utilize the `create_before_destroy` lifecycle directive within your user group resource definition.

### Implementing Lifecycle Changes
To implement the change, modify the user group resource in your Terraform configuration as follows:

```hcl
resource "newrelic_group" "viewer" {
#  Existing configuration ...

  lifecycle {
     create_before_destroy = true
  }
}
```
The `create_before_destroy = true` statement will ensure that Terraform updates the user group (e.g., removes the user) before attempting to destroy the user resource, thus preventing the error.

### Addition of New Users to a New Group

The following example illustrates the creation of a group using the `newrelic_group` resource, to which users created using the `newrelic_user` resource are added.

```hcl
data "newrelic_authentication_domain" "foo" {
  name = "Test Authentication Domain"
}

resource "newrelic_user" "foo" {
  name                     = "Test User One"
  email_id                 = "test_user_one@test.com"
  authentication_domain_id = data.newrelic_authentication_domain.foo.id
  user_type                = "CORE_USER_TIER"
}

resource "newrelic_user" "bar" {
  name                     = "Test User Two"
  email_id                 = "test_user_two@test.com"
  authentication_domain_id = data.newrelic_authentication_domain.foo.id
  user_type                = "BASIC_USER_TIER"
}

resource "newrelic_group" "foo" {
  name                     = "Test Group"
  authentication_domain_id = data.newrelic_authentication_domain.foo.id
  user_ids = [
    newrelic_user.foo.id,
    newrelic_user.bar.id,
  ]
}
```

### Addition of Existing Users to a New Group

The following example demonstrates the usage of the `newrelic_group` resource to create a group, wherein the `newrelic_user` data source is employed to associate existing users with the newly formed group.

```hcl
data "newrelic_authentication_domain" "foo" {
  name = "Test Authentication Domain"
}

data "newrelic_user" "foo" {
  authentication_domain_id = data.newrelic_authentication_domain.foo.id
  email_id                 = "test_user_one@test.com"
}

data "newrelic_user" "bar" {
  authentication_domain_id = data.newrelic_authentication_domain.foo.id
  name                     = "Test User Two"
}

resource "newrelic_group" "foo" {
  name                     = "Test Group"
  authentication_domain_id = data.newrelic_authentication_domain.foo.id
  user_ids = [
    data.newrelic_user.foo.id,
    data.newrelic_user.bar.id,
  ]
}
```

-> **NOTE** Please note that the addition of users to groups is only possible when both the group and the users to be added to it belong to the _same authentication domain_. If the group being created and the users being added to it belong to different authentication domains, an error indicating `user not found` or an equivalent error will be thrown.

## Import

A group can be imported using its ID. Example:

```shell
$ terraform import newrelic_group.foo <group_id>
```