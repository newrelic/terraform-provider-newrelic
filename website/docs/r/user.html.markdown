---
layout: 'newrelic'
page_title: 'New Relic: newrelic_user'
sidebar_current: 'docs-newrelic-resource-user'
description: |-
  Create and manage users in New Relic.
---

# Resource: newrelic\_user

The `newrelic_user` resource may be used to create, update and delete users in New Relic.

## Example Usage
```hcl
data "newrelic_authentication_domain" "foo" {
  name = "Test Authentication Domain"
}

resource "newrelic_user_management" "foo" {
  name                     = "Test New User"
  email                    = "test_user@test.com"
  authentication_domain_id = data.newrelic_authentication_domain.foo.id
  user_type                = "CORE_USER_TIER"
}
```

## Argument Reference
The following arguments are supported:

* `name` - (Required) The name of the user to be created.
* `email` - (Required) The email ID of the user to be created.
* `authentication_domain_id` - (Required) The ID of the authentication domain to which the user to be created would belong.
* `user_type` - (Optional) The tier to which the user to be created would belong. Accepted values for this argument are `BASIC_USER_TIER`, `CORE_USER_TIER`, or `FULL_USER_TIER`. If not specified in the configuration, the argument would default to `BASIC_USER_TIER`.

## Attributes Reference
In addition to the attributes listed above, the following attribute is also exported by this resource:

* `id` - The ID of the created user.

## Import
A user can be imported using its ID. Example:

```shell
$ terraform import newrelic_user.foo 1999999999
```